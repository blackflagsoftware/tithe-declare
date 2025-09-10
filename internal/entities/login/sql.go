package login

import (
	"context"
	"fmt"
	"time"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	usql "github.com/blackflagsoftware/tithe-declare/internal/util/sql"
	"github.com/jmoiron/sqlx"
)

type (
	SQLLoginV1 struct {
		DB *sqlx.DB
	}
)

func InitSQLV1() *SQLLoginV1 {
	db := stor.InitStorage()
	return &SQLLoginV1{DB: db}
}

func (d *SQLLoginV1) Read(ctx context.Context, login *Login) error {
	sqlGet := `
		SELECT
			id,
			email_addr,
			first_name,
			last_name,
			active,
			set_pwd,
			created_at,
			updated_at
		FROM login WHERE id = $1`
	if errDB := d.DB.Get(login, sqlGet, login.Id); errDB != nil {
		return ae.DBError("Login Get: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLLoginV1) ReadAll(ctx context.Context, login *[]Login, param LoginParam) (int, error) {
	searchStmt, args := usql.BuildSearchString(param.Param, false)
	sqlSearch := fmt.Sprintf(`
		SELECT
			id,
			email_addr,
			first_name,
			last_name,
			active,
			set_pwd,
			created_at,
			updated_at
		FROM login
		%s
		ORDER BY %s %s`, searchStmt, param.Sort, param.PaginationString)
	sqlSearch = d.DB.Rebind(sqlSearch)
	if errDB := d.DB.Select(login, sqlSearch, args...); errDB != nil {
		return 0, ae.DBError("Login Search: unable to select records.", errDB)
	}
	sqlCount := fmt.Sprintf(`
		SELECT
			COUNT(*)
		FROM login
		%s`, searchStmt)
	var count int
	sqlCount = d.DB.Rebind(sqlCount)
	if errDB := d.DB.Get(&count, sqlCount, args...); errDB != nil {
		return 0, ae.DBError("login Search: unable to select count.", errDB)
	}
	return count, nil
}

func (d *SQLLoginV1) Create(ctx context.Context, login *Login, res ResetRequest) (err error) {
	txn := d.DB.MustBegin()
	defer usql.TxnFinish(txn, &err)

	sqlPost := `
		INSERT INTO login (
			id,
			email_addr,
			first_name,
			last_name,
			pwd,
			active,
			set_pwd,
			created_at,
			updated_at
		) VALUES (
			:id,
			:email_addr,
			:first_name,
			:last_name,
			:pwd,
			:active,
			:set_pwd,
			:created_at,
			:updated_at
		)`
	_, errDB := txn.NamedExec(sqlPost, login)
	if errDB != nil {
		err = ae.DBError("Login Post: unable to insert record.", errDB)
		return err
	}

	sqlResetInsert := `INSERT INTO login_reset (login_id, reset_token, created_at) VALUES ($1, $2, $3)`
	if _, errDB := txn.Exec(sqlResetInsert, res.LoginId, res.ResetToken, res.CreatedAt); errDB != nil {
		err = ae.DBError("Create: reset request, unable to insert.", errDB)
		return
	}
	return err
}

func (d *SQLLoginV1) Update(ctx context.Context, login Login) error {
	sqlPatch := `
		UPDATE login SET
			email_addr = :email_addr,
			first_name = :first_name,
			last_name = :last_name,
			active = :active,
			updated_at = :updated_at
		WHERE id = :id`
	if _, errDB := d.DB.NamedExec(sqlPatch, login); errDB != nil {
		return ae.DBError("Login Patch: unable to update record.", errDB)
	}
	return nil
}

func (d *SQLLoginV1) UpdatePwd(ctx context.Context, login Login) (err error) {
	txn := d.DB.MustBegin()
	defer usql.TxnFinish(txn, &err)

	sqlPatch := `
		UPDATE login SET
			pwd = :pwd,
			updated_at = :updated_at,
			set_pwd = false
		WHERE id = :id`
	if _, errDB := txn.NamedExec(sqlPatch, login); errDB != nil {
		err = ae.DBError("Login UpdatePwd: unable to update record.", errDB)
		return
	}
	// update all reset records (tokens) if any
	now := time.Now().UTC()
	sqlReset := `UPDATE login_reset SET updated_at = $1 WHERE login_id = $2`
	if _, errDB := txn.Exec(sqlReset, now, login.Id); errDB != nil {
		err = ae.DBError("Login UpdatePwd: unable to set finish reset.", errDB)
		return
	}
	return nil
}

func (d *SQLLoginV1) Delete(ctx context.Context, login *Login) (err error) {
	txn := d.DB.MustBegin()
	defer usql.TxnFinish(txn, &err)

	sqlDelete := `
		DELETE FROM login WHERE id = $1`
	if _, errDB := txn.Exec(sqlDelete, login.Id); errDB != nil {
		return ae.DBError("Login Delete: unable to delete record.", errDB)
	}
	sqlRoleDelete := `DELETE FROM login_role WHERE login_id = $1`
	if _, errDB := txn.Exec(sqlRoleDelete, login.Id); errDB != nil {
		return ae.DBError("Login Delete: unable to delete login_role record.", errDB)
	}
	return nil
}

func (d *SQLLoginV1) GetByEmailAddr(ctx context.Context, login *Login) error {
	sqlGet := `SELECT id, pwd, active, set_pwd FROM login WHERE email_addr = $1`
	if errDB := d.DB.Get(login, sqlGet, login.EmailAddr); errDB != nil {
		return ae.DBError("GetByEmailAddr: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLLoginV1) GetResetRequest(ctx context.Context, resetRequest *ResetRequest) error {
	sqlGet := `SELECT login_id, reset_token, created_at FROM login_reset WHERE login_id = $1 AND reset_token = $2 AND updated_at IS NULL LIMIT 1`
	if errDB := d.DB.Get(resetRequest, sqlGet, resetRequest.LoginId, resetRequest.ResetToken); errDB != nil {
		return ae.DBError("GetByEmailAddr: unable to get record.", errDB)
	}
	return nil
}

func (d *SQLLoginV1) ProcessResetRequest(ctx context.Context, res *ResetRequest) (err error) {
	txn := d.DB.MustBegin()
	defer usql.TxnFinish(txn, &err)

	now := time.Now().UTC()
	sqlResetUpdate := `UPDATE login_reset SET updated_at = $1 WHERE login_id = $2`
	if _, errDB := txn.Exec(sqlResetUpdate, now, res.LoginId); errDB != nil {
		err = ae.DBError("Login Reset: unable to update.", errDB)
		return
	}
	sqlResetInsert := `INSERT INTO login_reset (login_Id, reset_token, created_at) VALUES ($1, $2, $3)`
	if _, errDB := txn.Exec(sqlResetInsert, res.LoginId, res.ResetToken, now); errDB != nil {
		err = ae.DBError("Login Reset: unable to insert.", errDB)
		return
	}
	// this assumes, from the check before, that the user was active
	sqlLoginUpdate := `UPDATE login SET set_pwd = true WHERE id = $1`
	if _, errDB := txn.Exec(sqlLoginUpdate, res.LoginId); errDB != nil {
		err = ae.DBError("Login Reset: unable to update login.", errDB)
		return
	}
	return err
}

func (d *SQLLoginV1) GetLoginRoles(ctx context.Context, loginId string, roles *[]string) error {
	sqlGet := `
		SELECT
			r.name
		FROM role AS r
		INNER JOIN login_role AS lr ON r.id = lr.role_id
		WHERE lr.login_id = $1`
	if errDB := d.DB.Select(roles, sqlGet, loginId); errDB != nil {
		return ae.DBError("GetLoginRoles: unable to get roles", errDB)
	}
	return nil
}

func (d *SQLLoginV1) WithRoles(ctx context.Context, login *[]LoginRoles) (int, error) {
	sqlLogin := `
		SELECT
			id,
			email_addr
		FROM login
		WHERE active = true
		ORDER BY email_addr`
	if errDB := d.DB.Select(login, sqlLogin); errDB != nil {
		return 0, ae.DBError("Login WithRoles: unable to select records.", errDB)
	}

	sqlRoles := `
		SELECT
			r.name
		FROM login_role AS lr
		INNER JOIN role AS r ON lr.role_id = r.id
		WHERE lr.login_id = $1`
	for i := range *login {
		roles := []string{}
		if errDB := d.DB.Select(&roles, sqlRoles, (*login)[i].LoginId); errDB != nil {
			return 0, ae.DBError("Login WithRoles: unable to select roles records.", errDB)
		}
		(*login)[i].Roles = roles
	}
	count := len(*login)
	return count, nil
}
