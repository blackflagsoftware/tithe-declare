package loginrole

import (
	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

type (
	LoginRole struct {
		LoginId null.String `db:"login_id" json:"login_id"`
		RoleId  null.String `db:"role_id" json:"role_id"`
	}

	LoginRoleUpdate struct {
		LoginId null.String `db:"login_id" json:"login_id"`
		RoleIds []string    `db:"role_ids" json:"role_ids"` // the list of role IDs to assign to the login
	}
	LoginRoleParam struct {
		// TODO: add any other custom params here
		h.Param
	}
)

const LoginRoleConst = "login_role"

func InitStorageV1() DataLoginRoleV1Adapter {
	return InitSQLV1()
}
