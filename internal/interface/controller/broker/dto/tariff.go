package dto

import "github.com/yeahyeahcore/redpanda-study/internal/models"

type Tariff struct {
	Privileges []models.Privilege `json:"privileges"`
}
