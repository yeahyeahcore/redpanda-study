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

	responses := receiver.client.ProduceSync(ctx, &kgo.Record{
		Key:   []byte(receiver.messageKey),
		Topic: receiver.topic,
		Value: bytes,
	})

	for _, response := range responses {
		if response.Err != nil {
			receiver.logger.Error("unable send tariff on <Send> of <TariffProducer>", zap.Error(response.Err))
			return err
		}
	}

	return nil
}

func (receiver *Producer) Close(_ context.Context) error {
	receiver.client.Close()

	return nil
}
