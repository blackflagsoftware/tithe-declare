package authclient

import (
	"context"
	"fmt"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	usql "github.com/blackflagsoftware/tithe-declare/internal/util/sql"
	"github.com/jmoiron/sqlx"
)

type (
	SQLAuthClientV1 struct {
		DB *sqlx.DB
	}
)

func InitSQLV1() *SQLAuthClientV1 {
	db := stor.InitStorage()
	return &SQLAuthClientV1{DB: db}
}

func (d *SQLAuthClientV1) Read(ctx context.Context, ac *AuthClient) error {
	sqlGet := `
		SELECT
			id,
			name,
			description,
			homepage_url,
			callback_url
		FROM auth_client WHERE id = $1`
	if errDB := d.DB.Get(ac, sqlGet, ac.Id, ac.Name, ac.Description, ac.HomepageUrl, ac.CallbackUrl); errDB != nil {
		return ae.DBError("AuthClient Get: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLAuthClientV1) ReadAll(ctx context.Context, ac *[]AuthClient, param AuthClientParam) (int, error) {
	searchStmt, args := usql.BuildSearchString(param.Param, false) // false => include the where clause, see internal/util/sql.go
	sqlSearch := fmt.Sprintf(`
		SELECT
			id,
			name,
			description,
			homepage_url,
			callback_url
		FROM auth_client
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(ac, sqlSearch, args...); errDB != nil {
		return 0, ae.DBError("AuthClient ReadAll: unable to select records.", errDB)
	}
	sqlCount := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM auth_client
		%s`, searchStmt)
	var count int
	sqlCount = d.DB.Rebind(sqlCount)
	if errDB := d.DB.Get(&count, sqlCount, args...); errDB != nil {
		return 0, ae.DBError("auth_client ReadAll: unable to select count.", errDB)
	}
	return count, nil
}

func (d *SQLAuthClientV1) Create(ctx context.Context, ac *AuthClient) error {
	sqlPost := `
		INSERT INTO auth_client (
			id,
			name,
			description,
			homepage_url,
			callback_url
		) VALUES (
			:id,
			:name,
			:description,
			:homepage_url,
			:callback_url
		)`
	_, errDB := d.DB.NamedExec(sqlPost, ac)
	if errDB != nil {
		return ae.DBError("AuthClient Post: unable to insert record.", errDB)
	}

	return nil
}

func (d *SQLAuthClientV1) Update(ctx context.Context, ac AuthClient) error {
	sqlPatch := `
		UPDATE auth_client SET
			id = :id,
			name = :name,
			description = :description,
			homepage_url = :homepage_url,
			callback_url = :callback_url
		WHERE id = :id`
	if _, errDB := d.DB.NamedExec(sqlPatch, ac); errDB != nil {
		return ae.DBError("AuthClient Patch: unable to update record.", errDB)
	}
	return nil
}

func (d *SQLAuthClientV1) Delete(ctx context.Context, ac *AuthClient) error {
	sqlDelete := `
		DELETE FROM auth_client WHERE id = $1`
	if _, errDB := d.DB.Exec(sqlDelete, ac.Id, ac.Name, ac.Description, ac.HomepageUrl, ac.CallbackUrl); errDB != nil {
		return ae.DBError("AuthClient Delete: unable to delete record.", errDB)
	}
	return nil
}
