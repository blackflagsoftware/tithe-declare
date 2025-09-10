package authauthorize

import (
	"context"
	"time"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	a "github.com/blackflagsoftware/tithe-declare/internal/audit"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	"github.com/blackflagsoftware/tithe-declare/internal/util/function"
)

//go:generate mockgen -source=domain.go -destination=mock.go -package=authauthorize
type (
	DataAuthAuthorizeV1Adapter interface {
		Read(context.Context, *AuthAuthorize) error
		ReadAll(context.Context, *[]AuthAuthorize, AuthAuthorizeParam) (int, error)
		Create(context.Context, *AuthAuthorize) error
		Update(context.Context, AuthAuthorize) error
		Delete(context.Context, *AuthAuthorize) error
	}

	DomainAuthAuthorizeV1 struct {
		dataAuthAuthorizeV1 DataAuthAuthorizeV1Adapter
		auditWriter         a.AuditAdapter
	}
)

func NewDomainAuthAuthorizeV1(caaV1 DataAuthAuthorizeV1Adapter) *DomainAuthAuthorizeV1 {
	aw := a.AuditInit()
	return &DomainAuthAuthorizeV1{dataAuthAuthorizeV1: caaV1, auditWriter: aw}
}

func (m *DomainAuthAuthorizeV1) Get(ctx context.Context, aa *AuthAuthorize) error {
	if aa.Id == "" {
		return ae.MissingParamError("Id")
	}
	return m.dataAuthAuthorizeV1.Read(ctx, aa)
}

func (m *DomainAuthAuthorizeV1) Search(ctx context.Context, aa *[]AuthAuthorize, param AuthAuthorizeParam) (int, error) {
	// the second argument (map[string]string) is a list of columns to use for filtering
	// the key matches the json struct tag, the value is the actual table column name (this should change if aliases are used in your query)
	param.Param.CalculateParam("client_id", map[string]string{"id": "id", "client_id": "client_id", "verifier": "verifier", "verifier_encode_method": "verifier_encode_method", "state": "state", "scope": "scope", "authorized_at": "authorized_at", "auth_code_at": "auth_code_at", "auth_code": "auth_code"})
	param.Param.PaginationString = stor.FormatPagination(param.Param.Limit, param.Param.Offset)

	return m.dataAuthAuthorizeV1.ReadAll(ctx, aa, param)
}

func (m *DomainAuthAuthorizeV1) Post(ctx context.Context, aa *AuthAuthorize) error {
	if !aa.ClientId.Valid {
		return ae.MissingParamError("ClientId")
	}
	if aa.ClientId.Valid && len(aa.ClientId.ValueOrZero()) > 32 {
		return ae.StringLengthError("ClientId", 32)
	}
	if aa.VerifierEncodeMethod.Valid && len(aa.VerifierEncodeMethod.ValueOrZero()) > 10 {
		return ae.StringLengthError("VerifierEncodeMethod", 10)
	}
	if aa.State.Valid && len(aa.State.ValueOrZero()) > 100 {
		return ae.StringLengthError("State", 100)
	}
	if aa.Scope.Valid && len(aa.Scope.ValueOrZero()) > 256 {
		return ae.StringLengthError("Scope", 256)
	}
	if !aa.AuthorizedAt.Valid {
		return ae.MissingParamError("AuthorizedAt")
	}
	if aa.AuthCode.Valid && len(aa.AuthCode.ValueOrZero()) > 256 {
		return ae.StringLengthError("AuthCode", 256)
	}
	aa.Id = function.GenerateRandomString(32)
	if err := m.dataAuthAuthorizeV1.Create(ctx, aa); err != nil {
		return err
	}
	go a.AuditCreate(m.auditWriter, *aa, AuthAuthorizeConst, a.KeysToString("id", aa.Id))
	return nil
}

func (m *DomainAuthAuthorizeV1) Patch(ctx context.Context, aaIn AuthAuthorize) error {
	aa := &AuthAuthorize{Id: aaIn.Id}
	errGet := m.dataAuthAuthorizeV1.Read(ctx, aa)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]any)
	// ClientId
	if aaIn.ClientId.Valid {
		if aaIn.ClientId.Valid && len(aaIn.ClientId.ValueOrZero()) > 32 {
			return ae.StringLengthError("ClientId", 32)
		}
		existingValues["client_id"] = aa.ClientId.String
		aa.ClientId = aaIn.ClientId
	}
	// Verifier
	if aaIn.Verifier.Valid {
		existingValues["verifier"] = aa.Verifier.String
		aa.Verifier = aaIn.Verifier
	}
	// VerifierEncodeMethod
	if aaIn.VerifierEncodeMethod.Valid {
		if aaIn.VerifierEncodeMethod.Valid && len(aaIn.VerifierEncodeMethod.ValueOrZero()) > 10 {
			return ae.StringLengthError("VerifierEncodeMethod", 10)
		}
		existingValues["verifier_encode_method"] = aa.VerifierEncodeMethod.String
		aa.VerifierEncodeMethod = aaIn.VerifierEncodeMethod
	}
	// State
	if aaIn.State.Valid {
		if aaIn.State.Valid && len(aaIn.State.ValueOrZero()) > 100 {
			return ae.StringLengthError("State", 100)
		}
		existingValues["state"] = aa.State.String
		aa.State = aaIn.State
	}
	// Scope
	if aaIn.Scope.Valid {
		if aaIn.Scope.Valid && len(aaIn.Scope.ValueOrZero()) > 256 {
			return ae.StringLengthError("Scope", 256)
		}
		existingValues["scope"] = aa.Scope.String
		aa.Scope = aaIn.Scope
	}
	// AuthorizedAt
	if aaIn.AuthorizedAt.Valid {
		existingValues["authorized_at"] = aa.AuthorizedAt.Time.Format(time.RFC3339)
		aa.AuthorizedAt = aaIn.AuthorizedAt
	}
	// AuthCodeAt
	if aaIn.AuthCodeAt.Valid {
		existingValues["auth_code_at"] = aa.AuthCodeAt.Time.Format(time.RFC3339)
		aa.AuthCodeAt = aaIn.AuthCodeAt
	}
	// AuthCode
	if aaIn.AuthCode.Valid {
		if aaIn.AuthCode.Valid && len(aaIn.AuthCode.ValueOrZero()) > 256 {
			return ae.StringLengthError("AuthCode", 256)
		}
		existingValues["auth_code"] = aa.AuthCode.String
		aa.AuthCode = aaIn.AuthCode
	}
	if err := m.dataAuthAuthorizeV1.Update(ctx, *aa); err != nil {
		return err
	}
	go a.AuditPatch(m.auditWriter, *aa, AuthAuthorizeConst, a.KeysToString("id", aa.Id), existingValues)
	return nil
}

func (m *DomainAuthAuthorizeV1) Delete(ctx context.Context, aa *AuthAuthorize) error {
	if aa.Id == "" {
		return ae.MissingParamError("Id")
	}
	if err := m.dataAuthAuthorizeV1.Delete(ctx, aa); err != nil {
		return err
	}
	go a.AuditDelete(m.auditWriter, *aa, AuthAuthorizeConst, a.KeysToString("id", aa.Id))
	return nil
}
