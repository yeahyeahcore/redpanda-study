package main

import (
	"log"

	"github.com/yeahyeahcore/redpanda-study/internal/app"
	"github.com/yeahyeahcore/redpanda-study/internal/config"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		log.Fatalf("failed to init zap logger: %v", err)
	}

	config, err := config.Initialize("config.dev.json", ".env.test")
	if err != nil {
		logger.Fatal("failed to init config", zap.Error(err))
	}

	if err := app.Run(config, logger); err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}
