package env

import (
	"context"
	"fmt"
	"log"

	envLoader "github.com/joho/godotenv"
	envParser "github.com/sethvargo/go-envconfig"
)

func Parse[T interface{}](envFilePath string) (*T, error) {
	var configuration T

	if err := envLoader.Load(envFilePath); err != nil {
		log.Printf("load local env file error: %s", err.Error())
	}
	if err := envParser.Process(context.Background(), &configuration); err != nil {
		return nil, fmt.Errorf("parsing env error: %s", err.Error())
	}

	return &configuration, nil
}
