package role

import (
	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

type (
	Role struct {
		Id          string      `db:"id" json:"id"`
		Name        null.String `db:"name" json:"name"`
		Description null.String `db:"description" json:"description"`
	}

	RoleParam struct {
		// TODO: add any other custom params here
		h.Param
	}
)

const RoleConst = "role"

func InitStorageV1() DataRoleV1Adapter {
	return InitSQLV1()
}
