package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/StanislavSizhuk/go_to_do/internal/core/logger"
	core_pgx_pool "github.com/StanislavSizhuk/go_to_do/internal/core/repository/posgres/pool/pgx"
	core_http_middleware "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/middleware"
	core_http_server "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/StanislavSizhuk/go_to_do/internal/feature/tasks/repository/postgres"
	tasks_service "github.com/StanislavSizhuk/go_to_do/internal/feature/tasks/service"
	tasks_transport_http "github.com/StanislavSizhuk/go_to_do/internal/feature/tasks/transport/http"
	users_postgres_repository "github.com/StanislavSizhuk/go_to_do/internal/feature/user/repository/postgres"
	users_service "github.com/StanislavSizhuk/go_to_do/internal/feature/user/service"
	users_transport_http "github.com/StanislavSizhuk/go_to_do/internal/feature/user/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	time.Local = timeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())

	if err != nil {
		fmt.Println("failed to init application layer: ", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", timeZone))

	logger.Debug("init database pool")

	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init database pool: ", zap.Error(err))

	}
	defer pool.Close()

	logger.Debug("initialize feature ", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initialize HTTP server ")

	httpServer := core_http_server.NewHTTpServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)

	httpServer.RegisterAPIRoutes(apiVersionRouterV1)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
