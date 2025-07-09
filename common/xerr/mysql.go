package xerr

import (
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	MysqlCodeDuplicate = 1062
)

func AssertMysqlErrorCode(code uint16, err error) bool {
	me, ok := errors.Cause(err).(*mysql.MySQLError)
	if !ok {
		return false
	}
	return me.Number == code
}
