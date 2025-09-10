package authclientsecret

import (
	"context"
	"fmt"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	usql "github.com/blackflagsoftware/tithe-declare/internal/util/sql"
	"github.com/jmoiron/sqlx"
)

type (
	SQLAuthClientSecretV1 struct {
		DB *sqlx.DB
	}
)

func InitSQLV1() *SQLAuthClientSecretV1 {
	db := stor.InitStorage()
	return &SQLAuthClientSecretV1{DB: db}
}

func (d *SQLAuthClientSecretV1) Read(ctx context.Context, au *AuthClientSecret) error {
	sqlGet := `
		SELECT
			client_id,
			secret
		FROM auth_client_secret WHERE client_id = $1 and secret = $2`
	if errDB := d.DB.Get(au, sqlGet, au.ClientId, au.Secret); errDB != nil {
		return ae.DBError("AuthClientSecret Get: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLAuthClientSecretV1) ReadAll(ctx context.Context, au *[]AuthClientSecret, param AuthClientSecretParam) (int, error) {
	searchStmt, args := usql.BuildSearchString(param.Param, false)
	sqlSearch := fmt.Sprintf(`
		SELECT
			client_id,
			secret
		FROM auth_client_secret
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(au, sqlSearch, args...); errDB != nil {
		return 0, ae.DBError("AuthClientSecret Search: unable to select records.", errDB)
	}
	sqlCount := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM auth_client_secret
		%s`, searchStmt)
	var count int
	sqlCount = d.DB.Rebind(sqlCount)
	if errDB := d.DB.Get(&count, sqlCount, args...); errDB != nil {
		return 0, ae.DBError("auth_client_secret Search: unable to select count.", errDB)
	}
	return count, nil
}

func (d *SQLAuthClientSecretV1) Create(ctx context.Context, au *AuthClientSecret) error {
	sqlPost := `
		INSERT INTO auth_client_secret (
			client_id,
			secret
		) VALUES (
			:client_id,
			:secret
		)`
	_, errDB := d.DB.NamedExec(sqlPost, au)
	if errDB != nil {
		return ae.DBError("AuthClientSecret Post: unable to insert record.", errDB)
	}

	return nil
}

func (d *SQLAuthClientSecretV1) Update(ctx context.Context, au AuthClientSecret) error {
	sqlPatch := `
		UPDATE auth_client_secret SET
			
		WHERE client_id = :client_idsecret = :secret`
	if _, errDB := d.DB.NamedExec(sqlPatch, au); errDB != nil {
		return ae.DBError("AuthClientSecret Patch: unable to update record.", errDB)
	}
	return nil
}

func (d *SQLAuthClientSecretV1) Delete(ctx context.Context, au *AuthClientSecret) error {
	sqlDelete := `
		DELETE FROM auth_client_secret WHERE client_id = $1 and secret = $2`
	if _, errDB := d.DB.Exec(sqlDelete, au.ClientId, au.Secret); errDB != nil {
		return ae.DBError("AuthClientSecret Delete: unable to delete record.", errDB)
	}
	return nil
}

func (d *SQLAuthClientSecretV1) ReadByIdAndSecret(ctx context.Context, au *AuthClientSecret) error {
	sqlGet := `
		SELECT
			id,
			client_id,
			secret,
			active
		FROM auth_client_secret WHERE client_id = $1 AND secret = $2`
	if errDB := d.DB.Get(au, sqlGet, au.ClientId, au.Secret); errDB != nil {
		return ae.DBError("AuthClientSecret Get: unable to get record.", errDB)
	}
	return nil
}
