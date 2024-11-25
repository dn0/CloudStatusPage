package msg

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type dummyPublisher struct {
	wg      sync.WaitGroup
	topic   string
	timeout time.Duration
}

func newDummyPublisher(cfg *PubsubPublisherConfig, topicID string) *dummyPublisher {
	return &dummyPublisher{
		topic:   topicID,
		timeout: cfg.PubsubPublishTimeout,
	}
}

func (p *dummyPublisher) Close() {
	p.wg.Wait()
}

func (p *dummyPublisher) Publish(ctx context.Context, msg *Message) {
	_ = p.publish(ctx, msg, false)
}

func (p *dummyPublisher) PublishWait(ctx context.Context, msg *Message) error {
	return p.publish(ctx, msg, true)
}

func (p *dummyPublisher) publish(ctx context.Context, msg *Message, block bool) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()
		defer cancel()

		slog.Debug("Published message", "topic", p.topic, "msg_attrs", msg.Attributes, "msg_data", msg.Data)
	}()

	if block {
		<-ctx.Done()
	}

	return nil
}
