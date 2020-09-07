package temporal_experiments

// Move is moving-units payload.
type Move struct {
	Type        string `json:"type"`
	Value       int64  `json:"value"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
}

// DataUpdatedMessage provides information about type and scope of changed inventory for company.
type DataUpdatedMessage struct {
	Type      string           `json:"type"`
	CompanyID string           `json:"companyId"`
	Scope     DataUpdatedScope `json:"scope"`
}

// DataUpdatedScope defines the scope of changed inventory.
type DataUpdatedScope struct {
	GroupIDList   []string `json:"groupIds,omitempty"`
	PackageIDList []string `json:"packageIds,omitempty"`
	ProductIDList []string `json:"productIds,omitempty"`
}
