package aws

import (
	"context"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"

	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/aws/s3"
	"cspage/pkg/msg"
)

func NewS3ObjectProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.AWS],
	topic *msg.PubsubTopic,
	awsConfig *awsSDK.Config,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.S3ObjectProbeInterval),
		s3.NewObjectProbe(cfg, awsConfig),
	)
}

func NewS3BucketProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.AWS],
	topic *msg.PubsubTopic,
	awsConfig *awsSDK.Config,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.S3BucketProbeInterval),
		s3.NewBucketProbe(cfg, awsConfig),
	)
}
