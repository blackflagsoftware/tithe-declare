package authclientcallback

import (
	"context"

	a "github.com/blackflagsoftware/tithe-declare/internal/audit"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
)

//go:generate mockgen -source=domain.go -destination=mock.go -package=authclientcallback
type (
	DataAuthClientCallbackV1Adapter interface {
		Read(context.Context, *AuthClientCallback) error
		ReadAll(context.Context, *[]AuthClientCallback, AuthClientCallbackParam) (int, error)
		Create(context.Context, *AuthClientCallback) error
		Update(context.Context, AuthClientCallback) error
		Delete(context.Context, *AuthClientCallback) error
	}

	DomainAuthClientCallbackV1 struct {
		dataAuthClientCallbackV1 DataAuthClientCallbackV1Adapter
		auditWriter              a.AuditAdapter
	}
)

func NewDomainAuthClientCallbackV1(cauV1 DataAuthClientCallbackV1Adapter) *DomainAuthClientCallbackV1 {
	aw := a.AuditInit()
	return &DomainAuthClientCallbackV1{dataAuthClientCallbackV1: cauV1, auditWriter: aw}
}

func (m *DomainAuthClientCallbackV1) Get(ctx context.Context, au *AuthClientCallback) error {

	return m.dataAuthClientCallbackV1.Read(ctx, au)
}

func (m *DomainAuthClientCallbackV1) Search(ctx context.Context, au *[]AuthClientCallback, param AuthClientCallbackParam) (int, error) {
	// the second argument (map[string]string) is a list of columns to use for filtering
	// the key matches the json struct tag, the value is the actual table column name (this should change if aliases are used in your query)
	param.Param.CalculateParam("client_id", map[string]string{"client_id": "client_id", "callback_url": "callback_url"})
	param.Param.PaginationString = stor.FormatPagination(param.Param.Limit, param.Param.Offset)

	return m.dataAuthClientCallbackV1.ReadAll(ctx, au, param)
}

func (m *DomainAuthClientCallbackV1) Post(ctx context.Context, au *AuthClientCallback) error {

	if err := m.dataAuthClientCallbackV1.Create(ctx, au); err != nil {
		return err
	}
	go a.AuditCreate(m.auditWriter, *au, AuthClientCallbackConst, a.KeysToString("client_id", au.ClientId, "callback_url", au.CallbackUrl))
	return nil
}

func (m *DomainAuthClientCallbackV1) Patch(ctx context.Context, auIn AuthClientCallback) error {
	au := &AuthClientCallback{ClientId: auIn.ClientId, CallbackUrl: auIn.CallbackUrl}
	errGet := m.dataAuthClientCallbackV1.Read(ctx, au)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]any)

	if err := m.dataAuthClientCallbackV1.Update(ctx, *au); err != nil {
		return err
	}
	go a.AuditPatch(m.auditWriter, *au, AuthClientCallbackConst, a.KeysToString("client_id", au.ClientId, "callback_url", au.CallbackUrl), existingValues)
	return nil
}

func (m *DomainAuthClientCallbackV1) Delete(ctx context.Context, au *AuthClientCallback) error {

	if err := m.dataAuthClientCallbackV1.Delete(ctx, au); err != nil {
		return err
	}
	go a.AuditDelete(m.auditWriter, *au, AuthClientCallbackConst, a.KeysToString("client_id", au.ClientId, "callback_url", au.CallbackUrl))
	return nil
}
