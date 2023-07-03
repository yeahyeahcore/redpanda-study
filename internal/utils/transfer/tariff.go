package transfer

import (
	"github.com/yeahyeahcore/redpanda-study/internal/models"
	"github.com/yeahyeahcore/redpanda-study/internal/proto/tariff"
)

func TariffProtoToModel(tariff *tariff.Tariff) *models.Tariff {
	if tariff == nil {
		return &models.Tariff{}
	}

	return &models.Tariff{Privileges: PrivilegeArrayProtoToModel(tariff.Privileges)}
}

func TariffModelToProto(tariffModel *models.Tariff) *tariff.Tariff {
	if tariffModel == nil {
		return &tariff.Tariff{}
	}

	return &tariff.Tariff{Privileges: PrivilegeArrayModelToProto(tariffModel.Privileges)}
}
