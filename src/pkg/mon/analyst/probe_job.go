package analyst

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"time"

	"cspage/pkg/data"
	"cspage/pkg/pb"
	"cspage/pkg/worker"
)

const (
	probeMaxFramesPerTask          = 10
	probeMaxFrameSizeDiv           = 4
	probeMinStates          uint32 = 3
	probeAlertTriggerSlow          = "zscore(probe.latency)=%.5f > %.2f"
	probeAlertTriggerFailed        = "probe.status=%s"
)

var (
	errFirstSuccessProbeIsNil = errors.New("firstSuccessProbe is nil")
	errFirstFailureProbeIsNil = errors.New("firstFailureProbe is nil")
	errFirstHighProbeIsNil    = errors.New("firstHighProbe is nil")
	errGetProbeAction         = errors.New("failed to get probe action")
)

type probeState struct {
	key   string // key in probeStates
	count uint32
	frame time.Time

	failedProbeAlert  *data.Alert
	failureCount      uint32
	successCount      uint32
	firstFailureProbe *probeResult
	firstSuccessProbe *probeResult

	slowProbeAlert     *data.Alert
	slowProbeThreshold time.Duration
	highCount          uint32
	lowCount           uint32
	firstHighProbe     *probeResult
	firstLowProbe      *probeResult
}

type probeStates map[string]*probeState

type probeJob struct {
	analystJob

	checkWindow        time.Duration
	failureThreshold   uint32
	zscoreWindow       time.Duration
	zscoreThresholdOn  float32
	zscoreThresholdOff float32
	distanceFromAvgDiv time.Duration
	distanceFromAvgMin time.Duration
}

func (j *probeJob) Do(ctx context.Context, tick worker.Tick) {
	j.analystJob.do(ctx, tick, j.analyze)
}

func (j *probeJob) analyze(ctx context.Context, checkpoint time.Time) (time.Time, *data.Notification, bool) {
	frames := []time.Time{checkpoint}

	for checkpoint.Sub(j.checkpoint) > j.checkWindow/probeMaxFrameSizeDiv {
		checkpoint = checkpoint.Add(-j.checkWindow / probeMaxFrameSizeDiv)
		frames = append(frames, checkpoint)
	}

	totalFrames := len(frames)
	if totalFrames > 1 {
		slices.Reverse(frames)
		if totalFrames > probeMaxFramesPerTask {
			frames = frames[:probeMaxFramesPerTask]
		}
		checkpoint = frames[len(frames)-1]
		j.taskLog.Info(
			"Probe analysis has some catching up to do",
			"frames", frames,
			"total", totalFrames,
			"last_checkpoint", checkpoint,
		)
	}

	var notifyAlerts data.Alerts
	var ok bool

	for _, frame := range frames {
		var alerts data.Alerts
		alerts, ok = j.analyzeFrame(ctx, frame, notifyAlerts)
		notifyAlerts = append(notifyAlerts, alerts...)
		if !ok {
			break
		}
	}

	//goland:noinspection GoDfaNilDereference
	return checkpoint, notifyAlerts.ToNotification(), ok
}

func (j *probeJob) analyzeFrame(ctx context.Context, checkpoint time.Time, newAlerts data.Alerts) (data.Alerts, bool) {
	probeResults, err := j.getProbeResults(ctx, checkpoint)
	if err != nil {
		j.taskLog.Error("Could not fetch probe results", "err", err)
		return nil, false
	}

	if len(probeResults) == 0 {
		j.taskLog.Debug("No new probe results")
		return nil, true
	}

	openAlerts, err := j.analystJob.getOpenAlerts(ctx, pb.AlertType_PROBE, j.probe.Name)
	if err != nil {
		j.taskLog.Error("Could not fetch open alerts", "err", err)
		return nil, false
	}
	// Append new & uncommitted alerts from previous frame so that we don't open the same alert twice
	for _, alert := range newAlerts {
		openAlerts[alert.Id] = alert
	}

	j.taskLog.Debug("Probe analysis", "results", len(probeResults), "alerts", len(openAlerts))
	notifyAlerts, err := j.analyzeProbe(ctx, checkpoint, probeResults, openAlerts)
	if err != nil {
		j.taskLog.Error("Could not analyze probe", "err", err)
	}

	return notifyAlerts, err == nil
}

