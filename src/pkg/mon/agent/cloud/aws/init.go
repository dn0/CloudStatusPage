package aws

import (
	"context"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"

	"cspage/pkg/mon/agent"
)

func LoadConfig(ctx context.Context, cfg *agent.CloudConfig[agent.AWS]) *awsSDK.Config {
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(cfg.Env.Region))
	if err != nil {
		agent.Die("Could not load AWS config", "err", err)
	}
	return &awsCfg
}
