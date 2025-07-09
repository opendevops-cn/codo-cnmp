// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProxyAgent is the golang structure for table proxy_agent.
type ProxyAgent struct {
	Id        uint64      `json:"id"         orm:"id"         ` //
	CreatedAt *gtime.Time `json:"created_at" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updated_at" orm:"updated_at" ` //
	DeletedAt *gtime.Time `json:"deleted_at" orm:"deleted_at" ` //
	Name      string      `json:"name"       orm:"name"       ` //
	AgentId   string      `json:"agent_id"   orm:"agent_id"   ` //
}