func (j *probeJob) analyzeProbe(
	ctx context.Context,
	checkpoint time.Time,
	probeResults []*probeResult,
	openAlerts map[string]*data.Alert,
) (data.Alerts, error) {
	var notifyAlerts data.Alerts
	var errs []error
	pstates := probeStates{}

	// Associate open alerts with probe states
	for _, alert := range openAlerts {
		state := pstates.getState(checkpoint, alert.CloudRegion, alert.ProbeAction)
		state.addAlert(alert, j.zscoreThresholdOff)
	}

	// Go through each row = p(robe) and update state
	for _, pr := range probeResults {
		state := pstates.getState(checkpoint, pr.CloudRegion, pr.Action)
		state.set(pr, j.zscoreThresholdOn, j.distanceFromAvgDiv, j.distanceFromAvgMin)
	}

	// Go through every state and evaluate => open or close alerts
	for _, state := range pstates {
		if state.count < probeMinStates {
			j.taskLog.With(state.logAttrs()).Debug("Skipping probe state analysis")
			continue
		}

		alert, err := j.evaluateProbeState(ctx, checkpoint, state)
		if alert != nil {
			notifyAlerts = append(notifyAlerts, alert)
		}
		if err != nil {
			j.taskLog.With(state.logAttrs()).Error("Could not evaluate probe state", "err", err)
			errs = append(errs, err)
		}
	}

	return notifyAlerts, errors.Join(errs...)
}

//nolint:funlen,gocognit,cyclop,nilnil,lll // This has to be complex.
func (j *probeJob) evaluateProbeState(ctx context.Context, checkpoint time.Time, state *probeState) (*data.Alert, error) {
	var notifyAlert *data.Alert
	var err error

	switch {
	case state.failureCount >= j.failureThreshold && state.failedProbeAlert != nil: // Existing alert for failed probe
		if notifyAlert, err = j.updateOpenAlert(ctx, checkpoint, state, state.firstSuccessProbe, state.failedProbeAlert); err != nil {
			err = fmt.Errorf("%s: failed to update failed probe alert: %w", state.key, err)
		}
	case state.failureCount >= j.failureThreshold && state.failedProbeAlert == nil:
		if state.firstFailureProbe == nil {
			err = fmt.Errorf("%s: unexpected state: %w", state.key, errFirstFailureProbeIsNil)
			return nil, err
		}
		if state.firstFailureProbe.AlertId != "" {
			return nil, nil
		}
		if j.isSilentTime(ctx, checkpoint, state, state.firstFailureProbe) {
			return nil, nil
		}
		if state.failedProbeAlert, err = j.createFailedProbeAlert(ctx, checkpoint, state, state.firstFailureProbe); err != nil {
			err = fmt.Errorf("%s: failed to create failed probe alert: %w", state.key, err)
		}
		notifyAlert = state.failedProbeAlert
	case state.failureCount < j.failureThreshold && state.failedProbeAlert != nil:
		if state.firstSuccessProbe == nil {
			err = fmt.Errorf("%s: failed to close failed probe alert: %w", state.key, errFirstSuccessProbeIsNil)
			return nil, err
		}
		if state.failedProbeAlert.TimeBegin.After(state.firstSuccessProbe.JobTime) {
			return nil, nil
		}
		if err = j.closeProbeAlert(ctx, checkpoint, state, state.firstSuccessProbe, state.failedProbeAlert); err != nil {
			err = fmt.Errorf("%s: failed to close failed probe alert: %w", state.key, err)
		}
		notifyAlert = state.failedProbeAlert
		state.failedProbeAlert = nil
	case state.highCount >= state.lowCount && state.slowProbeAlert != nil: // existing alert for slow probe
		if notifyAlert, err = j.updateOpenAlert(ctx, checkpoint, state, state.firstLowProbe, state.slowProbeAlert); err != nil {
			err = fmt.Errorf("%s: failed to update slow probe alert: %w", state.key, err)
		}
	case state.highCount > state.lowCount && state.slowProbeAlert == nil:
		if state.firstHighProbe == nil {
			err = fmt.Errorf("%s: unexpected state: %w", state.key, errFirstHighProbeIsNil)
			return nil, err
		}
		if state.firstHighProbe.AlertId != "" {
			return nil, nil
		}
		if j.isSilentTime(ctx, checkpoint, state, state.firstHighProbe) {
			return nil, nil
		}
		if state.slowProbeAlert, err = j.createSlowProbeAlert(ctx, checkpoint, state, state.firstHighProbe); err != nil {
			err = fmt.Errorf("%s: failed to create slow probe alert: %w", state.key, err)
		}
		notifyAlert = state.slowProbeAlert
	case state.highCount < state.lowCount && state.slowProbeAlert != nil:
		if state.firstLowProbe == nil {
			// can't close the alert if last probe in this frame is not low; let's wait for more probe results
			return nil, nil
		}
		if state.slowProbeAlert.TimeBegin.After(state.firstLowProbe.JobTime) {
			return nil, nil
		}
		if err = j.closeProbeAlert(ctx, checkpoint, state, state.firstLowProbe, state.slowProbeAlert); err != nil {
			err = fmt.Errorf("%s: failed to close slow probe alert: %w", state.key, err)
		}
		notifyAlert = state.slowProbeAlert
		state.slowProbeAlert = nil
	}

	return notifyAlert, err
}

