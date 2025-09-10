package emailreminder

import (
	"context"
	"encoding/json"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	p "github.com/blackflagsoftware/tithe-declare/pkg/proto"
)

type (
	EmailReminderGrpc struct {
		p.UnimplementedEmailReminderServiceServer
		domainEmailReminder DomainEmailReminderV1
	}
)

func NewEmailReminderGrpc(mema DomainEmailReminderV1) *EmailReminderGrpc {
	return &EmailReminderGrpc{domainEmailReminder: mema}
}

func (a *EmailReminderGrpc) GetEmailReminder(ctx context.Context, in *p.EmailReminderIDIn) (*p.EmailReminderResponse, error) {
	result := &p.Result{Success: false}
	response := &p.EmailReminderResponse{Result: result}
	ema := &EmailReminder{Id: int(in.Id)}
	if err := a.domainEmailReminder.Get(ctx, ema); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	var err error
	response.EmailReminder, err = translateOut(ema)
	if err != nil {
		return response, err
	}
	response.Result.Success = true
	return response, nil
}

func (a *EmailReminderGrpc) SearchEmailReminder(ctx context.Context, in *p.EmailReminder) (*p.EmailReminderRepeatResponse, error) {
	emailReminderParam := EmailReminderParam{}
	result := &p.Result{Success: false}
	response := &p.EmailReminderRepeatResponse{Result: result}
	emas := &[]EmailReminder{}
	if _, err := a.domainEmailReminder.Search(ctx, emas, emailReminderParam); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	for _, a := range *emas {
		protoEmailReminder, err := translateOut(&a)
		if err != nil {
			return response, err
		}
		response.EmailReminder = append(response.EmailReminder, protoEmailReminder)
	}
	response.Result.Success = true
	return response, nil
}

func (a *EmailReminderGrpc) CreateEmailReminder(ctx context.Context, in *p.EmailReminder) (*p.EmailReminderResponse, error) {
	result := &p.Result{Success: false}
	response := &p.EmailReminderResponse{Result: result}
	ema, err := translateIn(in)
	if err != nil {
		return response, err
	}
	if err := a.domainEmailReminder.Post(ctx, ema); err != nil {
		response.Result.Error = err.Error()
		return response, err
	}
	var errTranslate error
	response.EmailReminder, errTranslate = translateOut(ema)
	if errTranslate != nil {
		return response, errTranslate
	}
	response.Result.Success = true
	return response, nil
}

func (a *EmailReminderGrpc) UpdateEmailReminder(ctx context.Context, in *p.EmailReminder) (*p.Result, error) {
	response := &p.Result{Success: false}
	ema, err := translateIn(in)
	if err != nil {
		return response, err
	}
	if err := a.domainEmailReminder.Patch(ctx, *ema); err != nil {
		response.Error = err.Error()
		return response, err
	}
	response.Success = true
	return response, nil
}

func (a *EmailReminderGrpc) DeleteEmailReminder(ctx context.Context, in *p.EmailReminderIDIn) (*p.Result, error) {
	response := &p.Result{Success: false}
	ema := &EmailReminder{Id: int(in.Id)}
	if err := a.domainEmailReminder.Delete(ctx, ema); err != nil {
		response.Error = err.Error()
		return response, err
	}
	response.Success = true
	return response, nil
}

func translateOut(ema *EmailReminder) (*p.EmailReminder, error) {
	protoEmailReminder := p.EmailReminder{}
	protoEmailReminder.Id = int64(ema.Id)
	protoEmailReminder.Email = ema.Email.String
	return &protoEmailReminder, nil
}

func translateIn(in *p.EmailReminder) (*EmailReminder, error) {
	ema := EmailReminder{}
	ema.Id = int(in.Id)
	ema.Email.Scan(in.Email)
	return &ema, nil
}

// found these are slower; deprecated; keep them, just in case
func translateJsonOut(ema *EmailReminder) (*p.EmailReminder, error) {
	protoEmailReminder := p.EmailReminder{}
	outBytes, err := json.Marshal(ema)
	if err != nil {
		return &protoEmailReminder, ae.GeneralError("Unable to encode from EmailReminder", err)
	}
	err = json.Unmarshal(outBytes, &protoEmailReminder)
	if err != nil {
		return &protoEmailReminder, ae.GeneralError("Unable to decode to proto.EmailReminder", err)
	}
	return &protoEmailReminder, nil
}

func translateJsonIn(in *p.EmailReminder) (*EmailReminder, error) {
	ema := EmailReminder{}
	outBytes, err := json.Marshal(in)
	if err != nil {
		return &ema, ae.GeneralError("Unable to encode from proto.EmailReminder", err)
	}
	err = json.Unmarshal(outBytes, &ema)
	if err != nil {
		return &ema, ae.GeneralError("Unable to decode to EmailReminder", err)
	}
	return &ema, nil
}
