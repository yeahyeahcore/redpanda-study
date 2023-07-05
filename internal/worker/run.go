package worker

import (
	"context"

	"github.com/yeahyeahcore/redpanda-study/internal/initialize"
)

func Run(ctx context.Context, workers *initialize.Workers) {
	go workers.TariffBroker.Run(ctx)
}
