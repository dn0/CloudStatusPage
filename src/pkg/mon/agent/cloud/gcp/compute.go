package gcp

import (
	"context"

	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/gcp/compute"
	"cspage/pkg/mon/agent/cloud/gcp/vpc"
	"cspage/pkg/msg"
)

func NewComputeVMStandardProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.GCP],
	topic *msg.PubsubTopic,
) *agent.Worker {
	pingJob := vpc.NewIntraPingProbeJob(cfg, topic, cfg.Cloud.ComputeVMPrefix+"-TBD")
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeLongInterval(cfg.Cloud.ComputeVMProbeInterval),
		compute.NewVMProbe(cfg, compute.VMProvisioningModelStandard, pingJob),
	)
}

func NewComputeVMSpotProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.GCP],
	topic *msg.PubsubTopic,
) *agent.Worker {
	pingJob := vpc.NewIntraPingProbeJob(cfg, topic, cfg.Cloud.ComputeVMSpotPrefix+"-TBD")
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeLongInterval(cfg.Cloud.ComputeVMSpotProbeInterval),
		compute.NewVMProbe(cfg, compute.VMProvisioningModelSpot, pingJob),
	)
}

func NewComputeVMMetadataProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.GCP],
	topic *msg.PubsubTopic,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.ComputeVMMetadataProbeInterval),
		compute.NewVMMetadataProbe(cfg),
	)
}

func NewComputeDiskSnapshotProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.GCP],
	topic *msg.PubsubTopic,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeLongInterval(cfg.Cloud.ComputeDiskSnapshotProbeInterval),
		compute.NewDiskSnapshotProbe(cfg),
	)
}
