package models

type Privilege struct {
	PrivilegeID string `json:"privilegeId"`
	ServiceKey  string `json:"serviceKey"`
	Type        string `json:"type"`
	Value       string `json:"value"`
}
