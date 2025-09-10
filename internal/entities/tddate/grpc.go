package tddate

import (
	"context"
	"encoding/json"
	"time"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	p "github.com/blackflagsoftware/tithe-declare/pkg/proto"
)

type (
	TdDateGrpc struct {
		p.UnimplementedTdDateServiceServer
		domainTdDate DomainTdDateV1
	}
)

func NewTdDateGrpc(mtd_ DomainTdDateV1) *TdDateGrpc {
	return &TdDateGrpc{domainTdDate: mtd_}
}

func (a *TdDateGrpc) GetTdDate(ctx context.Context, in *p.TdDateIDIn) (*p.TdDateResponse, error) {
	result := &p.Result{Success: false}
	response := &p.TdDateResponse{Result: result}
	td_ := &TdDate{Id: int(in.Id)}
	if err := a.domainTdDate.Get(ctx, td_); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	var err error
	response.TdDate, err = translateOut(td_)
	if err != nil {
		return response, err
	}
	response.Result.Success = true
	return response, nil
}

func (a *TdDateGrpc) SearchTdDate(ctx context.Context, in *p.TdDate) (*p.TdDateRepeatResponse, error) {
	tdDateParam := TdDateParam{}
	result := &p.Result{Success: false}
	response := &p.TdDateRepeatResponse{Result: result}
	td_s := &[]TdDate{}
	if _, err := a.domainTdDate.Search(ctx, td_s, tdDateParam); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	for _, a := range *td_s {
		protoTdDate, err := translateOut(&a)
		if err != nil {
			return response, err
		}
		response.TdDate = append(response.TdDate, protoTdDate)
	}
	response.Result.Success = true
	return response, nil
}

func (a *TdDateGrpc) CreateTdDate(ctx context.Context, in *p.TdDate) (*p.TdDateResponse, error) {
	result := &p.Result{Success: false}
	response := &p.TdDateResponse{Result: result}
	td_, err := translateIn(in)
	if err != nil {
		return response, err
	}
	if err := a.domainTdDate.Post(ctx, td_); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	var errTranslate error
	response.TdDate, errTranslate = translateOut(td_)
	if errTranslate != nil {
		return response, errTranslate
	}
	response.Result.Success = true
	return response, nil
}

func (a *TdDateGrpc) UpdateTdDate(ctx context.Context, in *p.TdDate) (*p.Result, error) {
	response := &p.Result{Success: false}
	td_, err := translateIn(in)
	if err != nil {
		return response, err
	}
	if err := a.domainTdDate.Patch(ctx, *td_); err != nil {
		response.Error = err.Error()
		return response, err
	}
	response.Success = true
	return response, nil
}

func (a *TdDateGrpc) DeleteTdDate(ctx context.Context, in *p.TdDateIDIn) (*p.Result, error) {
	response := &p.Result{Success: false}
	td_ := &TdDate{Id: int(in.Id)}
	if err := a.domainTdDate.Delete(ctx, td_); err != nil {
		response.Error = err.Error()
		return response, err
	}
	response.Success = true
	return response, nil
}

func translateOut(td_ *TdDate) (*p.TdDate, error) {
	protoTdDate := p.TdDate{}
	protoTdDate.Id = int64(td_.Id)
	protoTdDate.DateValue = td_.DateValue.Time.Format(time.RFC3339)
	protoTdDate.Hold = td_.Hold.Time.Format(time.RFC3339)
	protoTdDate.Confirm = td_.Confirm.Time.Format(time.RFC3339)
	protoTdDate.Name = td_.Name.String
	protoTdDate.Phone = td_.Phone.String
	protoTdDate.Email = td_.Email.String
	return &protoTdDate, nil
}

func translateIn(in *p.TdDate) (*TdDate, error) {
	td_ := TdDate{}
	td_.Id = int(in.Id)
	td_.DateValue.Scan(in.DateValue)
	td_.Hold.Scan(in.Hold)
	td_.Confirm.Scan(in.Confirm)
	td_.Name.Scan(in.Name)
	td_.Phone.Scan(in.Phone)
	td_.Email.Scan(in.Email)
	return &td_, nil
}

// found these are slower; deprecated; keep them, just in case
func translateJsonOut(td_ *TdDate) (*p.TdDate, error) {
	protoTdDate := p.TdDate{}
	outBytes, err := json.Marshal(td_)
	if err != nil {
		return &protoTdDate, ae.GeneralError("Unable to encode from TdDate", err)
	}
	err = json.Unmarshal(outBytes, &protoTdDate)
	if err != nil {
		return &protoTdDate, ae.GeneralError("Unable to decode to proto.TdDate", err)
	}
	return &protoTdDate, nil
}

func translateJsonIn(in *p.TdDate) (*TdDate, error) {
	td_ := TdDate{}
	outBytes, err := json.Marshal(in)
	if err != nil {
		return &td_, ae.GeneralError("Unable to encode from proto.TdDate", err)
	}
	err = json.Unmarshal(outBytes, &td_)
	if err != nil {
		return &td_, ae.GeneralError("Unable to decode to TdDate", err)
	}
	return &td_, nil
}
