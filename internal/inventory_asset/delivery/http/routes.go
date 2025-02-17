package http

import (
	"inventory_assets/internal/inventory_asset"
	"inventory_assets/internal/middleware"

	"github.com/labstack/echo/v4"
)

// Map inventory assets routes
func MapInventoryAssetRoutes(inventoryAssetGroup *echo.Group, handler inventory_asset.Handlers, mw *middleware.MiddlewareManager) {
	inventoryAssetGroup.GET("/:id", handler.GetInventoryAssetsById())
	inventoryAssetGroup.GET("", handler.GetInventoryAssets())
}
