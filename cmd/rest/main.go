package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/blackflagsoftware/tithe-declare/config"
	ae "github.com/blackflagsoftware/tithe-declare/internal/api_error"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/auth"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/authauthorize"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/authclient"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/authclientcallback"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/authclientsecret"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/authrefresh"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/emailreminder"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/login"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/loginreset"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/loginrole"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/registerroute"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/role"
	"github.com/blackflagsoftware/tithe-declare/internal/entities/tddate"
	mid "github.com/blackflagsoftware/tithe-declare/internal/middleware"
	l "github.com/blackflagsoftware/tithe-declare/internal/middleware/logging"
	rt "github.com/blackflagsoftware/tithe-declare/internal/middleware/route"
	mig "github.com/blackflagsoftware/tithe-declare/tools/migration/src"
	"github.com/labstack/echo-contrib/echoprometheus"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	setPidFile()

	// argument flag
	var restPort string
	flag.StringVar(&restPort, "restPort", "", "the port number used for the REST listener")

	flag.Parse()

	if restPort == "" {
		restPort = config.Srv.RestPort
	}

	migration()

	e := echo.New()
	e.HTTPErrorHandler = ae.ErrorHandler // set echo's error handler
	if !strings.Contains(config.Srv.Env, "prod") {
		l.Default.Infoln("Logging set to debug...")
		e.Debug = true
		e.Use(l.DebugHandler)
	}
	e.Use(
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		}), // TODO: may not need for production, remove if you don't need
		l.Handler,
	)
	if config.Srv.EnableMetrics {
		e.Use(echoprometheus.NewMiddleware("TITHE_DECLARE"))
		e.GET("/metrics", echoprometheus.NewHandler())
	}

	// set all non-endpoints here
	e.Static("/", config.Srv.DocumentDir)
	fmt.Println("serving static files from: ", config.Srv.DocumentDir)
	e.HEAD("/status", ServerStatus) // for traditional server check
	e.GET("/liveness", Liveness)    // for k8s liveness

	InitializeRoutes()
	RegisterRoutes(e)

	go func() {
		// if err := e.StartTLS(fmt.Sprintf(":%s", restPort), "", ""); err != nil && err != http.ErrServerClosed { // or TLS, supplying the cert/key files as needed
		if err := e.Start(fmt.Sprintf(":%s", restPort)); err != nil && err != http.ErrServerClosed {
			l.Default.Printf("graceful server stop with error: %s", err)
		}
	}()

	ctx, cancelCheck := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Minute):
				sql := tddate.InitSQLV1()
				mgr := tddate.NewDomainTdDateV1(sql)
				mgr.CheckHoldConfirm(ctx)
			}
		}
	}(ctx)

	ctx, cancelEmail := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Minute):
				sql := emailreminder.InitSQLV1()
				mgr := emailreminder.NewDomainEmailReminderV1(sql)
				if err := mgr.SendEmail(ctx); err != nil {
					l.Default.Printf("error sending email reminders: %s", err)
				}
			}
		}
	}(ctx)

	// main server wait to exit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	cancelCheck()
	cancelEmail()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		l.Default.Printf("graceful shutdown with error: %s", err)
	}
}

func setPidFile() {
	// purpose: to set the starting applications pid number to file
	if pidFile, err := os.Create(config.Srv.PidPath); err != nil {
		l.Default.Panicln("Unable to create pid file...")
	} else if _, err := pidFile.Write((fmt.Appendf([]byte{}, "%d", os.Getpid()))); err != nil {
		l.Default.Panicln("Unable to write pid to file...")
	}
}

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the TITHE_DECLARE API")
}

func ServerStatus(c echo.Context) error {
	c.Response().Header().Add("TITHE_DECLARE", config.Srv.AppVersion)
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

func Liveness(c echo.Context) error {
	return c.String(http.StatusOK, "live")
}

func InitializeRoutes() {
	regDomain := registerroute.InitializeRegisterRouteV1()
	role.InitializeRoleV1()
	rt.InitRoute(regDomain)
	authclientsecret.InitializeAuthClientSecretV1()
	authclientcallback.InitializeAuthClientCallbackV1()
	authrefresh.InitializeAuthRefreshV1()
	login.InitializeLoginV1()
	loginreset.InitializeLoginResetV1()
	loginrole.InitializeLoginRoleV1()
	authauthorize.InitializeAuthAuthorizeV1()
	authclient.InitializeAuthClientV1()
	auth.InitializeAuthV1()
	tddate.InitializeTdDateV1()
	emailreminder.InitializeEmailReminderV1()
}

func RegisterRoutes(e *echo.Echo) {
	// register all routes here
	routeGroup := e.Group("")
	routeGroup.Use(mid.VersionHandler)
	routeGroup.Use(middleware.BasicAuthWithConfig(mid.BasicAuthConfig())) // uncomment to use AuthBasic, see internal/middleware/auth.go for more info
	additionalMiddlewareSetup(routeGroup)
	registerroute.RegisterRegisterRoute(routeGroup)
	role.RegisterRole(routeGroup)
	authclientsecret.RegisterAuthClientSecret(routeGroup)
	authclientcallback.RegisterAuthClientCallback(routeGroup)
	authrefresh.RegisterAuthRefresh(routeGroup)
	login.RegisterLogin(routeGroup)
	loginreset.RegisterLoginReset(routeGroup)
	loginrole.RegisterLoginRole(routeGroup)
	authauthorize.RegisterAuthAuthorize(routeGroup)
	authclient.RegisterAuthClient(routeGroup)
	auth.RegisterAuth(routeGroup)
	tddate.RegisterTdDate(routeGroup)
	emailreminder.RegisterEmailReminder(routeGroup)
}

func additionalMiddlewareSetup(rg *echo.Group) {
	rg.Use(echojwt.WithConfig(mid.AuthConfig())) // by default this will use JWT authentication, see internal/middleware/auth.go for more info
	rg.Use(mid.AuthorizationHandler)             // this will need to be called after the AuthConfig middleware
}

func migration() {
	if config.Mig.Enable {
		err := os.MkdirAll(config.Mig.Dir, 0744)
		if err != nil {
			l.Default.Printf("Unable to make scripts/migrations directory structure: %s\n", err)
		}
		c := mig.Connection{
			Host:           config.DB.SqlitePath,
			MigrationPath:  config.Mig.Dir,
			SkipInitialize: config.Mig.SkipInit,
			Engine:         mig.EngineType(config.DB.Engine),
		}
		if err := mig.StartMigration(c); err != nil {
			l.Default.Panicf("Migration failed due to: %s", err)
		}
	}
}
