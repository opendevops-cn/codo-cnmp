// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Node is the golang structure of table node for DAO operations like Where/Data.
type Node struct {
	g.Meta            `orm:"table:node, do:true"`
	Id                interface{} //
	Name              interface{} //
	Conditions        interface{} //
	Capacity          interface{} //
	Allocatable       interface{} //
	Addresses         interface{} //
	CreationTimestamp interface{} //
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
	DeletedAt         *gtime.Time //
	ClusterId         interface{} //
	CpuUsage          interface{} //
	MemoryUsage       interface{} //
	Status            interface{} //
	Labels            interface{} //
	Annotations       interface{} //
	NodeInfo          interface{} //
	Roles             interface{} //
	Uid               interface{} //
	ResourceVersion   interface{} //
	HealthState       interface{} //
	Spec              interface{} //
}
