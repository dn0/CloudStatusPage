package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"cloud.google.com/go/pubsub"

	"cspage/pkg/config"
	"cspage/pkg/msg"
)

const (
	ConsumerJobPrefix          = "consumer:"
	messageReceiveDelayWarning = 10 * time.Second
)

type consumer struct {
	job        ConsumerJob
	timeout    time.Duration
	errTimeout error
	sub        *pubsub.Subscription
	taskCtx    context.Context
}

func NewConsumer(
	ctx context.Context,
	cfg *config.BaseConfig,
	job ConsumerJob,
	sub *msg.PubsubSubscription,
	loopDesc ...any,
) *Daemon {
	loopDesc = append([]any{
		"job", job.String(),
		"subscription", sub.ID,
	}, loopDesc...)
	worker := &consumer{
		job:        job,
		timeout:    cfg.WorkerTaskTimeout,
		errTimeout: fmt.Errorf("%s: %w", job.String(), errConsumerTaskTimeoutExceeded),
		sub:        sub.Subscription,
	}
	return newDaemon(ctx, cfg, job, sub.Subscription != nil, worker.loop, loopDesc...)
}

func (c *consumer) loop(loopCtx, taskCtx context.Context) {
	c.taskCtx = taskCtx
	err := c.sub.Receive(loopCtx, c.process)
	if err != nil {
		config.Die("Consumer received a non-retryable error", "err", err)
	}
}

//nolint:contextcheck // The context is inherited from runnerCtx.
func (c *consumer) process(_ context.Context, message *msg.Message) {
	log := slog.With("job", c.job.String(), "msg_id", message.ID, "msg_attrs", message.Attributes)
	now := time.Now()
	delay := now.Sub(message.PublishTime)
	log.Debug("Received message", "delay", delay, "attempt", message.DeliveryAttempt)

	if delay > messageReceiveDelayWarning {
		log.Warn(
			"Received message with substantial delay of "+delay.String(),
			"delay", delay,
			"attempt", message.DeliveryAttempt,
		)
	}

	taskCtx, taskCancel := context.WithTimeoutCause(c.taskCtx, c.timeout, c.errTimeout)
	err := c.job.Process(taskCtx, message)
	taskCancel()

	if err == nil {
		message.Ack()
		log.Debug("Processed message", "took", time.Since(now))
	} else {
		message.Nack()
		log.Error("Failed to process message", "took", time.Since(now), "err", err)
	}
}
