// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Cluster is the golang structure of table cluster for DAO operations like Where/Data.
type Cluster struct {
	g.Meta        `orm:"table:cluster, do:true"`
	Id            interface{} //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	DeletedAt     *gtime.Time //
	Name          interface{} //
	Description   interface{} //
	ImportType    interface{} //
	ImportDetail  interface{} //
	Status        interface{} //
	ServerVersion interface{} //
	Platform      interface{} //
	BuildDate     interface{} //
	ExtInfo       interface{} //
	NodeState     interface{} //
	HealthState   interface{} //
	CpuUsage      interface{} //
	MemoryUsage   interface{} //
	CpuTotal      interface{} //
	MemoryTotal   interface{} //
	NodeCount     interface{} //
	ClusterState  interface{} //
	Uid           interface{} //
	Idip          interface{} //
	AppId         interface{} //
	AppSecret     interface{} //
	Ops           interface{} //
	DstAgentId    interface{} //
	ConnetType    interface{} //
	MeshAddr      interface{} //
	Links         interface{} //
}
