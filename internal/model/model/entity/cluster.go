// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Cluster is the golang structure for table cluster.
type Cluster struct {
	Id            uint64      `json:"id"             orm:"id"             ` //
	CreatedAt     *gtime.Time `json:"created_at"     orm:"created_at"     ` //
	UpdatedAt     *gtime.Time `json:"updated_at"     orm:"updated_at"     ` //
	DeletedAt     *gtime.Time `json:"deleted_at"     orm:"deleted_at"     ` //
	Name          string      `json:"name"           orm:"name"           ` //
	Description   string      `json:"description"    orm:"description"    ` //
	ImportType    int         `json:"import_type"    orm:"import_type"    ` //
	ImportDetail  string      `json:"import_detail"  orm:"import_detail"  ` //
	Status        int         `json:"status"         orm:"status"         ` //
	ServerVersion string      `json:"server_version" orm:"server_version" ` //
	Platform      string      `json:"platform"       orm:"platform"       ` //
	BuildDate     string      `json:"build_date"     orm:"build_date"     ` //
	ExtInfo       string      `json:"ext_info"       orm:"ext_info"       ` //
	NodeState     int         `json:"node_state"     orm:"node_state"     ` //
	HealthState   string      `json:"health_state"   orm:"health_state"   ` //
	CpuUsage      float64     `json:"cpu_usage"      orm:"cpu_usage"      ` //
	MemoryUsage   float64     `json:"memory_usage"   orm:"memory_usage"   ` //
	CpuTotal      float64     `json:"cpu_total"      orm:"cpu_total"      ` //
	MemoryTotal   float64     `json:"memory_total"   orm:"memory_total"   ` //
	NodeCount     uint        `json:"node_count"     orm:"node_count"     ` //
	ClusterState  int         `json:"cluster_state"  orm:"cluster_state"  ` //
	Uid           string      `json:"uid"            orm:"uid"            ` //
	Idip          string      `json:"idip"           orm:"idip"           ` //
	AppId         string      `json:"app_id"         orm:"app_id"         ` //
	AppSecret     string      `json:"app_secret"     orm:"app_secret"     ` //
	Ops           string      `json:"ops"            orm:"ops"            ` //
	DstAgentId    int64       `json:"dst_agent_id"   orm:"dst_agent_id"   ` //
	ConnetType    uint        `json:"connet_type"    orm:"connet_type"    ` //
	MeshAddr      string      `json:"mesh_addr"      orm:"mesh_addr"      ` //
	Links         string      `json:"links"          orm:"links"          ` //
}
