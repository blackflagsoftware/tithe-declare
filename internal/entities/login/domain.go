package login

import (
	"context"
	"net/mail"
	"time"

	"github.com/blackflagsoftware/tithe-declare/config"
	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	a "github.com/blackflagsoftware/tithe-declare/internal/audit"
	mid "github.com/blackflagsoftware/tithe-declare/internal/middleware"
	"github.com/blackflagsoftware/tithe-declare/internal/util"
	"github.com/blackflagsoftware/tithe-declare/internal/util/email"
	"github.com/blackflagsoftware/tithe-declare/internal/util/function"
	"gopkg.in/guregu/null.v3"
)

//go:generate mockgen -source=domain.go -destination=mock.go -package=login
type (
	DataLoginV1Adapter interface {
		Read(context.Context, *Login) error
		ReadAll(context.Context, *[]Login, LoginParam) (int, error)
		Create(context.Context, *Login, ResetRequest) error
		Update(context.Context, Login) error
		UpdatePwd(context.Context, Login) error
		Delete(context.Context, *Login) error
		GetByEmailAddr(context.Context, *Login) error
		GetResetRequest(context.Context, *ResetRequest) error
		ProcessResetRequest(context.Context, *ResetRequest) error
		GetLoginRoles(context.Context, string, *[]string) error
		WithRoles(context.Context, *[]LoginRoles) (int, error)
	}

	DomainLoginV1 struct {
		dataLoginV1 DataLoginV1Adapter
		auditWriter a.AuditAdapter
		emailer     email.Emailer
	}
)

func NewDomainLoginV1(clog DataLoginV1Adapter) *DomainLoginV1 {
	aw := a.AuditInit()
	em := email.EmailInit()
	return &DomainLoginV1{dataLoginV1: clog, auditWriter: aw, emailer: em}
}

func (m *DomainLoginV1) Get(ctx context.Context, login *Login) error {
	if login.Id == "" {
		return ae.MissingParamError("Id")
	}
	return m.dataLoginV1.Read(ctx, login)
}

func (m *DomainLoginV1) Search(ctx context.Context, login *[]Login, param LoginParam) (int, error) {
	param.Param.CalculateParam("email_addr", map[string]string{"email_addr": "email_addr", "pwd": "pwd", "active": "active", "set_pwd": "set_pwd", "created_at": "created_at", "updated_at": "updated_at"})

	return m.dataLoginV1.ReadAll(ctx, login, param)
}

func (m *DomainLoginV1) Post(ctx context.Context, login *Login) error {
	if !login.EmailAddr.Valid {
		return ae.MissingParamError("EmailAddress")
	}
	if login.EmailAddr.Valid && len(login.EmailAddr.ValueOrZero()) > 100 {
		return ae.StringLengthError("EmailAddress", 100)
	}
	if _, err := mail.ParseAddress(login.EmailAddr.String); err != nil {
		return ae.EmailValidError(err.Error())
	}
	// check if email is already used before
	logDup := &Login{EmailAddr: login.EmailAddr}
	if err := m.dataLoginV1.GetByEmailAddr(ctx, logDup); err != nil {
		title := err.(ae.ApiError).BodyError().Title
		if title != "No Results Error" {
			return err
		}
	}
	if logDup.Id != "" {
		return ae.DuplicateEmailError(logDup.EmailAddr.String)
	}
	// set this to empty string
	login.Pwd = null.NewString("", true)
	login.SetPwd = null.BoolFrom(true)
	login.Active = null.BoolFrom(true)
	login.CreatedAt.Scan(time.Now().UTC())
	login.Id = function.GenerateUUID()
	resetRequest := ResetRequest{LoginId: login.Id, ResetToken: function.GenerateUUID(), CreatedAt: time.Now().UTC()}
	if err := m.dataLoginV1.Create(ctx, login, resetRequest); err != nil {
		return err
	}
	go m.emailer.SendReset(ctx, login.EmailAddr.String, resetRequest.ResetToken)
	go a.AuditCreate(m.auditWriter, *login, LoginConst, a.KeysToString("id", login.Id))
	return nil
}

