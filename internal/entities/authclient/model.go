package authclient

import (
	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

type (
	AuthClient struct {
		Id          string      `db:"id" json:"id"`
		Name        null.String `db:"name" json:"name"`
		Description null.String `db:"description" json:"description"`
		HomepageUrl null.String `db:"homepage_url" json:"homepage_url"`
		CallbackUrl null.String `db:"callback_url" json:"callback_url"`
	}

	AuthClientParam struct {
		// TODO: add any other custom params here
		h.Param
	}
)

const AuthClientConst = "auth_client"

func InitStorageV1() DataAuthClientV1Adapter {
	return InitSQLV1()
}
