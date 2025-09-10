package authrefresh

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDomainAuthRefreshV1_Get(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	mockDataAuthRefresh := NewMockDataAuthRefreshV1Adapter(ctrl)

	tests := []struct {
		name    string
		ar      *AuthRefresh
		wantErr bool
		calls   []*gomock.Call
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DomainAuthRefreshV1{dataAuthRefreshV1: mockDataAuthRefresh}
			err := m.Get(ctx, tt.ar)
			if !tt.wantErr {
				assert.Nil(t, err, "DomainAuthRefreshV1.Get().%s => expected not error; got: %s", tt.name, err)
			}
			if tt.wantErr {
				assert.NotNil(t, err, "DomainAuthRefreshV1.Get().%s => expected error: got nil", tt.name)
			}
		})
	}
}

func TestDomainAuthRefreshV1_Post(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	mockDataAuthRefresh := NewMockDataAuthRefreshV1Adapter(ctrl)

	tests := []struct {
		name    string
		ar      *AuthRefresh
		wantErr bool
		calls   []*gomock.Call
	}{
		{
			"successful",
			&AuthRefresh{ClientId: "a", Token: "a", CreatedAt: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)},
			false,
			[]*gomock.Call{mockDataAuthRefresh.EXPECT().Create(ctx, gomock.Any()).Return(nil).AnyTimes()},
		},
		{
			"failed - clientId",
			&AuthRefresh{ClientId: "", Token: "a", CreatedAt: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)},
			true,
			[]*gomock.Call{},
		},
		{
			"failed - length clientId",
			&AuthRefresh{ClientId: "0123456789012345678901234567890123456789", Token: "a", CreatedAt: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)},
			true,
			[]*gomock.Call{},
		},
		{
			"failed - token",
			&AuthRefresh{ClientId: "a", Token: "", CreatedAt: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)},
			true,
			[]*gomock.Call{},
		},
		{
			"failed - length token",
			&AuthRefresh{ClientId: "a", Token: "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789", CreatedAt: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)},
			true,
			[]*gomock.Call{},
		},
		{
			"failed - createdAt",
			&AuthRefresh{ClientId: "a", Token: "a", CreatedAt: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)},
			true,
			[]*gomock.Call{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DomainAuthRefreshV1{dataAuthRefreshV1: mockDataAuthRefresh}
			err := m.Post(ctx, tt.ar)
			if !tt.wantErr {
				assert.Nil(t, err, "DomainAuthRefreshV1.Create().%s => expected not error; got: %s", tt.name, err)
			}
			if tt.wantErr {
				assert.NotNil(t, err, "DomainAuthRefreshV1.Create().%s => expected error: got nil", tt.name)
			}
		})
	}
}

func TestDomainAuthRefreshV1_Patch(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	mockDataAuthRefresh := NewMockDataAuthRefreshV1Adapter(ctrl)

	tests := []struct {
		name    string
		body    AuthRefresh
		wantErr bool
		calls   []*gomock.Call
	}{
		{
			"successful",
			AuthRefresh{},
			false,
			[]*gomock.Call{
				mockDataAuthRefresh.EXPECT().Read(ctx, gomock.Any()).Return(nil),
				mockDataAuthRefresh.EXPECT().Update(ctx, gomock.Any()).Return(nil),
			},
		},
		{
			"invalid id",
			AuthRefresh{},
			true,
			[]*gomock.Call{
				mockDataAuthRefresh.EXPECT().Read(ctx, gomock.Any()).Return(fmt.Errorf("missing record")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DomainAuthRefreshV1{dataAuthRefreshV1: mockDataAuthRefresh}
			err := m.Patch(ctx, tt.body)
			if !tt.wantErr {
				assert.Nil(t, err, "DomainAuthRefreshV1.Update().%s => expected not error; got: %s", tt.name, err)
			}
			if tt.wantErr {
				assert.NotNil(t, err, "DomainAuthRefreshV1.Update().%s => expected error: got nil", tt.name)
			}
		})
	}
}

func TestDomainAuthRefreshV1_Delete(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	mockDataAuthRefresh := NewMockDataAuthRefreshV1Adapter(ctrl)

	tests := []struct {
		name    string
		ar      *AuthRefresh
		wantErr bool
		calls   []*gomock.Call
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DomainAuthRefreshV1{dataAuthRefreshV1: mockDataAuthRefresh}
			err := m.Delete(ctx, tt.ar)
			if !tt.wantErr {
				assert.Nil(t, err, "DomainAuthRefreshV1.Delete().%s => expected not error; got: %s", tt.name, err)
			}
			if tt.wantErr {
				assert.NotNil(t, err, "DomainAuthRefreshV1.Delete().%s => expected error: got nil", tt.name)
			}
		})
	}
}
