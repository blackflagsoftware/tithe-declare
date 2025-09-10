package emailreminder

import (
	"context"
	"net/http"
	"strconv"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	r "github.com/blackflagsoftware/tithe-declare/internal/middleware/route"
	"github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"github.com/labstack/echo/v4"
)

type (
	RestEmailReminderV1 struct{}
)

var (
	restV1   RestEmailReminderV1
	domainV1 *DomainEmailReminderV1
)

func InitializeEmailReminderV1() *DomainEmailReminderV1 {
	storV1 := InitStorageV1()
	domainV1 = NewDomainEmailReminderV1(storV1)
	restV1 = *NewRestEmailReminderV1()
	return domainV1
}

func RegisterEmailReminder(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/email-reminder/:id", Get)
	r.RegisterAndAdd(eg, http.MethodPost, "/email-reminder/search", Search)
	r.RegisterAndAdd(eg, http.MethodPost, "/email-reminder", Post)
	r.RegisterAndAdd(eg, http.MethodPatch, "/email-reminder", Patch)
	r.RegisterAndAdd(eg, http.MethodDelete, "/email-reminder/:id", Delete)
}

func Get(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.Get(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func Search(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.Search(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func Post(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.Post(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func Patch(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.Patch(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func Delete(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.Delete(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

// V1
func NewRestEmailReminderV1() *RestEmailReminderV1 {
	return &RestEmailReminderV1{}
}

func (h *RestEmailReminderV1) Get(c echo.Context) error {
	ctx := context.Background()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	emailReminder := &EmailReminder{Id: int(id)}
	if err := domainV1.Get(ctx, emailReminder); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *emailReminder, nil)
}

func (h *RestEmailReminderV1) Search(c echo.Context) error {
	ctx := context.Background()
	param := EmailReminderParam{}
	if err := c.Bind(&param); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	emailReminders := &[]EmailReminder{}
	totalCount, err := domainV1.Search(ctx, emailReminders, param)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *emailReminders, &totalCount)
}

func (h *RestEmailReminderV1) Post(c echo.Context) error {
	ctx := context.Background()
	ema := EmailReminder{}
	if err := c.Bind(&ema); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Post(ctx, &ema); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, ema, nil)
}

func (h *RestEmailReminderV1) Patch(c echo.Context) error {
	ctx := context.Background()
	ema := EmailReminder{}
	if err := c.Bind(&ema); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Patch(ctx, ema); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestEmailReminderV1) Delete(c echo.Context) error {
	ctx := context.Background()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	emailReminder := &EmailReminder{Id: int(id)}
	if err := domainV1.Delete(ctx, emailReminder); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}
