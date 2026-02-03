package api

import (
	"fmt"
	"net/http"
	"strings"

	"DiscordBotAgent/internal/api/apierror"
	"DiscordBotAgent/internal/core/module_manager"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) registerRoutes(port string) {
	s.router.GET(
		"/swagger", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
		},
	)

	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", port))
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/health", s.handleHealth)
		v1.GET("/modules", s.handleGetModules)
		v1.GET("/modules/detail", s.handleGetModuleDetail)
	}
}

// @Summary Health check
// @Description Get current server status
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/health [get]
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// @Summary Get modules list
// @Description Get list of registered modules with optional status filtering
// @Tags modules
// @Produce json
// @Param status query string false "Filter by status (enabled, disabled, error, dependency_disabled)"
// @Success 200 {array} module_manager.ModuleInfo
// @Router /api/v1/modules [get]
func (s *Server) handleGetModules(c *gin.Context) {
	statusFilter := strings.ToLower(c.Query("status"))
	allModules := s.mm.GetAllModules()

	if statusFilter == "" {
		c.JSON(http.StatusOK, allModules)
		return
	}

	filtered := make([]module_manager.ModuleInfo, 0)
	for _, m := range allModules {
		if strings.ToLower(string(m.Status)) == statusFilter {
			filtered = append(filtered, m)
		}
	}

	c.JSON(http.StatusOK, filtered)
}

// @Summary Get module detail
// @Description Get detailed information about a specific module using query parameter
// @Tags modules
// @Produce json
// @Param name query string true "Module Name"
// @Success 200 {object} module_manager.ModuleInfo
// @Failure 400 {object} apierror.ErrorResponse
// @Failure 404 {object} apierror.ErrorResponse
// @Router /api/v1/modules/detail [get]
func (s *Server) handleGetModuleDetail(c *gin.Context) {
	name := c.Query("name")

	if name == "" {
		apierror.Abort(c, apierror.Errors.INVALID_REQUEST.WithMeta("query parameter 'name' is required"))
		return
	}

	info, ok := s.mm.GetModuleInfo(name)
	if !ok {
		apierror.Abort(c, apierror.Errors.MODULE_NOT_FOUND)
		return
	}

	c.JSON(http.StatusOK, info)
}
