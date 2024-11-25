package vpc

import (
	"cspage/pkg/mon/agent"
	VPC "cspage/pkg/mon/agent/cloud/common/vpc"
)

const (
	pingHostName = "cloudstatus.page Google LB or CDN"
	pingHostAddr = "cloudstatus.page"

	pingProbeName        = "dummy_vpc" // doc="Dummy VPC"
	pingProbeAction64B   = 901         // name="vpc.inter.Ping64B" doc="64 bytes ICMP ping between internet hosts" url=""
	pingProbeAction1500B = 902         // name="vpc.inter.Ping1500B" doc="1500 bytes ICMP ping between internet hosts" url=""
)

//nolint:mnd // Random IDs and strings for testing.
func NewAgentPingProbe[T agent.Dummy](cfg *agent.Config) *VPC.AgentPingProbe[T] {
	return VPC.NewAgentPingProbe[T](
		cfg,
		pingProbeName,
		map[string]*VPC.PingAction{
			"europe-nitra1":      {Id: 9991, Size: VPC.PingSize64B},
			"europe-bratislava1": {Id: 9992, Size: VPC.PingSize64B},
		},
	)
}

func NewInternetPingProbe[T agent.Dummy](cfg *agent.Config) *VPC.PingProbe[T] {
	return VPC.NewPingProbe[T](&cfg.VPCIntraPing, pingProbeName, pingHostName, pingHostAddr,
		&VPC.PingAction{
			Id:   pingProbeAction64B,
			Size: VPC.PingSize64B,
		},
		&VPC.PingAction{
			Id:   pingProbeAction1500B,
			Size: VPC.PingSize1500B,
		},
	)
}
