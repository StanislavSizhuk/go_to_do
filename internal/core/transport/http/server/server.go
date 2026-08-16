package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/StanislavSizhuk/go_to_do/internal/core/logger"
	core_http_middleware "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger

	middleware []core_http_middleware.Middleware
}

func NewHTTpServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware, 
) *HTTPServer {
	return &HTTPServer{
		mux:    http.NewServeMux(),
		config: config,
		log:    log,
		middleware: middleware,
	}
}

func (h *HTTPServer) RegisterAPIRoutes(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		h.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

func (h *HTTPServer) Run(ctx context.Context) error {
mux :=  core_http_middleware.ChainMiddleware(h.mux, h.middleware...)

	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: mux,

	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)
		h.log.Warn("starting HTTP server", zap.String("addr", h.config.Addr))
		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("HTTP server listen and serve: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("shutting down HTTP server")
		shutdowmCtx, cancel := context.WithTimeout(
			context.Background(),
			h.config.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdowmCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("HTTP server shutdown: %w", err)
		}

		h.log.Warn("HTTP server shutdown complete")
	}
	return nil
}