// isSilentTime = true means that we should not create new alerts.
func (j *probeJob) isSilentTime(
	ctx context.Context,
	checkpoint time.Time,
	state *probeState,
	firstKOProbe *probeResult,
) bool {
	// Very artificial check that should not be here. It exists only because of Azure :(
	agentUptime, err := j.getJobAgentUptime(ctx, checkpoint, firstKOProbe.JobId)
	if err != nil {
		j.taskLog.With(state.logAttrs()).Error(
			"Unexpected situation! Could not get job agent uptime",
			"job_id", firstKOProbe.JobId,
			"err", err,
		)
		return true
	}
	silent := agentUptime < j.cfg.ProbeAlertSilenceAfterAgentStart

	if silent {
		j.taskLog.With(state.logAttrs()).Debug(
			"No new alerts because the relevant agent is too fresh",
			"job_id", firstKOProbe.JobId,
			"agent_uptime", agentUptime,
		)
	}
	return silent
}

func (j *probeJob) updateOpenAlert(
	ctx context.Context,
	checkpoint time.Time,
	state *probeState,
	firstOKProbe *probeResult,
	alert *data.Alert,
) (*data.Alert, error) {
	if alert.TimeCheck.After(checkpoint) || alert.TimeCheck == checkpoint {
		//nolint:nilnil // This is OK.
		return nil, nil
	}

	alert.TimeCheck = checkpoint

	// This represents nil. It's also important to set the alert's time_end back to nil if the issue reappeared.
	timeEnd := time.Time{}
	// The alert is still open, but we'll set its time_end to first good probe if it's after alert's time_begin,
	// which should be somehow reflected in the UI (i.e. the alert should be shown as open)
	if firstOKProbe != nil && firstOKProbe.JobTime.After(alert.TimeBegin) {
		timeEnd = firstOKProbe.JobTime
	}

	alertTimeEnd := time.Time{}
	if alert.TimeEnd != nil {
		alertTimeEnd = *alert.TimeEnd
	}

	if timeEnd != alertTimeEnd {
		if timeEnd.IsZero() {
			alert.TimeEnd = nil
		} else {
			alert.TimeEnd = &timeEnd
		}
		j.taskLog.With(state.logAttrs()).Debug("Updating open alert's time_end", "alert", alert)
		// We should also notify all consumers
		return alert, j.updateAlert(ctx, alert, true)
	}

	return nil, j.updateAlert(ctx, alert, false)
}

func (j *probeJob) createFailedProbeAlert(
	ctx context.Context,
	checkpoint time.Time,
	state *probeState,
	firstFailedProbe *probeResult,
) (*data.Alert, error) {
	alert := &data.Alert{
		Type: pb.AlertType(firstFailedProbe.Status),
		Data: &data.AlertData{
			Trigger:           fmt.Sprintf(probeAlertTriggerFailed, pb.AlertType_name[int32(firstFailedProbe.Status)]),
			ProbeError:        firstFailedProbe.Error,
			ProbeLatencyValue: firstFailedProbe.Latency,
		},
	}
	return alert, j.createProbeAlert(ctx, checkpoint, state, firstFailedProbe, alert)
}

func (j *probeJob) createSlowProbeAlert(
	ctx context.Context,
	checkpoint time.Time,
	state *probeState,
	firstSlowProbe *probeResult,
) (*data.Alert, error) {
	alert := &data.Alert{
		Type: pb.AlertType_PROBE_SLOW,
		Data: &data.AlertData{
			Trigger:            fmt.Sprintf(probeAlertTriggerSlow, firstSlowProbe.Zscore, j.zscoreThresholdOn),
			ProbeLatencyValue:  firstSlowProbe.Latency,
			ProbeLatencyAvg:    firstSlowProbe.LatencyAvg,
			ProbeLatencySD:     firstSlowProbe.LatencySD,
			ProbeLatencyZscore: firstSlowProbe.Zscore,
		},
	}
	return alert, j.createProbeAlert(ctx, checkpoint, state, firstSlowProbe, alert)
}

