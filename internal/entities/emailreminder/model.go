package emailreminder

import (
	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

type (
	EmailReminder struct {
		Id    int         `db:"id" json:"id"`
		Email null.String `db:"email" json:"email"`
	}

	EmailReminderParam struct {
		// TODO: add any other custom params here
		h.Param
	}
)

const EmailReminderConst = "email_reminder"

func InitStorageV1() DataEmailReminderV1Adapter {
	return InitSQLV1()
}
