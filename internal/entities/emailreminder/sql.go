package emailreminder

import (
	"context"
	"fmt"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	usql "github.com/blackflagsoftware/tithe-declare/internal/util/sql"
	"github.com/jmoiron/sqlx"
)

type (
	SQLEmailReminderV1 struct {
		DB *sqlx.DB
	}
)

func InitSQLV1() *SQLEmailReminderV1 {
	db := stor.InitStorage()
	return &SQLEmailReminderV1{DB: db}
}

func (d *SQLEmailReminderV1) Read(ctx context.Context, ema *EmailReminder) error {
	sqlGet := `
		SELECT
			id,
			email
		FROM email_reminder WHERE id = $1`
	if errDB := d.DB.Get(ema, sqlGet, ema.Id, ema.Email); errDB != nil {
		return ae.DBError("EmailReminder Get: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLEmailReminderV1) ReadAll(ctx context.Context, ema *[]EmailReminder, param EmailReminderParam) (int, error) {
	searchStmt, args := usql.BuildSearchString(param.Param, false) // false => include the where clause, see internal/util/sql.go
	sqlSearch := fmt.Sprintf(`
		SELECT
			id,
			email
		FROM email_reminder
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(ema, sqlSearch, args...); errDB != nil {
		return 0, ae.DBError("EmailReminder ReadAll: unable to select records.", errDB)
	}
	sqlCount := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM email_reminder
		%s`, searchStmt)
	var count int
	sqlCount = d.DB.Rebind(sqlCount)
	if errDB := d.DB.Get(&count, sqlCount, args...); errDB != nil {
		return 0, ae.DBError("email_reminder ReadAll: unable to select count.", errDB)
	}
	return count, nil
}

func (d *SQLEmailReminderV1) Create(ctx context.Context, ema *EmailReminder) error {
	count, errCount := d.count()
	if errCount != nil {
		return errCount
	}
	ema.Id = count
	sqlPost := `
		INSERT INTO email_reminder (
			id,
			email
		) VALUES (
			:id,
			:email
		)`
	_, errDB := d.DB.NamedExec(sqlPost, ema)
	if errDB != nil {
		return ae.DBError("EmailReminder Post: unable to insert record.", errDB)
	}

	return nil
}

func (d *SQLEmailReminderV1) Update(ctx context.Context, ema EmailReminder) error {
	sqlPatch := `
		UPDATE email_reminder SET
			id = :id,
			email = :email
		WHERE id = :id`
	if _, errDB := d.DB.NamedExec(sqlPatch, ema); errDB != nil {
		return ae.DBError("EmailReminder Patch: unable to update record.", errDB)
	}
	return nil
}

func (d *SQLEmailReminderV1) Delete(ctx context.Context, ema *EmailReminder) error {
	sqlDelete := `
		DELETE FROM email_reminder WHERE id = $1`
	if _, errDB := d.DB.Exec(sqlDelete, ema.Id, ema.Email); errDB != nil {
		return ae.DBError("EmailReminder Delete: unable to delete record.", errDB)
	}
	return nil
}

func (d *SQLEmailReminderV1) count() (int, error) {
	count := 0
	if errDB := d.DB.Get(&count, "SELECT COALESCE(MAX(id), 0) FROM email_reminder"); errDB != nil {
		return 0, ae.DBError("EmailReminder count: unable to get count.", errDB)
	}
	return count + 1, nil
}
