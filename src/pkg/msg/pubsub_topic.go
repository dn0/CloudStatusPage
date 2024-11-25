package msg

import (
	"log/slog"

	"cloud.google.com/go/pubsub"
)

type PubsubTopic struct {
	ID    string
	Topic *pubsub.Topic
}

func ClosePubsubTopic(topic *PubsubTopic) {
	slog.Debug("Stopping pubsub topic...", "topic", topic.ID)
	if topic.Topic != nil {
		topic.Topic.Stop()
	}
}

func NewPubsubTopic(client *pubsub.Client, id string) *PubsubTopic {
	return &PubsubTopic{
		ID:    id,
		Topic: newPubsubTopic(client, id),
	}
}

func newPubsubTopic(client *pubsub.Client, id string) *pubsub.Topic {
	if client == nil {
		return nil
	}
	return client.Topic(id)
}
