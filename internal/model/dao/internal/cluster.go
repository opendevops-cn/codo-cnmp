// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ClusterDao is the data access object for table cluster.
type ClusterDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns ClusterColumns // columns contains all the column names of Table for convenient usage.
}

// ClusterColumns defines and stores column names for table cluster.
type ClusterColumns struct {
	Id            string //
	CreatedAt     string //
	UpdatedAt     string //
	DeletedAt     string //
	Name          string //
	Description   string //
	ImportType    string //
	ImportDetail  string //
	Status        string //
	ServerVersion string //
	Platform      string //
	BuildDate     string //
	ExtInfo       string //
	NodeState     string //
	HealthState   string //
	CpuUsage      string //
	MemoryUsage   string //
	CpuTotal      string //
	MemoryTotal   string //
	NodeCount     string //
	ClusterState  string //
	Uid           string //
	Idip          string //
	AppId         string //
	AppSecret     string //
	Ops           string //
	DstAgentId    string //
	ConnetType    string //
	MeshAddr      string //
	Links         string //
}

// clusterColumns holds the columns for table cluster.
var clusterColumns = ClusterColumns{
	Id:            "id",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
	Name:          "name",
	Description:   "description",
	ImportType:    "import_type",
	ImportDetail:  "import_detail",
	Status:        "status",
	ServerVersion: "server_version",
	Platform:      "platform",
	BuildDate:     "build_date",
	ExtInfo:       "ext_info",
	NodeState:     "node_state",
	HealthState:   "health_state",
	CpuUsage:      "cpu_usage",
	MemoryUsage:   "memory_usage",
	CpuTotal:      "cpu_total",
	MemoryTotal:   "memory_total",
	NodeCount:     "node_count",
	ClusterState:  "cluster_state",
	Uid:           "uid",
	Idip:          "idip",
	AppId:         "app_id",
	AppSecret:     "app_secret",
	Ops:           "ops",
	DstAgentId:    "dst_agent_id",
	ConnetType:    "connet_type",
	MeshAddr:      "mesh_addr",
	Links:         "links",
}

// NewClusterDao creates and returns a new DAO object for table data access.
func NewClusterDao() *ClusterDao {
	return &ClusterDao{
		group:   "default",
		table:   "cluster",
		columns: clusterColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ClusterDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ClusterDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ClusterDao) Columns() ClusterColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ClusterDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ClusterDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ClusterDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
