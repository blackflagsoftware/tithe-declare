package loginreset

import (
	"context"
	"time"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	a "github.com/blackflagsoftware/tithe-declare/internal/audit"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
)

//go:generate mockgen -source=domain.go -destination=mock.go -package=loginreset
type (
	DataLoginResetV1Adapter interface {
		Read(context.Context, *LoginReset) error
		ReadAll(context.Context, *[]LoginReset, LoginResetParam) (int, error)
		Create(context.Context, *LoginReset) error
		Update(context.Context, LoginReset) error
		Delete(context.Context, *LoginReset) error
	}

	DomainLoginResetV1 struct {
		dataLoginResetV1 DataLoginResetV1Adapter
		auditWriter      a.AuditAdapter
	}
)

func NewDomainLoginResetV1(cloV1 DataLoginResetV1Adapter) *DomainLoginResetV1 {
	aw := a.AuditInit()
	return &DomainLoginResetV1{dataLoginResetV1: cloV1, auditWriter: aw}
}

func (m *DomainLoginResetV1) Get(ctx context.Context, lo *LoginReset) error {

	return m.dataLoginResetV1.Read(ctx, lo)
}

func (m *DomainLoginResetV1) Search(ctx context.Context, lo *[]LoginReset, param LoginResetParam) (int, error) {
	// the second argument (map[string]string) is a list of columns to use for filtering
	// the key matches the json struct tag, the value is the actual table column name (this should change if aliases are used in your query)
	param.Param.CalculateParam("login_id", map[string]string{"login_id": "login_id", "reset_token": "reset_token", "created_at": "created_at", "updated_at": "updated_at"})
	param.Param.PaginationString = stor.FormatPagination(param.Param.Limit, param.Param.Offset)

	return m.dataLoginResetV1.ReadAll(ctx, lo, param)
}

func (m *DomainLoginResetV1) Post(ctx context.Context, lo *LoginReset) error {
	if !lo.CreatedAt.Valid {
		return ae.MissingParamError("CreatedAt")
	}
	if !lo.UpdatedAt.Valid {
		return ae.MissingParamError("UpdatedAt")
	}
	lo.CreatedAt.Scan(time.Now().UTC())
	if err := m.dataLoginResetV1.Create(ctx, lo); err != nil {
		return err
	}
	go a.AuditCreate(m.auditWriter, *lo, LoginResetConst, a.KeysToString("login_id", lo.LoginId, "reset_token", lo.ResetToken))
	return nil
}

func (m *DomainLoginResetV1) Patch(ctx context.Context, loIn LoginReset) error {
	lo := &LoginReset{LoginId: loIn.LoginId, ResetToken: loIn.ResetToken}
	errGet := m.dataLoginResetV1.Read(ctx, lo)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]any)

	lo.UpdatedAt.Scan(time.Now().UTC())
	if err := m.dataLoginResetV1.Update(ctx, *lo); err != nil {
		return err
	}
	go a.AuditPatch(m.auditWriter, *lo, LoginResetConst, a.KeysToString("login_id", lo.LoginId, "reset_token", lo.ResetToken), existingValues)
	return nil
}

func (m *DomainLoginResetV1) Delete(ctx context.Context, lo *LoginReset) error {

	if err := m.dataLoginResetV1.Delete(ctx, lo); err != nil {
		return err
	}
	go a.AuditDelete(m.auditWriter, *lo, LoginResetConst, a.KeysToString("login_id", lo.LoginId, "reset_token", lo.ResetToken))
	return nil
}
