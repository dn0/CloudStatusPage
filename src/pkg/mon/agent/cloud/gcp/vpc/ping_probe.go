package vpc

import (
	"cspage/pkg/mon/agent"
	VPC "cspage/pkg/mon/agent/cloud/common/vpc"
	"cspage/pkg/msg"
)

const (
	pingProbeName = "gcp_vpc" // doc="Google Cloud VPC"
)

//nolint:mnd // Documentation is used by `make db/sql`.
func NewIntraPingProbeJob(
	cfg *agent.GCPConfig,
	topic *msg.PubsubTopic,
	host string,
) *VPC.PingProbeJob[agent.GCP] {
	return VPC.NewPingProbeJob(cfg, topic, VPC.NewPingProbe[agent.GCP](
		&cfg.Config.VPCIntraPing,
		pingProbeName,
		host,
		host,
		&VPC.PingAction{
			Id:   91, // name="vpc.intra.Ping64B" doc="64 bytes ICMP ping within a VPC network subnet" url=""
			Size: VPC.PingSize64B,
		},
		&VPC.PingAction{
			Id:   92, // name="vpc.intra.Ping1500B" doc="1500 bytes ICMP ping within a VPC network subnet" url=""
			Size: VPC.PingSize1500B,
		},
	))
}

//nolint:lll,mnd // Documentation is used by `make db/sql`.
func NewInterPingProbe[T agent.GCP](cfg *agent.GCPConfig) *VPC.AgentPingProbe[T] {
	return VPC.NewAgentPingProbe[T](
		&cfg.Config,
		pingProbeName,
		map[string]*VPC.PingAction{
			"us-east1":                {Id: 9201, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-east1" doc="64 bytes ICMP ping sent to VPC network subnet in us-east1" url=""
			"us-east4":                {Id: 9202, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-east4" doc="64 bytes ICMP ping sent to VPC network subnet in us-east4" url=""
			"us-east5":                {Id: 9203, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-east5" doc="64 bytes ICMP ping sent to VPC network subnet in us-east5" url=""
			"us-west1":                {Id: 9204, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-west1" doc="64 bytes ICMP ping sent to VPC network subnet in us-west1" url=""
			"us-west2":                {Id: 9205, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-west2" doc="64 bytes ICMP ping sent to VPC network subnet in us-west2" url=""
			"us-west3":                {Id: 9206, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-west3" doc="64 bytes ICMP ping sent to VPC network subnet in us-west3" url=""
			"us-west4":                {Id: 9207, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-west4" doc="64 bytes ICMP ping sent to VPC network subnet in us-west4" url=""
			"us-central1":             {Id: 9208, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-central1" doc="64 bytes ICMP ping sent to VPC network subnet in us-central1" url=""
			"us-south1":               {Id: 9209, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B us-south1" doc="64 bytes ICMP ping sent to VPC network subnet in us-south1" url=""
			"northamerica-northeast1": {Id: 9210, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B northamerica-northeast1" doc="64 bytes ICMP ping sent to VPC network subnet in northamerica-northeast1" url=""
			"northamerica-northeast2": {Id: 9211, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B northamerica-northeast2" doc="64 bytes ICMP ping sent to VPC network subnet in northamerica-northeast2" url=""
			"southamerica-east1":      {Id: 9231, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B southamerica-east1" doc="64 bytes ICMP ping sent to VPC network subnet in southamerica-east1" url=""
			"southamerica-west1":      {Id: 9232, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B southamerica-west1" doc="64 bytes ICMP ping sent to VPC network subnet in southamerica-west1" url=""
			"europe-west1":            {Id: 9241, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B europe-west1" doc="64 bytes ICMP ping sent to VPC network subnet in europe-west1" url=""
			"europe-west2":            {Id: 9242, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B europe-west2" doc="64 bytes ICMP ping sent to VPC network subnet in europe-west2" url=""
			"europe-west3":            {Id: 9243, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B europe-west3" doc="64 bytes ICMP ping sent to VPC network subnet in europe-west3" url=""
			"europe-west4":            {Id: 9244, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B europe-west4" doc="64 bytes ICMP ping sent to VPC network subnet in europe-west4" url=""
			"europe-west6":            {Id: 9245, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B europe-west6" doc="64 bytes ICMP ping sent to VPC network subnet in europe-west6" url=""
			"europe-west8":            {Id: 9246, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B europe-west8" doc="64 bytes ICMP ping sent to VPC network subnet in europe-west8" url=""
			"europe-west9":            {Id: 9247, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B europe-west9" doc="64 bytes ICMP ping sent to VPC network subnet in europe-west9" url=""
			"europe-west10":           {Id: 9248, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B europe-west10" doc="64 bytes ICMP ping sent to VPC network subnet in europe-west10" url=""
			"europe-west12":           {Id: 9249, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B europe-west12" doc="64 bytes ICMP ping sent to VPC network subnet in europe-west12" url=""
			"europe-north1":           {Id: 9250, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B europe-north1" doc="64 bytes ICMP ping sent to VPC network subnet in europe-north1" url=""
			"europe-southwest1":       {Id: 9251, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B europe-southwest1" doc="64 bytes ICMP ping sent to VPC network subnet in europe-southwest1" url=""
			"europe-central2":         {Id: 9252, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B europe-central2" doc="64 bytes ICMP ping sent to VPC network subnet in europe-central2" url=""
			"australia-southeast1":    {Id: 9271, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B australia-southeast1" doc="64 bytes ICMP ping sent to VPC network subnet in australia-southeast1" url=""
			"australia-southeast2":    {Id: 9272, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B australia-southeast2" doc="64 bytes ICMP ping sent to VPC network subnet in australia-southeast2" url=""
			"asia-east1":              {Id: 9273, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B asia-east1" doc="64 bytes ICMP ping sent to VPC network subnet in asia-east1" url=""
			"asia-east2":              {Id: 9274, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B asia-east2" doc="64 bytes ICMP ping sent to VPC network subnet in asia-east2" url=""
			"asia-northeast1":         {Id: 9275, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B asia-northeast1" doc="64 bytes ICMP ping sent to VPC network subnet in asia-northeast1" url=""
			"asia-northeast2":         {Id: 9276, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B asia-northeast2" doc="64 bytes ICMP ping sent to VPC network subnet in asia-northeast2" url=""
			"asia-northeast3":         {Id: 9277, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B asia-northeast3" doc="64 bytes ICMP ping sent to VPC network subnet in asia-northeast3" url=""
			"asia-south1":             {Id: 9278, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B asia-south1" doc="64 bytes ICMP ping sent to VPC network subnet in asia-south1" url=""
			"asia-south2":             {Id: 9279, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B asia-south2" doc="64 bytes ICMP ping sent to VPC network subnet in asia-south2" url=""
			"asia-southeast1":         {Id: 9280, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B asia-southeast1" doc="64 bytes ICMP ping sent to VPC network subnet in asia-southeast1" url=""
			"asia-southeast2":         {Id: 9281, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B asia-southeast2" doc="64 bytes ICMP ping sent to VPC network subnet in asia-southeast2" url=""
			"me-central1":             {Id: 9282, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B me-central1" doc="64 bytes ICMP ping sent to VPC network subnet in me-central1" url=""
			"me-central2":             {Id: 9283, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B me-central2" doc="64 bytes ICMP ping sent to VPC network subnet in me-central2" url=""
			"me-west1":                {Id: 9284, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B me-west1" doc="64 bytes ICMP ping sent to VPC network subnet in me-west1" url=""
			"africa-south1":           {Id: 9285, Size: VPC.PingSize64B}, // name="vpc.inter.Ping64B africa-south1" doc="64 bytes ICMP ping sent to VPC network subnet in africa-south1" url=""
		},
	)
}
