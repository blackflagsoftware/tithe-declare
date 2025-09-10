package login

import (
	"context"
	"net/http"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	r "github.com/blackflagsoftware/tithe-declare/internal/middleware/route"
	"github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"github.com/labstack/echo/v4"
)

type (
	RestLoginV1 struct{}
)

var (
	restV1   RestLoginV1
	domainV1 *DomainLoginV1
)

func InitializeLoginV1() *DomainLoginV1 {
	storV1 := InitSQLV1()
	domainV1 = NewDomainLoginV1(storV1)
	restV1 = *NewRestLoginV1()
	return domainV1
}

func RegisterLogin(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/login/:id", Get)
	r.RegisterAndAdd(eg, http.MethodPost, "/login/search", Search)
	r.RegisterAndAdd(eg, http.MethodPost, "/login", Post)
	r.RegisterAndAdd(eg, http.MethodPatch, "/login", Patch)
	r.RegisterAndAdd(eg, http.MethodDelete, "/login/:id", Delete)
	r.RegisterAndAdd(eg, http.MethodPatch, "/login/pwd", PatchPwd)
	r.RegisterAndAdd(eg, http.MethodPost, "/login/verify", Verify)
	r.RegisterAndAdd(eg, http.MethodGet, "/login/roles", WithRoles)
	r.RegisterAndAdd(eg, http.MethodPost, "/login/reset/pwd", PostPwd)
	r.RegisterAndAdd(eg, http.MethodGet, "/login/forgot-password/:email_addr", ProcessResetRequest)
	r.RegisterAndAdd(eg, http.MethodPost, "/login/sign-in", SignIn)
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

func PatchPwd(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.PatchPwd(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func Verify(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.Verify(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func WithRoles(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.WithRoles(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func PostPwd(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.PostPwd(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func ProcessResetRequest(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.ProcessResetRequest(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func SignIn(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.SignIn(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

// V1
func NewRestLoginV1() *RestLoginV1 {
	return &RestLoginV1{}
}

func (h *RestLoginV1) Get(c echo.Context) error {
	ctx := context.TODO()
	id := c.Param("id")
	login := &Login{Id: id}
	if err := domainV1.Get(ctx, login); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *login, nil)
}

func (h *RestLoginV1) Search(c echo.Context) error {
	ctx := context.TODO()
	param := LoginParam{}
	if err := c.Bind(&param); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	logins := &[]Login{}
	totalCount, err := domainV1.Search(ctx, logins, param)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *logins, &totalCount)
}

func (h *RestLoginV1) Post(c echo.Context) error {
	ctx := context.TODO()
	login := Login{}
	if err := c.Bind(&login); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Post(ctx, &login); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 201, login, nil)
}

func (h *RestLoginV1) Patch(c echo.Context) error {
	ctx := context.TODO()
	login := Login{}
	if err := c.Bind(&login); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.Patch(ctx, login); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestLoginV1) PatchPwd(c echo.Context) error {
	ctx := context.TODO()
	login := Login{}
	if err := c.Bind(&login); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.PatchPwd(ctx, login); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestLoginV1) Delete(c echo.Context) error {
	ctx := context.TODO()
	id := c.Param("id")
	login := &Login{Id: id}
	if err := domainV1.Delete(ctx, login); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestLoginV1) PostPwd(c echo.Context) error {
	ctx := context.TODO()
	pwd := PasswordReset{}
	if err := c.Bind(&pwd); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if err := domainV1.PwdReset(ctx, pwd); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestLoginV1) ProcessResetRequest(c echo.Context) error {
	ctx := context.TODO()
	emailAddr := c.Param("email_addr")
	resetRequest := &ResetRequest{EmailAddr: emailAddr}
	if err := domainV1.ProcessResetRequest(ctx, resetRequest); err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *RestLoginV1) SignIn(c echo.Context) error {
	ctx := context.TODO()
	login := Login{}
	if err := c.Bind(&login); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	token, err := domainV1.SignIn(ctx, login)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *RestLoginV1) Verify(c echo.Context) error {
	// if it made it this far, the jwt made it past the middleware auth layer
	return c.JSON(http.StatusOK, echo.Map{"status": "successful"})
}

func (h *RestLoginV1) WithRoles(c echo.Context) error {
	ctx := context.TODO()
	login := &[]LoginRoles{}
	totalCount, err := domainV1.WithRoles(ctx, login)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	return handler.FormatResponse(c, 200, *login, &totalCount)
}