func (j *probeJob) createProbeAlert(
	ctx context.Context,
	checkpoint time.Time,
	state *probeState,
	firstKOProbe *probeResult,
	alert *data.Alert,
) error {
	alert.JobId = firstKOProbe.JobId
	alert.TimeBegin = firstKOProbe.JobTime
	alert.TimeEnd = nil
	alert.TimeCheck = checkpoint
	alert.Status = pb.AlertStatus_ALERT_OPEN
	alert.CloudRegion = firstKOProbe.CloudRegion
	alert.ProbeAction = firstKOProbe.Action
	alert.ProbeName = j.probe.Name
	alert.Data.ServiceName = j.probe.Config.ServiceName
	alert.Data.ServiceGroup = j.probe.Config.ServiceGroup
	alert.Data.ProbeDescription = j.probe.Description
	action, ok := j.probe.Config.ActionGet(alert.ProbeAction)
	if !ok {
		return fmt.Errorf("%w %d", errGetProbeAction, alert.ProbeAction)
	}
	alert.Data.ProbeActionName = action.Name
	alert.Data.ProbeActionTitle = action.Title()
	j.taskLog.With(state.logAttrs()).Warn("NEW alert for probe", "alert", alert)
	return j.createAlert(ctx, alert)
}

func (j *probeJob) closeProbeAlert(
	ctx context.Context,
	checkpoint time.Time,
	state *probeState,
	firstOKProbe *probeResult,
	alert *data.Alert,
) error {
	alert.Status = pb.AlertStatus_ALERT_CLOSED_AUTO
	alert.TimeCheck = checkpoint
	alert.TimeEnd = &firstOKProbe.JobTime
	j.taskLog.With(state.logAttrs()).Warn("Closing alert for probe", "alert", alert)
	return j.updateAlert(ctx, alert, true)
}

func (s *probeState) logAttrs() slog.Attr {
	return slog.Group(
		"probe_state",
		"name", s.key,
		"frame", s.frame,
		"count", s.count,
		"success_count", s.successCount,
		"failure_count", s.failureCount,
		"low_count", s.lowCount,
		"high_count", s.highCount,
		"slow_probe_threshold", s.slowProbeThreshold,
	)
}

func (s *probeState) addAlert(alert *data.Alert, zscoreThresholdOff float32) {
	//nolint:exhaustive // Ping alert type is used by ping analyst.
	switch alert.Type {
	case pb.AlertType_PROBE_FAILURE, pb.AlertType_PROBE_TIMEOUT:
		s.failedProbeAlert = alert
	case pb.AlertType_PROBE_SLOW:
		s.slowProbeAlert = alert
		s.slowProbeThreshold = alert.Data.ProbeLatencyThreshold(float64(zscoreThresholdOff))
	}
}

//nolint:varnamelen,cyclop // Using pr(probeResult) is OK in this context.
func (s *probeState) set(pr *probeResult, zscoreThresholdOn float32, distanceFromAvgDiv, distanceFromAvgMin time.Duration) {
	s.count++

	switch pr.Status {
	case pb.ResultStatus_RESULT_UNKNOWN:
		return
	case pb.ResultStatus_RESULT_FAILURE, pb.ResultStatus_RESULT_TIMEOUT:
		s.failureCount++
		s.firstSuccessProbe = nil // Moves first success probe after last failed probe
		if s.firstFailureProbe == nil {
			s.firstFailureProbe = pr
		}
	case pb.ResultStatus_RESULT_SUCCESS:
		s.successCount++
		if s.firstSuccessProbe == nil {
			s.firstSuccessProbe = pr
		}

		// SlowProbeThreshold is either a large constant or zscoreOff * stddev + average latency of
		// probeResult that opened the currently active alert; This condition is here because
		// zscore can become lower after an alert is open for a longer time.
		// The second condition is the main one for detecting new anomalies; Besides Zscore
		// we also check whether the current probeResult distance from mean is greater than an absolute value and
		// grater the mean latency divided by a tunable parameter (usually 3 which gives us 1/3 of the mean latency),
		if pr.Latency > s.slowProbeThreshold ||
			(pr.Zscore >= zscoreThresholdOn &&
				pr.Distance >= distanceFromAvgMin &&
				pr.Distance >= pr.LatencyAvg/distanceFromAvgDiv) {
			s.highCount++
			s.firstLowProbe = nil // Moves first "good" probe after last "slow" probe
			if s.firstHighProbe == nil {
				s.firstHighProbe = pr
			}
		} else {
			s.lowCount++
			if s.firstLowProbe == nil {
				s.firstLowProbe = pr
			}
		}
	}
}

func (p probeStates) getState(frame time.Time, cloudRegion string, probeAction uint32) *probeState {
	key := cloudRegion + "|" + strconv.Itoa(int(probeAction))
	state, exists := p[key]
	if !exists {
		state = &probeState{
			key:                key,
			frame:              frame,
			slowProbeThreshold: time.Hour,
		}
		p[key] = state
	}
	return state
}
