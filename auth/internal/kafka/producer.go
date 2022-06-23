package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Producer - kafka producer.
type Producer interface {
	SendMessage(context.Context, protoreflect.ProtoMessage) error
	Close() error
}

// producer - producer implementation.
type producer struct {
	writer *kafka.Writer
}

// NewProducer - ConsumerImpl constructor.
func NewProducer(topic string, brokers []string, groupID string) Producer {
	return &producer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
			Dialer: &kafka.Dialer{
				Timeout:   10 * time.Second,
				DualStack: true,
			},
		}),
	}
}

func (p *producer) SendMessage(ctx context.Context, msg protoreflect.ProtoMessage) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
}

func (p *producer) Close() error {
	return p.writer.Close()
}
