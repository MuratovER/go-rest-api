//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package inventory_asset

import (
	"context"
	"inventory_assets/internal/models"
	"inventory_assets/pkg/utils"
)

// Inventory Asset use case
type UseCase interface {
	GetInventoryAssetsById(ctx context.Context, pq *utils.PaginationQuery, Id int) (*models.InventoryAssetList, error)
	GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error)
}
