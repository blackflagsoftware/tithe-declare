package authclientsecret

import (
	"context"
	"net/http"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	r "github.com/blackflagsoftware/tithe-declare/internal/middleware/route"
	"github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"github.com/labstack/echo/v4"
)

type (
	RestAuthClientSecretV1 struct{}
)

var (
	restV1   RestAuthClientSecretV1
	domainV1 *DomainAuthClientSecretV1
)

func InitializeAuthClientSecretV1() *DomainAuthClientSecretV1 {
	storV1 := InitStorageV1()
	domainV1 = NewDomainAuthClientSecretV1(storV1)
	restV1 = *NewRestAuthClientSecretV1()
	return domainV1
}

func RegisterAuthClientSecret(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/auth-client-secret/:client_id/secret/:secret", Get)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth-authorize/search", Search)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth-client-secret", Post)
	r.RegisterAndAdd(eg, http.MethodPatch, "/auth-client-secret", Patch)
	r.RegisterAndAdd(eg, http.MethodDelete, "/auth-client-secret/:client_id/secret/:secret", Delete)
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
func NewRestAuthClientSecretV1() *RestAuthClientSecretV1 {
	return &RestAuthClientSecretV1{}
}

func (h *RestAuthClientSecretV1) Get(c echo.Context) error {
	ctx := context.TODO()
	authClientSecret := &AuthClientSecret{}
	if err := domainV1.Get(ctx, authClientSecret); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *authClientSecret, nil)
}

func (h *RestAuthClientSecretV1) Search(c echo.Context) error {
	ctx := context.TODO()
	param := AuthClientSecretParam{}
	if err := c.Bind(&param); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	authClientSecrets := &[]AuthClientSecret{}
	totalCount, err := domainV1.Search(ctx, authClientSecrets, param)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *authClientSecrets, &totalCount)
}

func (h *RestAuthClientSecretV1) Post(c echo.Context) error {
	ctx := context.TODO()
	au := AuthClientSecret{}
	if err := c.Bind(&au); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Post(ctx, &au); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, au, nil)
}

func (h *RestAuthClientSecretV1) Patch(c echo.Context) error {
	ctx := context.TODO()
	au := AuthClientSecret{}
	if err := c.Bind(&au); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Patch(ctx, au); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestAuthClientSecretV1) Delete(c echo.Context) error {
	ctx := context.TODO()
	authClientSecret := &AuthClientSecret{}
	if err := domainV1.Delete(ctx, authClientSecret); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}
