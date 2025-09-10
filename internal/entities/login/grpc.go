package login

import (
	"context"
	"encoding/json"
	"time"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	p "github.com/blackflagsoftware/tithe-declare/pkg/proto"
)

type (
	LoginGrpc struct {
		p.UnimplementedLoginServiceServer
		domainLogin DomainLoginV1
	}
)

func NewLoginGrpc(mlog DomainLoginV1) *LoginGrpc {
	return &LoginGrpc{domainLogin: mlog}
}

func (a *LoginGrpc) GetLogin(ctx context.Context, in *p.LoginIDIn) (*p.LoginResponse, error) {
	result := &p.Result{Success: false}
	response := &p.LoginResponse{Result: result}
	fin := &Login{Id: in.Id}
	if err := a.domainLogin.Get(ctx, fin); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	var err error
	response.Login, err = translateOut(fin)
	if err != nil {
		return response, err
	}
	response.Result.Success = true
	return response, nil
}

func (a *LoginGrpc) SearchLogin(ctx context.Context, in *p.Login) (*p.LoginRepeatResponse, error) {
	loginParam := LoginParam{}
	result := &p.Result{Success: false}
	response := &p.LoginRepeatResponse{Result: result}
	fins := &[]Login{}
	if _, err := a.domainLogin.Search(ctx, fins, loginParam); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	for _, a := range *fins {
		protoLogin, err := translateOut(&a)
		if err != nil {
			return response, err
		}
		response.Login = append(response.Login, protoLogin)
	}
	response.Result.Success = true
	return response, nil
}

func (a *LoginGrpc) CreateLogin(ctx context.Context, in *p.Login) (*p.LoginResponse, error) {
	result := &p.Result{Success: false}
	response := &p.LoginResponse{Result: result}
	fin, err := translateIn(in)
	if err != nil {
		return response, err
	}
	if err := a.domainLogin.Post(ctx, fin); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	var errTranslate error
	response.Login, errTranslate = translateOut(fin)
	if errTranslate != nil {
		return response, errTranslate
	}
	response.Result.Success = true
	return response, nil
}

func (a *LoginGrpc) UpdateLogin(ctx context.Context, in *p.Login) (*p.Result, error) {
	response := &p.Result{Success: false}
	fin, err := translateIn(in)
	if err != nil {
		return response, err
	}
	if err := a.domainLogin.Patch(ctx, *fin); err != nil {
		response.Error = err.Error()
		return response, err
	}
	response.Success = true
	return response, nil
}

func (a *LoginGrpc) DeleteLogin(ctx context.Context, in *p.LoginIDIn) (*p.Result, error) {
	response := &p.Result{Success: false}
	fin := &Login{Id: in.Id}
	if err := a.domainLogin.Delete(ctx, fin); err != nil {
		response.Error = err.Error()
		return response, err
	}
	response.Success = true
	return response, nil
}

func translateOut(fin *Login) (*p.Login, error) {
	protoLogin := p.Login{}
	protoLogin.Id = fin.Id
	protoLogin.EmailAddr = fin.EmailAddr.String
	protoLogin.Pwd = fin.Pwd.String
	protoLogin.Active = fin.Active.Bool
	protoLogin.SetPwd = fin.SetPwd.Bool
	protoLogin.CreatedAt = fin.CreatedAt.Time.Format(time.RFC3339)
	protoLogin.UpdatedAt = fin.UpdatedAt.Time.Format(time.RFC3339)
	return &protoLogin, nil
}

func translateIn(in *p.Login) (*Login, error) {
	fin := Login{}
	fin.Id = in.Id
	fin.EmailAddr.Scan(in.EmailAddr)
	fin.Pwd.Scan(in.Pwd)
	fin.Active.Scan(in.Active)
	fin.SetPwd.Scan(in.SetPwd)
	fin.CreatedAt.Scan(in.CreatedAt)
	fin.UpdatedAt.Scan(in.UpdatedAt)
	return &fin, nil
}

// found these are slower; deprecated; keep them, just in case
func translateJsonOut(fin *Login) (*p.Login, error) {
	protoLogin := p.Login{}
	outBytes, err := json.Marshal(fin)
	if err != nil {
		return &protoLogin, ae.GeneralError("Unable to encode from Login", err)
	}
	err = json.Unmarshal(outBytes, &protoLogin)
	if err != nil {
		return &protoLogin, ae.GeneralError("Unable to decode to proto.Login", err)
	}
	return &protoLogin, nil
}

func translateJsonIn(in *p.Login) (*Login, error) {
	fin := Login{}
	outBytes, err := json.Marshal(in)
	if err != nil {
		return &fin, ae.GeneralError("Unable to encode from proto.Login", err)
	}
	err = json.Unmarshal(outBytes, &fin)
	if err != nil {
		return &fin, ae.GeneralError("Unable to decode to Login", err)
	}
	return &fin, nil
}
