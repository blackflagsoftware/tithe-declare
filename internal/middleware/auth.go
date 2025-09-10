package middleware

import (
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/blackflagsoftware/tithe-declare/config"
	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	l "github.com/blackflagsoftware/tithe-declare/internal/middleware/logging"
	r "github.com/blackflagsoftware/tithe-declare/internal/middleware/route"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	CustomClaims struct {
		Roles []string `json:"roles"`
		jwt.RegisteredClaims
	}
)

// basic auth
// eg.Use(middleware.BasicAuthWithConfig(m.BasicAuthConfig()))
// also uncomment BasicAuthUser and BasicAuthPwd in config/config.go
func BasicAuthConfig() middleware.BasicAuthConfig {
	return middleware.BasicAuthConfig{
		Skipper:   SkipperBasicFunc(),
		Validator: BasicAuthFunc(),
	}
}

func BasicAuthFunc() middleware.BasicAuthValidator {
	return func(userName, userPwd string, c echo.Context) (bool, error) {
		// uncomment if using this feature
		if config.A.BasicAuthUser == userName && config.A.BasicAuthPwd == userPwd {
			c.Set("authenticated", "basic")
			c.Set("roles", []string{"admin"})
		}
		return true, nil
	}
}

func SkipperBasicFunc() func(echo.Context) bool {
	return func(c echo.Context) bool {
		bearer := "Bearer"
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		l := len(bearer)
		if len(auth) > l+1 && strings.EqualFold(auth[:l], bearer) {
			return true
		}
		return false
	}
}

// if you want to use this conjuction with the basic auth
// daisy chain the Echo Group eg.Use() calls with this one first
// eg.Use(middleware.BasicAuthWithConfig(m.BasicAuthConfig()))
// eg.Use(echojwt.WithConfig(m.AuthConfig()))
func AuthConfig() echojwt.Config {
	keyContent, err := base64.StdEncoding.DecodeString(config.A.AuthSecret)
	if err != nil {
		l.Default.Println("Unable to DecodeString for auth secret", err)
		return echojwt.Config{
			Skipper: SkipperJWTFunc(),
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(CustomClaims)
			},
			KeyFunc: func(token *jwt.Token) (any, error) {
				publicKey, err := base64.StdEncoding.DecodeString(config.A.AuthPublic)
				if err != nil {
					return nil, fmt.Errorf("unable to decode public string: %w", err)
				}
				return publicKey, nil
			},
		}
	}
	return echojwt.Config{
		Skipper: SkipperJWTFunc(),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(CustomClaims)
		},
		SigningKey: keyContent,
	}
}

func AuthBuild(loginId string, roles []string) (string, error) {
	// build claims
	now := time.Now().UTC()
	hours := config.A.GetExpiresAtDuration()
	expiresAt := now.Add(time.Duration(hours) * time.Hour).UTC()
	claims := &CustomClaims{
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        loginId,
		},
	}
	// determine which alg to use and then create the jwt
	keyContent, err := base64.StdEncoding.DecodeString(config.A.AuthSecret)
	if err != nil {
		l.Default.Println("Unable to DecodeString for auth secret", err)
		return "", ae.GeneralError("auth: unable to sign token", nil)
	}
	alg := config.A.AuthAlg
	var (
		secretByte    any
		signingMethod jwt.SigningMethod
		ok            bool
	)
	switch alg {
	case "HMAC":
		secretByte = keyContent
		signingMethod = jwt.SigningMethodHS256 // doesn't like the 512
	case "ECDSA":
		block, _ := pem.Decode(keyContent)
		if block == nil || !strings.Contains(block.Type, "PRIVATE KEY") {
			l.Default.Print("Failed to decode PEM block containing private key")
			return "", ae.GeneralError("auth: unable to sign token", nil)
		}
		var err error
		secretByte, err = x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			l.Default.Println("Unable to parse private key:", err)
			return "", ae.GeneralError("auth: unable to sign token", nil)
		}
		signingMethod = jwt.SigningMethodES512
	case "RSA", "EdDSA":
		block, _ := pem.Decode(keyContent)
		if block == nil || !strings.Contains(block.Type, "PRIVATE KEY") {
			l.Default.Print("Failed to decode PEM block containing private key")
			return "", ae.GeneralError("auth: unable to sign token", nil)
		}
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			l.Default.Println("Unable to parse private key:", err)
			return "", ae.GeneralError("auth: unable to sign token", nil)
		}
		if alg == "RSA" {
			signingMethod = jwt.SigningMethodRS512
			secretByte, ok = parsedKey.(*rsa.PrivateKey)
			if !ok {
				l.Default.Println("Not a RSA key")
				return "", ae.GeneralError("auth: unable to sign token", nil)
			}
		}
		if alg == "EdDSA" {
			signingMethod = jwt.SigningMethodEdDSA
			secretByte, ok = parsedKey.(ed25519.PrivateKey)
			if !ok {
				l.Default.Print("Not an ecdsa key")
				return "", ae.GeneralError("auth: unable to sign token", nil)
			}
		}
	}
	token := jwt.NewWithClaims(signingMethod, claims)
	tokenStr, err := token.SignedString(secretByte)
	if err != nil {
		l.Default.Println("AuthBuild: error getting signed token:", err)
		return "", ae.GeneralError("auth: unable to sign token", nil)
	}
	return tokenStr, nil
}

func SkipperJWTFunc() func(echo.Context) bool {
	return func(c echo.Context) bool {
		if c.Get("authenticated") == "basic" {
			return true
		}
		uriPath := c.Request().URL.Path
		method := c.Request().Method
		restrictedRoles, err := r.GetRolesForRegisterRoute(method, uriPath)
		if err != nil {
			l.Default.Printf("error getting roles for route: %s, error: %v", uriPath, err)
			return false
		}
		if len(restrictedRoles) == 0 {
			return true // unrestricted route
		}
		return false
	}
}

func AuthorizationHandler(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		uriPath := c.Request().URL.Path
		method := c.Request().Method
		restrictedRoles, err := r.GetRolesForRegisterRoute(method, uriPath)
		if err != nil {
			err = ae.GeneralError(fmt.Sprintf("error getting roles for route: %s", uriPath), err)
			return
		}
		if len(restrictedRoles) == 0 {
			// unrestricted route
			if err := h(c); err != nil {
				c.Error(err)
			}
			return
		}
		userRoles := []string{}
		if c.Get("authenticated") == "basic" {
			// get them from the context because they are set by basic auth
			var ok bool
			userRoles, ok = c.Get("roles").([]string)
			if !ok {
				return ae.AuthorizationError("BasicAuth: No roles assigned to the user")
			}
		} else {
			// get them from the jwt
			user, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return ae.AuthorizationError("JWT malformed: No roles present")
			}
			claims, ok := user.Claims.(*CustomClaims)
			if !ok {
				return ae.MissingParamError("JWT Claims malformed: No roles present")
			}
			userRoles = claims.Roles
		}
		for _, role := range slices.All(userRoles) {
			if slices.Contains(restrictedRoles, role) {
				if err := h(c); err != nil {
					c.Error(err)
				}
				return
			}
		}
		err = ae.AuthorizationError("Insufficient role to access route")
		return
	}
}
