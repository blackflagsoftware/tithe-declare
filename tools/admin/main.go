package main

import (
	"context"
	"fmt"
	"os"

	"github.com/blackflagsoftware/tithe-declare/config"
	l "github.com/blackflagsoftware/tithe-declare/internal/entities/login"
	lr "github.com/blackflagsoftware/tithe-declare/internal/entities/loginrole"
	r "github.com/blackflagsoftware/tithe-declare/internal/entities/role"
	"github.com/blackflagsoftware/tithe-declare/internal/util/handler"
	"gopkg.in/guregu/null.v3"
)

// this will populate the admin user
// and populate the reset pwd process

func main() {
	ctx := context.TODO()
	login := &l.Login{EmailAddr: null.StringFrom(config.E.AdminEmail), SetPwd: null.BoolFrom(true)}
	dLog := l.InitializeLoginV1()
	if err := dLog.Post(ctx, login); err != nil {
		fmt.Println("Creating admin error:", err)
		os.Exit(1)
	}
	dRol := r.InitializeRoleV1()
	roles := []r.Role{}
	param := r.RoleParam{Param: handler.Param{Search: handler.Search{Filters: []handler.Filter{{Column: "name", Compare: "=", Value: "admin"}}}}}
	if _, errDB := dRol.Search(ctx, &roles, param); errDB != nil {
		fmt.Println("Creating admin unable to get admin role id:", errDB)
		os.Exit(1)
	}
	if len(roles) == 0 {
		fmt.Println("Creating admin unable to get admin role id, search failed")
		os.Exit(1)
	}
	adminId := roles[0].Id
	// insert login_role
	dLogRol := lr.InitializeLoginRoleV1()
	loginRole := &lr.LoginRole{LoginId: null.StringFrom(login.Id), RoleId: null.StringFrom(adminId)}
	if err := dLogRol.Post(ctx, loginRole); err != nil {
		fmt.Println("Create admin unable to set login role:", err)
		os.Exit(1)
	}
	os.Exit(0)
}
