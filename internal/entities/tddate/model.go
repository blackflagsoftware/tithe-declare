package tddate

import (
	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

type (
	TdDate struct {
		Id        int         `db:"id" json:"id"`
		DateValue null.Time   `db:"date_value" json:"date_value"`
		Hold      null.Time   `db:"hold" json:"hold"`
		Confirm   null.Time   `db:"confirm" json:"confirm"`
		Name      null.String `db:"name" json:"name"`
		Phone     null.String `db:"phone" json:"phone"`
		Email     null.String `db:"email" json:"email"`
	}

	TdDateParam struct {
		// TODO: add any other custom params here
		h.Param
	}

	TdDateBlock struct {
		NewDate   null.String `json:"new_date"`
		StartTime null.String `json:"start_time"`
		EndTime   null.String `json:"end_time"`
	}

	CurrentDateTime struct {
		DayAndTimes map[string][]string `json:"day_and_times"`
	}

	CheckHoldTimeRequest struct {
		Date string `json:"date"`
		Time string `json:"time"`
	}

	ConfirmRequest struct {
		TdDate
		CheckHoldTimeRequest
	}
)

const TdDateConst = "td_date"

func InitStorageV1() DataTdDateV1Adapter {
	return InitSQLV1()
}
