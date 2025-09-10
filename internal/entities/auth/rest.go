package auth

import (
	"context"
	"net/http"
	"net/url"
	"time"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	au "github.com/blackflagsoftware/tithe-declare/internal/entities/authauthorize"
	r "github.com/blackflagsoftware/tithe-declare/internal/middleware/route"
	"github.com/blackflagsoftware/tithe-declare/internal/util/function"
	"github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v3"
)

type (
	RestAuthV1 struct{}
)

var (
	restV1   RestAuthV1
	domainV1 *DomainAuthV1
)

func InitializeAuthV1() *DomainAuthV1 {
	domainV1 = NewDomainAuthV1()
	restV1 = *NewRestAuthV1()
	return domainV1
}

func RegisterAuth(eg *echo.Group) {
	r.RegisterAndAdd(eg, http.MethodGet, "/auth/oauth2/authorize", OAuth2Authorize)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth/oauth2/verify-consent", OAuth2VerifyConsent)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth/oauth2/sign-in", OAuth2SignIn)
	r.RegisterAndAdd(eg, http.MethodPost, "/auth/oauth2/token-exchange", OAuth2Exchange)
}

func OAuth2Authorize(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.OAuth2Authorize(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func OAuth2VerifyConsent(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.OAuth2VerifyConsent(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func OAuth2SignIn(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.OAuth2SignIn(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

func OAuth2Exchange(c echo.Context) error {
	version := c.Get("version").(string)
	if version == "v1" {
		return restV1.OAuth2Exchange(c)
	}
	return handler.FormatResponseWithError(c, ae.RouteVersionNotFoundError(c.Request().URL.Path, version))
}

// V1
func NewRestAuthV1() *RestAuthV1 {
	return &RestAuthV1{}
}

func (h *RestAuthV1) Redirect(c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, "https://localhost:3001/redirect")
}

func (h *RestAuthV1) OAuth2Authorize(c echo.Context) error {
	// return the login/consent form code
	// do a redirect and call our "consent" page, store this on that page and send when calling "sign-in"
	ctx := context.TODO()
	auth := OAuthLogin{}
	if err := c.Bind(&auth); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	if auth.ResponseType != "code" {
		responseErr := ae.GeneralError("Invalid response_type, 'code' is the only valid option", nil)
		return handler.FormatResponseWithError(c, responseErr)
	}
	redirectPath := url.URL{}
	redirectPath.Scheme = "https"
	clientName, err := domainV1.GetClientId(ctx, auth.ClientId)
	if err != nil {
		responseErr := ae.GeneralError("Not a valid client", err)
		return handler.FormatResponseWithError(c, responseErr)
	}
	query := c.Request().URL.Query()
	query.Add("client_id", auth.ClientId)
	query.Add("client_name", clientName)
	query.Add("scope", auth.Scope)
	query.Add("state", auth.State)
	query.Add("redirect_uri", auth.RedirectUri)
	query.Add("code_challenge", auth.CodeChallenge)
	query.Add("code_challenge_method", auth.CodeChallengeMethod)

	redirectPath.Opaque = "//localhost/consent"
	redirectPath.RawQuery = query.Encode()
	// payload := map[string]string{"consent_form_path": redirectPath.String()}

	return c.Redirect(http.StatusSeeOther, redirectPath.String())
	// return c.JSON(http.StatusOK, handler.FormatResponse(c, payload, nil, nil))
}

func (h *RestAuthV1) OAuth2SignIn(c echo.Context) error {
	// this is called from the login form or consent form
	// the consent form will validate the response here and redirect if needed
	ctx := context.TODO()
	login := OAuthLogin{}
	if err := c.Bind(&login); err != nil {
		// since this is coming from our consent form, rare but avoidable
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}
	authorizeId := function.GenerateRandomString(32) // or authorization_code
	authAuthorize := au.AuthAuthorize{
		Id:                   authorizeId,
		ClientId:             null.StringFrom(login.ClientId),
		Verifier:             null.StringFrom(login.CodeChallenge),
		VerifierEncodeMethod: null.StringFrom(login.CodeChallengeMethod),
		State:                null.StringFrom(login.State),
		Scope:                null.StringFrom(login.Scope),
		AuthorizedAt:         null.TimeFrom(time.Now().UTC()),
	}

	response, err := domainV1.OAuthSignIn(ctx, login, authAuthorize)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}
	payload := map[string]string{"code": response.AuthCode, "access_token": response.AccessToken}
	return handler.FormatResponse(c, 200, payload, nil)
}

func (h *RestAuthV1) OAuth2Exchange(c echo.Context) error {
	// grant_type=authorization_code
	ctx := context.TODO()
	oAuthToken := OAuthToken{}
	if err := c.Bind(&oAuthToken); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}

	authToken, err := domainV1.OAuthExchange(ctx, oAuthToken)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}

	return handler.FormatResponse(c, 200, authToken, nil)
}

func (h *RestAuthV1) OAuth2VerifyConsent(c echo.Context) error {
	// grant_type=authorization_code
	ctx := context.TODO()
	oAuthToken := OAuthToken{}
	if err := c.Bind(&oAuthToken); err != nil {
		bindErr := ae.BindError(err)
		return handler.FormatResponseWithError(c, bindErr)
	}

	authToken, err := domainV1.OAuthExchange(ctx, oAuthToken)
	if err != nil {
		apiError := err.(ae.ApiError)
		return handler.FormatResponseWithError(c, apiError)
	}

	return handler.FormatResponse(c, 200, authToken, nil)
}
