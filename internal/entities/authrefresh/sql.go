package authrefresh

import (
	"context"
	"fmt"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	usql "github.com/blackflagsoftware/tithe-declare/internal/util/sql"
	"github.com/jmoiron/sqlx"
)

type (
	SQLAuthRefreshV1 struct {
		DB  *sqlx.DB
		Txn *sqlx.Tx
	}
)

func InitSQLV1() *SQLAuthRefreshV1 {
	db := stor.InitStorage()
	return &SQLAuthRefreshV1{DB: db}
}

func (d *SQLAuthRefreshV1) Read(ctx context.Context, ar *AuthRefresh) error {
	sqlGet := `
		SELECT
			client_id,
			token,
			created_at
		FROM auth_refresh WHERE client_id = $1 and token = $2`
	if errDB := d.DB.Get(ar, sqlGet, ar.ClientId, ar.Token); errDB != nil {
		return ae.DBError("AuthRefresh Get: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLAuthRefreshV1) ReadAll(ctx context.Context, ar *[]AuthRefresh, param AuthRefreshParam) (int, error) {
	searchStmt, args := usql.BuildSearchString(param.Param, false)
	sqlSearch := fmt.Sprintf(`
		SELECT
			client_id,
			token,
			created_at
		FROM auth_refresh
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(ar, sqlSearch, args...); errDB != nil {
		return 0, ae.DBError("AuthRefresh Search: unable to select records.", errDB)
	}
	sqlCount := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM auth_refresh
		%s`, searchStmt)
	var count int
	sqlCount = d.DB.Rebind(sqlCount)
	if errDB := d.DB.Get(&count, sqlCount, args...); errDB != nil {
		return 0, ae.DBError("auth_refresh Search: unable to select count.", errDB)
	}
	return count, nil
}

func (d *SQLAuthRefreshV1) Create(ctx context.Context, ar *AuthRefresh) error {
	sqlPost := `
		INSERT INTO auth_refresh (
			client_id,
			token,
			created_at
		) VALUES (
			:client_id,
			:token,
			:created_at
		)`
	_, errDB := d.DB.NamedExec(sqlPost, ar)
	if errDB != nil {
		return ae.DBError("AuthRefresh Post: unable to insert record.", errDB)
	}

	return nil
}

func (d *SQLAuthRefreshV1) CreateTxn(ar *AuthRefresh) error {
	if d.Txn == nil {
		return ae.GeneralError("Unexpected Error", fmt.Errorf("CreateTxn: txn was not set"))
	}
	sqlPost := `
		INSERT INTO auth_refresh (
			login_id,
			token,
			active,
			created_at
		) VALUES (
			:login_id,
			:token,
			:active,
			:created_at
		)`
	_, errDB := d.Txn.NamedExec(sqlPost, ar)
	if errDB != nil {
		return ae.DBError("AuthRefresh Post: unable to insert record.", errDB)
	}

	return nil
}

func (d *SQLAuthRefreshV1) Update(ctx context.Context, ar AuthRefresh) error {
	sqlPatch := `
		UPDATE auth_refresh SET
			created_at = :created_at
		WHERE client_id = :client_idtoken = :token`
	if _, errDB := d.DB.NamedExec(sqlPatch, ar); errDB != nil {
		return ae.DBError("AuthRefresh Patch: unable to update record.", errDB)
	}
	return nil
}

func (d *SQLAuthRefreshV1) DeactiveAllTxn(ar AuthRefresh) error {
	if d.Txn == nil {
		return ae.GeneralError("Unexpected Error", fmt.Errorf("DeactiveAllTxn: txn was not set"))
	}
	sqlPatch := `
		UPDATE auth_refresh SET
			active = :active,
			created_at = :created_at
		WHERE login_id = :login_id`
	if _, errDB := d.Txn.NamedExec(sqlPatch, ar); errDB != nil {
		return ae.DBError("AuthRefresh DeactiveAllTxn: unable to update records.", errDB)
	}
	return nil
}

func (d *SQLAuthRefreshV1) Delete(ctx context.Context, ar *AuthRefresh) error {
	sqlDelete := `
		DELETE FROM auth_refresh WHERE client_id = $1 and token = $2`
	if _, errDB := d.DB.Exec(sqlDelete, ar.ClientId, ar.Token); errDB != nil {
		return ae.DBError("AuthRefresh Delete: unable to delete record.", errDB)
	}
	return nil
}

func (d *SQLAuthRefreshV1) CycleRefreshToken(ctx context.Context, refreshOld, refreshNew AuthRefresh) (err error) {
	txn := d.DB.MustBegin()
	d.Txn = txn
	defer usql.TxnFinish(txn, &err)

	err = d.DeactiveAllTxn(refreshOld)
	if err != nil {
		return err
	}
	return d.CreateTxn(&refreshNew)
}
