package dataprocess

import (
	"encoding/json"
)

// ToParameterValue Convert WorkflowParameters to JSON
func (v WorkflowParameters) ToParameterValue() string {
	data, _ := json.Marshal(v.WorkflowParameter)
	return string(data)
}

// ToParameterValue Convert MetaQueryAggregations to JSON
func (v MetaQueryAggregations) ToParameterValue() string {
	data, _ := json.Marshal(v.Aggregations)
	return string(data)
}

// ToParameterValue Convert WithFields to JSON
func (v WithFields) ToParameterValue() string {
	data, _ := json.Marshal(v.WithField)
	return string(data)
}

// ToParameterValue Convert MetaQueryMediaTypes to JSON
func (v MetaQueryMediaTypes) ToParameterValue() string {
	data, _ := json.Marshal(v.MediaTypes)
	return string(data)
}

// ToParameterValue Convert SmartClusterRules to JSON
func (v SmartClusterRules) ToParameterValue() string {
	data, _ := json.Marshal(v.Rules)
	return string(data)
}

// ToParameterValue Convert SmartClusterNotification to JSON
func (v SmartClusterNotification) ToParameterValue() string {
	data, _ := json.Marshal(v)
	return string(data)
}

// ToParameterValue Convert DatasetConfig to JSON
func (v DatasetConfig) ToParameterValue() string {
	data, _ := json.Marshal(v)
	return string(data)
}

// ToParameterValue Convert SimpleQuery to JSON
func (v SimpleQuery) ToParameterValue() string {
	data, _ := json.Marshal(v)
	return string(data)
}
