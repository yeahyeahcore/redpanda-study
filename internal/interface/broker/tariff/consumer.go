package tariff

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/yeahyeahcore/redpanda-study/internal/models"
	"github.com/yeahyeahcore/redpanda-study/internal/proto/tariff"
	"github.com/yeahyeahcore/redpanda-study/internal/utils/transfer"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type ConsumerDeps struct {
	Logger     *zap.Logger
	Brokers    []string
	Topic      string
	GroupID    string
	MessageKey string
}

type Consumer struct {
	logger     *zap.Logger
	client     *kgo.Client
	topic      string
	messageKey string
}

func NewConsumer(deps ConsumerDeps) (*Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(deps.Brokers...),
		kgo.ConsumerGroup(deps.GroupID),
		kgo.ConsumeTopics(deps.Topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client:     client,
		topic:      deps.Topic,
		messageKey: deps.MessageKey,
	}, nil
}

func (receiver *Consumer) FetchMessages(ctx context.Context) []models.Privilege {
	fetches := receiver.client.PollFetches(ctx)
	iter := fetches.RecordIter()

	tariffs := make([]models.Tariff, 0)
	privileges := make([]models.Privilege, 0)

	for !iter.Done() {
		record := iter.Next()
		tariff := tariff.Tariff{}

		fmt.Println(string(record.Key))

		if err := proto.Unmarshal(record.Value, &tariff); err != nil {
			receiver.logger.Error(fmt.Sprintf("error decoding message: %v\n", err))
			continue
		}

		tariffs = append(tariffs, *transfer.TariffProtoToModel(&tariff))
	}

	for _, tariff := range tariffs {
		privileges = append(privileges, tariff.Privileges...)
	}

	return privileges
}

func (receiver *Consumer) Close(_ context.Context) error {
	receiver.client.Close()

	return nil
}
