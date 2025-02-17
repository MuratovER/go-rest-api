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

// Inventory Asset UseCase
type inventoryAssetUC struct {
	cfg                *config.Config
	inventoryAssetRepo inventory_asset.Repository
	logger             logger.Logger
}

// Inventory Asset UseCase constructor
func NewInventoryAssetUseCase(cfg *config.Config, inventoryAssetRepo inventory_asset.Repository, logger logger.Logger) *inventoryAssetUC {
	return &inventoryAssetUC{cfg: cfg, inventoryAssetRepo: inventoryAssetRepo, logger: logger}
}

// Get inventory assets by id
func (u *inventoryAssetUC) GetInventoryAssetsById(ctx context.Context, pq *utils.PaginationQuery, Id int) (*models.InventoryAssetList, error) {
	n, err := u.inventoryAssetRepo.GetInventoryAssetsById(ctx, pq, Id)
	if err != nil {

		return nil, httpErrors.NewInventoryAssetNotFoundError(errors.WithMessage(err, "Inventoy asset with this uid not found."))
	}

	return n, nil
}

// Get inventory assets
func (u *inventoryAssetUC) GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error) {

	return u.inventoryAssetRepo.GetInventoryAssets(ctx, pq)

}
