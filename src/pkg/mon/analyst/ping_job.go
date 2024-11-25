package analyst

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cspage/pkg/data"
	"cspage/pkg/pb"
	"cspage/pkg/worker"
)

type pingJob struct {
	analystJob

	durationBuffer time.Duration
	alertTrigger   string
}

func (j *pingJob) Do(ctx context.Context, tick worker.Tick) {
	j.analystJob.do(ctx, tick, j.analyze)
}

func (j *pingJob) analyze(ctx context.Context, checkpoint time.Time) (time.Time, *data.Notification, bool) {
	ok := true

	openAlerts, errr := j.analystJob.getOpenAlerts(ctx, pb.AlertType_PING, j.probe.Name)
	if errr != nil {
		j.taskLog.Error("Could not fetch open alerts", "err", errr)
		return checkpoint, nil, false
	}

	notifyAlerts, err := j.analyzeMissingPings(ctx, checkpoint, openAlerts)
	if err != nil {
		j.taskLog.Error("Could not analyze missing pings", "err", err)
		ok = false
	}

	return checkpoint, notifyAlerts.ToNotification(), ok
}

func (j *pingJob) analyzeMissingPings(
	ctx context.Context,
	checkpoint time.Time,
	openAlerts map[string]*data.Alert,
) (data.Alerts, error) {
	var notifyAlerts data.Alerts
	var errs []error

	missingPings, errr := j.getMissingPings(ctx, checkpoint)
	if errr != nil {
		return notifyAlerts, fmt.Errorf("failed to fetch missing pings: %w", errr)
	}

	//nolint:varnamelen // Using mp is OK in this context.
	for _, mp := range missingPings {
		j.taskLog.Debug("Missing ping analysis", "missing_ping", mp)
		notifyAlert, err := j.evaluateMissingPing(ctx, checkpoint, openAlerts, mp)
		if notifyAlert != nil {
			notifyAlerts = append(notifyAlerts, notifyAlert)
		}
		if err != nil {
			j.taskLog.Error("Could not evaluate missing pings", "missing_ping", mp, "err", err)
			errs = append(errs, err)
		}
	}

	// OpenAlerts should be empty by now, if not, we need to close each alert manually.
	// This can happen when a missing ping is delivered late or the missing ping window
	// is shorter than 2 x PingLagThreshold.
	for _, alert := range openAlerts {
		// Let's close alerts this way only after PingLagThreshold has passed since last check
		if checkpoint.Sub(alert.TimeCheck) < j.cfg.PingLagThreshold {
			continue
		}

		notifyAlert, err := j.closeOpenAlert(ctx, checkpoint, nil, alert)
		if notifyAlert != nil {
			notifyAlerts = append(notifyAlerts, notifyAlert)
		}
		if err != nil {
			j.taskLog.Error("Could not close missing pings alert", "alert", alert, "err", err)
			errs = append(errs, err)
		}
	}

	return notifyAlerts, errors.Join(errs...)
}

//nolint:varnamelen // Using mp is OK in this context.
func (j *pingJob) evaluateMissingPing(
	ctx context.Context,
	checkpoint time.Time,
	openAlerts map[string]*data.Alert,
	mp *missingPing,
) (*data.Alert, error) {
	if mp.JobId == nil {
		// Very old missing ping that couldn't be paired to job_id as we don't see behind the edge.
		// The only thing we can do is to match it with an open alert.
		for _, alert := range openAlerts {
			if alert.CloudRegion == mp.CloudRegion && alert.TimeBegin.Before(mp.TimeBegin) {
				mp.JobId = &alert.JobId
				mp.AlertId = &alert.Id
				break
			}
		}
	}

	if mp.JobId == nil {
		j.taskLog.Debug("Missing ping with NULL job ID => skipping", "missing_ping", mp)
		//nolint:nilnil // The caller can handle this.
		return nil, nil // Still no job ID, can't work with this missing ping
	}

	var notifyAlert *data.Alert
	var err error

	if mp.AlertId == nil {
		notifyAlert, err = j.createPingAlert(ctx, checkpoint, mp)
	} else if alert, exists := openAlerts[*mp.AlertId]; exists { // Existing open alert
		if checkpoint.After(mp.TimeEnd) {
			notifyAlert, err = j.closeOpenAlert(ctx, checkpoint, mp, alert)
		} else {
			err = j.updateOpenAlert(ctx, checkpoint, alert)
		}
		delete(openAlerts, alert.Id)
	} // Else existing closed (historical) alert => ignore

	return notifyAlert, err
}

//nolint:varnamelen // Using mp is OK in this context.
func (j *pingJob) createPingAlert(
	ctx context.Context,
	checkpoint time.Time,
	mp *missingPing,
) (*data.Alert, error) {
	// New missing ping
	alert := &data.Alert{
		JobId:       *mp.JobId,
		TimeBegin:   mp.TimeBegin,
		TimeEnd:     nil,
		TimeCheck:   checkpoint,
		Type:        pb.AlertType_PING_MISSING,
		Status:      pb.AlertStatus_ALERT_OPEN,
		CloudRegion: mp.CloudRegion,
		ProbeName:   j.probe.Name,
		ProbeAction: 0,
		Data: &data.AlertData{
			ProbeDescription: j.probe.Description,
			Trigger:          j.alertTrigger,
			ServiceName:      j.probe.Config.ServiceName,
			ServiceGroup:     j.probe.Config.ServiceGroup,
		},
	}

	if checkpoint.After(mp.TimeEnd) {
		// Older issue => we can backfill it
		alert.Status = pb.AlertStatus_ALERT_CLOSED_AUTO
		alert.TimeEnd = &mp.End
	}

	j.taskLog.Warn("NEW alert for missing pings", "alert", alert)
	err := j.createAlert(ctx, alert)
	if err != nil {
		err = fmt.Errorf("failed to create alert for missing pings: %w", err)
	}

	return alert, err
}

//nolint:varnamelen // Using mp is OK in this context.
func (j *pingJob) closeOpenAlert(
	ctx context.Context,
	checkpoint time.Time,
	mp *missingPing,
	alert *data.Alert,
) (*data.Alert, error) {
	firstPingSince, err := j.getFirstPingSince(ctx, alert.TimeBegin, alert.CloudRegion)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch last ping in %s since %s: %w", alert.CloudRegion, alert.TimeBegin, err)
	}

	if firstPingSince.IsZero() {
		j.taskLog.Debug("Alert for missing pings is still active", "alert", alert)
		//nolint:nilnil // The caller can handle this.
		return nil, nil
	}

	if mp == nil {
		alert.TimeEnd = &firstPingSince
	} else {
		alert.TimeEnd = &mp.End
		if mp.End != firstPingSince {
			j.taskLog.Error(
				"Unexpected difference between missing pings' end time and first ping since outage",
				"alert", alert,
				"missing_ping", mp,
				"first_ping", firstPingSince,
			)
		}
	}

	alert.Status = pb.AlertStatus_ALERT_CLOSED_AUTO
	alert.TimeCheck = checkpoint
	j.taskLog.Warn("Closing alert for missing pings", "alert", alert)

	err = j.updateAlert(ctx, alert, true)
	if err != nil {
		err = fmt.Errorf("failed to close alert for missing pings: %w", err)
	}

	return alert, err
}

func (j *pingJob) updateOpenAlert(ctx context.Context, checkpoint time.Time, alert *data.Alert) error {
	alert.TimeCheck = checkpoint
	if err := j.updateAlert(ctx, alert, false); err != nil {
		return fmt.Errorf("failed to update alert for missing pings: %w", err)
	}
	return nil
}
