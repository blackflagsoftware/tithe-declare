package authclientcallback

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
	RestAuthClientCallbackV1 struct{}
)

var (
	restV1   RestAuthClientCallbackV1
	domainV1 *DomainAuthClientCallbackV1
)

func InitializeAuthClientCallbackV1() *DomainAuthClientCallbackV1 {
	storV1 := InitStorageV1()
	domainV1 = NewDomainAuthClientCallbackV1(storV1)
	restV1 = *NewRestAuthClientCallbackV1()
	return domainV1
}

func RegisterAuthClientCallback(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/auth-client-callback/:client_id/callback_url/:callback_url", Get)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth-client-callback/search", Search)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth-client-callback", Post)
	r.RegisterAndAdd(eg, http.MethodPatch, "/auth-client-callback", Patch)
	r.RegisterAndAdd(eg, http.MethodDelete, "/auth-client-callback/:client_id/callback_url/:callback_url", Delete)
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
func NewRestAuthClientCallbackV1() *RestAuthClientCallbackV1 {
	return &RestAuthClientCallbackV1{}
}

func (h *RestAuthClientCallbackV1) Get(c echo.Context) error {
	ctx := context.Background()
	client_id := c.Param("client_id")
	callback_url := c.Param("callback_url")
	authClientCallback := &AuthClientCallback{ClientId: null.StringFrom(client_id), CallbackUrl: null.StringFrom(callback_url)}
	if err := domainV1.Get(ctx, authClientCallback); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *authClientCallback, nil)
}

func (h *RestAuthClientCallbackV1) Search(c echo.Context) error {
	ctx := context.Background()
	param := AuthClientCallbackParam{}
	if err := c.Bind(&param); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	authClientCallbacks := &[]AuthClientCallback{}
	totalCount, err := domainV1.Search(ctx, authClientCallbacks, param)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *authClientCallbacks, &totalCount)
}

func (h *RestAuthClientCallbackV1) Post(c echo.Context) error {
	ctx := context.Background()
	au := AuthClientCallback{}
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

func (h *RestAuthClientCallbackV1) Patch(c echo.Context) error {
	ctx := context.Background()
	au := AuthClientCallback{}
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

func (h *RestAuthClientCallbackV1) Delete(c echo.Context) error {
	ctx := context.Background()
	client_id := c.Param("client_id")
	callback_url := c.Param("callback_url")
	authClientCallback := &AuthClientCallback{ClientId: null.StringFrom(client_id), CallbackUrl: null.StringFrom(callback_url)}
	if err := domainV1.Delete(ctx, authClientCallback); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}
