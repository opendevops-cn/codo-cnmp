// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CodoClusterV1Dao is the data access object for table codo_cluster_v1.
type CodoClusterV1Dao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns CodoClusterV1Columns // columns contains all the column names of Table for convenient usage.
}

// CodoClusterV1Columns defines and stores column names for table codo_cluster_v1.
type CodoClusterV1Columns struct {
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
}

// codoClusterV1Columns holds the columns for table codo_cluster_v1.
var codoClusterV1Columns = CodoClusterV1Columns{
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
}

// NewCodoClusterV1Dao creates and returns a new DAO object for table data access.
func NewCodoClusterV1Dao() *CodoClusterV1Dao {
	return &CodoClusterV1Dao{
		group:   "default",
		table:   "codo_cluster_v1",
		columns: codoClusterV1Columns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CodoClusterV1Dao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CodoClusterV1Dao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CodoClusterV1Dao) Columns() CodoClusterV1Columns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CodoClusterV1Dao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CodoClusterV1Dao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CodoClusterV1Dao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
