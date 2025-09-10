package loginrole

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
	RestLoginRoleV1 struct{}
)

var (
	restV1   RestLoginRoleV1
	domainV1 *DomainLoginRoleV1
)

func InitializeLoginRoleV1() *DomainLoginRoleV1 {
	storV1 := InitStorageV1()
	domainV1 = NewDomainLoginRoleV1(storV1)
	restV1 = *NewRestLoginRoleV1()
	return domainV1
}

func RegisterLoginRole(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/login-role/:login_id/role_id/:role_id", Get)
	r.RegisterAndAdd(eg, http.MethodPost, "/login-role/search", Search)
	r.RegisterAndAdd(eg, http.MethodPost, "/login-role", Post)
	r.RegisterAndAdd(eg, http.MethodPost, "/login-role/bulk", Bulk)
	r.RegisterAndAdd(eg, http.MethodPatch, "/login-role", Patch)
	r.RegisterAndAdd(eg, http.MethodDelete, "/login-role/:login_id/role_id/:role_id", Delete)
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

func Bulk(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.Bulk(c)
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
func NewRestLoginRoleV1() *RestLoginRoleV1 {
	return &RestLoginRoleV1{}
}

func (h *RestLoginRoleV1) Get(c echo.Context) error {
	ctx := context.Background()
	login_id := c.Param("login_id")
	role_id := c.Param("role_id")
	loginRole := &LoginRole{LoginId: null.StringFrom(login_id), RoleId: null.StringFrom(role_id)}
	if err := domainV1.Get(ctx, loginRole); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *loginRole, nil)
}

func (h *RestLoginRoleV1) Search(c echo.Context) error {
	ctx := context.Background()
	param := LoginRoleParam{}
	if err := c.Bind(&param); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	loginRoles := &[]LoginRole{}
	totalCount, err := domainV1.Search(ctx, loginRoles, param)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *loginRoles, &totalCount)
}

func (h *RestLoginRoleV1) Post(c echo.Context) error {
	ctx := context.Background()
	lr := LoginRole{}
	if err := c.Bind(&lr); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Post(ctx, &lr); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, lr, nil)
}

func (h *RestLoginRoleV1) Bulk(c echo.Context) error {
	ctx := context.Background()
	lr := LoginRoleUpdate{}
	if err := c.Bind(&lr); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Bulk(ctx, &lr); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, lr, nil)
}

func (h *RestLoginRoleV1) Patch(c echo.Context) error {
	ctx := context.Background()
	lr := LoginRole{}
	if err := c.Bind(&lr); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Patch(ctx, lr); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestLoginRoleV1) Delete(c echo.Context) error {
	ctx := context.Background()
	login_id := c.Param("login_id")
	role_id := c.Param("role_id")
	loginRole := &LoginRole{LoginId: null.StringFrom(login_id), RoleId: null.StringFrom(role_id)}
	if err := domainV1.Delete(ctx, loginRole); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}
