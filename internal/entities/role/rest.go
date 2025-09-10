package role

import (
	"context"
	"net/http"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	r "github.com/blackflagsoftware/tithe-declare/internal/middleware/route"
	"github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"github.com/labstack/echo/v4"
)

type (
	RestRoleV1 struct{}
)

var (
	restV1   RestRoleV1
	domainV1 *DomainRoleV1
)

func InitializeRoleV1() *DomainRoleV1 {
	storV1 := InitStorageV1()
	domainV1 = NewDomainRoleV1(storV1)
	restV1 = *NewRestRoleV1()
	return domainV1
}

func RegisterRole(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/role/:id", Get)
	r.RegisterAndAdd(eg, http.MethodPost, "/role/search", Search)
	r.RegisterAndAdd(eg, http.MethodPost, "/role", Post)
	r.RegisterAndAdd(eg, http.MethodPatch, "/role", Patch)
	r.RegisterAndAdd(eg, http.MethodDelete, "/role/:id", Delete)
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
func NewRestRoleV1() *RestRoleV1 {
	return &RestRoleV1{}
}

func (h *RestRoleV1) Get(c echo.Context) error {
	ctx := context.Background()
	id := c.Param("id")
	role := &Role{Id: id}
	if err := domainV1.Get(ctx, role); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *role, nil)
}

func (h *RestRoleV1) Search(c echo.Context) error {
	ctx := context.Background()
	param := RoleParam{}
	if err := c.Bind(&param); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	roles := &[]Role{}
	totalCount, err := domainV1.Search(ctx, roles, param)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *roles, &totalCount)
}

func (h *RestRoleV1) Post(c echo.Context) error {
	ctx := context.Background()
	rol := Role{}
	if err := c.Bind(&rol); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Post(ctx, &rol); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, rol, nil)
}

func (h *RestRoleV1) Patch(c echo.Context) error {
	ctx := context.Background()
	rol := Role{}
	if err := c.Bind(&rol); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Patch(ctx, rol); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestRoleV1) Delete(c echo.Context) error {
	ctx := context.Background()
	id := c.Param("id")
	role := &Role{Id: id}
	if err := domainV1.Delete(ctx, role); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}
