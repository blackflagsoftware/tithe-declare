package role

import (
	"context"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	a "github.com/blackflagsoftware/tithe-declare/internal/audit"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	"github.com/blackflagsoftware/tithe-declare/internal/util/function"
)

//go:generate mockgen -source=domain.go -destination=mock.go -package=role
type (
	DataRoleV1Adapter interface {
		Read(context.Context, *Role) error
		ReadAll(context.Context, *[]Role, RoleParam) (int, error)
		Create(context.Context, *Role) error
		Update(context.Context, Role) error
		Delete(context.Context, *Role) error
	}

	DomainRoleV1 struct {
		dataRoleV1  DataRoleV1Adapter
		auditWriter a.AuditAdapter
	}
)

func NewDomainRoleV1(crolV1 DataRoleV1Adapter) *DomainRoleV1 {
	aw := a.AuditInit()
	return &DomainRoleV1{dataRoleV1: crolV1, auditWriter: aw}
}

func (m *DomainRoleV1) Get(ctx context.Context, rol *Role) error {
	if rol.Id == "" {
		return ae.MissingParamError("Id")
	}
	return m.dataRoleV1.Read(ctx, rol)
}

func (m *DomainRoleV1) Search(ctx context.Context, rol *[]Role, param RoleParam) (int, error) {
	// the second argument (map[string]string) is a list of columns to use for filtering
	// the key matches the json struct tag, the value is the actual table column name (this should change if aliases are used in your query)
	param.Param.CalculateParam("name", map[string]string{"id": "id", "name": "name", "description": "description"})
	param.Param.PaginationString = stor.FormatPagination(param.Param.Limit, param.Param.Offset)

	return m.dataRoleV1.ReadAll(ctx, rol, param)
}

func (m *DomainRoleV1) Post(ctx context.Context, rol *Role) error {
	if !rol.Name.Valid {
		return ae.MissingParamError("Name")
	}
	if rol.Name.Valid && len(rol.Name.ValueOrZero()) > 50 {
		return ae.StringLengthError("Name", 50)
	}
	rol.Id = function.GenerateRandomString(12)
	if err := m.dataRoleV1.Create(ctx, rol); err != nil {
		return err
	}
	go a.AuditCreate(m.auditWriter, *rol, RoleConst, a.KeysToString("id", rol.Id))
	return nil
}

func (m *DomainRoleV1) Patch(ctx context.Context, rolIn Role) error {
	rol := &Role{Id: rolIn.Id}
	errGet := m.dataRoleV1.Read(ctx, rol)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]any)
	// Name
	if rolIn.Name.Valid {
		if rolIn.Name.Valid && len(rolIn.Name.ValueOrZero()) > 50 {
			return ae.StringLengthError("Name", 50)
		}
		existingValues["name"] = rol.Name.String
		rol.Name = rolIn.Name
	}
	// Description
	if rolIn.Description.Valid {
		existingValues["description"] = rol.Description.String
		rol.Description = rolIn.Description
	}
	if err := m.dataRoleV1.Update(ctx, *rol); err != nil {
		return err
	}
	go a.AuditPatch(m.auditWriter, *rol, RoleConst, a.KeysToString("id", rol.Id), existingValues)
	return nil
}

func (m *DomainRoleV1) Delete(ctx context.Context, rol *Role) error {
	if rol.Id == "" {
		return ae.MissingParamError("Id")
	}
	if err := m.dataRoleV1.Delete(ctx, rol); err != nil {
		return err
	}
	go a.AuditDelete(m.auditWriter, *rol, RoleConst, a.KeysToString("id", rol.Id))
	return nil
}
