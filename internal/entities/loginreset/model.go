package loginreset

import (
	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

type (
	LoginReset struct {
		LoginId    null.String `db:"login_id" json:"login_id"`
		ResetToken null.String `db:"reset_token" json:"reset_token"`
		CreatedAt  null.Time   `db:"created_at" json:"created_at"`
		UpdatedAt  null.Time   `db:"updated_at" json:"updated_at"`
	}

	LoginResetParam struct {
		// TODO: add any other custom params here
		h.Param
	}
)

const LoginResetConst = "login_reset"

func InitStorageV1() DataLoginResetV1Adapter {
	return InitSQLV1()
}
