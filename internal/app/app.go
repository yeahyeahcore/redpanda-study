package app

import (
	"context"
	"errors"
	"os/signal"
	"syscall"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/yeahyeahcore/redpanda-study/internal/config"
	"github.com/yeahyeahcore/redpanda-study/internal/initialize"
	"github.com/yeahyeahcore/redpanda-study/internal/server"
	"github.com/yeahyeahcore/redpanda-study/internal/worker"
	"github.com/yeahyeahcore/redpanda-study/pkg/closer"
	"github.com/yeahyeahcore/redpanda-study/pkg/kafka"
	"go.uber.org/zap"
)

const (
	shutdownTimeout = 5 * time.Second
)

func Run(config *config.Config, logger *zap.Logger) (appErr error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			appErr = errors.New("recovered panic on application run")
			logger.Error("recovered panic on application run", zap.Any("recovered", recovered))
		}
	}()

	closerUtil := closer.New()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := kafka.Initialize(ctx, config.Service.Kafka.Brokers, []string{
		config.Service.Kafka.Tariff.Topic,
	}); err != nil {
		return err
	}

	kafkaTariffClient, err := kgo.NewClient(
		kgo.SeedBrokers(config.Service.Kafka.Brokers...),
		kgo.ConsumerGroup(config.Service.Kafka.ConsumerGroupID),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		kgo.DefaultProduceTopic(config.Service.Kafka.Tariff.Topic),
		kgo.ConsumeTopics(config.Service.Kafka.Tariff.Topic),
	)
	if err != nil {
		return err
	}

	brokers := initialize.NewBrokers(initialize.BrokersDeps{
		Logger:       logger,
		TariffClient: kafkaTariffClient,
	})

	services := initialize.NewServices(initialize.ServicesDeps{
		Logger:  logger,
		Brokers: *brokers,
	})

	workers := initialize.NewWorkers(initialize.WorkersDeps{
		Logger:   logger,
		Services: *services,
	})

	controllers := initialize.NewControllers(initialize.ControllersDeps{
		Logger:   logger,
		Services: *services,
	})

	serverHTTP := server.NewHTTP(server.DepsHTTP{
		Logger: logger,
	})

	serverHTTP.Register(controllers)

	go serverHTTP.Run(&config.HTTP)
	go worker.Run(ctx, workers)

	closerUtil.Add(serverHTTP.Stop)
	closerUtil.Add(closer.Wrap(kafkaTariffClient.Close))

	<-ctx.Done()

	logger.Info("shutting down app gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	defer cancel()

	if err := closerUtil.Close(shutdownCtx); err != nil {
		logger.Error("closer error", zap.Error(err))
		return err
	}

	return nil
}
