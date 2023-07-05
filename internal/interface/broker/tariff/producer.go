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
	Logger *zap.Logger
	Client *kgo.Client
	Topic  string
}

type Producer struct {
	logger *zap.Logger
	client *kgo.Client
	topic  string
}

func NewProducer(deps ProducerDeps) *Producer {
	return &Producer{
		logger: deps.Logger,
		client: deps.Client,
		topic:  deps.Topic,
	}
}

func (receiver *Producer) Send(ctx context.Context, tariff *models.Tariff) error {
	bytes, err := proto.Marshal(transfer.TariffModelToProto(tariff))
	if err != nil {
		return err
	}

	responses := receiver.client.ProduceSync(ctx, &kgo.Record{
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
