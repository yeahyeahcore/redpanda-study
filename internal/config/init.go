package config

import (
	"github.com/yeahyeahcore/redpanda-study/pkg/env"
	"github.com/yeahyeahcore/redpanda-study/pkg/json"
)

func Initialize(pathJSON, pathENV string) (*Config, error) {
	config, err := env.Parse[Config](pathENV)
	if err != nil {
		return nil, err
	}

	config, err = json.Read[Config](pathJSON)
	if err != nil {
		return nil, err
	}

	return config, nil
}
