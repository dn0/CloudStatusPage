package servicebus

import (
	"context"
	"log/slog"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"

	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

//nolint:lll // Documentation is used by `make db/sql`.
const (
	queueMessageProbeSleep         = 10 * time.Millisecond            // Sleep between send and receive. Subtracted from the final message latency.
	queueMessageProbeName          = "azure_servicebus_queue_message" // doc="Azure Service Bus Queue Message"
	queueMessageProbeActionSend    = 60                               // name="servicebus.queue.message.send" doc="Sends a couple of small messages to a queue" url="https://learn.microsoft.com/en-us/rest/api/servicebus/send-message-to-queue"
	queueMessageProbeActionReceive = 70                               // name="servicebus.queue.message.receive" doc="Receives all messages sent by `servicebus.queue.message.send` and removes them from queue" url="https://learn.microsoft.com/en-us/rest/api/servicebus/receive-and-delete-message-destructive-read"
	queueMessageProbeActionLatency = 80                               // name="servicebus.queue.message.latency" doc="Average delay between sending and receiving of one message" url=""
)

type QueueMessageProbe[T agent.Azure] struct {
	message  *message
	warmedUp bool
}

func NewQueueMessageProbe[T agent.Azure](cfg *agent.AzureConfig, client *azservicebus.Client) *QueueMessageProbe[T] {
	return &QueueMessageProbe[T]{
		message: &message{
			namespace: cfg.Cloud.ServiceBusNamespace,
			queue:     cfg.Cloud.ServiceBusQueueName,
			client:    client,
			timeout:   cfg.ProbeTimeout,
		},
	}
}

func (p *QueueMessageProbe[T]) String() string {
	return queueMessageProbeName
}

func (p *QueueMessageProbe[T]) Start(_ context.Context) {
	var err error
	if p.message.sender, err = p.message.client.NewSender(p.message.queue, nil); err != nil {
		agent.DieLog(p.log(), "Could not initialize Azure service bus queue sender", "err", err)
	}
	if p.message.receiver, err = p.message.client.NewReceiverForQueue(p.message.queue, nil); err != nil {
		agent.DieLog(p.log(), "Could not initialize Azure service bus queue receiver", "err", err)
	}
	p.log().Info("Probe initialized")
}

func (p *QueueMessageProbe[T]) Do(ctx context.Context) []*pb.Result {
	res := []*pb.Result{
		pb.NewResult(queueMessageProbeActionSend),
		pb.NewResult(queueMessageProbeActionReceive),
		pb.NewResult(queueMessageProbeActionLatency),
	}
	p.message.new()

	res[0].Store(p.sendMessages(ctx))
	if res[0].Failed() {
		return res
	}

	time.Sleep(queueMessageProbeSleep)

	res[1].Store(p.receiveMessages(ctx))

	if p.message.receivedAll() {
		res[2].Store(p.getLatency(), nil)
	}

	if !p.warmedUp {
		p.log().Debug("Cold start => discarding all results")
		p.warmedUp = true
		return nil
	}

	return res
}

func (p *QueueMessageProbe[T]) Stop(ctx context.Context) {
	_ = p.message.sender.Close(ctx)
	_ = p.message.receiver.Close(ctx)
}

//nolint:dupl // We can tolerate the duplicate code between receiveMessages and sendMessages.
func (p *QueueMessageProbe[T]) sendMessages(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Sending messages to service bus queue...")
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
			p.log().Debug("Sent message to service bus queue", "num", num, "took", ret.Took)
		} else {
			p.log().Error("Could not send message to service bus queue", "num", num, "took", ret.Took, "err", err)
			break
		}
	}

	retAll.Took /= time.Duration(num)
	p.log().Debug("Sent messages to service bus queue", "sent", num, "took", retAll.Took)
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return retAll, err
}

//nolint:dupl // We can tolerate the duplicate code between receiveMessages and sendMessages.
func (p *QueueMessageProbe[T]) receiveMessages(ctx context.Context) (pb.ResultTime, error) {
	p.log().Debug("Receiving messages from service bus queue...")
	retAll := pb.EmptyResultTime()
	var err error
	var num int

	for !p.message.receivedAll() {
		var ret pb.ResultTime
		ret, err = pb.Timeit(p.message.receive, ctx)
		num++
		//goland:noinspection GoDfaErrorMayBeNotNil
		retAll.Took += ret.Took
		if err == nil {
			p.log().Debug("Received message from service bus queue", "num", num, "took", ret.Took)
		} else {
			p.log().Error("Could not receive message from service bus queue", "num", num, "took", ret.Took, "err", err)
			break
		}
	}

	retAll.Took /= time.Duration(num)
	p.log().Debug("Received messages from service bus queue", "received", num, "took", retAll.Took)
	//nolint:wrapcheck // A probe error is expected to be properly wrapped by the lowest probe method.
	return retAll, err
}

func (p *QueueMessageProbe[T]) getLatency() pb.ResultTime {
	ret := pb.EmptyResultTime()
	if len(p.message.received) == 0 {
		return ret
	}

	var sum time.Duration
	for _, took := range p.message.received {
		sum += took - queueMessageProbeSleep
	}

	ret.Took = sum / time.Duration(len(p.message.received))
	p.log().Debug("Average message send/receive latency", "took", ret.Took)
	return ret
}

func (p *QueueMessageProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "servicebus", p.message.String(), "count", messageCount)
}
