package dto

import "github.com/yeahyeahcore/redpanda-study/internal/models"

type Message struct {
	PrivilegeArray []models.Privilege `json:"privilegeArray"`
}
