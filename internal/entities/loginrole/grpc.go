package loginrole

import (
	"context"
	"encoding/json"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	p "github.com/blackflagsoftware/tithe-declare/pkg/proto"
	"gopkg.in/guregu/null.v3"
)

type (
	LoginRoleGrpc struct {
		p.UnimplementedLoginRoleServiceServer
		domainLoginRole DomainLoginRoleV1
	}
)

func NewLoginRoleGrpc(mlr DomainLoginRoleV1) *LoginRoleGrpc {
	return &LoginRoleGrpc{domainLoginRole: mlr}
}

func (a *LoginRoleGrpc) GetLoginRole(ctx context.Context, in *p.LoginRoleIDIn) (*p.LoginRoleResponse, error) {
	result := &p.Result{Success: false}
	response := &p.LoginRoleResponse{Result: result}
	lr := &LoginRole{LoginId: null.StringFrom(in.LoginId), RoleId: null.StringFrom(in.RoleId)}
	if err := a.domainLoginRole.Get(ctx, lr); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	var err error
	response.LoginRole, err = translateOut(lr)
	if err != nil {
		return response, err
	}
	response.Result.Success = true
	return response, nil
}

func (a *LoginRoleGrpc) SearchLoginRole(ctx context.Context, in *p.LoginRole) (*p.LoginRoleRepeatResponse, error) {
	loginRoleParam := LoginRoleParam{}
	result := &p.Result{Success: false}
	response := &p.LoginRoleRepeatResponse{Result: result}
	lrs := &[]LoginRole{}
	if _, err := a.domainLoginRole.Search(ctx, lrs, loginRoleParam); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	for _, a := range *lrs {
		protoLoginRole, err := translateOut(&a)
		if err != nil {
			return response, err
		}
		response.LoginRole = append(response.LoginRole, protoLoginRole)
	}
	response.Result.Success = true
	return response, nil
}

func (a *LoginRoleGrpc) CreateLoginRole(ctx context.Context, in *p.LoginRole) (*p.LoginRoleResponse, error) {
	result := &p.Result{Success: false}
	response := &p.LoginRoleResponse{Result: result}
	lr, err := translateIn(in)
	if err != nil {
		return response, err
	}
	if err := a.domainLoginRole.Post(ctx, lr); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	response.LoginRole, err = translateOut(lr)
	if err != nil {
		return response, err
	}
	response.Result.Success = true
	return response, nil
}

func (a *LoginRoleGrpc) BulkLoginRole(ctx context.Context, in *p.LoginRoleUpdate) (*p.LoginRoleUpdateResponse, error) {
	result := &p.Result{Success: false}
	response := &p.LoginRoleUpdateResponse{Result: result}
	lr := &LoginRoleUpdate{}
	lr.LoginId.Scan(in.LoginId)
	lr.RoleIds = in.RoleIds
	if err := a.domainLoginRole.Bulk(ctx, lr); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	// response.LoginRoleUpdate = lr TODO: if we really want to send back data, will need to add TranslateOutBulk (or something similar)
	response.Result.Success = true
	return response, nil
}

func (a *LoginRoleGrpc) UpdateLoginRole(ctx context.Context, in *p.LoginRole) (*p.Result, error) {
	response := &p.Result{Success: false}
	lr, err := translateIn(in)
	if err != nil {
		return response, err
	}
	if err := a.domainLoginRole.Patch(ctx, *lr); err != nil {
		response.Error = err.Error()
		return response, err
	}
	response.Success = true
	return response, nil
}

func (a *LoginRoleGrpc) DeleteLoginRole(ctx context.Context, in *p.LoginRoleIDIn) (*p.Result, error) {
	response := &p.Result{Success: false}
	lr := &LoginRole{LoginId: null.StringFrom(in.LoginId), RoleId: null.StringFrom(in.RoleId)}
	if err := a.domainLoginRole.Delete(ctx, lr); err != nil {
		response.Error = err.Error()
		return response, err
	}
	response.Success = true
	return response, nil
}

func translateOut(lr *LoginRole) (*p.LoginRole, error) {
	protoLoginRole := p.LoginRole{}
	protoLoginRole.LoginId = lr.LoginId.String
	protoLoginRole.RoleId = lr.RoleId.String
	return &protoLoginRole, nil
}

func translateIn(in *p.LoginRole) (*LoginRole, error) {
	lr := LoginRole{}
	lr.LoginId.Scan(in.LoginId)
	lr.RoleId.Scan(in.RoleId)
	return &lr, nil
}

// found these are slower; deprecated; keep them, just in case
func translateJsonOut(lr *LoginRole) (*p.LoginRole, error) {
	protoLoginRole := p.LoginRole{}
	outBytes, err := json.Marshal(lr)
	if err != nil {
		return &protoLoginRole, ae.GeneralError("Unable to encode from LoginRole", err)
	}
	err = json.Unmarshal(outBytes, &protoLoginRole)
	if err != nil {
		return &protoLoginRole, ae.GeneralError("Unable to decode to proto.LoginRole", err)
	}
	return &protoLoginRole, nil
}

func translateJsonIn(in *p.LoginRole) (*LoginRole, error) {
	lr := LoginRole{}
	outBytes, err := json.Marshal(in)
	if err != nil {
		return &lr, ae.GeneralError("Unable to encode from proto.LoginRole", err)
	}
	err = json.Unmarshal(outBytes, &lr)
	if err != nil {
		return &lr, ae.GeneralError("Unable to decode to LoginRole", err)
	}
	return &lr, nil
}
