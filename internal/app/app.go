package app

import (
	"context"
	"errors"
	"garantex-monitor/config"
	"garantex-monitor/gen/pb"
	"garantex-monitor/internal/controller"
	"garantex-monitor/internal/health"
	"garantex-monitor/internal/service"
	"garantex-monitor/internal/storage"
	"garantex-monitor/pkg/metrics"
	"garantex-monitor/pkg/tracer"
	"net/http"

	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
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
	// Настройка експортера и провайдера для трейсов
	tr, err := tracer.NewTracerProvider(context.Background(), a.conf.TraceURL, a.conf.Name)
	if err != nil {
		a.logger.Fatal(
			"failed to create tracer provider",
			zap.String("trace url", a.conf.TraceURL),
			zap.String("service name", a.conf.Name),
			zap.Error(err))
	}
	otel.SetTracerProvider(tr)

	// Регистрация всех метрик и создание http сервера для prometheus
	metrics.MustAllRegister()

	//
	//
	//
	//
	//
	//
	//
	//
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
		if !errors.Is(err, migrate.ErrNoChange) {
			a.logger.Fatal("failed to up migrations", zap.Error(err))
		}
		a.logger.Warn("migrate", zap.Error(err))
	}

	storage := storage.NewGMStorage(conn)
	service := service.NewGMService(storage, a.logger)
	srv := controller.NewGMController(service, a.logger, "tracerName")
	healhCheck := health.NewGMHealhCheck(conn, a.logger)

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	pb.RegisterGarantexMonitorServer(s, srv)
	grpc_health_v1.RegisterHealthServer(s, healhCheck)

	lis, err := net.Listen("tcp", a.conf.Host+":"+a.conf.Port)
	if err != nil {
		a.logger.Fatal(
			"failed to create tcp listener",
			zap.String("host", a.conf.Host),
			zap.String("port", a.conf.Port),
			zap.Error(err))
	}

	a.connDB = conn
	a.grpcSrv = s
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

	// HTTP Server
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := http.ListenAndServe(":8080", mux)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("failed to serve http server", zap.Error(err))
		}
	}()

	// gRPC Server
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.logger.Debug("start server", zap.String("address", a.lis.Addr().String()))
		if err := a.grpcSrv.Serve(a.lis); err != nil {
			a.logger.Fatal("failed to serve gRPC server", zap.Error(err))
		}
	}()

	//

	wg.Wait()
}
