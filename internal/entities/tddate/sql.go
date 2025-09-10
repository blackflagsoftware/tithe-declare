package tddate

import (
	"context"
	"fmt"
	"time"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	usql "github.com/blackflagsoftware/tithe-declare/internal/util/sql"
	"github.com/jmoiron/sqlx"
)

type (
	SQLTdDateV1 struct {
		DB *sqlx.DB
	}
)

func InitSQLV1() *SQLTdDateV1 {
	db := stor.InitStorage()
	return &SQLTdDateV1{DB: db}
}

func (d *SQLTdDateV1) Read(ctx context.Context, td_ *TdDate) error {
	sqlGet := `
		SELECT
			id,
			date_value,
			hold,
			confirm,
			name,
			phone,
			email
		FROM td_date WHERE id = $1`
	if errDB := d.DB.Get(td_, sqlGet, td_.Id, td_.DateValue, td_.Hold, td_.Confirm, td_.Name, td_.Phone, td_.Email); errDB != nil {
		return ae.DBError("TdDate Get: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLTdDateV1) ReadAll(ctx context.Context, td_ *[]TdDate, param TdDateParam) (int, error) {
	searchStmt, args := usql.BuildSearchString(param.Param, false) // false => include the where clause, see internal/util/sql.go
	sqlSearch := fmt.Sprintf(`
		SELECT
			id,
			date_value,
			hold,
			confirm,
			name,
			phone,
			email
		FROM td_date
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(td_, sqlSearch, args...); errDB != nil {
		return 0, ae.DBError("TdDate ReadAll: unable to select records.", errDB)
	}
	sqlCount := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM td_date
		%s`, searchStmt)
	var count int
	sqlCount = d.DB.Rebind(sqlCount)
	if errDB := d.DB.Get(&count, sqlCount, args...); errDB != nil {
		return 0, ae.DBError("td_date ReadAll: unable to select count.", errDB)
	}
	return count, nil
}

func (d *SQLTdDateV1) Create(ctx context.Context, td_ *TdDate) error {
	count, err := d.count()
	if err != nil {
		return err
	}
	td_.Id = count
	sqlPost := `
		INSERT INTO td_date (
			id,
			date_value,
			hold,
			confirm,
			name,
			phone,
			email
		) VALUES (
		 	:id,
			:date_value,
			:hold,
			:confirm,
			:name,
			:phone,
			:email
		)`
	_, errDB := d.DB.NamedExec(sqlPost, td_)
	if errDB != nil {
		return ae.DBError("TdDate Post: unable to insert record.", errDB)
	}

	return nil
}

func (d *SQLTdDateV1) Update(ctx context.Context, td_ TdDate) error {
	sqlPatch := `
		UPDATE td_date SET
			date_value = :date_value,
			hold = :hold,
			confirm = :confirm,
			name = :name,
			phone = :phone,
			email = :email
		WHERE id = :id`
	if _, errDB := d.DB.NamedExec(sqlPatch, td_); errDB != nil {
		return ae.DBError("TdDate Patch: unable to update record.", errDB)
	}
	return nil
}

func (d *SQLTdDateV1) Delete(ctx context.Context, td_ *TdDate) error {
	sqlDelete := `
		DELETE FROM td_date WHERE id = $1`
	if _, errDB := d.DB.Exec(sqlDelete, td_.Id, td_.DateValue, td_.Hold, td_.Confirm, td_.Name, td_.Phone, td_.Email); errDB != nil {
		return ae.DBError("TdDate Delete: unable to delete record.", errDB)
	}
	return nil
}

func (d *SQLTdDateV1) GetCurrentDays(ctx context.Context, dates *[]time.Time, param TdDateParam) error {
	searchStmt, args := usql.BuildSearchString(param.Param, false) // false => include the where clause, see internal/util/sql.go
	sqlSearch := fmt.Sprintf(`
		SELECT
			date_value
		FROM td_date
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(dates, sqlSearch, args...); errDB != nil {
		return ae.DBError("TdDate GetCurrentDays: unable to select records.", errDB)
	}
	return nil
}

// checks if the given dateTime is available to be held (i.e. not already held)
// if available, will hold the time
func (d *SQLTdDateV1) CheckSetHoldTime(ctx context.Context, dateTime time.Time) error {
	exists := true
	sqlExists := `
		SELECT EXISTS (
			SELECT 1 FROM td_date WHERE date_value = $1 AND hold IS NOT NULL
		)`
	if errDB := d.DB.Get(&exists, sqlExists, dateTime); errDB != nil {
		return ae.DBError("TdDate CheckHoldTime: unable to check hold time.", errDB)
	}
	if exists {
		return ae.HoldError()
	}
	sqlHold := `
		UPDATE td_date SET
			hold = $1
		WHERE date_value = $2 AND hold IS NULL`
	if _, errDB := d.DB.Exec(sqlHold, time.Now().UTC(), dateTime); errDB != nil {
		return ae.DBError("TdDate CheckHoldTime: unable to hold time.", errDB)
	}
	return nil
}

func (d *SQLTdDateV1) Confirm(ctx context.Context, dtDate TdDate) error {
	sqlConfirm := `
		UPDATE td_date SET
			confirm = :confirm,
			name = :name,
			email = :email,
			phone = :phone
		WHERE date_value = :date_value`
	if _, errDB := d.DB.NamedExec(sqlConfirm, dtDate); errDB != nil {
		return ae.DBError("TdDate Confirm: unable to confirm time.", errDB)
	}
	return nil
}

func (d *SQLTdDateV1) count() (int, error) {
	count := 0
	sqlCount := `
		SELECT
			COALESCE(MAX(id), 0) AS id
		FROM td_date`
	if errDB := d.DB.Get(&count, sqlCount); errDB != nil {
		return 0, ae.DBError("TdDate GetCount: unable to select count.", errDB)
	}
	return count + 1, nil
}
