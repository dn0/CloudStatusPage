package analyst

import (
	"context"
	"errors"
	"time"

	"cspage/pkg/data"
	"cspage/pkg/worker"
)

type incidentJob struct {
	analystJob
}

func (j *incidentJob) Do(ctx context.Context, tick worker.Tick) {
	j.analystJob.do(ctx, tick, j.analyze)
}

func (j *incidentJob) analyze(ctx context.Context, checkpoint time.Time) (time.Time, *data.Notification, bool) {
	ok := true

	recentAlerts, errr := j.getRecentAlerts(ctx)
	if errr != nil {
		j.taskLog.Error("Could not fetch recent alerts", "err", errr)
		return checkpoint, nil, false
	}

	notifyIncidents, err := j.analyzeRecentAlerts(ctx, checkpoint, recentAlerts)
	if err != nil {
		j.taskLog.Error("Could not analyze recent alerts", "err", err)
		ok = false
	}

	return checkpoint, notifyIncidents.ToNotification(), ok
}

func (j *incidentJob) analyzeRecentAlerts(
	ctx context.Context,
	checkpoint time.Time,
	recentAlerts []*data.ExtendedAlert,
) (data.Incidents, error) {
	var notifyIncidents data.Incidents
	var errs []error
	clusters := data.ClusterAlerts(j.cloud, recentAlerts, data.ClusterMaxDistance)

	for _, cluster := range clusters {
		if len(cluster.Alerts) <= 1 {
			continue
		}

		var inc *data.Incident
		var msg string
		if cluster.Incident.Id == "" {
			msg = "NEW incident for alerts"
			inc = &data.Incident{}
		} else if cluster.Incident.Outdated {
			msg = "Updating incident for alerts"
			var err error
			if inc, err = data.GetIncident(ctx, j.dbc.Read, j.cloud, cluster.Incident.Id); err != nil {
				j.taskLog.Error("Could not fetch incident", "incident", inc, "err", err)
				errs = append(errs, err)
				continue
			}
		}

		if inc == nil {
			continue
		}

		err := inc.CreateOrUpdateFromAlerts(ctx, j.dbtx, j.cloud, checkpoint, cluster.Alerts)
		if err == nil {
			j.taskLog.Warn(msg, "incident", inc, "alert_cluster", cluster)
			notifyIncidents = append(notifyIncidents, inc)
		} else {
			j.taskLog.Error("Could not "+msg, "incident", inc, "alert_cluster", cluster, "err", err)
			errs = append(errs, err)
		}
	}

	return notifyIncidents, errors.Join(errs...)
}
