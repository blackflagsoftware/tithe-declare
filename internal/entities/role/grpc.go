package role

import (
	"context"
	"encoding/json"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	p "github.com/blackflagsoftware/tithe-declare/pkg/proto"
)

type (
	RoleGrpc struct {
		p.UnimplementedRoleServiceServer
		domainRole DomainRoleV1
	}
)

func NewRoleGrpc(mrol DomainRoleV1) *RoleGrpc {
	return &RoleGrpc{domainRole: mrol}
}

func (a *RoleGrpc) GetRole(ctx context.Context, in *p.RoleIDIn) (*p.RoleResponse, error) {
	result := &p.Result{Success: false}
	response := &p.RoleResponse{Result: result}
	rol := &Role{Id: in.Id}
	if err := a.domainRole.Get(ctx, rol); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	var err error
	response.Role, err = translateOut(rol)
	if err != nil {
		return response, err
	}
	response.Result.Success = true
	return response, nil
}

func (a *RoleGrpc) SearchRole(ctx context.Context, in *p.Role) (*p.RoleRepeatResponse, error) {
	roleParam := RoleParam{}
	result := &p.Result{Success: false}
	response := &p.RoleRepeatResponse{Result: result}
	rols := &[]Role{}
	if _, err := a.domainRole.Search(ctx, rols, roleParam); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	for _, a := range *rols {
		protoRole, err := translateOut(&a)
		if err != nil {
			return response, err
		}
		response.Role = append(response.Role, protoRole)
	}
	response.Result.Success = true
	return response, nil
}

func (a *RoleGrpc) CreateRole(ctx context.Context, in *p.Role) (*p.RoleResponse, error) {
	result := &p.Result{Success: false}
	response := &p.RoleResponse{Result: result}
	rol, err := translateIn(in)
	if err != nil {
		return response, err
	}
	if err := a.domainRole.Post(ctx, rol); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	var errTranslate error
	response.Role, errTranslate = translateOut(rol)
	if errTranslate != nil {
		return response, errTranslate
	}
	response.Result.Success = true
	return response, nil
}

func (a *RoleGrpc) UpdateRole(ctx context.Context, in *p.Role) (*p.Result, error) {
	response := &p.Result{Success: false}
	rol, err := translateIn(in)
	if err != nil {
		return response, err
	}
	if err := a.domainRole.Patch(ctx, *rol); err != nil {
		response.Error = err.Error()
		return response, err
	}
	response.Success = true
	return response, nil
}

func (a *RoleGrpc) DeleteRole(ctx context.Context, in *p.RoleIDIn) (*p.Result, error) {
	response := &p.Result{Success: false}
	rol := &Role{Id: in.Id}
	if err := a.domainRole.Delete(ctx, rol); err != nil {
		response.Error = err.Error()
		return response, err
	}
	response.Success = true
	return response, nil
}

func translateOut(rol *Role) (*p.Role, error) {
	protoRole := p.Role{}
	protoRole.Id = rol.Id
	protoRole.Name = rol.Name.String
	protoRole.Description = rol.Description.String
	return &protoRole, nil
}

func translateIn(in *p.Role) (*Role, error) {
	rol := Role{}
	rol.Id = in.Id
	rol.Name.Scan(in.Name)
	rol.Description.Scan(in.Description)
	return &rol, nil
}

// found these are slower; deprecated; keep them, just in case
func translateJsonOut(rol *Role) (*p.Role, error) {
	protoRole := p.Role{}
	outBytes, err := json.Marshal(rol)
	if err != nil {
		return &protoRole, ae.GeneralError("Unable to encode from Role", err)
	}
	err = json.Unmarshal(outBytes, &protoRole)
	if err != nil {
		return &protoRole, ae.GeneralError("Unable to decode to proto.Role", err)
	}
	return &protoRole, nil
}

func translateJsonIn(in *p.Role) (*Role, error) {
	rol := Role{}
	outBytes, err := json.Marshal(in)
	if err != nil {
		return &rol, ae.GeneralError("Unable to encode from proto.Role", err)
	}
	err = json.Unmarshal(outBytes, &rol)
	if err != nil {
		return &rol, ae.GeneralError("Unable to decode to Role", err)
	}
	return &rol, nil
}
