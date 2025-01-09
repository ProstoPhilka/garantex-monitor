package app

import (
	"garantex-monitor/config"
	"garantex-monitor/internal/controller"
	"garantex-monitor/internal/service"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"

	pb "garantex-monitor/gen/grpc"

	"go.uber.org/zap"
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
}

func NewApp(conf *config.Config, logger *zap.Logger) *App {
	return &App{
		conf:   conf,
		logger: logger,
		Sig:    make(chan os.Signal, 1),
	}
}

func (a *App) Bootstrap(options ...interface{}) Runner {
	a.logger.Debug("bootstrapping application")
	service := service.NewGarantexMonitor(a.logger)
	controller := controller.NewController(service, a.logger)

	a.logger.Debug("creating server")
	a.grpcSrv = grpc.NewServer()
	srv := controller
	pb.RegisterGarantexMonitorServer(a.grpcSrv, srv)

	a.logger.Debug("creating tcp listener", zap.String("host", a.conf.Host), zap.String("port", a.conf.Port))
	lis, err := net.Listen("tcp", a.conf.Host+":"+a.conf.Port)
	if err != nil {
		a.logger.Fatal("failed to create tcp listener", zap.Error(err))
	}
	a.lis = lis

	return a
}

func (a *App) Run() {
	wg := sync.WaitGroup{}

	// Graceful
	wg.Add(1)
	go func() {
		defer wg.Done()
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
