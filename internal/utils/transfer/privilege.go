package transfer

import (
	"github.com/yeahyeahcore/redpanda-study/internal/models"
	"github.com/yeahyeahcore/redpanda-study/internal/proto/tariff"
)

func PrivilegeProtoToModel(privilege *tariff.Privilege) *models.Privilege {
	if privilege == nil {
		return &models.Privilege{}
	}

	return &models.Privilege{
		PrivilegeID: privilege.GetPrivilegeID(),
		ServiceKey:  privilege.GetServiceKey(),
		Type:        privilege.GetType(),
		Value:       privilege.GetValue(),
	}
}

func PrivilegeArrayProtoToModel(privileges []*tariff.Privilege) []models.Privilege {
	privilegesModel := make([]models.Privilege, len(privileges))

	for index, privilege := range privileges {
		privilegesModel[index] = *PrivilegeProtoToModel(privilege)
	}

	return privilegesModel
}

func PrivilegeModelToProto(privilege *models.Privilege) *tariff.Privilege {
	if privilege == nil {
		return &tariff.Privilege{}
	}

	return &tariff.Privilege{
		PrivilegeID: privilege.PrivilegeID,
		ServiceKey:  privilege.ServiceKey,
		Type:        privilege.Type,
		Value:       privilege.Value,
	}
}

func PrivilegeArrayModelToProto(privileges []models.Privilege) []*tariff.Privilege {
	privilegesProto := make([]*tariff.Privilege, len(privileges))

	for index, privilege := range privileges {
		privilegeCopy := privilege
		privilegesProto[index] = PrivilegeModelToProto(&privilegeCopy)
	}

	return privilegesProto
}
