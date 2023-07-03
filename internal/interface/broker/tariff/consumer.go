package tariff

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/yeahyeahcore/redpanda-study/internal/interface/broker/tariff/dto"
	"github.com/yeahyeahcore/redpanda-study/internal/models"
	"github.com/yeahyeahcore/redpanda-study/pkg/json"
	"go.uber.org/zap"
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

	messages := make([]dto.Message, 0)
	privilegies := make([]models.Privilege, 0)

	for !iter.Done() {
		record := iter.Next()

		fmt.Println(string(record.Key))

		message, err := json.Unmarshal[dto.Message](record.Value)
		if err != nil {
			receiver.logger.Error(fmt.Sprintf("error decoding message: %v\n", err))
			continue
		}

		messages = append(messages, *message)
	}

	for _, message := range messages {
		privilegies = append(privilegies, message.PrivilegeArray...)
	}

	return privilegies
}

func (receiver *Consumer) Close(_ context.Context) error {
	receiver.client.Close()

	return nil
}
