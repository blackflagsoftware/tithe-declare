package auth

import (
	"gopkg.in/guregu/null.v3"
)

type (
	OAuthAuthResponse struct {
		ClientName string `json:"client_name"`
		AuthId     string `json:"auth_id"`
	}

	OAuthLogin struct {
		EmailAddr           null.String `db:"email_addr" json:"email_address"`
		Pwd                 null.String `db:"pwd" json:"password"`
		ClientId            string      `json:"client_id"`
		RedirectUri         string      `json:"redirect_uri"`
		Scope               string      `json:"scope"`
		State               string      `json:"state"`
		CodeChallenge       string      `json:"code_challenge"`
		CodeChallengeMethod string      `json:"code_challenge_method"`
		ResponseType        string      `json:"response_type"`
	}

	OAuthResponse struct {
		State       string `json:"state"`
		AuthCode    string `json:"code"`
		AccessToken string `json:"access_token"`
		RedirectUrl string `json:"redirect_url"`
	}

	OAuthToken struct {
		TokenType    string `json:"token_type"`
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"` // in seconds
		Scope        string `json:"scope"`      // comma delimited list
		RefreshToken string `json:"refresh_token"`
		GrantType    string `json:"grant_type,omitempty"` // authorization_code || refresh_token
		CodeVerifier string `json:"code_verifier,omitempty"`
		ClientId     string `json:"client_id,omitempty"`
		ClientSecret string `json:"client_secret,omitempty"`
		Code         string `json:"code,omitempty"`
		RedirectUri  string `json:"redirect_uri,omitempty"`
	}
)
