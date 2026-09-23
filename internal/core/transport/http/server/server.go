package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/MlkProduction/golang-project/internal/core/logger"
	core_http_middleware "github.com/MlkProduction/golang-project/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux         *http.ServeMux // мультиплексер, перенаправить мидлваре в обработчик
	config      Config	
	log         *core_logger.Logger // 
	middleweare []core_http_middleware.Middleware
}

func NewHTTPserver(config Config, log *core_logger.Logger, middleweare ...core_http_middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:         http.NewServeMux(),
		config:      config,
		log:         log,
		middleweare: middleweare,
	}
}

func (h *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		h.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}

func (h *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(h.mux, h.middleweare...)

	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)
		h.log.Warn("start HTTP server", zap.String("addr", h.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), h.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("Shutdown HTTP server:%w", err)
		}
		h.log.Warn("HTTP server stopped")
	}
	return nil
}
