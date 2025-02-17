package inventory_asset

import "github.com/labstack/echo/v4"

// Inventory Assets HTTP Handlers interface
type Handlers interface {
	GetInventoryAssetsById() echo.HandlerFunc
	GetInventoryAssets() echo.HandlerFunc
}
