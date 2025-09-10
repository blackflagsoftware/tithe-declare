package authauthorize

import (
	"context"
	"net/http"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	r "github.com/blackflagsoftware/tithe-declare/internal/middleware/route"
	"github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"github.com/labstack/echo/v4"
)

type (
	RestAuthAuthorizeV1 struct{}
)

var (
	restV1   RestAuthAuthorizeV1
	domainV1 *DomainAuthAuthorizeV1
)

func InitializeAuthAuthorizeV1() *DomainAuthAuthorizeV1 {
	storV1 := InitStorageV1()
	domainV1 = NewDomainAuthAuthorizeV1(storV1)
	restV1 = *NewRestAuthAuthorizeV1()
	return domainV1
}

func RegisterAuthAuthorize(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/auth-authorize/:id", Get)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth-authorize/search", Search)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth-authorize", Post)
	r.RegisterAndAdd(eg, http.MethodPatch, "/auth-authorize", Patch)
	r.RegisterAndAdd(eg, http.MethodDelete, "/auth-authorize/:id", Delete)
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
func NewRestAuthAuthorizeV1() *RestAuthAuthorizeV1 {
	return &RestAuthAuthorizeV1{}
}

func (h *RestAuthAuthorizeV1) Get(c echo.Context) error {
	ctx := context.Background()
	id := c.Param("id")
	authAuthorize := &AuthAuthorize{Id: id}
	if err := domainV1.Get(ctx, authAuthorize); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *authAuthorize, nil)
}

func (h *RestAuthAuthorizeV1) Search(c echo.Context) error {
	ctx := context.Background()
	param := AuthAuthorizeParam{}
	if err := c.Bind(&param); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	authAuthorizes := &[]AuthAuthorize{}
	totalCount, err := domainV1.Search(ctx, authAuthorizes, param)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *authAuthorizes, &totalCount)
}

func (h *RestAuthAuthorizeV1) Post(c echo.Context) error {
	ctx := context.Background()
	aa := AuthAuthorize{}
	if err := c.Bind(&aa); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Post(ctx, &aa); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, aa, nil)
}

func (h *RestAuthAuthorizeV1) Patch(c echo.Context) error {
	ctx := context.Background()
	aa := AuthAuthorize{}
	if err := c.Bind(&aa); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Patch(ctx, aa); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestAuthAuthorizeV1) Delete(c echo.Context) error {
	ctx := context.Background()
	id := c.Param("id")
	authAuthorize := &AuthAuthorize{Id: id}
	if err := domainV1.Delete(ctx, authAuthorize); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}
