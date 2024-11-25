package aws

import (
	"context"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"

	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/aws/ec2"
	"cspage/pkg/mon/agent/cloud/aws/vpc"
	"cspage/pkg/msg"
)

func NewEC2VMStandardProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.AWS],
	topic *msg.PubsubTopic,
	awsConfig *awsSDK.Config,
) *agent.Worker {
	pingJob := vpc.NewIntraPingProbeJob(cfg, topic, cfg.Cloud.EC2VMPrefix+"-TBD")
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeLongInterval(cfg.Cloud.EC2VMProbeInterval),
		ec2.NewVMProbe(cfg, awsConfig, ec2.VMProvisioningModelStandard, pingJob),
	)
}

func NewEC2VMSpotProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.AWS],
	topic *msg.PubsubTopic,
	awsConfig *awsSDK.Config,
) *agent.Worker {
	pingJob := vpc.NewIntraPingProbeJob(cfg, topic, cfg.Cloud.EC2VMSpotPrefix+"-TBD")
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeLongInterval(cfg.Cloud.EC2VMSpotProbeInterval),
		ec2.NewVMProbe(cfg, awsConfig, ec2.VMProvisioningModelSpot, pingJob),
	)
}

func NewEC2VMMetadataProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.AWS],
	topic *msg.PubsubTopic,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.EC2VMMetadataProbeInterval),
		ec2.NewVMMetadataProbe(cfg),
	)
}

func NewEC2EBSSnapshotProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.AWS],
	topic *msg.PubsubTopic,
	awsConfig *awsSDK.Config,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeLongInterval(cfg.Cloud.EC2EBSSnapshotProbeInterval),
		ec2.NewEBSSnapshotProbe(cfg, awsConfig),
	)
}
