package authclient

import (
	"context"
	"net/http"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	r "github.com/blackflagsoftware/tithe-declare/internal/middleware/route"
	"github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"github.com/labstack/echo/v4"
)

type (
	RestAuthClientV1 struct{}
)

var (
	restV1   RestAuthClientV1
	domainV1 *DomainAuthClientV1
)

func InitializeAuthClientV1() *DomainAuthClientV1 {
	storV1 := InitStorageV1()
	domainV1 = NewDomainAuthClientV1(storV1)
	restV1 = *NewRestAuthClientV1()
	return domainV1
}

func RegisterAuthClient(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/auth-client/:id", Get)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth-client/search", Search)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth-client", Post)
	r.RegisterAndAdd(eg, http.MethodPatch, "/auth-client", Patch)
	r.RegisterAndAdd(eg, http.MethodDelete, "/auth-client/:id", Delete)
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
func NewRestAuthClientV1() *RestAuthClientV1 {
	return &RestAuthClientV1{}
}

func (h *RestAuthClientV1) Get(c echo.Context) error {
	ctx := context.Background()
	id := c.Param("id")
	authClient := &AuthClient{Id: id}
	if err := domainV1.Get(ctx, authClient); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *authClient, nil)
}

func (h *RestAuthClientV1) Search(c echo.Context) error {
	ctx := context.Background()
	param := AuthClientParam{}
	if err := c.Bind(&param); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	authClients := &[]AuthClient{}
	totalCount, err := domainV1.Search(ctx, authClients, param)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *authClients, &totalCount)
}

func (h *RestAuthClientV1) Post(c echo.Context) error {
	ctx := context.Background()
	ac := AuthClient{}
	if err := c.Bind(&ac); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Post(ctx, &ac); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, ac, nil)
}

func (h *RestAuthClientV1) Patch(c echo.Context) error {
	ctx := context.Background()
	ac := AuthClient{}
	if err := c.Bind(&ac); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Patch(ctx, ac); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestAuthClientV1) Delete(c echo.Context) error {
	ctx := context.Background()
	id := c.Param("id")
	authClient := &AuthClient{Id: id}
	if err := domainV1.Delete(ctx, authClient); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}
