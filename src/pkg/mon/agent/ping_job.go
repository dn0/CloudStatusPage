package agent

import (
	"context"
	"log/slog"

	"cspage/pkg/msg"
	"cspage/pkg/pb"
	"cspage/pkg/sysinfo"
	"cspage/pkg/worker"
)

const (
	pingJobName = worker.TickerJobPrefix + pb.JobNamePing
)

type sysInfo struct {
	Host *sysinfo.OSInfo `json:"host"`
	VM   *vm             `json:"vm"`
}

type pingJob struct {
	cfg       *Config
	publisher msg.Publisher
}

func (j *pingJob) String() string {
	return pingJobName
}

func (j *pingJob) PreStart(ctx context.Context) {
	if err := j.publishAgentStartMessage(ctx); err != nil {
		Die("Could not publish AGENT_START message", "err", err)
	}
}

func (j *pingJob) Start(_ context.Context) {}

func (j *pingJob) Do(ctx context.Context, tick worker.Tick) {
	if ctx.Err() != nil {
		return // Task was interrupted => no ping
	}
	//nolint:contextcheck // Publish & forget => not using task context.
	j.publishPingMessage(tick)
}

func (j *pingJob) Stop(ctx context.Context) {
	defer j.publisher.Close()

	if err := j.publishAgentStopMessage(ctx); err != nil {
		slog.Error("Could not publish AGENT_STOP message", "err", err)
	}
}

func (j *pingJob) Shutdown(_ error) {
	j.publishAgentStoppingMessage()
}

func (j *pingJob) publishAgentStartMessage(ctx context.Context) error {
	//nolint:contextcheck // Sysinfo is using a separate context.
	message := newAgentMessage(&j.cfg.Env, pb.AgentAction_AGENT_START, &sysInfo{
		Host: sysinfo.GetOSInfo(context.Background()),
		VM:   &j.cfg.Env.VM,
	})
	//nolint:wrapcheck // This error is properly logged above when this method is called.
	return j.publisher.PublishWait(ctx, msg.NewMessage(message, msg.NewAttrs(
		msg.TypeAgent,
		j.cfg.Env.Cloud,
		j.cfg.Env.Region,
	)))
}

func (j *pingJob) publishAgentStopMessage(ctx context.Context) error {
	message := newAgentMessage(&j.cfg.Env, pb.AgentAction_AGENT_STOP, nil)
	//nolint:wrapcheck // This error is properly logged above when this method is called.
	return j.publisher.PublishWait(ctx, msg.NewMessage(message, msg.NewAttrs(
		msg.TypeAgent,
		j.cfg.Env.Cloud,
		j.cfg.Env.Region,
	)))
}

func (j *pingJob) publishAgentStoppingMessage() {
	message := newAgentMessage(&j.cfg.Env, pb.AgentAction_AGENT_STOPPING, nil)
	// NOTE: this means publish & forget as we are not waiting for the result
	//       therefore we can't use the task context here
	j.publisher.Publish(context.Background(), msg.NewMessage(message, msg.NewAttrs(
		msg.TypeAgent,
		j.cfg.Env.Cloud,
		j.cfg.Env.Region,
	)))
}

func (j *pingJob) publishPingMessage(tick worker.Tick) {
	message := &pb.Ping{
		Job:     newJobMessage(&j.cfg.Env, &tick, pb.JobNamePing),
		Sysstat: sysinfo.GetSysStat(context.Background()),
	}
	// NOTE: this means publish & forget as we are not waiting for the result
	//       therefore we can't use the task context here
	j.publisher.Publish(context.Background(), msg.NewMessage(message, msg.NewAttrs(
		msg.TypePing,
		j.cfg.Env.Cloud,
		j.cfg.Env.Region,
	)))
}
