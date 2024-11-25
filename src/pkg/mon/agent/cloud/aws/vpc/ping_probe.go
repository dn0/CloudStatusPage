package vpc

import (
	"cspage/pkg/mon/agent"
	VPC "cspage/pkg/mon/agent/cloud/common/vpc"
	"cspage/pkg/msg"
)

const (
	pingProbeName = "aws_vpc" // doc="Amazon VPC"
)

//nolint:mnd // Documentation is used by `make db/sql`.
func NewIntraPingProbeJob(
	cfg *agent.AWSConfig,
	topic *msg.PubsubTopic,
	host string,
) *VPC.PingProbeJob[agent.AWS] {
	return VPC.NewPingProbeJob(cfg, topic, VPC.NewPingProbe[agent.AWS](
		&cfg.Config.VPCIntraPing,
		pingProbeName,
		host,
		host,
		&VPC.PingAction{
			Id:   91, // name="vpc.intra.Ping64B" doc="64 bytes ICMP ping within a VPC subnet" url=""
			Size: VPC.PingSize64B,
		},
		&VPC.PingAction{
			Id:   92, // name="vpc.intra.Ping1500B" doc="1500 bytes ICMP ping within a VPC subnet" url=""
			Size: VPC.PingSize1500B,
		},
	))
}

//nolint:lll,mnd // Documentation is used by `make db/sql`.
func NewInterPingProbe[T agent.AWS](cfg *agent.AWSConfig) *VPC.AgentPingProbe[T] {
	return VPC.NewAgentPingProbe[T](
		&cfg.Config,
		pingProbeName,
		map[string]*VPC.PingAction{
			"us-east-1":      {Id: 9201, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-east-1" doc="64 bytes ICMP ping sent to VPC in us-east-1" url=""
			"us-east-2":      {Id: 9202, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-east-2" doc="64 bytes ICMP ping sent to VPC in us-east-2" url=""
			"us-west-1":      {Id: 9203, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-west-1" doc="64 bytes ICMP ping sent to VPC in us-west-1" url=""
			"us-west-2":      {Id: 9204, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-west-2" doc="64 bytes ICMP ping sent to VPC in us-west-2" url=""
			"ca-west-1":      {Id: 9205, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ca-west-1" doc="64 bytes ICMP ping sent to VPC in ca-west-1" url=""
			"ca-central-1":   {Id: 9206, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ca-central-1" doc="64 bytes ICMP ping sent to VPC in ca-central-1" url=""
			"mx-central-1":   {Id: 9207, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B mx-central-1" doc="64 bytes ICMP ping sent to VPC in mx-central-1" url=""
			"sa-east-1":      {Id: 9231, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B sa-east-1" doc="64 bytes ICMP ping sent to VPC in sa-east-1" url=""
			"eu-west-1":      {Id: 9241, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B eu-west-1" doc="64 bytes ICMP ping sent to VPC in eu-west-1" url=""
			"eu-west-2":      {Id: 9242, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B eu-west-2" doc="64 bytes ICMP ping sent to VPC in eu-west-2" url=""
			"eu-west-3":      {Id: 9243, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B eu-west-3" doc="64 bytes ICMP ping sent to VPC in eu-west-3" url=""
			"eu-north-1":     {Id: 9244, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B eu-north-1" doc="64 bytes ICMP ping sent to VPC in eu-north-1" url=""
			"eu-south-1":     {Id: 9245, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B eu-south-1" doc="64 bytes ICMP ping sent to VPC in eu-south-1" url=""
			"eu-south-2":     {Id: 9246, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B eu-south-2" doc="64 bytes ICMP ping sent to VPC in eu-south-2" url=""
			"eu-central-1":   {Id: 9247, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B eu-central-1" doc="64 bytes ICMP ping sent to VPC in eu-central-1" url=""
			"eu-central-2":   {Id: 9248, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B eu-central-2" doc="64 bytes ICMP ping sent to VPC in eu-central-2" url=""
			"ap-northeast-1": {Id: 9271, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-northeast-1" doc="64 bytes ICMP ping sent to VPC in ap-northeast-1" url=""
			"ap-northeast-2": {Id: 9272, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-northeast-2" doc="64 bytes ICMP ping sent to VPC in ap-northeast-2" url=""
			"ap-northeast-3": {Id: 9273, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-northeast-3" doc="64 bytes ICMP ping sent to VPC in ap-northeast-3" url=""
			"ap-east-1":      {Id: 9274, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-east-1" doc="64 bytes ICMP ping sent to VPC in ap-east-1" url=""
			"ap-southeast-1": {Id: 9275, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-southeast-1" doc="64 bytes ICMP ping sent to VPC in ap-southeast-1" url=""
			"ap-southeast-2": {Id: 9276, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-southeast-2" doc="64 bytes ICMP ping sent to VPC in ap-southeast-2" url=""
			"ap-southeast-3": {Id: 9277, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-southeast-3" doc="64 bytes ICMP ping sent to VPC in ap-southeast-3" url=""
			"ap-southeast-4": {Id: 9278, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-southeast-4" doc="64 bytes ICMP ping sent to VPC in ap-southeast-4" url=""
			"ap-southeast-5": {Id: 9279, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-southeast-5" doc="64 bytes ICMP ping sent to VPC in ap-southeast-5" url=""
			"ap-southeast-6": {Id: 9280, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-southeast-6" doc="64 bytes ICMP ping sent to VPC in ap-southeast-6" url=""
			"ap-southeast-7": {Id: 9281, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-southeast-7" doc="64 bytes ICMP ping sent to VPC in ap-southeast-7" url=""
			"ap-south-1":     {Id: 9282, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-south-1" doc="64 bytes ICMP ping sent to VPC in ap-south-1" url=""
			"ap-south-2":     {Id: 9283, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B ap-south-2" doc="64 bytes ICMP ping sent to VPC in ap-south-2" url=""
			"me-south-1":     {Id: 9284, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B me-south-1" doc="64 bytes ICMP ping sent to VPC in me-south-1" url=""
			"me-south-2":     {Id: 9285, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B me-south-2" doc="64 bytes ICMP ping sent to VPC in me-south-2" url=""
			"me-central-1":   {Id: 9286, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B me-central-1" doc="64 bytes ICMP ping sent to VPC in me-central-1" url=""
			"il-central-1":   {Id: 9287, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B il-central-1" doc="64 bytes ICMP ping sent to VPC in il-central-1" url=""
			"af-south-1":     {Id: 9288, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B af-south-1" doc="64 bytes ICMP ping sent to VPC in af-south-1" url=""
		},
	)
}