// Patch only allows to update the email or active, see PatchPwd to update the pwd
func (m *DomainLoginV1) Patch(ctx context.Context, logIn Login) error {
	login := &Login{Id: logIn.Id}
	errGet := m.dataLoginV1.Read(ctx, login)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]interface{})
	// EmailAddr
	if logIn.EmailAddr.Valid {
		if logIn.EmailAddr.Valid && len(logIn.EmailAddr.ValueOrZero()) > 100 {
			return ae.StringLengthError("EmailAddress", 100)
		}
		existingValues["email_addr"] = login.EmailAddr.String
		login.EmailAddr = logIn.EmailAddr
	}
	// FirstName
	if logIn.FirstName.Valid {
		if logIn.FirstName.Valid && len(logIn.FirstName.ValueOrZero()) > 50 {
			return ae.StringLengthError("FirstName", 50)
		}
		existingValues["first_name"] = login.FirstName.String
		login.FirstName = logIn.FirstName
	}
	// LastName
	if logIn.LastName.Valid {
		if logIn.LastName.Valid && len(logIn.LastName.ValueOrZero()) > 100 {
			return ae.StringLengthError("LastName", 100)
		}
		existingValues["last_name"] = login.LastName.String
		login.LastName = logIn.LastName
	}
	// Active
	if logIn.Active.Valid {
		existingValues["active"] = login.Active.Bool
		login.Active = logIn.Active
	}

	login.UpdatedAt.Scan(time.Now().UTC())
	if err := m.dataLoginV1.Update(ctx, *login); err != nil {
		return err
	}
	go a.AuditPatch(m.auditWriter, *login, LoginConst, a.KeysToString("Id", login.Id), existingValues)
	return nil
}

// PatchPwd only allows to change the pwd, see Patch to update the email
func (m *DomainLoginV1) PatchPwd(ctx context.Context, logIn Login) error {
	if logIn.Id == "" {
		return ae.MissingParamError("Id")
	}
	login := &Login{Id: logIn.Id}
	errGet := m.dataLoginV1.Read(ctx, login)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]interface{})
	if logIn.Pwd.Valid {
		if login.Pwd.Valid && len(logIn.Pwd.ValueOrZero()) > 72 {
			return ae.StringLengthError("Pwd", 72)
		}
		existingValues["pwd"] = login.Pwd.String
	}
	if err := util.PasswordValidator(logIn.Pwd.String, logIn.ConfirmPwd.String); err != nil {
		return err
	}
	hashPwd, errHash := util.EncryptPassword(logIn.Pwd.String)
	if errHash != nil {
		return errHash
	}
	login.Pwd.Scan(hashPwd)
	if err := m.dataLoginV1.UpdatePwd(ctx, *login); err != nil {
		return err
	}
	go a.AuditPatch(m.auditWriter, *login, LoginConst, a.KeysToString("id", login.Id), existingValues)
	return nil
}

func (m *DomainLoginV1) Delete(ctx context.Context, login *Login) error {
	if login.Id == "" {
		return ae.MissingParamError("Id")
	}
	if err := m.dataLoginV1.Delete(ctx, login); err != nil {
		return err
	}
	go a.AuditDelete(m.auditWriter, *login, LoginConst, a.KeysToString("id", login.Id))
	return nil
}

