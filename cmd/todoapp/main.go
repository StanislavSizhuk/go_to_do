package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/StanislavSizhuk/go_to_do/internal/core/logger"
	core_postgres_pool "github.com/StanislavSizhuk/go_to_do/internal/core/repository/posgres/pool"
	core_http_middleware "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/middleware"
	core_http_server "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/server"
	users_postgres_repository "github.com/StanislavSizhuk/go_to_do/internal/feature/user/repository/postgres"
	users_service "github.com/StanislavSizhuk/go_to_do/internal/feature/user/service"
	users_transport_http "github.com/StanislavSizhuk/go_to_do/internal/feature/user/transport/http"
	"go.uber.org/zap"
)

func main() {

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

	logger.Debug("init database pool")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init database pool: ", zap.Error(err))

	}
	defer pool.Close()

	logger.Debug("initialize feature ", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)

	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initialize HTTP server ")

	httpServer := core_http_server.NewHTTpServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)



	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
