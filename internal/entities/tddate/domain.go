package tddate

import (
	"context"
	"time"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	a "github.com/blackflagsoftware/tithe-declare/internal/audit"
	"github.com/blackflagsoftware/tithe-declare/internal/middleware/logging"
	stor "github.com/blackflagsoftware/tithe-declare/internal/storage"
	h "github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

//go:generate mockgen -source=domain.go -destination=mock.go -package=tddate
type (
	DataTdDateV1Adapter interface {
		Read(context.Context, *TdDate) error
		ReadAll(context.Context, *[]TdDate, TdDateParam) (int, error)
		Create(context.Context, *TdDate) error
		Update(context.Context, TdDate) error
		Delete(context.Context, *TdDate) error
		GetCurrentDays(context.Context, *[]time.Time, TdDateParam) error
		CheckSetHoldTime(context.Context, time.Time) error
		Confirm(context.Context, TdDate) error
	}

	DomainTdDateV1 struct {
		dataTdDateV1 DataTdDateV1Adapter
		auditWriter  a.AuditAdapter
	}
)

func NewDomainTdDateV1(ctd_V1 DataTdDateV1Adapter) *DomainTdDateV1 {
	aw := a.AuditInit()
	return &DomainTdDateV1{dataTdDateV1: ctd_V1, auditWriter: aw}
}

func (m *DomainTdDateV1) Get(ctx context.Context, td_ *TdDate) error {
	if td_.Id < 1 {
		return ae.MissingParamError("Id")
	}
	return m.dataTdDateV1.Read(ctx, td_)
}

func (m *DomainTdDateV1) Search(ctx context.Context, td_ *[]TdDate, param TdDateParam) (int, error) {
	// the second argument (map[string]string) is a list of columns to use for filtering
	// the key matches the json struct tag, the value is the actual table column name (this should change if aliases are used in your query)
	param.Param.CalculateParam("date_value", map[string]string{"id": "id", "date_value": "date_value", "hold": "hold", "confirm": "confirm", "name": "name", "phone": "phone", "email": "email"})
	param.Param.PaginationString = stor.FormatPagination(param.Param.Limit, param.Param.Offset)

	return m.dataTdDateV1.ReadAll(ctx, td_, param)
}

func (m *DomainTdDateV1) Post(ctx context.Context, td_ *TdDate) error {
	if !td_.DateValue.Valid {
		return ae.MissingParamError("DateValue")
	}
	if td_.Name.Valid && len(td_.Name.ValueOrZero()) > 255 {
		return ae.StringLengthError("Name", 255)
	}
	if td_.Phone.Valid && len(td_.Phone.ValueOrZero()) > 255 {
		return ae.StringLengthError("Phone", 255)
	}
	if td_.Email.Valid && len(td_.Email.ValueOrZero()) > 100 {
		return ae.StringLengthError("Email", 100)
	}
	if err := m.dataTdDateV1.Create(ctx, td_); err != nil {
		return err
	}
	go a.AuditCreate(m.auditWriter, *td_, TdDateConst, a.KeysToString("id", td_.Id))
	return nil
}

func (m *DomainTdDateV1) Patch(ctx context.Context, td_In TdDate) error {
	td_ := &TdDate{Id: td_In.Id}
	errGet := m.dataTdDateV1.Read(ctx, td_)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]any)
	// DateValue
	if td_In.DateValue.Valid {
		existingValues["date_value"] = td_.DateValue.Time.Format(time.RFC3339)
		td_.DateValue = td_In.DateValue
	}
	// Hold
	if td_In.Hold.Valid {
		existingValues["hold"] = td_.Hold.Time.Format(time.RFC3339)
		td_.Hold = td_In.Hold
	}
	// Confirm
	if td_In.Confirm.Valid {
		existingValues["confirm"] = td_.Confirm.Time.Format(time.RFC3339)
		td_.Confirm = td_In.Confirm
	}
	// Name
	if td_In.Name.Valid {
		if td_In.Name.Valid && len(td_In.Name.ValueOrZero()) > 255 {
			return ae.StringLengthError("Name", 255)
		}
		existingValues["name"] = td_.Name.String
		td_.Name = td_In.Name
	}
	// Phone
	if td_In.Phone.Valid {
		if td_In.Phone.Valid && len(td_In.Phone.ValueOrZero()) > 255 {
			return ae.StringLengthError("Phone", 255)
		}
		existingValues["phone"] = td_.Phone.String
		td_.Phone = td_In.Phone
	}
	// Email
	if td_In.Email.Valid {
		if td_In.Email.Valid && len(td_In.Email.ValueOrZero()) > 100 {
			return ae.StringLengthError("Email", 100)
		}
		existingValues["email"] = td_.Email.String
		td_.Email = td_In.Email
	}
	if err := m.dataTdDateV1.Update(ctx, *td_); err != nil {
		return err
	}
	go a.AuditPatch(m.auditWriter, *td_, TdDateConst, a.KeysToString("id", td_.Id), existingValues)
	return nil
}

func (m *DomainTdDateV1) Delete(ctx context.Context, td_ *TdDate) error {
	if td_.Id < 1 {
		return ae.MissingParamError("Id")
	}
	if err := m.dataTdDateV1.Delete(ctx, td_); err != nil {
		return err
	}
	go a.AuditDelete(m.auditWriter, *td_, TdDateConst, a.KeysToString("id", td_.Id))
	return nil
}

