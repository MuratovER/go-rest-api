package usecase

import (
	"context"
	"inventory_assets/config"
	"inventory_assets/internal/inventory_asset"
	"inventory_assets/internal/models"
	"inventory_assets/pkg/httpErrors"
	"inventory_assets/pkg/logger"
	"inventory_assets/pkg/utils"

	"github.com/pkg/errors"
)

// inventoryAssetUC represents the inventory asset use case
type inventoryAssetUC struct {
	cfg                *config.Config
	inventoryAssetRepo inventory_asset.Repository
	logger             logger.Logger
}

// NewInventoryAssetUseCase creates a new inventory asset use case instance
func NewInventoryAssetUseCase(cfg *config.Config, inventoryAssetRepo inventory_asset.Repository, logger logger.Logger) *inventoryAssetUC {
	return &inventoryAssetUC{
		cfg:                cfg,
		inventoryAssetRepo: inventoryAssetRepo,
		logger:             logger,
	}
}

// GetInventoryAssetsById retrieves inventory assets by ID with pagination
func (u *inventoryAssetUC) GetInventoryAssetsById(ctx context.Context, pq *utils.PaginationQuery, id int) (*models.InventoryAssetList, error) {
	u.logger.Infof("Getting inventory asset by ID: %d", id)

	result, err := u.inventoryAssetRepo.GetInventoryAssetsById(ctx, pq, id)
	if err != nil {
		u.logger.Errorf("Failed to get inventory asset by ID %d: %v", id, err)
		return nil, httpErrors.NewInventoryAssetNotFoundError(errors.WithMessage(err, "inventory asset with this id not found"))
	}

	u.logger.Infof("Successfully retrieved inventory asset by ID: %d", id)
	return result, nil
}

// GetInventoryAssets retrieves all inventory assets with pagination
func (u *inventoryAssetUC) GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error) {
	u.logger.Infof("Getting inventory assets with pagination: page=%d, size=%d", pq.GetPage(), pq.GetSize())

	result, err := u.inventoryAssetRepo.GetInventoryAssets(ctx, pq)
	if err != nil {
		u.logger.Errorf("Failed to get inventory assets: %v", err)
		return nil, errors.Wrap(err, "failed to get inventory assets")
	}

	u.logger.Infof("Successfully retrieved %d inventory assets", len(*result.Items))
	return result, nil
}
