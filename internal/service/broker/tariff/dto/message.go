package dto

type Privilege struct {
	PrivilegeID string `json:"privilegeId"`
	ServiceKey  string `json:"serviceKey"`
	Type        string `json:"type"`
	Value       string `json:"value"`
}

type Message struct {
	PrivilegeArray []Privilege `json:"privilegeArray"`
}
