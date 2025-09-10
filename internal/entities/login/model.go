package login

import (
	"time"

	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

type (
	Login struct {
		Id         string      `db:"id" json:"id"`
		EmailAddr  null.String `db:"email_addr" json:"email_address"`
		FirstName  null.String `db:"first_name" json:"first_name"`
		LastName   null.String `db:"last_name" json:"last_name"`
		Pwd        null.String `db:"pwd" json:"password"`
		ConfirmPwd null.String `json:"confirm_password,omitempty"`
		Active     null.Bool   `db:"active" json:"active"`
		SetPwd     null.Bool   `db:"set_pwd" json:"set_password"`
		CreatedAt  null.Time   `db:"created_at" json:"created_at"`
		UpdatedAt  null.Time   `db:"updated_at" json:"updated_at"`
	}

	LoginParam struct {
		// TODO: add any other custom params here
		h.Param
	}

	ResetRequest struct {
		EmailAddr  string    `db:"email_address" json:"email_address"`
		LoginId    string    `db:"login_id"`
		ResetToken string    `db:"reset_token"`
		CreatedAt  time.Time `db:"created_at"`
	}

	PasswordReset struct {
		EmailAddr  string      `json:"email_address"`
		ResetToken string      `json:"reset_token"`
		Pwd        null.String `json:"password"`
		ConfirmPwd null.String `json:"confirm_password"`
	}

	LoginRoles struct {
		LoginId   string   `db:"id" json:"id"`
		EmailAddr string   `db:"email_addr" json:"email_address"`
		Roles     []string `db:"roles" json:"roles"`
	}
)

const LoginConst = "login"

func InitStorage() DataLoginV1Adapter {
	return InitSQLV1()
}
