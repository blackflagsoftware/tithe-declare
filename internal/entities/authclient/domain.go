package authclient

import (
	"context"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	a "github.com/blackflagsoftware/tithe-declare/internal/audit"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	"github.com/blackflagsoftware/tithe-declare/internal/util/function"
)

//go:generate mockgen -source=domain.go -destination=mock.go -package=authclient
type (
	DataAuthClientV1Adapter interface {
		Read(context.Context, *AuthClient) error
		ReadAll(context.Context, *[]AuthClient, AuthClientParam) (int, error)
		Create(context.Context, *AuthClient) error
		Update(context.Context, AuthClient) error
		Delete(context.Context, *AuthClient) error
	}

	DomainAuthClientV1 struct {
		dataAuthClientV1 DataAuthClientV1Adapter
		auditWriter      a.AuditAdapter
	}
)

func NewDomainAuthClientV1(cacV1 DataAuthClientV1Adapter) *DomainAuthClientV1 {
	aw := a.AuditInit()
	return &DomainAuthClientV1{dataAuthClientV1: cacV1, auditWriter: aw}
}

func (m *DomainAuthClientV1) Get(ctx context.Context, ac *AuthClient) error {
	if ac.Id == "" {
		return ae.MissingParamError("Id")
	}
	return m.dataAuthClientV1.Read(ctx, ac)
}

func (m *DomainAuthClientV1) Search(ctx context.Context, ac *[]AuthClient, param AuthClientParam) (int, error) {
	// the second argument (map[string]string) is a list of columns to use for filtering
	// the key matches the json struct tag, the value is the actual table column name (this should change if aliases are used in your query)
	param.Param.CalculateParam("name", map[string]string{"id": "id", "name": "name", "description": "description", "homepage_url": "homepage_url", "callback_url": "callback_url"})
	param.Param.PaginationString = stor.FormatPagination(param.Param.Limit, param.Param.Offset)

	return m.dataAuthClientV1.ReadAll(ctx, ac, param)
}

func (m *DomainAuthClientV1) Post(ctx context.Context, ac *AuthClient) error {
	if !ac.Name.Valid {
		return ae.MissingParamError("Name")
	}
	if ac.Name.Valid && len(ac.Name.ValueOrZero()) > 100 {
		return ae.StringLengthError("Name", 100)
	}
	if ac.Description.Valid && len(ac.Description.ValueOrZero()) > 1000 {
		return ae.StringLengthError("Description", 1000)
	}
	if !ac.HomepageUrl.Valid {
		return ae.MissingParamError("HomepageUrl")
	}
	if ac.HomepageUrl.Valid && len(ac.HomepageUrl.ValueOrZero()) > 500 {
		return ae.StringLengthError("HomepageUrl", 500)
	}
	if !ac.CallbackUrl.Valid {
		return ae.MissingParamError("CallbackUrl")
	}
	if ac.CallbackUrl.Valid && len(ac.CallbackUrl.ValueOrZero()) > 500 {
		return ae.StringLengthError("CallbackUrl", 500)
	}
	ac.Id = function.GenerateRandomString(32)
	if err := m.dataAuthClientV1.Create(ctx, ac); err != nil {
		return err
	}
	go a.AuditCreate(m.auditWriter, *ac, AuthClientConst, a.KeysToString("id", ac.Id))
	return nil
}

func (m *DomainAuthClientV1) Patch(ctx context.Context, acIn AuthClient) error {
	ac := &AuthClient{Id: acIn.Id}
	errGet := m.dataAuthClientV1.Read(ctx, ac)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]any)
	// Name
	if acIn.Name.Valid {
		if acIn.Name.Valid && len(acIn.Name.ValueOrZero()) > 100 {
			return ae.StringLengthError("Name", 100)
		}
		existingValues["name"] = ac.Name.String
		ac.Name = acIn.Name
	}
	// Description
	if acIn.Description.Valid {
		if acIn.Description.Valid && len(acIn.Description.ValueOrZero()) > 1000 {
			return ae.StringLengthError("Description", 1000)
		}
		existingValues["description"] = ac.Description.String
		ac.Description = acIn.Description
	}
	// HomepageUrl
	if acIn.HomepageUrl.Valid {
		if acIn.HomepageUrl.Valid && len(acIn.HomepageUrl.ValueOrZero()) > 500 {
			return ae.StringLengthError("HomepageUrl", 500)
		}
		existingValues["homepage_url"] = ac.HomepageUrl.String
		ac.HomepageUrl = acIn.HomepageUrl
	}
	// CallbackUrl
	if acIn.CallbackUrl.Valid {
		if acIn.CallbackUrl.Valid && len(acIn.CallbackUrl.ValueOrZero()) > 500 {
			return ae.StringLengthError("CallbackUrl", 500)
		}
		existingValues["callback_url"] = ac.CallbackUrl.String
		ac.CallbackUrl = acIn.CallbackUrl
	}
	if err := m.dataAuthClientV1.Update(ctx, *ac); err != nil {
		return err
	}
	go a.AuditPatch(m.auditWriter, *ac, AuthClientConst, a.KeysToString("id", ac.Id), existingValues)
	return nil
}

func (m *DomainAuthClientV1) Delete(ctx context.Context, ac *AuthClient) error {
	if ac.Id == "" {
		return ae.MissingParamError("Id")
	}
	if err := m.dataAuthClientV1.Delete(ctx, ac); err != nil {
		return err
	}
	go a.AuditDelete(m.auditWriter, *ac, AuthClientConst, a.KeysToString("id", ac.Id))
	return nil
}
