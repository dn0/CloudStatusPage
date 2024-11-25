package pubsub

import (
	"context"
	"log/slog"
	"time"

	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	messageProbeCount         = 5
	messageProbeSleep         = 10 * time.Millisecond // Sleep between send and receive. Subtracted from the final message latency.
	messageProbeSleepInit     = 90 * time.Millisecond // Sleep between initial warm-up and actual a run.
	messageProbeName          = "gcp_pubsub_message"  // doc="Google Cloud Pub/Sub Message"
	messageProbeActionSend    = 60                    // name="pubsub.message.send" doc="Publishes a couple of small messages to a topic" url="https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics/publish"
	messageProbeActionReceive = 70                    // name="pubsub.message.receive" doc="Receives and acknowledges all messages published by `pubsub.message.send`" url="https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.subscriptions/pull"
	messageProbeActionLatency = 80                    // name="pubsub.message.latency" doc="Average delay between sending and receiving of one message" url=""
)

type MessageProbe[T agent.GCP] struct {
	message *message
}

func NewMessageProbe[T agent.GCP](cfg *agent.GCPConfig) *MessageProbe[T] {
	return &MessageProbe[T]{
		message: &message{
			count:          messageProbeCount,
			projectID:      cfg.Cloud.PubsubProject,
			region:         cfg.Env.Region,
			topicID:        cfg.Cloud.PubsubTopic,
			subscriptionID: cfg.Cloud.PubsubSubscription,
			timeout:        cfg.ProbeTimeout,
		},
	}
}

func (p *MessageProbe[T]) String() string {
	return messageProbeName
}

func (p *MessageProbe[T]) Start(_ context.Context) {
	p.log().Info("Probe initialized")
}

func (p *MessageProbe[T]) Do(ctx context.Context) []*pb.Result {
	res := []*pb.Result{
		pb.NewResult(messageProbeActionSend),
		pb.NewResult(messageProbeActionReceive),
		pb.NewResult(messageProbeActionLatency),
	}

	if err := p.message.init(ctx); err != nil {
		p.log().Error("Could not initialize pubsub probe", "err", err)
		res[0].Store(pb.EmptyResultTime(), err)
		return res
	}
	defer p.message.close()
	p.message.new()

	// warm-up
	if ret, err := p.sendMessages(ctx); err != nil {
		p.log().Error("Initial message sending warmup failed", "err", err)
		res[0].Store(ret, err)
		return res
	}
	time.Sleep(messageProbeSleep)
	if ret, err := p.receiveMessages(ctx); err != nil {
		p.log().Error("Initial message receiving warmup failed", "err", err)
		res[1].Store(ret, err)
		return res
	}

	time.Sleep(messageProbeSleepInit)
	p.message.new()

	res[0].Store(p.sendMessages(ctx))
	if res[0].Failed() {
		return res
	}
	time.Sleep(messageProbeSleep)
	res[1].Store(p.receiveMessages(ctx))
	if p.message.receivedAll() {
		res[2].Store(p.getLatency(), nil)
	}

	return res
}

func (p *MessageProbe[T]) Stop(_ context.Context) {}

func (p *MessageProbe[T]) sendMessages(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Publishing messages to topic...")
	retAll := pb.EmptyResultTime()
	var err error
	var num int

	for !p.message.sentAll() {
		var ret pb.ResultTime
		ret, err = pb.Timeit(p.message.send, ctx)
		num++
		//goland:noinspection GoDfaErrorMayBeNotNil
		retAll.Took += ret.Took
		if err == nil {
			p.log().Debug("Published message to topic", "num", num, "took", ret.Took)
		} else {
			p.log().Error("Could not publish message to topic", "num", num, "took", ret.Took, "err", err)
			break
		}
	}

	retAll.Took /= time.Duration(num)
	p.log().Debug("Published messages to topic", "published", num, "took", retAll.Took)
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return retAll, err
}

func (p *MessageProbe[T]) receiveMessages(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Receiving messages from subscription...")
	ret, err := pb.Timeit(p.message.receive, ctx)
	//goland:noinspection GoDfaErrorMayBeNotNil
	ret.Took /= time.Duration(len(p.message.received))
	if err == nil {
		p.log().Debug("Received all messages from subscription", "took", ret.Took)
	} else {
		p.log().Error("Could not receive all messages from subscription", "took", ret.Took, "err", err)
	}

	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return ret, err
}

func (p *MessageProbe[T]) getLatency() pb.ResultTime {
	ret := pb.EmptyResultTime()
	if len(p.message.received) == 0 {
		return ret
	}

	var sum time.Duration
	for _, took := range p.message.received {
		sum += took - messageProbeSleep
	}

	ret.Took = sum / time.Duration(len(p.message.received))
	p.log().Debug("Average message send/receive latency", "took", ret.Took)
	return ret
}

func (p *MessageProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "pubsub", p.message.String(), "count", p.message.count)
}
