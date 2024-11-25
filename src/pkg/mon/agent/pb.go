package agent

import (
	"encoding/json"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"cspage/pkg/pb"
	"cspage/pkg/worker"
)

func newAgentMessage(env *env, action pb.AgentAction, sysInfo *sysInfo) *pb.Agent {
	var sysInfoJSON []byte
	ipAddress := ""
	if sysInfo != nil {
		var err error
		sysInfoJSON, err = json.Marshal(sysInfo)
		if err != nil {
			Die("Could not marshal sysinfo for AGENT message", "err", err)
		}
		ipAddress = sysInfo.Host.PrimaryIPAddress()
	}

	return &pb.Agent{
		Id:          env.AgentID,
		Action:      action,
		Time:        timestamppb.New(time.Now().Round(time.Microsecond)),
		Version:     env.Version,
		Hostname:    env.Hostname,
		IpAddress:   ipAddress,
		CloudRegion: env.Region,
		CloudZone:   env.Zone,
		Sysinfo:     sysInfoJSON,
	}
}

func newJobMessage(env *env, tick *worker.Tick, name string) *pb.Job {
	if tick.Took == 0 {
		tick.Took = time.Since(tick.Time)
	}
	return &pb.Job{
		AgentId: env.AgentID,
		Id:      tick.ID,
		Time:    timestamppb.New(tick.Time.Round(time.Microsecond)),
		Drift:   durationpb.New(tick.Drift),
		Took:    durationpb.New(tick.Took),
		Name:    name,
		Errors:  0,
	}
}
