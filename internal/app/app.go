package app

import (
	"context"
	"garantex-monitor/config"
	"garantex-monitor/gen/pb"
	"garantex-monitor/internal/controller"
	"garantex-monitor/internal/health"
	"garantex-monitor/internal/service"
	"garantex-monitor/internal/storage"

	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Application interface {
	Runner
	Bootstraper
}

type Runner interface {
	Run()
}

type Bootstraper interface {
	Bootstrap(options ...interface{}) Runner
}

type App struct {
	conf    *config.Config
	logger  *zap.Logger
	grpcSrv *grpc.Server
	lis     net.Listener
	Sig     chan os.Signal
	connDB  *pgx.Conn
}

func NewApp(conf *config.Config, logger *zap.Logger) *App {
	return &App{
		conf:   conf,
		logger: logger,
		Sig:    make(chan os.Signal, 1),
	}
}

func (a *App) Bootstrap(options ...interface{}) Runner {
	conn, err := pgx.Connect(context.Background(), a.conf.DB)
	if err != nil {
		a.logger.Fatal(
			"failed to create database connection",
			zap.String("databaseURL", a.conf.DB),
			zap.Error(err))
	}

	m, err := migrate.New("file://./migrations", a.conf.DB)
	if err != nil {
		a.logger.Fatal(
			"failed to create migrator",
			zap.String("sourceURL", "file://./migrations"),
			zap.String("databaseURL", a.conf.DB),
			zap.Error(err))
	}
	if err = m.Up(); err != nil {
		//if errors.Is(err, )
		a.logger.Debug("failed to up migrations", zap.Error(err))
	}

	storage := storage.NewGMStorage(conn)
	service := service.NewGMService(storage, a.logger)
	controller := controller.NewGMController(service, a.logger)
	healhCheck := health.NewGMHealhCheck(conn, a.logger)

	grpcSrv := grpc.NewServer()
	srv := controller
	pb.RegisterGarantexMonitorServer(grpcSrv, srv)
	grpc_health_v1.RegisterHealthServer(grpcSrv, healhCheck)

	lis, err := net.Listen("tcp", a.conf.Host+":"+a.conf.Port)
	if err != nil {
		a.logger.Fatal(
			"failed to create tcp listener",
			zap.String("host", a.conf.Host),
			zap.String("port", a.conf.Port),
			zap.Error(err))
	}

	a.connDB = conn
	a.grpcSrv = grpcSrv
	a.lis = lis

	return a
}

func (a *App) Run() {
	wg := sync.WaitGroup{}

	// Graceful
	wg.Add(1)
	go func() {
		// close DB connection
		defer func() {
			if err := a.connDB.Close(context.Background()); err != nil {
				a.logger.Error("failed to close database connection", zap.Error(err))
			}
			wg.Done()
		}()
		signal.Notify(a.Sig, syscall.SIGINT, syscall.SIGTERM)
		sig := <-a.Sig
		a.logger.Info("signal interrupt recieved", zap.Stringer("os_signal", sig))
		a.grpcSrv.GracefulStop()
	}()

	// Server
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.logger.Debug("start server", zap.String("address", a.lis.Addr().String()))
		if err := a.grpcSrv.Serve(a.lis); err != nil {
			a.logger.Fatal("failed to serve gRPC server", zap.Error(err))
		}
	}()

	wg.Wait()
}
