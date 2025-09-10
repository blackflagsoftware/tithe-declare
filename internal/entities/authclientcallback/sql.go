package authclientcallback

import (
	"context"
	"fmt"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	usql "github.com/blackflagsoftware/tithe-declare/internal/util/sql"
	"github.com/jmoiron/sqlx"
)

type (
	SQLAuthClientCallbackV1 struct {
		DB *sqlx.DB
	}
)

func InitSQLV1() *SQLAuthClientCallbackV1 {
	db := stor.InitStorage()
	return &SQLAuthClientCallbackV1{DB: db}
}

func (d *SQLAuthClientCallbackV1) Read(ctx context.Context, au *AuthClientCallback) error {
	sqlGet := `
		SELECT
			client_id,
			callback_url
		FROM auth_client_callback WHERE client_id = $1 AND callback_url = $2`
	if errDB := d.DB.Get(au, sqlGet, au.ClientId, au.CallbackUrl); errDB != nil {
		return ae.DBError("AuthClientCallback Get: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLAuthClientCallbackV1) ReadAll(ctx context.Context, au *[]AuthClientCallback, param AuthClientCallbackParam) (int, error) {
	searchStmt, args := usql.BuildSearchString(param.Param, false) // false => include the where clause, see internal/util/sql.go
	sqlSearch := fmt.Sprintf(`
		SELECT
			client_id,
			callback_url
		FROM auth_client_callback
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(au, sqlSearch, args...); errDB != nil {
		return 0, ae.DBError("AuthClientCallback ReadAll: unable to select records.", errDB)
	}
	sqlCount := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM auth_client_callback
		%s`, searchStmt)
	var count int
	sqlCount = d.DB.Rebind(sqlCount)
	if errDB := d.DB.Get(&count, sqlCount, args...); errDB != nil {
		return 0, ae.DBError("auth_client_callback ReadAll: unable to select count.", errDB)
	}
	return count, nil
}

func (d *SQLAuthClientCallbackV1) Create(ctx context.Context, au *AuthClientCallback) error {
	sqlPost := `
		INSERT INTO auth_client_callback (
			client_id,
			callback_url
		) VALUES (
			:client_id,
			:callback_url
		)`
	_, errDB := d.DB.NamedExec(sqlPost, au)
	if errDB != nil {
		return ae.DBError("AuthClientCallback Post: unable to insert record.", errDB)
	}

	return nil
}

func (d *SQLAuthClientCallbackV1) Update(ctx context.Context, au AuthClientCallback) error {
	sqlPatch := `
		UPDATE auth_client_callback SET
			client_id = :client_id,
			callback_url = :callback_url
		WHERE client_id = :client_id AND callback_url = :callback_url`
	if _, errDB := d.DB.NamedExec(sqlPatch, au); errDB != nil {
		return ae.DBError("AuthClientCallback Patch: unable to update record.", errDB)
	}
	return nil
}

func (d *SQLAuthClientCallbackV1) Delete(ctx context.Context, au *AuthClientCallback) error {
	sqlDelete := `
		DELETE FROM auth_client_callback WHERE client_id = $1 AND callback_url = $2`
	if _, errDB := d.DB.Exec(sqlDelete, au.ClientId, au.CallbackUrl); errDB != nil {
		return ae.DBError("AuthClientCallback Delete: unable to delete record.", errDB)
	}
	return nil
}
