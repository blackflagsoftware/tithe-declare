package main

import (
	"net"
	"os"

	"github.com/blackflagsoftware/tithe-declare/config"
	ema "github.com/blackflagsoftware/tithe-declare/internal/entities/emailreminder"
	log "github.com/blackflagsoftware/tithe-declare/internal/entities/login"
	lr "github.com/blackflagsoftware/tithe-declare/internal/entities/loginrole"
	rol "github.com/blackflagsoftware/tithe-declare/internal/entities/role"
	td_ "github.com/blackflagsoftware/tithe-declare/internal/entities/tddate"
	l "github.com/blackflagsoftware/tithe-declare/internal/middleware/logging"
	pb "github.com/blackflagsoftware/tithe-declare/pkg/proto"
	mig "github.com/blackflagsoftware/tithe-declare/tools/migration/src"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
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
	tcpListener, err := net.Listen("tcp", ":"+config.Srv.GrpcPort)
	if err != nil {
		l.Default.Panic("Unable to start GRPC port:", err)
	}
	defer tcpListener.Close()
	s := grpc.NewServer()

	registerServices(s)

	reflection.Register(s)
	l.Default.Printf("Starting GRPC server on port: %s...\n", config.Srv.GrpcPort)
	s.Serve(tcpListener)
}

func registerServices(s *grpc.Server) {
	// Role
	drol := rol.InitializeRoleV1()
	hrol := rol.NewRoleGrpc(*drol)
	pb.RegisterRoleServiceServer(s, hrol)
	// Login
	dlog := log.InitializeLoginV1()
	hlog := log.NewLoginGrpc(*dlog)
	pb.RegisterLoginServiceServer(s, hlog)
	// LoginRole
	dlr := lr.InitializeLoginRoleV1()
	hlr := lr.NewLoginRoleGrpc(*dlr)
	pb.RegisterLoginRoleServiceServer(s, hlr)
	// TdDate
	dtd_ := td_.InitializeTdDateV1()
	htd_ := td_.NewTdDateGrpc(*dtd_)
	pb.RegisterTdDateServiceServer(s, htd_)
	// EmailReminder
	dema := ema.InitializeEmailReminderV1()
	hema := ema.NewEmailReminderGrpc(*dema)
	pb.RegisterEmailReminderServiceServer(s, hema)
}
