package repository

import (
	"context"
	"inventory_assets/internal/inventory_asset"
	"inventory_assets/internal/models"
	"inventory_assets/pkg/utils"

	"inventory_assets/pkg/logger"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Inventory Asset Repository
type inventoryAssetRepo struct {
	db     *gorm.DB
	logger logger.Logger
}

// Inventory Asser Repository constructor
func NewInventoryAssetRepository(db *gorm.DB, logger logger.Logger) inventory_asset.Repository {
	return &inventoryAssetRepo{db: db, logger: logger}
}

// Get inventory assets by id
func (r *inventoryAssetRepo) GetInventoryAssetsById(ctx context.Context, pq *utils.PaginationQuery, Id int) (*models.InventoryAssetList, error) {
	var totalCount int64
	var inventory_assets []models.InventoryAsset

	if err := r.db.Find(&inventory_assets).Count(&totalCount); err.Error != nil {
		return nil, errors.Wrap(err.Error, "inventoryAssetRepo.GetInventoryAssets.GetContext.totalCount")
	}

	r.db.Offset(pq.GetOffset()).Limit(pq.GetLimit()).Find(&inventory_assets, Id)
	if len(inventory_assets) == 0 {
		r.logger.Error("Inventoy asset with this uid not found.")
		return nil, errors.New("Inventory_asset with this uid not found")
	}

	return &models.InventoryAssetList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		Items:      &inventory_assets,
	}, nil
}

// Get inventory assets
func (r *inventoryAssetRepo) GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error) {
	var totalCount int64
	var inventory_assets []models.InventoryAsset

	if err := r.db.Find(&inventory_assets).Count(&totalCount); err.Error != nil {
		return nil, errors.Wrap(err.Error, "inventoryAssetRepo.GetInventoryAssets.GetContext.totalCount")
	}

	err := r.db.Offset(pq.GetOffset()).Limit(pq.GetLimit()).Find(&inventory_assets)
	if err.Error != nil {
		return nil, errors.Wrap(err.Error, "inventoryAssetRepo.GetInventoryAssets.QueryxContext")
	}

	return &models.InventoryAssetList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		Items:      &inventory_assets,
	}, nil
}
