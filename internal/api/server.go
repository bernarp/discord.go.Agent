package api

import (
	"context"
	"errors"
	"net/http"

	"DiscordBotAgent/internal/core/module_manager"
	"DiscordBotAgent/internal/core/zap_logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	log    *zap_logger.Logger
	mm     *module_manager.Manager
	router *gin.Engine
	srv    *http.Server
}

func New(
	log *zap_logger.Logger,
	mm *module_manager.Manager,
) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	s := &Server{
		log:    log,
		mm:     mm,
		router: router,
	}

	router.Use(s.loggerMiddleware())
	router.Use(gin.Recovery())

	return s
}

func (s *Server) Start(port string) error {
	s.registerRoutes(port)

	s.srv = &http.Server{
		Addr:    ":" + port,
		Handler: s.router,
	}

	s.log.Info("http server starting", zap.String("port", port))

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("http server failed", zap.Error(err))
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}
	s.log.Info("shutting down http server")
	return s.srv.Shutdown(ctx)
}
