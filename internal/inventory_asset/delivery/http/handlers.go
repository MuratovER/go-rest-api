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

// Inventory assets handlers
type inventoryAssetHandlers struct {
	cfg              *config.Config
	inventoryassetUC inventory_asset.UseCase
	logger           logger.Logger
}

// NewInventoryAssetHandlers Inventory Asset handlers constructor
func NewInventoryAssetHandlers(cfg *config.Config, inventoryassetUC inventory_asset.UseCase, logger logger.Logger) inventory_asset.Handlers {

	return &inventoryAssetHandlers{cfg: cfg, inventoryassetUC: inventoryassetUC, logger: logger}

}

// GetByID godoc
// @Summary Get by id inventory assets
// @Description Get by id inventory assets handler
// @Tags Inventory Assets
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.InventoryAssetList
// @Router /api/v1/inventory-assets/{id} [get]
func (handlers inventoryAssetHandlers) GetInventoryAssetsById() echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(echoContext), "inventoryAssetHandlers.GetByID")
		defer span.Finish()

		inventoryAssetId, err := strconv.Atoi(echoContext.Param("id"))
		if err != nil {
			utils.LogResponseError(echoContext, handlers.logger, err)
			return echoContext.JSON(httpErrors.ErrorResponse(err))
		}
		pq, err := utils.GetPaginationFromCtx(echoContext)
		if err != nil {
			utils.LogResponseError(echoContext, handlers.logger, err)
			return echoContext.JSON(httpErrors.ErrorResponse(err))

		}

		inventoryAssetById, err := handlers.inventoryassetUC.GetInventoryAssetsById(ctx, pq, inventoryAssetId)
		if err != nil {
			utils.LogResponseError(echoContext, handlers.logger, err)
			return echoContext.JSON(httpErrors.ErrorResponse(err))
		}

		return echoContext.JSON(http.StatusOK, inventoryAssetById)
	}
}

// GetInventoryAssets godoc
// @Summary Get all inventory assets
// @Description Get all invntory assets with pagination
// @Tags Inventory Assets
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} models.InventoryAssetList
// @Router /api/v1/inventory-assets [get]
func (handlers inventoryAssetHandlers) GetInventoryAssets() echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(echoContext), "inventoryAssetHandlers.GetInventoryAssets")
		defer span.Finish()

		pq, err := utils.GetPaginationFromCtx(echoContext)
		if err != nil {
			utils.LogResponseError(echoContext, handlers.logger, err)
			return echoContext.JSON(httpErrors.ErrorResponse(err))
		}

		inventoryAssetsList, err := handlers.inventoryassetUC.GetInventoryAssets(ctx, pq)
		if err != nil {
			utils.LogResponseError(echoContext, handlers.logger, err)
			return echoContext.JSON(httpErrors.ErrorResponse(err))
		}

		return echoContext.JSON(http.StatusOK, inventoryAssetsList)
	}
}
