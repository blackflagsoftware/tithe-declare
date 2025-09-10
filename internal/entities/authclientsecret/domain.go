package authclientsecret

import (
	"context"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	a "github.com/blackflagsoftware/tithe-declare/internal/audit"
)

//go:generate mockgen -source=domain.go -destination=mock.go -package=authclientsecret
type (
	DataAuthClientSecretV1Adapter interface {
		Read(context.Context, *AuthClientSecret) error
		ReadAll(context.Context, *[]AuthClientSecret, AuthClientSecretParam) (int, error)
		Create(context.Context, *AuthClientSecret) error
		Update(context.Context, AuthClientSecret) error
		Delete(context.Context, *AuthClientSecret) error
		ReadByIdAndSecret(context.Context, *AuthClientSecret) error
	}

	DomainAuthClientSecretV1 struct {
		dataAuthClientSecretV1 DataAuthClientSecretV1Adapter
		auditWriter            a.AuditAdapter
	}
)

func NewDomainAuthClientSecretV1(cacsV1 DataAuthClientSecretV1Adapter) *DomainAuthClientSecretV1 {
	aw := a.AuditInit()
	return &DomainAuthClientSecretV1{dataAuthClientSecretV1: cacsV1, auditWriter: aw}
}

func (d *DomainAuthClientSecretV1) Get(ctx context.Context, acs *AuthClientSecret) error {
	return d.dataAuthClientSecretV1.Read(ctx, acs)
}

func (d *DomainAuthClientSecretV1) Search(ctx context.Context, acs *[]AuthClientSecret, param AuthClientSecretParam) (int, error) {
	param.Param.CalculateParam("", map[string]string{})
	return d.dataAuthClientSecretV1.ReadAll(ctx, acs, param)
}

func (d *DomainAuthClientSecretV1) Post(ctx context.Context, acs *AuthClientSecret) error {
	if err := d.dataAuthClientSecretV1.Create(ctx, acs); err != nil {
		return err
	}
	go a.AuditCreate(d.auditWriter, *acs, AuthClientSecretConst, a.KeysToString("client_id", acs.ClientId, "secret", acs.Secret))
	return nil
}

func (d *DomainAuthClientSecretV1) Patch(ctx context.Context, auIn AuthClientSecret) error {
	acs := &AuthClientSecret{ClientId: auIn.ClientId, Secret: auIn.Secret}
	errGet := d.dataAuthClientSecretV1.Read(ctx, acs)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]interface{})

	if err := d.dataAuthClientSecretV1.Update(ctx, *acs); err != nil {
		return err
	}
	go a.AuditPatch(d.auditWriter, *acs, AuthClientSecretConst, a.KeysToString("client_id", acs.ClientId, "secret", acs.Secret), existingValues)
	return nil
}

func (d *DomainAuthClientSecretV1) Delete(ctx context.Context, acs *AuthClientSecret) error {
	if err := d.dataAuthClientSecretV1.Delete(ctx, acs); err != nil {
		return err
	}
	go a.AuditDelete(d.auditWriter, *acs, AuthClientSecretConst, a.KeysToString("client_id", acs.ClientId, "secret", acs.Secret))
	return nil
}

func (d *DomainAuthClientSecretV1) GetByIdAndSecret(ctx context.Context, acs *AuthClientSecret) error {
	if !acs.ClientId.Valid || acs.ClientId.String == "" {
		return ae.MissingParamError("ClientId")
	}
	if !acs.Secret.Valid || acs.Secret.String == "" {
		return ae.MissingParamError("Secret")
	}
	return d.dataAuthClientSecretV1.ReadByIdAndSecret(ctx, acs)
}
