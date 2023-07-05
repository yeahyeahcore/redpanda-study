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
	Logger *zap.Logger
	Client *kgo.Client
}

type Consumer struct {
	logger *zap.Logger
	client *kgo.Client
}

func NewConsumer(deps ConsumerDeps) *Consumer {
	return &Consumer{
		logger: deps.Logger,
		client: deps.Client,
	}
}

func (receiver *Consumer) FetchPrivileges(ctx context.Context) []models.Privilege {
	fetches := receiver.client.PollFetches(ctx)
	iter := fetches.RecordIter()

	tariffs := make([]models.Tariff, 0)
	privileges := make([]models.Privilege, 0)

	for !iter.Done() {
		record := iter.Next()
		tariff := tariff.Tariff{}

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
