package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	wbKafka "github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

type Kafka struct {
	Cons *wbKafka.Consumer
}

func New(brokers []string, topic string, groupID string) *Kafka {
	kafkaCons := wbKafka.NewConsumer(brokers, topic, groupID)

	return &Kafka{
		Cons: kafkaCons,
	}
}

func (k *Kafka) Consume(ctx context.Context, out chan<- kafka.Message, strategy retry.Strategy) {
	go func() {
		defer close(out)
		for {
			msg, err := k.Cons.Fetch(ctx)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					zlog.Logger.Warn().Err(err).Msg("kafka.go - failed to fetch message from kafka")
					continue
				}
			}
			if len(msg.Value) == 0 && len(msg.Key) == 0 {
				continue
			}

			select {
			case out <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()
}
