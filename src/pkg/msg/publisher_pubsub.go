package msg

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
)

type pubsubPublisher struct {
	wg      sync.WaitGroup
	topic   *pubsub.Topic
	timeout time.Duration
}

func newPubsubPublisher(cfg *PubsubPublisherConfig, topic *pubsub.Topic) *pubsubPublisher {
	return &pubsubPublisher{
		topic:   topic,
		timeout: cfg.PubsubPublishTimeout,
	}
}

func (p *pubsubPublisher) Close() {
	p.wg.Wait()
}

func (p *pubsubPublisher) Publish(ctx context.Context, msg *Message) {
	_ = p.publish(ctx, msg, false)
}

func (p *pubsubPublisher) PublishWait(ctx context.Context, msg *Message) error {
	return p.publish(ctx, msg, true)
}

func (p *pubsubPublisher) publish(ctx context.Context, msg *Message, block bool) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	resultChan := make(chan error, 1)
	result := p.topic.Publish(ctx, msg)
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()
		defer cancel()

		id, err := result.Get(ctx)
		if err == nil {
			slog.Debug("Published message", "topic", p.topic.ID(), "msg_attrs", msg.Attributes, "msg_data", msg.Data, "id", id)
		} else {
			slog.Error("Could not publish", "topic", p.topic.ID(), "msg_attrs", msg.Attributes, "msg_data", msg.Data, "err", err)
		}
		resultChan <- err
	}()

	if block {
		return <-resultChan
	}

	return nil
}