func (m *DomainTdDateV1) CreateBlock(ctx context.Context, block TdDateBlock) error {
	if !block.NewDate.Valid {
		return ae.MissingParamError("NewDate")
	}
	if !block.StartTime.Valid {
		return ae.MissingParamError("StartTime")
	}
	if !block.EndTime.Valid {
		return ae.MissingParamError("EndTime")
	}
	// new_date should be in YYYY-MM-DD format
	// start_time and end_time should be in HH:MM format (24 hour clock)
	// with start_time, increment by 15 minutes until end_time is reached
	layoutDateTime := "2006-01-02 15:04"
	startDT := block.NewDate.String + " " + block.StartTime.String
	endDT := block.NewDate.String + " " + block.EndTime.String
	startTime, errStart := time.Parse(layoutDateTime, startDT)
	if errStart != nil {
		return ae.ParseError("StartTime not in correct format")
	}
	endTime, errEnd := time.Parse(layoutDateTime, endDT)
	if errEnd != nil {
		return ae.ParseError("EndTime not in correct format")
	}
	t := startTime
	for {
		if t.After(endTime) {
			break
		}
		if t.Equal(endTime) {
			break
		}
		td_ := TdDate{DateValue: null.TimeFrom(t)}
		if err := m.dataTdDateV1.Create(ctx, &td_); err != nil {
			return err
		}
		go a.AuditCreate(m.auditWriter, td_, TdDateConst, a.KeysToString("id", td_.Id))
		t = t.Add(15 * time.Minute)
	}
	return nil
}

func (m *DomainTdDateV1) GetCurrentDays(ctx context.Context, dayWithTimes map[string][]string) error {
	param := TdDateParam{
		Param: h.Param{
			Search: h.Search{
				Filters: []h.Filter{
					{Column: "date_value", Compare: ">", Value: time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")},
					{Column: "hold", Compare: "NULL", Value: nil},
				},
				Sort: "date_value",
			},
		},
	}
	param.Param.CalculateParam("date_value", map[string]string{"id": "id", "date_value": "date_value", "hold": "hold", "confirm": "confirm", "name": "name", "phone": "phone", "email": "email"})
	dateTimes := []time.Time{}
	if err := m.dataTdDateV1.GetCurrentDays(ctx, &dateTimes, param); err != nil {
		return err
	}
	for _, day := range dateTimes {
		// dayOnly := day[:10] // YYYY-MM-DD
		// hourStr := day[11:13]
		// hourOnly, err := strconv.Atoi(hourStr)
		// if err != nil {
		// 	return ae.ParseError("DateValue not in correct format")
		// }
		// minOnly := day[14:16]
		// amPM := "AM"
		// if hourOnly > 11 {
		// 	amPM = "PM"
		// }
		// if hourOnly > 12 {
		// 	hourInt := (int(hourStr[0]-'0') * 10) + int(hourStr[1]-'0') - 12
		// 	hourStr = "" + string(rune(hourInt/10+'0')) + string(rune(hourInt%10+'0'))
		// }
		// dayWithTimes[dayOnly] = append(dayWithTimes[dayOnly], hourStr+":"+minOnly+" "+amPM)
		dayOnly := day.Format("2006-01-02")
		dayWithTimes[dayOnly] = append(dayWithTimes[dayOnly], day.Format("03:04 PM"))
	}
	return nil
}
func (m *DomainTdDateV1) CheckSetHoldTime(ctx context.Context, checkHold CheckHoldTimeRequest) error {
	dt, err := formatDateTime(checkHold)
	if err != nil {
		return err
	}
	return m.dataTdDateV1.CheckSetHoldTime(ctx, dt)
}

func (m *DomainTdDateV1) Confirm(ctx context.Context, confirm ConfirmRequest) error {
	dt, err := formatDateTime(confirm.CheckHoldTimeRequest)
	if err != nil {
		return err
	}
	confirm.DateValue = null.TimeFrom(dt)
	confirm.Confirm = null.TimeFrom(time.Now().UTC())
	return m.dataTdDateV1.Confirm(ctx, confirm.TdDate)
}

func (m *DomainTdDateV1) CheckHoldConfirm(ctx context.Context) error {
	param := TdDateParam{
		Param: h.Param{
			Search: h.Search{
				Filters: []h.Filter{
					{Column: "hold", Compare: "NOT NULL", Value: nil},
					{Column: "confirm", Compare: "NULL", Value: nil},
				},
			},
		},
	}
	param.Param.CalculateParam("date_value", map[string]string{"id": "id", "date_value": "date_value", "hold": "hold", "confirm": "confirm", "name": "name", "phone": "phone", "email": "email"})
	tdDates := []TdDate{}
	if _, err := m.dataTdDateV1.ReadAll(ctx, &tdDates, param); err != nil {
		return err
	}
	tenMinsAgo := time.Now().UTC().Add(-10 * time.Minute)
	for _, td := range tdDates {
		if td.Hold.Valid && td.Hold.Time.Before(tenMinsAgo) {
			td.Hold = null.Time{}
			if err := m.dataTdDateV1.Update(ctx, td); err != nil {
				logging.Default.Println("Error releasing hold on td_date id", td.Id, ":", err)
			}
		}
	}
	return nil
}

func formatDateTime(checkHold CheckHoldTimeRequest) (time.Time, error) {
	// date should be in YYYY-MM-DD format
	// time should be in HH:MM AM/PM format
	layoutDateTime := "2006-01-02 03:04 PM"
	dateTime := checkHold.Date + " " + checkHold.Time
	dt, errParse := time.Parse(layoutDateTime, dateTime)
	if errParse != nil {
		return time.Time{}, ae.ParseError("Date or Time not in correct format")
	}
	return dt, nil
}
