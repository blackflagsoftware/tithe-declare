package loginrole

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
	SQLLoginRoleV1 struct {
		DB *sqlx.DB
	}
)

func InitSQLV1() *SQLLoginRoleV1 {
	db := stor.InitStorage()
	return &SQLLoginRoleV1{DB: db}
}

func (d *SQLLoginRoleV1) Read(ctx context.Context, lr *LoginRole) error {
	sqlGet := `
		SELECT
			login_id,
			role_id
		FROM login_role WHERE login_id = $1 AND role_id = $2`
	if errDB := d.DB.Get(lr, sqlGet, strings.ToLower(lr.LoginId.String), lr.RoleId); errDB != nil {
		return ae.DBError("LoginRole Get: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLLoginRoleV1) ReadAll(ctx context.Context, lr *[]LoginRole, param LoginRoleParam) (int, error) {
	searchStmt, args := usql.BuildSearchString(param.Param, false) // false => include the where clause, see internal/util/sql.go
	sqlSearch := fmt.Sprintf(`
		SELECT
			login_id,
			role_id
		FROM login_role
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(lr, sqlSearch, args...); errDB != nil {
		return 0, ae.DBError("LoginRole ReadAll: unable to select records.", errDB)
	}
	sqlCount := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM login_role
		%s`, searchStmt)
	var count int
	sqlCount = d.DB.Rebind(sqlCount)
	if errDB := d.DB.Get(&count, sqlCount, args...); errDB != nil {
		return 0, ae.DBError("login_role ReadAll: unable to select count.", errDB)
	}
	return count, nil
}

func (d *SQLLoginRoleV1) Create(ctx context.Context, lr *LoginRole) error {
	sqlPost := `
		INSERT INTO login_role (
			login_id,
			role_id
		) VALUES (
			:login_id,
			:role_id
		)`
	_, errDB := d.DB.NamedExec(sqlPost, lr)
	if errDB != nil {
		return ae.DBError("LoginRole Post: unable to insert record.", errDB)
	}

	return nil
}

func (d *SQLLoginRoleV1) Update(ctx context.Context, lr LoginRole) error {
	sqlPatch := `
		UPDATE login_role SET
			login_id = :login_id,
			role_id = :role_id
		WHERE login_id = :login_id AND role_id = :role_id`
	if _, errDB := d.DB.NamedExec(sqlPatch, lr); errDB != nil {
		return ae.DBError("LoginRole Patch: unable to update record.", errDB)
	}
	return nil
}

func (d *SQLLoginRoleV1) Delete(ctx context.Context, lr *LoginRole) error {
	sqlDelete := `
		DELETE FROM login_role WHERE login_id = $1 AND role_id = $2`
	if _, errDB := d.DB.Exec(sqlDelete, strings.ToLower(lr.LoginId.String), lr.RoleId); errDB != nil {
		return ae.DBError("LoginRole Delete: unable to delete record.", errDB)
	}
	return nil
}
