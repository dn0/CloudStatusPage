package common

import (
	"context"

	"cspage/pkg/http"
	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

func GetRunningAgents(ctx context.Context, cfg *agent.Config, client *http.Client) (pb.Agents, error) {
	url := cfg.SiteURL + "/cloud/" + cfg.Env.Cloud + "/agent"
	headers := map[string]string{
		"Authorization": "Basic " + cfg.SiteSecret,
	}
	agents, err := http.GetJSON[pb.Agents](ctx, client, url, headers)
	if err != nil {
		return nil, err
	}
	return *agents, nil
}
