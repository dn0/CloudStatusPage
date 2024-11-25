package msg

import (
	"cloud.google.com/go/pubsub"
)

type PubsubSubscription struct {
	ID           string
	Subscription *pubsub.Subscription
}

func NewPubsubSubscription(client *pubsub.Client, id string) *PubsubSubscription {
	return &PubsubSubscription{
		ID:           id,
		Subscription: newPubsubSubscription(client, id),
	}
}

func newPubsubSubscription(client *pubsub.Client, id string) *pubsub.Subscription {
	if client == nil {
		return nil
	}
	return client.Subscription(id)
}
