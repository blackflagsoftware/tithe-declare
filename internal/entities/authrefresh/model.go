package authrefresh

import (
	"time"

	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
)

type (
	AuthRefresh struct {
		ClientId  string    `db:"client_id" json:"client_id"`
		Token     string    `db:"token" json:"token"`
		Active    bool      `db:"active" json:"active"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
	}

	AuthRefreshParam struct {
		// TODO: add any other custom params here
		h.Param
	}
)

const AuthRefreshConst = "auth_refresh"

func InitStorageV1() DataAuthRefreshV1Adapter {
	return InitSQLV1()
}
