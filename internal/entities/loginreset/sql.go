package loginreset

import (
	"context"
	"fmt"
	"strings"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	usql "github.com/blackflagsoftware/tithe-declare/internal/util/sql"
	"github.com/jmoiron/sqlx"
)

type (
	SQLLoginResetV1 struct {
		DB *sqlx.DB
	}
)

func InitSQLV1() *SQLLoginResetV1 {
	db := stor.InitStorage()
	return &SQLLoginResetV1{DB: db}
}

func (d *SQLLoginResetV1) Read(ctx context.Context, lo *LoginReset) error {
	sqlGet := `
		SELECT
			login_id,
			reset_token,
			created_at,
			updated_at
		FROM login_reset WHERE login_id = $1 AND reset_token = $2`
	if errDB := d.DB.Get(lo, sqlGet, strings.ToLower(lo.LoginId.String), strings.ToLower(lo.ResetToken.String), lo.CreatedAt, lo.UpdatedAt); errDB != nil {
		return ae.DBError("LoginReset Get: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLLoginResetV1) ReadAll(ctx context.Context, lo *[]LoginReset, param LoginResetParam) (int, error) {
	searchStmt, args := usql.BuildSearchString(param.Param, false) // false => include the where clause, see internal/util/sql.go
	sqlSearch := fmt.Sprintf(`
		SELECT
			login_id,
			reset_token,
			created_at,
			updated_at
		FROM login_reset
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(lo, sqlSearch, args...); errDB != nil {
		return 0, ae.DBError("LoginReset ReadAll: unable to select records.", errDB)
	}
	sqlCount := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM login_reset
		%s`, searchStmt)
	var count int
	sqlCount = d.DB.Rebind(sqlCount)
	if errDB := d.DB.Get(&count, sqlCount, args...); errDB != nil {
		return 0, ae.DBError("login_reset ReadAll: unable to select count.", errDB)
	}
	return count, nil
}

func (d *SQLLoginResetV1) Create(ctx context.Context, lo *LoginReset) error {
	sqlPost := `
		INSERT INTO login_reset (
			login_id,
			reset_token,
			created_at,
			updated_at
		) VALUES (
			:login_id,
			:reset_token,
			:created_at,
			:updated_at
		)`
	_, errDB := d.DB.NamedExec(sqlPost, lo)
	if errDB != nil {
		return ae.DBError("LoginReset Post: unable to insert record.", errDB)
	}

	return nil
}

func (d *SQLLoginResetV1) Update(ctx context.Context, lo LoginReset) error {
	sqlPatch := `
		UPDATE login_reset SET
			login_id = :login_id,
			reset_token = :reset_token,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE login_id = :login_id AND reset_token = :reset_token`
	if _, errDB := d.DB.NamedExec(sqlPatch, lo); errDB != nil {
		return ae.DBError("LoginReset Patch: unable to update record.", errDB)
	}
	return nil
}

func (d *SQLLoginResetV1) Delete(ctx context.Context, lo *LoginReset) error {
	sqlDelete := `
		DELETE FROM login_reset WHERE login_id = $1 AND reset_token = $2`
	if _, errDB := d.DB.Exec(sqlDelete, strings.ToLower(lo.LoginId.String), strings.ToLower(lo.ResetToken.String), lo.CreatedAt, lo.UpdatedAt); errDB != nil {
		return ae.DBError("LoginReset Delete: unable to delete record.", errDB)
	}
	return nil
}
