package authclientsecret

import (
	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

type (
	AuthClientSecret struct {
		Id       string      `db:"id" json:"id"`
		ClientId null.String `db:"client_id" json:"client_id"`
		Secret   null.String `db:"secret" json:"secret"`
		Active   bool        `db:"active" json:"active"`
	}

	AuthClientSecretParam struct {
		// TODO: add any other custom params here
		h.Param
	}
)

const AuthClientSecretConst = "auth_client_secret"

func InitStorageV1() DataAuthClientSecretV1Adapter {
	return InitSQLV1()
}
