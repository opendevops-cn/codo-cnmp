// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CodoClusterDao is the data access object for table codo_cluster.
type CodoClusterDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns CodoClusterColumns // columns contains all the column names of Table for convenient usage.
}

// CodoClusterColumns defines and stores column names for table codo_cluster.
type CodoClusterColumns struct {
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
}

// codoClusterColumns holds the columns for table codo_cluster.
var codoClusterColumns = CodoClusterColumns{
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
}

// NewCodoClusterDao creates and returns a new DAO object for table data access.
func NewCodoClusterDao() *CodoClusterDao {
	return &CodoClusterDao{
		group:   "default",
		table:   "codo_cluster",
		columns: codoClusterColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CodoClusterDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CodoClusterDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CodoClusterDao) Columns() CodoClusterColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CodoClusterDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CodoClusterDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CodoClusterDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
