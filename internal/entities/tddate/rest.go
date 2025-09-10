package tddate

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
	RestTdDateV1 struct{}
)

var (
	restV1   RestTdDateV1
	domainV1 *DomainTdDateV1
)

func InitializeTdDateV1() *DomainTdDateV1 {
	storV1 := InitStorageV1()
	domainV1 = NewDomainTdDateV1(storV1)
	restV1 = *NewRestTdDateV1()
	return domainV1
}

func RegisterTdDate(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/td-date/:id", Get)
	r.RegisterAndAdd(eg, http.MethodPost, "/td-date/search", Search)
	r.RegisterAndAdd(eg, http.MethodPost, "/td-date", Post)
	r.RegisterAndAdd(eg, http.MethodPatch, "/td-date", Patch)
	r.RegisterAndAdd(eg, http.MethodDelete, "/td-date/:id", Delete)
	r.RegisterAndAdd(eg, http.MethodPost, "/td-date/block", CreateBlock)
	r.RegisterAndAdd(eg, http.MethodGet, "/td-date/current-days", GetCurrentDays)
	r.RegisterAndAdd(eg, http.MethodPost, "/td-date/check-hold-time", CheckHoldTime)
	r.RegisterAndAdd(eg, http.MethodPost, "/td-date/confirm", Confirm)
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

func CreateBlock(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.CreateBlock(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func GetCurrentDays(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.GetCurrentDays(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func CheckHoldTime(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.CheckSetHoldTime(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func Confirm(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.Confirm(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

// V1
func NewRestTdDateV1() *RestTdDateV1 {
	return &RestTdDateV1{}
}

func (h *RestTdDateV1) Get(c echo.Context) error {
	ctx := context.Background()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	tdDate := &TdDate{Id: int(id)}
	if err := domainV1.Get(ctx, tdDate); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *tdDate, nil)
}

func (h *RestTdDateV1) Search(c echo.Context) error {
	ctx := context.Background()
	param := TdDateParam{}
	if err := c.Bind(&param); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	tdDates := &[]TdDate{}
	totalCount, err := domainV1.Search(ctx, tdDates, param)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *tdDates, &totalCount)
}

func (h *RestTdDateV1) Post(c echo.Context) error {
	ctx := context.Background()
	td_ := TdDate{}
	if err := c.Bind(&td_); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Post(ctx, &td_); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, td_, nil)
}

func (h *RestTdDateV1) Patch(c echo.Context) error {
	ctx := context.Background()
	td_ := TdDate{}
	if err := c.Bind(&td_); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Patch(ctx, td_); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestTdDateV1) Delete(c echo.Context) error {
	ctx := context.Background()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	tdDate := &TdDate{Id: int(id)}
	if err := domainV1.Delete(ctx, tdDate); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestTdDateV1) CreateBlock(c echo.Context) error {
	ctx := context.Background()
	block := TdDateBlock{}
	if err := c.Bind(&block); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	err := domainV1.CreateBlock(ctx, block)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, nil, nil)
}

func (h *RestTdDateV1) GetCurrentDays(c echo.Context) error {
	ctx := context.Background()
	dayWithTimes := make(map[string][]string)
	if err := domainV1.GetCurrentDays(ctx, dayWithTimes); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, dayWithTimes, nil)
}

func (h *RestTdDateV1) CheckSetHoldTime(c echo.Context) error {
	ctx := context.Background()
	checkHold := CheckHoldTimeRequest{}
	if err := c.Bind(&checkHold); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	err := domainV1.CheckSetHoldTime(ctx, checkHold)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, nil, nil)
}

func (h *RestTdDateV1) Confirm(c echo.Context) error {
	ctx := context.Background()
	confirm := ConfirmRequest{}
	if err := c.Bind(&confirm); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	err := domainV1.Confirm(ctx, confirm)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, nil, nil)
}
