package authauthorize

import (
	"context"
	"fmt"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	usql "github.com/blackflagsoftware/tithe-declare/internal/util/sql"
	"github.com/jmoiron/sqlx"
)

type (
	SQLAuthAuthorizeV1 struct {
		DB *sqlx.DB
	}
)

func InitSQLV1() *SQLAuthAuthorizeV1 {
	db := stor.InitStorage()
	return &SQLAuthAuthorizeV1{DB: db}
}

func (d *SQLAuthAuthorizeV1) Read(ctx context.Context, aa *AuthAuthorize) error {
	sqlGet := `
		SELECT
			id,
			client_id,
			verifier,
			verifier_encode_method,
			state,
			scope,
			authorized_at,
			auth_code_at,
			auth_code
		FROM auth_authorize WHERE id = $1`
	if errDB := d.DB.Get(aa, sqlGet, aa.Id, aa.ClientId, aa.Verifier, aa.VerifierEncodeMethod, aa.State, aa.Scope, aa.AuthorizedAt, aa.AuthCodeAt, aa.AuthCode); errDB != nil {
		return ae.DBError("AuthAuthorize Get: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLAuthAuthorizeV1) ReadAll(ctx context.Context, aa *[]AuthAuthorize, param AuthAuthorizeParam) (int, error) {
	searchStmt, args := usql.BuildSearchString(param.Param, false) // false => include the where clause, see internal/util/sql.go
	sqlSearch := fmt.Sprintf(`
		SELECT
			id,
			client_id,
			verifier,
			verifier_encode_method,
			state,
			scope,
			authorized_at,
			auth_code_at,
			auth_code
		FROM auth_authorize
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(aa, sqlSearch, args...); errDB != nil {
		return 0, ae.DBError("AuthAuthorize ReadAll: unable to select records.", errDB)
	}
	sqlCount := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM auth_authorize
		%s`, searchStmt)
	var count int
	sqlCount = d.DB.Rebind(sqlCount)
	if errDB := d.DB.Get(&count, sqlCount, args...); errDB != nil {
		return 0, ae.DBError("auth_authorize ReadAll: unable to select count.", errDB)
	}
	return count, nil
}

func (d *SQLAuthAuthorizeV1) Create(ctx context.Context, aa *AuthAuthorize) error {
	sqlPost := `
		INSERT INTO auth_authorize (
			id,
			client_id,
			verifier,
			verifier_encode_method,
			state,
			scope,
			authorized_at,
			auth_code_at,
			auth_code
		) VALUES (
			:id,
			:client_id,
			:verifier,
			:verifier_encode_method,
			:state,
			:scope,
			:authorized_at,
			:auth_code_at,
			:auth_code
		)`
	_, errDB := d.DB.NamedExec(sqlPost, aa)
	if errDB != nil {
		return ae.DBError("AuthAuthorize Post: unable to insert record.", errDB)
	}

	return nil
}

func (d *SQLAuthAuthorizeV1) Update(ctx context.Context, aa AuthAuthorize) error {
	sqlPatch := `
		UPDATE auth_authorize SET
			id = :id,
			client_id = :client_id,
			verifier = :verifier,
			verifier_encode_method = :verifier_encode_method,
			state = :state,
			scope = :scope,
			authorized_at = :authorized_at,
			auth_code_at = :auth_code_at,
			auth_code = :auth_code
		WHERE id = :id`
	if _, errDB := d.DB.NamedExec(sqlPatch, aa); errDB != nil {
		return ae.DBError("AuthAuthorize Patch: unable to update record.", errDB)
	}
	return nil
}

func (d *SQLAuthAuthorizeV1) Delete(ctx context.Context, aa *AuthAuthorize) error {
	sqlDelete := `
		DELETE FROM auth_authorize WHERE id = $1`
	if _, errDB := d.DB.Exec(sqlDelete, aa.Id, aa.ClientId, aa.Verifier, aa.VerifierEncodeMethod, aa.State, aa.Scope, aa.AuthorizedAt, aa.AuthCodeAt, aa.AuthCode); errDB != nil {
		return ae.DBError("AuthAuthorize Delete: unable to delete record.", errDB)
	}
	return nil
}
