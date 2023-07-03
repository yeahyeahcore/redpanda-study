package tariff

import (
	"context"
	"encoding/json"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/yeahyeahcore/redpanda-study/internal/interface/broker/tariff/dto"
	"go.uber.org/zap"
)

type ProducerDeps struct {
	Logger     *zap.Logger
	Brokers    []string
	Topic      string
	MessageKey string
}

type Producer struct {
	logger     *zap.Logger
	client     *kgo.Client
	topic      string
	messageKey string
}

func NewProducer(deps ProducerDeps) (*Producer, error) {
	client, err := kgo.NewClient(kgo.SeedBrokers(deps.Brokers...))
	if err != nil {
		return nil, err
	}

	return &Producer{
		logger:     deps.Logger,
		client:     client,
		topic:      deps.Topic,
		messageKey: deps.MessageKey,
	}, nil
}

func (receiver *Producer) SendMessage(ctx context.Context, message *dto.Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	receiver.client.Produce(ctx, &kgo.Record{
		Key:   []byte(receiver.messageKey),
		Topic: receiver.topic,
		Value: bytes,
	}, nil)

	return nil
}

func (receiver *Producer) Close(_ context.Context) error {
	receiver.client.Close()

	return nil
}
