package authauthorize

import (
	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

type (
	AuthAuthorize struct {
		Id                   string      `db:"id" json:"id"`
		ClientId             null.String `db:"client_id" json:"client_id"`
		Verifier             null.String `db:"verifier" json:"verifier"`
		VerifierEncodeMethod null.String `db:"verifier_encode_method" json:"verifier_encode_method"`
		State                null.String `db:"state" json:"state"`
		Scope                null.String `db:"scope" json:"scope"`
		AuthorizedAt         null.Time   `db:"authorized_at" json:"authorized_at"`
		AuthCodeAt           null.Time   `db:"auth_code_at" json:"auth_code_at"`
		AuthCode             null.String `db:"auth_code" json:"auth_code"`
	}

	AuthAuthorizeParam struct {
		// TODO: add any other custom params here
		h.Param
	}
)

const AuthAuthorizeConst = "auth_authorize"

func InitStorageV1() DataAuthAuthorizeV1Adapter {
	return InitSQLV1()
}
