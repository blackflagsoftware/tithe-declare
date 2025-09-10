package authclientcallback

import (
	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

type (
	AuthClientCallback struct {
		ClientId    null.String `db:"client_id" json:"client_id"`
		CallbackUrl null.String `db:"callback_url" json:"callback_url"`
	}

	AuthClientCallbackParam struct {
		// TODO: add any other custom params here
		h.Param
	}
)

const AuthClientCallbackConst = "auth_client_callback"

func InitStorageV1() DataAuthClientCallbackV1Adapter {
	return InitSQLV1()
}
