package authrefresh

import (
	"context"
	"net/http"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	r "github.com/blackflagsoftware/tithe-declare/internal/middleware/route"
	"github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"github.com/labstack/echo/v4"
)

type (
	RestAuthRefreshV1 struct{}
)

var (
	restV1   RestAuthRefreshV1
	domainV1 *DomainAuthRefreshV1
)

func InitializeAuthRefreshV1() *DomainAuthRefreshV1 {
	storV1 := InitStorageV1()
	domainV1 = NewDomainAuthRefreshV1(storV1)
	restV1 = *NewRestAuthRefreshV1()
	return domainV1
}

func RegisterAuthRefresh(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/auth-refresh/:client_id/token/:token", Get)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth-refresh/search", Search)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth-refresh", Post)
	r.RegisterAndAdd(eg, http.MethodPatch, "/auth-refresh", Patch)
	r.RegisterAndAdd(eg, http.MethodDelete, "/auth-refresh/:client_id/token/:token", Delete)
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

func NewRestAuthRefreshV1() *RestAuthRefreshV1 {
	return &RestAuthRefreshV1{}
}

func (h *RestAuthRefreshV1) Get(c echo.Context) error {
	ctx := context.Background()
	client_id := c.Param("client_id")
	token := c.Param("token")
	authRefresh := &AuthRefresh{ClientId: client_id, Token: token}
	if err := domainV1.Get(ctx, authRefresh); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *authRefresh, nil)
}

func (h *RestAuthRefreshV1) Search(c echo.Context) error {
	ctx := context.Background()
	param := AuthRefreshParam{}
	if err := c.Bind(&param); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	authRefreshs := &[]AuthRefresh{}
	totalCount, err := domainV1.Search(ctx, authRefreshs, param)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *authRefreshs, &totalCount)
}

func (h *RestAuthRefreshV1) Post(c echo.Context) error {
	ctx := context.Background()
	authRefresh := AuthRefresh{}
	if err := c.Bind(&authRefresh); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Post(ctx, &authRefresh); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, authRefresh, nil)
}

func (h *RestAuthRefreshV1) Patch(c echo.Context) error {
	ctx := context.Background()
	ar := AuthRefresh{}
	if err := c.Bind(&ar); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Patch(ctx, ar); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestAuthRefreshV1) Delete(c echo.Context) error {
	ctx := context.Background()
	client_id := c.Param("client_id")
	token := c.Param("token")
	authRefresh := &AuthRefresh{ClientId: client_id, Token: token}
	if err := domainV1.Delete(ctx, authRefresh); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}
