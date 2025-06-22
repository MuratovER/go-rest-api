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

// InventoryAssetRepo represents the inventory asset repository
type inventoryAssetRepo struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewInventoryAssetRepository creates a new inventory asset repository instance
func NewInventoryAssetRepository(db *gorm.DB, logger logger.Logger) inventory_asset.Repository {
	return &inventoryAssetRepo{db: db, logger: logger}
}

// GetInventoryAssetsById retrieves inventory assets by ID with pagination
func (r *inventoryAssetRepo) GetInventoryAssetsById(ctx context.Context, pq *utils.PaginationQuery, id int) (*models.InventoryAssetList, error) {
	var totalCount int64
	var inventoryAssets []models.InventoryAsset

	// Count total records
	if err := r.db.Model(&models.InventoryAsset{}).Count(&totalCount).Error; err != nil {
		return nil, errors.Wrap(err, "inventoryAssetRepo.GetInventoryAssetsById.Count")
	}

	// Find specific record by ID
	if err := r.db.Where("id = ?", id).Find(&inventoryAssets).Error; err != nil {
		return nil, errors.Wrap(err, "inventoryAssetRepo.GetInventoryAssetsById.Find")
	}

	if len(inventoryAssets) == 0 {
		r.logger.Error("Inventory asset with this id not found")
		return nil, errors.New("inventory asset with this id not found")
	}

	return &models.InventoryAssetList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		Items:      &inventoryAssets,
	}, nil
}

// GetInventoryAssets retrieves all inventory assets with pagination
func (r *inventoryAssetRepo) GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error) {
	var totalCount int64
	var inventoryAssets []models.InventoryAsset

	// Count total records
	if err := r.db.Model(&models.InventoryAsset{}).Count(&totalCount).Error; err != nil {
		return nil, errors.Wrap(err, "inventoryAssetRepo.GetInventoryAssets.Count")
	}

	// Get paginated records
	if err := r.db.Offset(pq.GetOffset()).Limit(pq.GetLimit()).Find(&inventoryAssets).Error; err != nil {
		return nil, errors.Wrap(err, "inventoryAssetRepo.GetInventoryAssets.Find")
	}

	return &models.InventoryAssetList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		Items:      &inventoryAssets,
	}, nil
}
