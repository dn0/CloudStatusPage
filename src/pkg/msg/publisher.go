package msg

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, msg *Message)
	PublishWait(ctx context.Context, msg *Message) error
	Close()
}

func NewPublisher(cfg *PubsubPublisherConfig, topic *PubsubTopic) Publisher {
	if cfg.PubsubProjectID == "" {
		return newDummyPublisher(cfg, topic.ID)
	}
	return newPubsubPublisher(cfg, topic.Topic)
}
