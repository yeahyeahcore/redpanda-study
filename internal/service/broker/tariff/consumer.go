package tariff

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/yeahyeahcore/redpanda-study/internal/service/broker/tariff/dto"
	"go.uber.org/zap"
)

type ConsumerDeps struct {
	Logger  *zap.Logger
	Brokers []string
	Topic   string
	GroupID string
}

type Consumer struct {
	logger *zap.Logger
	client *kgo.Client
	topic  string
}

func NewConsumer(deps *ConsumerDeps) (*Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(deps.Brokers...),
		kgo.ConsumerGroup(deps.GroupID),
		kgo.ConsumeTopics(deps.Topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{client: client, topic: deps.Topic}, nil
}

func (receiver *Consumer) PrintMessages(ctx context.Context) {
	for {
		fetches := receiver.client.PollFetches(ctx)
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()

			var message dto.Message

			if err := json.Unmarshal(record.Value, &message); err != nil {
				receiver.logger.Error(fmt.Sprintf("error decoding message: %v\n", err))
				continue
			}

			for _, privilege := range message.PrivilegeArray {
				fmt.Printf("privilegeID: %s\n serviceKey: %s\n type: %s\n value: %s\n\n\n",
					privilege.PrivilegeID,
					privilege.ServiceKey,
					privilege.Value,
					privilege.Value,
				)
			}
		}
	}
}

func (receiver *Consumer) Close(_ context.Context) error {
	receiver.client.Close()

	return nil
}
