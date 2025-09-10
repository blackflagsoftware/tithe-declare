package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchBuilder_AppendCompare(t *testing.T) {
	type fields struct {
		Params []string
		Values []any
	}
	type args struct {
		param   string
		compare string
		value   any
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantParam []string
	}{
		{
			"success - 1",
			fields{
				[]string{},
				[]any{},
			},
			args{
				"id",
				"=",
				1,
			},
			[]string{"id = ?"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SearchBuilder{
				Params: tt.fields.Params,
				Values: tt.fields.Values,
			}
			s.AppendCompare(tt.args.param, tt.args.compare, tt.args.value)
			assert.Equal(t, tt.wantParam, s.Params, "params are not equal")
		})
	}
}

func TestSearchBuilder_AppendLike(t *testing.T) {
	type fields struct {
		Params []string
		Values []any
	}
	type args struct {
		param string
		value string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantParam []string
	}{
		{
			"successful - 1",
			fields{
				[]string{},
				[]any{},
			},
			args{
				"addr",
				"street",
			},
			[]string{"addr LIKE '%street%'"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SearchBuilder{
				Params: tt.fields.Params,
				Values: tt.fields.Values,
			}
			s.AppendLike(tt.args.param, tt.args.value)
			assert.Equal(t, tt.wantParam, s.Params, "params are not equal")
		})
	}
}

func TestSearchBuilder_AppendNull(t *testing.T) {
	type fields struct {
		Params []string
		Values []any
	}
	type args struct {
		param    string
		wantNull bool
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantParam []string
	}{
		{
			"successful - null",
			fields{
				[]string{},
				[]any{},
			},
			args{
				"name",
				true,
			},
			[]string{"name IS NULL"},
		},
		{
			"successful - not null",
			fields{
				[]string{},
				[]any{},
			},
			args{
				"name",
				false,
			},
			[]string{"name IS NOT NULL"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SearchBuilder{
				Params: tt.fields.Params,
				Values: tt.fields.Values,
			}
			s.AppendNull(tt.args.param, tt.args.wantNull)
			assert.Equal(t, tt.wantParam, s.Params, "params are not equal")
		})
	}
}

func TestSearchBuilder_String(t *testing.T) {
	type fields struct {
		Params []string
		Values []any
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"successful",
			fields{
				[]string{"name IS NOT NULL", "addr LIKE '%home%'"},
				[]any{},
			},
			"WHERE name IS NOT NULL\n\t\tAND addr LIKE '%home%'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SearchBuilder{
				Params: tt.fields.Params,
				Values: tt.fields.Values,
			}
			got := s.String()
			assert.Equal(t, tt.want, got, "output is not equal")
		})
	}
}

func TestSearchBuilder_AppendIn(t *testing.T) {
	type fields struct {
		Params []string
		Values []any
	}
	type args struct {
		param string
		value any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			"successful - empty",
			fields{
				[]string{},
				[]any{},
			},
			args{
				"addr",
				"street",
			},
			[]string{""},
		},
		{
			"successful - slice of string",
			fields{
				[]string{},
				[]any{},
			},
			args{
				"addr",
				[]string{"street", "home", "cow"},
			},
			[]string{"addr IN ('street', 'home', 'cow')"},
		},
		{
			"successful - slice of int",
			fields{
				[]string{},
				[]any{},
			},
			args{
				"addr",
				[]int{1, 101, 1000012},
			},
			[]string{"addr IN ('1', '101', '1000012')"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SearchBuilder{
				Params: tt.fields.Params,
				Values: tt.fields.Values,
			}
			s.AppendIn(tt.args.param, tt.args.value)
		})
	}
}
