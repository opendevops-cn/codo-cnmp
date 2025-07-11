// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"codo-cnmp/internal/model/dao/internal"
)

// internalUserGroupDao is internal type for wrapping internal DAO implements.
type internalUserGroupDao = *internal.UserGroupDao

// userGroupDao is the data access object for table user_group.
// You can define custom methods on it to extend its functionality as you wish.
type userGroupDao struct {
	internalUserGroupDao
}

var (
	// UserGroup is globally public accessible object for table user_group operations.
	UserGroup = userGroupDao{
		internal.NewUserGroupDao(),
	}
)

// Fill with you ideas below.
