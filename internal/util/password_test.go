package util

import (
	"testing"

	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	"github.com/stretchr/testify/assert"
)

func TestPasswordValidator(t *testing.T) {
	type args struct {
		pwd     string
		confirm string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{
			"successful",
			args{
				pwd:     "Test!User0",
				confirm: "Test!User0",
			},
			false,
			"",
		},
		{
			"mismatch error",
			args{
				pwd:     "Test!User0",
				confirm: "Test!User",
			},
			true,
			"Invalid password, reason: passwords do not match",
		},
		{
			"length error",
			args{
				pwd:     "Test!U0",
				confirm: "Test!U0",
			},
			true,
			"Invalid password, reason: must be > 8 characters long",
		},
		{
			"upper error",
			args{
				pwd:     "test!user0",
				confirm: "test!user0",
			},
			true,
			"Invalid password, reason: must have one uppercase letter",
		},
		{
			"number error",
			args{
				pwd:     "Test!user",
				confirm: "Test!user",
			},
			true,
			"Invalid password, reason: must have one number",
		},
		{
			"special error",
			args{
				pwd:     "Testuser0",
				confirm: "Testuser0",
			},
			true,
			"Invalid password, reason: must have one special character [@$!%*?]",
		},
		{
			"combo error",
			args{
				pwd:     "Testuser",
				confirm: "Testuser",
			},
			true,
			"Invalid password, reason: must have one number, must have one special character [@$!%*?]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PasswordValidator(tt.args.pwd, tt.args.confirm)
			if tt.wantErr {
				assert.NotNil(t, err, "TestPasswordValidator(); want error got nil")
			}
			if !tt.wantErr {
				assert.Nil(t, err, "TestPasswordValidator(); want nil got err: %s", err)
			}
			if err != nil {
				be := err.(ae.ApiError).BodyError()
				assert.Equal(t, tt.err, be.Detail, "TestPasswordValidator(); want error: %s: got: %s", tt.err, err.Error())
			}

		})
	}
}
