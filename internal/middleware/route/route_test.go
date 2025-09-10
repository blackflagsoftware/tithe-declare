package route

import "testing"

func Test_transformPathToRegex(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"No transformation",
			args{path: "/role/user"},
			"/role/user",
		},
		{
			"1 transformation - end",
			args{path: "/role/:id"},
			"/role/.+",
		},
		{
			"1 transformation - middle",
			args{path: "/role/:id/name"},
			"/role/.+/name",
		},
		{
			"multiple transformation",
			args{path: "/role/:id/name/:name"},
			"/role/.+/name/.+",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformPathToRegex(tt.args.path); got != tt.want {
				t.Errorf("transformPathToRegex() = %v, want %v", got, tt.want)
			}
		})
	}
}
