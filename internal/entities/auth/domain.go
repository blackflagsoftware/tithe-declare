package auth

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/blackflagsoftware/tithe-declare/config"
	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	aut "github.com/blackflagsoftware/tithe-declare/internal/entities/authauthorize"
	au "github.com/blackflagsoftware/tithe-declare/internal/entities/authclient"
	acc "github.com/blackflagsoftware/tithe-declare/internal/entities/authclientcallback"
	acs "github.com/blackflagsoftware/tithe-declare/internal/entities/authclientsecret"
	ar "github.com/blackflagsoftware/tithe-declare/internal/entities/authrefresh"
	l "github.com/blackflagsoftware/tithe-declare/internal/entities/login"
	"github.com/blackflagsoftware/tithe-declare/internal/middleware"
	"github.com/blackflagsoftware/tithe-declare/internal/util"
	"gopkg.in/guregu/null.v3"
)

type (
	DomainAuthV1 struct{}
)

func NewDomainAuthV1() *DomainAuthV1 {
	return &DomainAuthV1{}
}

func (d *DomainAuthV1) GetClientId(ctx context.Context, clientId string) (string, error) {
	// validate the client
	authClient := au.AuthClient{Id: clientId}
	acs := au.InitStorageV1()
	acm := au.NewDomainAuthClientV1(acs)
	if err := acm.Get(ctx, &authClient); err != nil {
		return "", err
	}
	return authClient.Name.String, nil
}

func (d *DomainAuthV1) OAuthSignIn(ctx context.Context, logIn OAuthLogin, authAuthorize aut.AuthAuthorize) (OAuthResponse, error) {
	oAuthResponse := OAuthResponse{}
	if !logIn.EmailAddr.Valid {
		return oAuthResponse, ae.MissingParamError("EmailAddress")
	}
	login := &l.Login{EmailAddr: logIn.EmailAddr}
	ls := l.InitSQLV1()
	err := ls.GetByEmailAddr(ctx, login)
	if err != nil {
		title := err.(ae.ApiError).BodyError().Title
		if title == "No Results Error" {
			return oAuthResponse, ae.EmailPasswordComboError()
		}
		return oAuthResponse, err
	}
	if err := util.CheckPassword(logIn.Pwd.String, login.Pwd.String); err != nil {
		return oAuthResponse, err
	}
	// make sure the redirect url is save on the client, another security check
	authClientCallback := acc.AuthClientCallback{ClientId: null.StringFrom(authAuthorize.ClientId.String), CallbackUrl: null.StringFrom("")} // TODO: need to figure this out, callback url
	accs := acc.InitStorageV1()
	accm := acc.NewDomainAuthClientCallbackV1(accs)
	if err := accm.Get(ctx, &authClientCallback); err != nil {
		return oAuthResponse, err
	}
	authAuthorize.AuthCodeAt = null.TimeFrom(time.Now().UTC())
	authAuthorize.ClientId = null.StringFrom(login.Id)
	// save off auth authorize record
	aas := aut.InitStorageV1()
	aam := aut.NewDomainAuthAuthorizeV1(aas)
	if err := aam.Post(ctx, &authAuthorize); err != nil {
		return oAuthResponse, err
	}
	accessToken, err := middleware.AuthBuild(authAuthorize.ClientId.String, []string{}) // TODO: add roles here for the user... ALSO: I don't think clientId is correct
	if err != nil {
		return oAuthResponse, ae.GeneralError("Unable to build access token", fmt.Errorf("Unable to build access token"))
	}

	oAuthResponse.AuthCode = authAuthorize.Id
	oAuthResponse.State = authAuthorize.State.String
	oAuthResponse.RedirectUrl = "" // TODO: same as above
	oAuthResponse.AccessToken = accessToken
	return oAuthResponse, nil
}

func (d *DomainAuthV1) OAuthExchange(ctx context.Context, authToken OAuthToken) (OAuthToken, error) {
	// check the client_secret is valid
	authClientSecret := acs.AuthClientSecret{ClientId: null.StringFrom(authToken.ClientId), Secret: null.StringFrom(authToken.ClientSecret)}
	acss := acs.InitStorageV1()
	acsm := acs.NewDomainAuthClientSecretV1(acss)
	if err := acsm.GetByIdAndSecret(ctx, &authClientSecret); err != nil {
		// TODO: send back different error if no record found
		return OAuthToken{}, err
	}
	if !authClientSecret.Secret.Valid {
		return OAuthToken{}, ae.GeneralError("Inactive secret", fmt.Errorf("Invalid secret"))
	}
	// get the loginId from AuthAuthorize
	authAuthorize := aut.AuthAuthorize{Id: authToken.Code}
	aas := aut.InitStorageV1()
	aam := aut.NewDomainAuthAuthorizeV1(aas)
	if err := aam.Get(ctx, &authAuthorize); err != nil {
		// TODO: send back different error if no record found
		return OAuthToken{}, ae.GeneralError("Bad AuthAuthorize", fmt.Errorf("Bad AuthAuthorize"))
	}
	// check the code
	if authToken.GrantType == "authorization_code" {
		return ExchangeAuthCode(ctx, authToken, authAuthorize)
	}
	if authToken.GrantType == "refresh_token" {
		return ExchangeRefreshToken(ctx, authToken, authAuthorize.ClientId.String)
	}
	return OAuthToken{}, fmt.Errorf("Invalid grant type")
}

