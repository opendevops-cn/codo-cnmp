// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProxyAgent is the golang structure of table proxy_agent for DAO operations like Where/Data.
type ProxyAgent struct {
	g.Meta    `orm:"table:proxy_agent, do:true"`
	Id        interface{} //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
	Name      interface{} //
	AgentId   interface{} //
}
