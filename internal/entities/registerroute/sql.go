package registerroute

import (
	"context"
	"fmt"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	usql "github.com/blackflagsoftware/tithe-declare/internal/util/sql"
	"github.com/jmoiron/sqlx"
)

type (
	SQLRegisterRouteV1 struct {
		DB *sqlx.DB
	}
)

func InitSQLV1() *SQLRegisterRouteV1 {
	db := stor.InitStorage()
	return &SQLRegisterRouteV1{DB: db}
}

func (d *SQLRegisterRouteV1) Read(ctx context.Context, reg *RegisterRoute) error {
	sqlGet := `
		SELECT
			raw_path,
			transformed_path,
			roles
		FROM register_route WHERE raw_path = $1`
	if errDB := d.DB.Get(reg, sqlGet, reg.RawPath, reg.TransformedPath, reg.Roles); errDB != nil {
		return ae.DBError("RegisterRoute Get: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLRegisterRouteV1) ReadAll(ctx context.Context, reg *[]RegisterRoute, param RegisterRouteParam) (int, error) {
	searchStmt, args := usql.BuildSearchString(param.Param, false) // false => include the where clause, see internal/util/sql.go
	sqlSearch := fmt.Sprintf(`
		SELECT
			raw_path,
			transformed_path,
			roles
		FROM register_route
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(reg, sqlSearch, args...); errDB != nil {
		return 0, ae.DBError("RegisterRoute ReadAll: unable to select records.", errDB)
	}
	sqlCount := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM register_route
		%s`, searchStmt)
	var count int
	sqlCount = d.DB.Rebind(sqlCount)
	if errDB := d.DB.Get(&count, sqlCount, args...); errDB != nil {
		return 0, ae.DBError("register_route ReadAll: unable to select count.", errDB)
	}
	return count, nil
}

func (d *SQLRegisterRouteV1) Create(ctx context.Context, reg *RegisterRoute) error {
	sqlPost := `
		INSERT INTO register_route (
			raw_path,
			transformed_path,
			roles
		) VALUES (
			:raw_path,
			:transformed_path,
			:roles
		)`
	_, errDB := d.DB.NamedExec(sqlPost, reg)
	if errDB != nil {
		return ae.DBError("RegisterRoute Post: unable to insert record.", errDB)
	}

	return nil
}

func (d *SQLRegisterRouteV1) Update(ctx context.Context, reg RegisterRoute) error {
	sqlPatch := `
		UPDATE register_route SET
			raw_path = :raw_path,
			transformed_path = :transformed_path,
			roles = :roles
		WHERE raw_path = :raw_path`
	if _, errDB := d.DB.NamedExec(sqlPatch, reg); errDB != nil {
		return ae.DBError("RegisterRoute Patch: unable to update record.", errDB)
	}
	return nil
}

func (d *SQLRegisterRouteV1) Delete(ctx context.Context, reg *RegisterRoute) error {
	sqlDelete := `
		DELETE FROM register_route WHERE raw_path = $1`
	if _, errDB := d.DB.Exec(sqlDelete, reg.RawPath, reg.TransformedPath, reg.Roles); errDB != nil {
		return ae.DBError("RegisterRoute Delete: unable to delete record.", errDB)
	}
	return nil
}