func ExchangeAuthCode(ctx context.Context, authToken OAuthToken, authAuthorize aut.AuthAuthorize) (OAuthToken, error) {
	// will need to check the timing on when that code was created
	now := time.Now().UTC()
	expires := authAuthorize.AuthCodeAt.Time.Add(time.Duration(config.A.GetAuthorizationExpires()) * time.Second)
	if now.After(expires) {
		return OAuthToken{}, ae.GeneralError("The Auth Code has expired", fmt.Errorf("The Auth Code has expired"))
	}
	// base64 encode the code_verifier (pkce) and verify from the earlier saved off code_challenge (pkce) using the challenge_method
	coded, err := PkceCodeChallengeCheck(authToken.CodeVerifier, authAuthorize.VerifierEncodeMethod.String)
	if err != nil {
		return OAuthToken{}, err
	}
	if coded != authAuthorize.Verifier.String {
		return OAuthToken{}, ae.GeneralError("Invalid PKCE value", fmt.Errorf("Invalid Pkce value"))
	}
	// authRefresh := ar.AuthRefresh{LoginId: loginId, Token: authToken.RefreshToken}
	ars := ar.InitStorageV1()
	arm := ar.NewDomainAuthRefreshV1(ars)
	refreshToken, err := arm.CycleRefreshToken(ctx, ar.AuthRefresh{ClientId: authAuthorize.ClientId.String})
	if err != nil {
		return OAuthToken{}, nil
	}
	// build return authToken
	authTokenNew := OAuthToken{}
	authTokenNew.AccessToken, err = middleware.AuthBuild(authAuthorize.ClientId.String, []string{}) // TODO: add roles here for the user
	if err != nil {
		return OAuthToken{}, ae.GeneralError("Unable to build access token", fmt.Errorf("Unable to build access token"))
	}
	authTokenNew.ExpiresIn = config.A.GetAuthorizationExpires()
	authTokenNew.RefreshToken = refreshToken
	authTokenNew.Scope = "" // TODO: fill in your scope(s) here
	authTokenNew.TokenType = "bearer"
	return authTokenNew, nil
}

func ExchangeRefreshToken(ctx context.Context, authToken OAuthToken, clientId string) (OAuthToken, error) {
	// grant_type=refresh_token
	// TODO: we will want to hash the refresh token on save and then do the same here when comparing
	// TODO: add a config setting to either -1 [keep forever]; 0 [always refresh]; >0 [expire time interval]
	authRefresh := ar.AuthRefresh{ClientId: clientId, Token: authToken.RefreshToken}
	ars := ar.InitStorageV1()
	arm := ar.NewDomainAuthRefreshV1(ars)
	if err := arm.Get(ctx, &authRefresh); err != nil {
		// TODO: send back different error
		return OAuthToken{}, err
	}
	if authRefresh.ClientId == "" {
		return OAuthToken{}, fmt.Errorf("Inactive refresh token")
	}
	var err error
	authToken.AccessToken, err = middleware.AuthBuild(authRefresh.ClientId, []string{}) // TODO: add roles here for the user
	if err != nil {
		return OAuthToken{}, fmt.Errorf("Unable to build access token")
	}
	refreshToken, err := arm.CycleRefreshToken(ctx, ar.AuthRefresh{ClientId: clientId}) // empty => old
	if err != nil {
		return OAuthToken{}, err
	}
	authToken.ExpiresIn = config.A.GetAuthorizationExpires()
	authToken.RefreshToken = refreshToken
	authToken.Scope = ""
	authToken.TokenType = "bearer"
	return authToken, nil
}

func PkceCodeChallengeCheck(code, method string) (string, error) {
	switch method {
	case "S256":
		sum := sha256.Sum256([]byte(code))
		coded := base64.StdEncoding.EncodeToString(sum[:32])
		return coded, nil
	case "S512":
		sum := sha512.Sum512([]byte(code))
		coded := base64.StdEncoding.EncodeToString(sum[:64])
		return coded, nil
	default:
		return "", fmt.Errorf("Invalid encoded method")
	}
}