// ResetPwd takes email, pwd, confirmPwd
func (m *DomainLoginV1) PwdReset(ctx context.Context, pwd PasswordReset) error {
	if pwd.EmailAddr == "" {
		return ae.MissingParamError("EmailAddress")
	}
	login := Login{EmailAddr: null.StringFrom(pwd.EmailAddr)}
	if errLogin := m.dataLoginV1.GetByEmailAddr(ctx, &login); errLogin != nil {
		return errLogin
	}
	// reset login is still valid
	resetRequest := ResetRequest{LoginId: login.Id, ResetToken: pwd.ResetToken}
	if errReset := m.dataLoginV1.GetResetRequest(ctx, &resetRequest); errReset != nil {
		title := errReset.(ae.ApiError).BodyError().Title
		if title == "No Results Error" {
			return ae.ResetTokenInvalidError()
		}
		return errReset
	}
	now := time.Now().UTC()
	expired := resetRequest.CreatedAt.Add(time.Duration(config.A.GetResetDuration()*24) * time.Hour)
	if now.After(expired) {
		return ae.ResetTokenInvalidError()
	}
	// makes sure the login record is in the correct state for pwd reset
	if !login.Active.Bool {
		return ae.LoginActiveError()
	}
	if !login.SetPwd.Bool {
		return ae.ResetTokenInvalidError()
	}
	// validates the pwd, save to storage if all is good
	existingValues := make(map[string]interface{})
	if pwd.Pwd.Valid {
		if pwd.Pwd.Valid && len(pwd.Pwd.ValueOrZero()) > 72 {
			return ae.StringLengthError("Pwd", 72)
		}
		existingValues["pwd"] = login.Pwd.String
	}
	if err := util.PasswordValidator(pwd.Pwd.String, pwd.ConfirmPwd.String); err != nil {
		return err
	}
	hashPwd, errHash := util.EncryptPassword(pwd.Pwd.String)
	if errHash != nil {
		return errHash
	}
	login.Pwd.Scan(hashPwd)
	if err := m.dataLoginV1.UpdatePwd(ctx, login); err != nil {
		return err
	}
	go a.AuditPatch(m.auditWriter, login, LoginConst, a.KeysToString("id", login.Id), existingValues)
	return nil
}

func (m *DomainLoginV1) ProcessResetRequest(ctx context.Context, res *ResetRequest) error {
	if res.EmailAddr == "" {
		return ae.MissingParamError("EmailAddress")
	}
	login := &Login{EmailAddr: null.StringFrom(res.EmailAddr)}
	err := m.dataLoginV1.GetByEmailAddr(ctx, login)
	if err != nil {
		title := err.(ae.ApiError).BodyError().Title
		if title == "No Results Error" {
			// no valid email, send back success
			return nil
		}
		return err // all other errors, at least let the caller know there was an issue
	}
	if login.Active.Valid && !login.Active.Bool {
		// TODO: email address was not active, report here
		return nil
	}
	res.LoginId = login.Id
	res.ResetToken = function.GenerateUUID()
	res.CreatedAt = time.Now().UTC()
	if err := m.dataLoginV1.ProcessResetRequest(ctx, res); err != nil {
		return err
	}
	go m.emailer.SendReset(ctx, login.EmailAddr.String, res.ResetToken)
	return nil
}

func (m *DomainLoginV1) SignIn(ctx context.Context, logIn Login) (string, error) {
	if !logIn.EmailAddr.Valid {
		return "", ae.MissingParamError("EmailAddress")
	}
	login := &Login{EmailAddr: logIn.EmailAddr}
	err := m.dataLoginV1.GetByEmailAddr(ctx, login)
	if err != nil {
		title := err.(ae.ApiError).BodyError().Title
		if title == "No Results Error" {
			return "", ae.EmailPasswordComboError()
		}
		return "", err
	}
	if err := util.CheckPassword(logIn.Pwd.String, login.Pwd.String); err != nil {
		return "", err
	}
	roles := []string{}
	if errRoles := m.dataLoginV1.GetLoginRoles(ctx, login.Id, &roles); errRoles != nil {
		title := errRoles.(ae.ApiError).BodyError().Title
		if title != "No Results Error" {
			return "", errRoles
		}
	}
	// the argument for AuthBuilder is a list of roles for this person
	token, err := mid.AuthBuild(login.Id, roles)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *DomainLoginV1) WithRoles(ctx context.Context, login *[]LoginRoles) (int, error) {
	return m.dataLoginV1.WithRoles(ctx, login)
}
