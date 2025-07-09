// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Node is the golang structure for table node.
type Node struct {
	Id                int         `json:"id"                 orm:"id"                 ` //
	Name              string      `json:"name"               orm:"name"               ` //
	Conditions        string      `json:"conditions"         orm:"conditions"         ` //
	Capacity          string      `json:"capacity"           orm:"capacity"           ` //
	Allocatable       string      `json:"allocatable"        orm:"allocatable"        ` //
	Addresses         string      `json:"addresses"          orm:"addresses"          ` //
	CreationTimestamp string      `json:"creation_timestamp" orm:"creation_timestamp" ` //
	CreatedAt         *gtime.Time `json:"created_at"         orm:"created_at"         ` //
	UpdatedAt         *gtime.Time `json:"updated_at"         orm:"updated_at"         ` //
	DeletedAt         *gtime.Time `json:"deleted_at"         orm:"deleted_at"         ` //
	ClusterId         uint64      `json:"cluster_id"         orm:"cluster_id"         ` //
	CpuUsage          float64     `json:"cpu_usage"          orm:"cpu_usage"          ` //
	MemoryUsage       float64     `json:"memory_usage"       orm:"memory_usage"       ` //
	Status            int         `json:"status"             orm:"status"             ` //
	Labels            string      `json:"labels"             orm:"labels"             ` //
	Annotations       string      `json:"annotations"        orm:"annotations"        ` //
	NodeInfo          string      `json:"node_info"          orm:"node_info"          ` //
	Roles             string      `json:"roles"              orm:"roles"              ` //
	Uid               string      `json:"uid"                orm:"uid"                ` //
	ResourceVersion   string      `json:"resource_version"   orm:"resource_version"   ` //
	HealthState       string      `json:"health_state"       orm:"health_state"       ` //
	Spec              string      `json:"spec"               orm:"spec"               ` //
}
