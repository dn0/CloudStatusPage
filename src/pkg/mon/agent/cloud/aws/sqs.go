package aws

import (
	"context"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"

	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/aws/sqs"
	"cspage/pkg/msg"
)

func NewSQSMessageProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.AWS],
	topic *msg.PubsubTopic,
	awsConfig *awsSDK.Config,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.SQSMessageInterval),
		sqs.NewMessageProbe(cfg, awsConfig),
	)
}
