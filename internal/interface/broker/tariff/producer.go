package tariff

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/yeahyeahcore/redpanda-study/internal/models"
	"github.com/yeahyeahcore/redpanda-study/internal/utils/transfer"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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

func (receiver *Producer) Send(ctx context.Context, tariff *models.Tariff) error {
	bytes, err := proto.Marshal(transfer.TariffModelToProto(tariff))
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
