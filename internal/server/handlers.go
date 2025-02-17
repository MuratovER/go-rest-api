package server

import (
	"inventory_assets/pkg/csrf"
	"inventory_assets/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"inventory_assets/internal/inventory_asset/repository"
	"inventory_assets/internal/inventory_asset/usecase"

	inventoryAssetHttp "inventory_assets/internal/inventory_asset/delivery/http"
	apiMiddlewares "inventory_assets/internal/middleware"
)

// Map Server Handlers
func (s *Server) MapHandlers(e *echo.Echo) error {

	// Init repositories
	inventoryAssetRepo := repository.NewInventoryAssetRepository(s.db, s.logger)

	// Init useCases
	inventoryAssetUC := usecase.NewInventoryAssetUseCase(s.cfg, inventoryAssetRepo, s.logger)

	// Init handlers
	inventoryAssetHandlers := inventoryAssetHttp.NewInventoryAssetHandlers(s.cfg, inventoryAssetUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())

	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	if s.cfg.Server.Debug {
		e.Use(mw.DebugMiddleware)
	}

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	inventoryAssetGroup := v1.Group("/inventory-assets")

	inventoryAssetHttp.MapInventoryAssetRoutes(inventoryAssetGroup, inventoryAssetHandlers, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
