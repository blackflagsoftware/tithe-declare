package loginreset

import (
	"context"
	"net/http"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	r "github.com/blackflagsoftware/tithe-declare/internal/middleware/route"
	"github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v3"
)

type (
	RestLoginResetV1 struct{}
)

var (
	restV1   RestLoginResetV1
	domainV1 *DomainLoginResetV1
)

func InitializeLoginResetV1() *DomainLoginResetV1 {
	storV1 := InitStorageV1()
	domainV1 = NewDomainLoginResetV1(storV1)
	restV1 = *NewRestLoginResetV1()
	return domainV1
}

func RegisterLoginReset(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/login-reset/:login_id/reset_token/:reset_token", Get)
	r.RegisterAndAdd(eg, http.MethodPost, "/login-reset/search", Search)
	r.RegisterAndAdd(eg, http.MethodPost, "/login-reset", Post)
	r.RegisterAndAdd(eg, http.MethodPatch, "/login-reset", Patch)
	r.RegisterAndAdd(eg, http.MethodDelete, "/login-reset/:login_id/reset_token/:reset_token", Delete)
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
func NewRestLoginResetV1() *RestLoginResetV1 {
	return &RestLoginResetV1{}
}

func (h *RestLoginResetV1) Get(c echo.Context) error {
	ctx := context.Background()
	login_id := c.Param("login_id")
	reset_token := c.Param("reset_token")
	loginReset := &LoginReset{LoginId: null.StringFrom(login_id), ResetToken: null.StringFrom(reset_token)}
	if err := domainV1.Get(ctx, loginReset); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *loginReset, nil)
}

func (h *RestLoginResetV1) Search(c echo.Context) error {
	ctx := context.Background()
	param := LoginResetParam{}
	if err := c.Bind(&param); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	loginResets := &[]LoginReset{}
	totalCount, err := domainV1.Search(ctx, loginResets, param)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *loginResets, &totalCount)
}

func (h *RestLoginResetV1) Post(c echo.Context) error {
	ctx := context.Background()
	lo := LoginReset{}
	if err := c.Bind(&lo); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Post(ctx, &lo); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, lo, nil)
}

func (h *RestLoginResetV1) Patch(c echo.Context) error {
	ctx := context.Background()
	lo := LoginReset{}
	if err := c.Bind(&lo); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Patch(ctx, lo); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestLoginResetV1) Delete(c echo.Context) error {
	ctx := context.Background()
	login_id := c.Param("login_id")
	reset_token := c.Param("reset_token")
	loginReset := &LoginReset{LoginId: null.StringFrom(login_id), ResetToken: null.StringFrom(reset_token)}
	if err := domainV1.Delete(ctx, loginReset); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}
