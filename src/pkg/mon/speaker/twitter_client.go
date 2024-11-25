package speaker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	gohttp "net/http"

	"github.com/dghubble/oauth1"

	"cspage/pkg/data"
	"cspage/pkg/http"
	"cspage/pkg/pb"
)

type tweetReply struct {
	InReplyToTweetID string `json:"in_reply_to_tweet_id,omitempty"`
}

type tweetRequest struct {
	Text  string      `json:"text"`
	Reply *tweetReply `json:"reply,omitempty"`
}

type tweetResponseData struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type tweetResponse struct {
	Data tweetResponseData `json:"data"`
}

type twitterClient struct {
	*http.Client
	cloud   string
	siteURL string
	headers map[string]string
}

func newTwitterClient(ctx context.Context, cfg *Config, cloud string) *twitterClient {
	config := oauth1.NewConfig(cfg.TwitterAPIKey, cfg.TwitterAPISecret)
	token := oauth1.NewToken(cfg.TwitterAccessToken, cfg.TwitterAccessSecret)
	return &twitterClient{
		Client:  http.NewClientFromClient(config.Client(ctx, token)),
		cloud:   cloud,
		siteURL: cfg.SiteURL,
		headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}
}

func (t *twitterClient) createTweetFromIncident(ctx context.Context, inc *data.Incident) (*tweetResponse, error) {
	tweetReq := t.newTweetRequest(inc)
	if tweetReq == nil {
		//nolint:nilnil // We don't need a sentinel error here.
		return nil, nil
	}

	body, err := json.Marshal(tweetReq)
	if err != nil {
		return nil, fmt.Errorf("incident=%s: failed to marshal tweet request: %w", inc.Id, err)
	}

	req, err := http.NewRequest(ctx, "POST", "https://api.x.com/2/tweets", t.headers, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("incident=%s: failed to create tweet request: %w", inc.Id, err)
	}

	res, err := t.Do(req)
	if err != nil {
		return nil, fmt.Errorf("incident=%s: failed to post new tweet: %w", inc.Id, err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != gohttp.StatusCreated {
		return nil, http.NewError(res)
	}

	tweetRes := &tweetResponse{}
	if err = decoder.Decode(tweetRes); err != nil {
		return nil, fmt.Errorf("incident=%s: failed to unmarshal tweet response: %w", inc.Id, err)
	}

	return tweetRes, nil
}

func (t *twitterClient) newTweetRequest(inc *data.Incident) *tweetRequest {
	req := &tweetRequest{}

	switch {
	case inc.Status == pb.IncidentStatus_INCIDENT_OPEN && inc.Data.TwitterId == "":
		req.Text = "📟 New incident: "
	case inc.Status == pb.IncidentStatus_INCIDENT_CLOSED && inc.Data.TwitterId != "":
		req.Text = "✅ Incident closed: "
		req.Reply = &tweetReply{InReplyToTweetID: inc.Data.TwitterId}
	default:
		return nil
	}

	switch inc.Severity {
	case pb.IncidentSeverity_INCIDENT_HIGH:
		req.Text += "🟥 "
	case pb.IncidentSeverity_INCIDENT_MEDIUM:
		req.Text += "🟧 "
	case pb.IncidentSeverity_INCIDENT_LOW:
		req.Text += "🟨 "
	case pb.IncidentSeverity_INCIDENT_NONE:
		req.Text += " "
	}

	req.Text += inc.Data.Summary + " " + t.siteURL + inc.DetailsURL(t.cloud)

	return req
}
