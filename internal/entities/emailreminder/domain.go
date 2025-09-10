package emailreminder

import (
	"context"
	"time"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	a "github.com/blackflagsoftware/tithe-declare/internal/audit"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/tddate"
	"github.com/blackflagsoftware/tithe-declare/internal/middleware/logging"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	"github.com/blackflagsoftware/tithe-declare/internal/util/email"
	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
)

//go:generate mockgen -source=domain.go -destination=mock.go -package=emailreminder
type (
	DataEmailReminderV1Adapter interface {
		Read(context.Context, *EmailReminder) error
		ReadAll(context.Context, *[]EmailReminder, EmailReminderParam) (int, error)
		Create(context.Context, *EmailReminder) error
		Update(context.Context, EmailReminder) error
		Delete(context.Context, *EmailReminder) error
	}

	DomainEmailReminderV1 struct {
		dataEmailReminderV1 DataEmailReminderV1Adapter
		auditWriter         a.AuditAdapter
		emailer             email.Emailer
	}
)

func NewDomainEmailReminderV1(cemaV1 DataEmailReminderV1Adapter) *DomainEmailReminderV1 {
	aw := a.AuditInit()
	em := email.EmailInit()
	return &DomainEmailReminderV1{dataEmailReminderV1: cemaV1, auditWriter: aw, emailer: em}
}

func (m *DomainEmailReminderV1) Get(ctx context.Context, ema *EmailReminder) error {
	if ema.Id < 1 {
		return ae.MissingParamError("Id")
	}
	return m.dataEmailReminderV1.Read(ctx, ema)
}

func (m *DomainEmailReminderV1) Search(ctx context.Context, ema *[]EmailReminder, param EmailReminderParam) (int, error) {
	// the second argument (map[string]string) is a list of columns to use for filtering
	// the key matches the json struct tag, the value is the actual table column name (this should change if aliases are used in your query)
	param.Param.CalculateParam("email", map[string]string{"id": "id", "email": "email"})
	param.Param.PaginationString = stor.FormatPagination(param.Param.Limit, param.Param.Offset)

	return m.dataEmailReminderV1.ReadAll(ctx, ema, param)
}

func (m *DomainEmailReminderV1) Post(ctx context.Context, ema *EmailReminder) error {
	if !ema.Email.Valid {
		return ae.MissingParamError("Email")
	}
	if ema.Email.Valid && len(ema.Email.ValueOrZero()) > 50 {
		return ae.StringLengthError("Email", 50)
	}
	if err := m.dataEmailReminderV1.Create(ctx, ema); err != nil {
		return err
	}
	go a.AuditCreate(m.auditWriter, *ema, EmailReminderConst, a.KeysToString("id", ema.Id))
	return nil
}

func (m *DomainEmailReminderV1) Patch(ctx context.Context, emaIn EmailReminder) error {
	ema := &EmailReminder{Id: emaIn.Id}
	errGet := m.dataEmailReminderV1.Read(ctx, ema)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]any)
	// Email
	if emaIn.Email.Valid {
		if emaIn.Email.Valid && len(emaIn.Email.ValueOrZero()) > 50 {
			return ae.StringLengthError("Email", 50)
		}
		existingValues["email"] = ema.Email.String
		ema.Email = emaIn.Email
	}
	if err := m.dataEmailReminderV1.Update(ctx, *ema); err != nil {
		return err
	}
	go a.AuditPatch(m.auditWriter, *ema, EmailReminderConst, a.KeysToString("id", ema.Id), existingValues)
	return nil
}

func (m *DomainEmailReminderV1) Delete(ctx context.Context, ema *EmailReminder) error {
	if ema.Id < 1 {
		return ae.MissingParamError("Id")
	}
	if err := m.dataEmailReminderV1.Delete(ctx, ema); err != nil {
		return err
	}
	go a.AuditDelete(m.auditWriter, *ema, EmailReminderConst, a.KeysToString("id", ema.Id))
	return nil
}

func (m *DomainEmailReminderV1) SendEmail(ctx context.Context) error {
	now := time.Now().UTC()
	// if now.Weekday() == time.Friday && now.Hour() == 23 && now.Minute() == 0 {
	if now.Weekday() == time.Wednesday && now.Hour() == 0 && now.Minute() == 15 {
		logging.Default.Println("It's Friday at 11:00 PM UTC, sending email reminders...")
		sql := tddate.InitSQLV1()
		tdDomain := tddate.NewDomainTdDateV1(sql)
		param := tddate.TdDateParam{
			Param: h.Param{
				Search: h.Search{
					Filters: []h.Filter{
						{Column: "hold", Compare: "NOT NULL", Value: nil},
						{Column: "confirm", Compare: "NOT NULL", Value: nil},
						{Column: "date_value", Compare: ">", Value: now},
						{Column: "date_value", Compare: "<", Value: now.AddDate(0, 0, 7)},
					},
				},
			},
		}
		param.Param.CalculateParam("date_value", map[string]string{"id": "id", "date_value": "date_value", "hold": "hold", "confirm": "confirm", "name": "name", "phone": "phone", "email": "email"})
		tdDates := []tddate.TdDate{}
		_, err := tdDomain.Search(ctx, &tdDates, param)
		if err != nil {
			return err
		}
		if len(tdDates) == 0 {
			logging.Default.Println("No upcoming declaration dates found, no emails to send.")
			return nil
		}
		emailBody := "The following upcoming declaration dates have been scheduled:\n\n"
		for _, td := range tdDates {
			emailBody += "- "
			if td.DateValue.Valid {
				emailBody += td.DateValue.Time.Format("01/02/2006")
			} else {
				emailBody += "No Date"
			}
			emailBody += " for "
			if td.Name.Valid {
				emailBody += td.Name.String
			} else {
				emailBody += "No Name"
			}
			if td.Email.Valid {
				emailBody += " (Email: " + td.Email.String + ")"
			}
			emailBody += "\n"
		}
		emails := []EmailReminder{}
		if _, err := m.Search(ctx, &emails, EmailReminderParam{}); err != nil {
			return err
		}
		if len(emails) == 0 {
			logging.Default.Println("No email reminders configured, no emails to send.")
			return nil
		}
		emailAddresses := []string{}
		for _, er := range emails {
			if er.Email.Valid {
				emailAddresses = append(emailAddresses, er.Email.String)
			}
		}
		m.emailer.SendReminder(ctx, emailAddresses, emailBody)
	}
	return nil
}
