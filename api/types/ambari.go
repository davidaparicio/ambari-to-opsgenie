package types

type (
	// ResponseAlert definition on JSON.
	ResponseAlert struct {
		Href  string `json:"href"`
		Items []Item `json:"items"`
	}

	// Item definition on JSON.ResponseAlert.
	Item struct {
		Href  string `json:"href"`
		Alert Alert  `json:"Alert"`
	}

	// Alert definition on JSON.ResponseAlert.Item.
	Alert struct {
		ClusterName   string `json:"cluster_name" mapstructure:"topic_prefix"`
		ComponentName string `json:"component_name" mapstructure:"component_name"`
		// https://go.googlesource.com/lint | We will not be adding pragmas
		// or other knobs to suppress specific warnings
		// OK: struct field DefinitionId should be DefinitionID
		DefinitionId   int    `json:"definition_id" mapstructure:"definition_id"`
		DefinitionName string `json:"definition_name" mapstructure:"definition_name"`
		DostName       string `json:"host_name" mapstructure:"host_name"`
		//OK: struct field Id should be ID
		Id                int    `json:"id"`
		Label             string `json:"label"`
		LatestTimestamp   int64  `json:"latest_timestamp" mapstructure:"latest_timestamp"`
		MaintenanceState  string `json:"maintenance_state" mapstructure:"maintenance_state"`
		OriginalTimestamp int64  `json:"original_timestamp" mapstructure:"original_timestamp"`
		Scope             string `json:"scope"`
		ServiceName       string `json:"service_name" mapstructure:"service_name"`
		State             string `json:"state" mapstructure:""`
		Text              string `json:"text"`
	}
)
