package registerroute

import (
	"context"
	"net/http"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	r "github.com/blackflagsoftware/tithe-declare/internal/middleware/route"
	"github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"github.com/labstack/echo/v4"
)

type (
	RestRegisterRouteV1 struct{}
)

var (
	restV1   RestRegisterRouteV1
	domainV1 *DomainRegisterRouteV1
)

func InitializeRegisterRouteV1() *DomainRegisterRouteV1 {
	storV1 := InitStorageV1()
	domainV1 = NewDomainRegisterRouteV1(storV1)
	restV1 = *NewRestRegisterRouteV1()
	return domainV1
}

func RegisterRegisterRoute(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/register-route/:raw_path", Get)
	r.RegisterAndAdd(eg, http.MethodPost, "/register-route/search", Search)
	r.RegisterAndAdd(eg, http.MethodPost, "/register-route", Post)
	r.RegisterAndAdd(eg, http.MethodPatch, "/register-route", Patch)
	r.RegisterAndAdd(eg, http.MethodDelete, "/register-route/:raw_path", Delete)
	r.RegisterAndAdd(eg, http.MethodPost, "/register-route/bulk", Bulk)
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

func Bulk(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.Bulk(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

// V1
func NewRestRegisterRouteV1() *RestRegisterRouteV1 {
	return &RestRegisterRouteV1{}
}

func (h *RestRegisterRouteV1) Get(c echo.Context) error {
	ctx := context.Background()
	raw_path := c.Param("raw_path")
	registerRoute := &RegisterRoute{RawPath: raw_path}
	if err := domainV1.Get(ctx, registerRoute); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *registerRoute, nil)
}

func (h *RestRegisterRouteV1) Search(c echo.Context) error {
	ctx := context.Background()
	param := RegisterRouteParam{}
	if err := c.Bind(&param); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	registerRoutes := &[]RegisterRoute{}
	totalCount, err := domainV1.Search(ctx, registerRoutes, param)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *registerRoutes, &totalCount)
}

func (h *RestRegisterRouteV1) Post(c echo.Context) error {
	ctx := context.Background()
	registerRoute := RegisterRoute{}
	if err := c.Bind(&registerRoute); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Post(ctx, &registerRoute); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, registerRoute, nil)
}

func (h *RestRegisterRouteV1) Patch(c echo.Context) error {
	ctx := context.Background()
	reg := RegisterRoute{}
	if err := c.Bind(&reg); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Patch(ctx, reg); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestRegisterRouteV1) Delete(c echo.Context) error {
	ctx := context.Background()
	raw_path := c.Param("raw_path")
	registerRoute := &RegisterRoute{RawPath: raw_path}
	if err := domainV1.Delete(ctx, registerRoute); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestRegisterRouteV1) Bulk(c echo.Context) error {
	ctx := context.Background()
	bulk := BulkRegisterRoute{}
	if err := c.Bind(&bulk); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Bulk(ctx, bulk); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}
