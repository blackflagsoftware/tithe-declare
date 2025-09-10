package authrefresh

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/blackflagsoftware/tithe-declare/config"
	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	a "github.com/blackflagsoftware/tithe-declare/internal/audit"
	"github.com/blackflagsoftware/tithe-declare/internal/util/function"
)

//go:generate mockgen -source=domain.go -destination=mock.go -package=authrefresh
type (
	DataAuthRefreshV1Adapter interface {
		Read(context.Context, *AuthRefresh) error
		ReadAll(context.Context, *[]AuthRefresh, AuthRefreshParam) (int, error)
		Create(context.Context, *AuthRefresh) error
		Update(context.Context, AuthRefresh) error
		Delete(context.Context, *AuthRefresh) error
		CycleRefreshToken(context.Context, AuthRefresh, AuthRefresh) error
	}

	DomainAuthRefreshV1 struct {
		dataAuthRefreshV1 DataAuthRefreshV1Adapter
		auditWriter       a.AuditAdapter
	}
)

func NewDomainAuthRefreshV1(car DataAuthRefreshV1Adapter) *DomainAuthRefreshV1 {
	aw := a.AuditInit()
	return &DomainAuthRefreshV1{dataAuthRefreshV1: car, auditWriter: aw}
}

func (d *DomainAuthRefreshV1) Get(ctx context.Context, ar *AuthRefresh) error {
	return d.dataAuthRefreshV1.Read(ctx, ar)
}

func (d *DomainAuthRefreshV1) Search(ctx context.Context, ar *[]AuthRefresh, param AuthRefreshParam) (int, error) {
	param.Param.CalculateParam("created_at", map[string]string{"created_at": "created_at"})

	return d.dataAuthRefreshV1.ReadAll(ctx, ar, param)
}

func (d *DomainAuthRefreshV1) Post(ctx context.Context, ar *AuthRefresh) error {
	if !ar.CreatedAt.IsZero() {
		return ae.MissingParamError("CreatedAt")
	}
	ar.CreatedAt = time.Now().UTC()
	if err := d.dataAuthRefreshV1.Create(ctx, ar); err != nil {
		return err
	}
	go a.AuditCreate(d.auditWriter, *ar, AuthRefreshConst, a.KeysToString("client_id", ar.ClientId, "token", ar.Token))
	return nil
}

func (d *DomainAuthRefreshV1) Patch(ctx context.Context, arIn AuthRefresh) error {
	ar := &AuthRefresh{ClientId: arIn.ClientId, Token: arIn.Token}
	errGet := d.dataAuthRefreshV1.Read(ctx, ar)
	if errGet != nil {
		return errGet
	}
	existingValues := make(map[string]interface{})

	if err := d.dataAuthRefreshV1.Update(ctx, *ar); err != nil {
		return err
	}
	go a.AuditPatch(d.auditWriter, *ar, AuthRefreshConst, a.KeysToString("client_id", ar.ClientId, "token", ar.Token), existingValues)
	return nil
}

func (d *DomainAuthRefreshV1) Delete(ctx context.Context, ar *AuthRefresh) error {
	if err := d.dataAuthRefreshV1.Delete(ctx, ar); err != nil {
		return err
	}
	go a.AuditDelete(d.auditWriter, *ar, AuthRefreshConst, a.KeysToString("client_id", ar.ClientId, "token", ar.Token))
	return nil
}

func (d *DomainAuthRefreshV1) CycleRefreshToken(ctx context.Context, authRefresh AuthRefresh) (string, error) {
	refreshTokenSize := 32
	expires := config.A.GetRefreshTokenExpires()
	refreshToken := function.GenerateRandomString(refreshTokenSize)
	buildNewAuth := authRefresh.Token == ""
	if buildNewAuth {
		authRefresh.Token = refreshToken
		authRefresh.Active = true
		authRefresh.CreatedAt = time.Now().UTC()
		if err := d.dataAuthRefreshV1.CycleRefreshToken(ctx, AuthRefresh{}, authRefresh); err != nil {
			return "", err
		}
		return refreshToken, nil
	}
	if err := d.dataAuthRefreshV1.Read(ctx, &authRefresh); err != nil {
		return "", fmt.Errorf("unable to read")
	}
	if !authRefresh.Active {
		// is duplicate
		return "", fmt.Errorf("duplicate")
	}
	if expires == -1 {
		// do nothing
		return authRefresh.Token, nil
	}
	expireDate := authRefresh.CreatedAt.Add(time.Duration(expires))
	if expireDate.After(time.Now().UTC()) {
		// the refresh token is still good, do nothing else
		return authRefresh.Token, nil
	}
	/* new auth
	deactivate any existing refresh tokens => update auth_refresh_token set active = false where login_id = :login_id
	generate new refresh token, insert into table => insert into auth_refresh_token
	done
	*/
	/* refresh
	based on authOld coming in
	determine if a new refresh token is needed:
		-1: send back old token
		0: always refresh
		>0: use this value and compare times, if expired, generate new token, if not, send back old
	call the db layer cycle refresh with oldauth and newauth
	*/
	// authRefresh => old refresh record
	hashedTokenOld := fmt.Sprintf("%x", sha256.Sum256([]byte(authRefresh.Token)))
	authRefresh.Token = hashedTokenOld
	authRefresh.Active = false
	// build new refresh token record
	hashedTokenNew := fmt.Sprintf("%x", sha256.Sum256([]byte(refreshToken)))
	authRefreshNew := AuthRefresh{ClientId: authRefresh.ClientId, Token: hashedTokenNew, Active: true, CreatedAt: time.Now().UTC()}
	if err := d.dataAuthRefreshV1.CycleRefreshToken(ctx, authRefresh, authRefreshNew); err != nil {
		return "", err
	}
	return refreshToken, nil
}
