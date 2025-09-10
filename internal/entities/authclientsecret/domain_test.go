package authclientsecret

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestDomainAuthClientSecretV1_Get(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	mockDataAuthClientSecret := NewMockDataAuthClientSecretV1Adapter(ctrl)

	tests := []struct {
		name    string
		au      *AuthClientSecret
		wantErr bool
		calls   []*gomock.Call
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DomainAuthClientSecretV1{dataAuthClientSecretV1: mockDataAuthClientSecret}
			err := m.Get(ctx, tt.au)
			if !tt.wantErr {
				assert.Nil(t, err, "DomainAuthClientSecretV1.Get().%s => expected not error; got: %s", tt.name, err)
			}
			if tt.wantErr {
				assert.NotNil(t, err, "DomainAuthClientSecretV1.Get().%s => expected error: got nil", tt.name)
			}
		})
	}
}

func TestDomainAuthClientSecretV1_Post(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	mockDataAuthClientSecret := NewMockDataAuthClientSecretV1Adapter(ctrl)

	tests := []struct {
		name    string
		au      *AuthClientSecret
		wantErr bool
		calls   []*gomock.Call
	}{
		{
			"successful",
			&AuthClientSecret{ClientId: null.NewString("a", true), Secret: null.NewString("a", true)},
			false,
			[]*gomock.Call{mockDataAuthClientSecret.EXPECT().Create(ctx, gomock.Any()).Return(nil).AnyTimes()},
		},
		{
			"failed - clientId",
			&AuthClientSecret{ClientId: null.NewString("a", false), Secret: null.NewString("a", true)},
			true,
			[]*gomock.Call{},
		},
		{
			"failed - length clientId",
			&AuthClientSecret{ClientId: null.NewString("0123456789012345678901234567890123456789", true), Secret: null.NewString("a", true)},
			true,
			[]*gomock.Call{},
		},
		{
			"failed - secret",
			&AuthClientSecret{ClientId: null.NewString("a", true), Secret: null.NewString("a", false)},
			true,
			[]*gomock.Call{},
		},
		{
			"failed - length secret",
			&AuthClientSecret{ClientId: null.NewString("a", true), Secret: null.NewString("01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789", true)},
			true,
			[]*gomock.Call{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DomainAuthClientSecretV1{dataAuthClientSecretV1: mockDataAuthClientSecret}
			err := m.Post(ctx, tt.au)
			if !tt.wantErr {
				assert.Nil(t, err, "DomainAuthClientSecretV1.Create().%s => expected not error; got: %s", tt.name, err)
			}
			if tt.wantErr {
				assert.NotNil(t, err, "DomainAuthClientSecretV1.Create().%s => expected error: got nil", tt.name)
			}
		})
	}
}

func TestDomainAuthClientSecretV1_Patch(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	mockDataAuthClientSecret := NewMockDataAuthClientSecretV1Adapter(ctrl)

	tests := []struct {
		name    string
		body    AuthClientSecret
		wantErr bool
		calls   []*gomock.Call
	}{
		{
			"successful",
			AuthClientSecret{},
			false,
			[]*gomock.Call{
				mockDataAuthClientSecret.EXPECT().Read(ctx, gomock.Any()).Return(nil),
				mockDataAuthClientSecret.EXPECT().Update(ctx, gomock.Any()).Return(nil),
			},
		},
		{
			"invalid id",
			AuthClientSecret{},
			true,
			[]*gomock.Call{
				mockDataAuthClientSecret.EXPECT().Read(ctx, gomock.Any()).Return(fmt.Errorf("missing record")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DomainAuthClientSecretV1{dataAuthClientSecretV1: mockDataAuthClientSecret}
			err := m.Patch(ctx, tt.body)
			if !tt.wantErr {
				assert.Nil(t, err, "DomainAuthClientSecretV1.Update().%s => expected not error; got: %s", tt.name, err)
			}
			if tt.wantErr {
				assert.NotNil(t, err, "DomainAuthClientSecretV1.Update().%s => expected error: got nil", tt.name)
			}
		})
	}
}

func TestDomainAuthClientSecretV1_Delete(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	mockDataAuthClientSecret := NewMockDataAuthClientSecretV1Adapter(ctrl)

	tests := []struct {
		name    string
		au      *AuthClientSecret
		wantErr bool
		calls   []*gomock.Call
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DomainAuthClientSecretV1{dataAuthClientSecretV1: mockDataAuthClientSecret}
			err := m.Delete(ctx, tt.au)
			if !tt.wantErr {
				assert.Nil(t, err, "DomainAuthClientSecretV1.Delete().%s => expected not error; got: %s", tt.name, err)
			}
			if tt.wantErr {
				assert.NotNil(t, err, "DomainAuthClientSecretV1.Delete().%s => expected error: got nil", tt.name)
			}
		})
	}
}
