package http

import (
	"inventory_assets/config"
	"inventory_assets/internal/inventory_asset"

	_ "inventory_assets/internal/models"
	"inventory_assets/pkg/httpErrors"
	"inventory_assets/pkg/logger"
	"inventory_assets/pkg/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

// inventoryAssetHandlers represents the inventory asset HTTP handlers
type inventoryAssetHandlers struct {
	cfg              *config.Config
	inventoryAssetUC inventory_asset.UseCase
	logger           logger.Logger
}

// NewInventoryAssetHandlers creates a new inventory asset handlers instance
func NewInventoryAssetHandlers(cfg *config.Config, inventoryAssetUC inventory_asset.UseCase, logger logger.Logger) inventory_asset.Handlers {
	return &inventoryAssetHandlers{
		cfg:              cfg,
		inventoryAssetUC: inventoryAssetUC,
		logger:           logger,
	}
}

// GetInventoryAssetsById godoc
// @Summary Get inventory asset by ID
// @Description Retrieves a specific inventory asset by its ID
// @Tags Inventory Assets
// @Accept json
// @Produce json
// @Param id path int true "Inventory asset ID"
// @Success 200 {object} models.InventoryAssetList
// @Failure 400 {object} httpErrors.RestError
// @Failure 404 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /api/v1/inventory-assets/{id} [get]
func (h *inventoryAssetHandlers) GetInventoryAssetsById() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "inventoryAssetHandlers.GetInventoryAssetsById")
		defer span.Finish()

		// Parse ID parameter
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			h.logger.Errorf("Invalid ID parameter: %v", err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		// Get pagination parameters
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			h.logger.Errorf("Invalid pagination parameters: %v", err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		// Get inventory asset by ID
		result, err := h.inventoryAssetUC.GetInventoryAssetsById(ctx, pq, id)
		if err != nil {
			h.logger.Errorf("Failed to get inventory asset by ID %d: %v", id, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		h.logger.Infof("Successfully retrieved inventory asset by ID: %d", id)
		return c.JSON(http.StatusOK, result)
	}
}

// GetInventoryAssets godoc
// @Summary Get all inventory assets
// @Description Retrieves all inventory assets with pagination support
// @Tags Inventory Assets
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of elements per page" default(10)
// @Param orderBy query string false "Sort field"
// @Success 200 {object} models.InventoryAssetList
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /api/v1/inventory-assets [get]
func (h *inventoryAssetHandlers) GetInventoryAssets() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "inventoryAssetHandlers.GetInventoryAssets")
		defer span.Finish()

		// Get pagination parameters
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			h.logger.Errorf("Invalid pagination parameters: %v", err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		// Get all inventory assets
		result, err := h.inventoryAssetUC.GetInventoryAssets(ctx, pq)
		if err != nil {
			h.logger.Errorf("Failed to get inventory assets: %v", err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		h.logger.Infof("Successfully retrieved %d inventory assets", len(*result.Items))
		return c.JSON(http.StatusOK, result)
	}
}
