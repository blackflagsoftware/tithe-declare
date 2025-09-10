package role

import (
	"context"
	"fmt"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	usql "github.com/blackflagsoftware/tithe-declare/internal/util/sql"
	"github.com/jmoiron/sqlx"
)

type (
	SQLRoleV1 struct {
		DB *sqlx.DB
	}
)

func InitSQLV1() *SQLRoleV1 {
	db := stor.InitStorage()
	return &SQLRoleV1{DB: db}
}

func (d *SQLRoleV1) Read(ctx context.Context, rol *Role) error {
	sqlGet := `
		SELECT
			id,
			name,
			description
		FROM role WHERE id = $1`
	if errDB := d.DB.Get(rol, sqlGet, rol.Id, rol.Name, rol.Description); errDB != nil {
		return ae.DBError("Role Get: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLRoleV1) ReadAll(ctx context.Context, rol *[]Role, param RoleParam) (int, error) {
	searchStmt, args := usql.BuildSearchString(param.Param, false) // false => include the where clause, see internal/util/sql.go
	sqlSearch := fmt.Sprintf(`
		SELECT
			id,
			name,
			description
		FROM role
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(rol, sqlSearch, args...); errDB != nil {
		return 0, ae.DBError("Role ReadAll: unable to select records.", errDB)
	}
	sqlCount := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM role
		%s`, searchStmt)
	var count int
	sqlCount = d.DB.Rebind(sqlCount)
	if errDB := d.DB.Get(&count, sqlCount, args...); errDB != nil {
		return 0, ae.DBError("role ReadAll: unable to select count.", errDB)
	}
	return count, nil
}

func (d *SQLRoleV1) Create(ctx context.Context, rol *Role) error {
	sqlPost := `
		INSERT INTO role (
			id,
			name,
			description
		) VALUES (
			:id,
			:name,
			:description
		)`
	_, errDB := d.DB.NamedExec(sqlPost, rol)
	if errDB != nil {
		return ae.DBError("Role Post: unable to insert record.", errDB)
	}

	return nil
}

func (d *SQLRoleV1) Update(ctx context.Context, rol Role) error {
	sqlPatch := `
		UPDATE role SET
			id = :id,
			name = :name,
			description = :description
		WHERE id = :id`
	if _, errDB := d.DB.NamedExec(sqlPatch, rol); errDB != nil {
		return ae.DBError("Role Patch: unable to update record.", errDB)
	}
	return nil
}

func (d *SQLRoleV1) Delete(ctx context.Context, rol *Role) error {
	sqlDelete := `
		DELETE FROM role WHERE id = $1`
	if _, errDB := d.DB.Exec(sqlDelete, rol.Id, rol.Name, rol.Description); errDB != nil {
		return ae.DBError("Role Delete: unable to delete record.", errDB)
	}
	return nil
}
