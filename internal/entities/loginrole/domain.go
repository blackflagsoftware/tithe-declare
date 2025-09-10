package loginrole

import (
	"context"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	a "github.com/blackflagsoftware/tithe-declare/internal/audit"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	"github.com/blackflagsoftware/tithe-declare/internal/util/function"
	"github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

//go:generate mockgen -source=domain.go -destination=mock.go -package=loginrole
type (
	DataLoginRoleV1Adapter interface {
		Read(context.Context, *LoginRole) error
		ReadAll(context.Context, *[]LoginRole, LoginRoleParam) (int, error)
		Create(context.Context, *LoginRole) error
		Update(context.Context, LoginRole) error
		Delete(context.Context, *LoginRole) error
	}

	DomainLoginRoleV1 struct {
		dataLoginRoleV1 DataLoginRoleV1Adapter
		auditWriter     a.AuditAdapter
	}
)

func NewDomainLoginRoleV1(clrV1 DataLoginRoleV1Adapter) *DomainLoginRoleV1 {
	aw := a.AuditInit()
	return &DomainLoginRoleV1{dataLoginRoleV1: clrV1, auditWriter: aw}
}

func (m *DomainLoginRoleV1) Get(ctx context.Context, lr *LoginRole) error {

	return m.dataLoginRoleV1.Read(ctx, lr)
}

func (m *DomainLoginRoleV1) Search(ctx context.Context, lr *[]LoginRole, param LoginRoleParam) (int, error) {
	// the second argument (map[string]string) is a list of columns to use for filtering
	// the key matches the json struct tag, the value is the actual table column name (this should change if aliases are used in your query)
	param.Param.CalculateParam("login_id", map[string]string{"login_id": "login_id", "role_id": "role_id"})
	param.Param.PaginationString = stor.FormatPagination(param.Param.Limit, param.Param.Offset)

	return m.dataLoginRoleV1.ReadAll(ctx, lr, param)
}

func (m *DomainLoginRoleV1) Post(ctx context.Context, lr *LoginRole) error {
	if !lr.LoginId.Valid || lr.LoginId.String == "" {
		return ae.MissingParamError("LoginId")
	}
	return m.dataLoginRoleV1.Create(ctx, lr)
}

func (m *DomainLoginRoleV1) Bulk(ctx context.Context, lr *LoginRoleUpdate) error {
	if !lr.LoginId.Valid || lr.LoginId.String == "" {
		return ae.MissingParamError("LoginId")
	}
	if len(lr.RoleIds) == 0 {
		return ae.MissingParamError("RoleIds")
	}
	oldLoginRole := []LoginRole{}
	if _, err := m.Search(ctx, &oldLoginRole, LoginRoleParam{
		handler.Param{
			Search: handler.Search{
				Filters: []handler.Filter{
					{
						Column:  "login_id",
						Value:   lr.LoginId.String,
						Compare: "=",
					},
				},
			},
		},
	}); err != nil {
		return err
	}
	oldRoles := make([]string, len(oldLoginRole))
	for i, v := range oldLoginRole {
		oldRoles[i] = v.RoleId.String
	}
	functionAdd := function.ArrayDiff(oldRoles, lr.RoleIds)
	functionDelete := function.ArrayDiff(lr.RoleIds, oldRoles)
	for _, roleId := range functionAdd {
		lrAdd := &LoginRole{LoginId: lr.LoginId, RoleId: null.StringFrom(roleId)}
		if err := m.dataLoginRoleV1.Create(ctx, lrAdd); err != nil {
			return err
		}
	}
	for _, roleId := range functionDelete {
		lrDelete := &LoginRole{LoginId: lr.LoginId, RoleId: null.StringFrom(roleId)}
		if err := m.dataLoginRoleV1.Delete(ctx, lrDelete); err != nil {
			return err
		}
	}
	go a.AuditCreate(m.auditWriter, *lr, LoginRoleConst, a.KeysToString())
	return nil
}

func (m *DomainLoginRoleV1) Patch(ctx context.Context, lrIn LoginRole) error {
	lr := &LoginRole{LoginId: lrIn.LoginId, RoleId: lrIn.RoleId}
	errGet := m.dataLoginRoleV1.Read(ctx, lr)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]any)

	if err := m.dataLoginRoleV1.Update(ctx, *lr); err != nil {
		return err
	}
	go a.AuditPatch(m.auditWriter, *lr, LoginRoleConst, a.KeysToString("login_id", lr.LoginId, "role_id", lr.RoleId), existingValues)
	return nil
}

func (m *DomainLoginRoleV1) Delete(ctx context.Context, lr *LoginRole) error {

	if err := m.dataLoginRoleV1.Delete(ctx, lr); err != nil {
		return err
	}
	go a.AuditDelete(m.auditWriter, *lr, LoginRoleConst, a.KeysToString("login_id", lr.LoginId, "role_id", lr.RoleId))
	return nil
}
