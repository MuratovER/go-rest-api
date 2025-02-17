//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package inventory_asset

import (
	"context"

	"inventory_assets/internal/models"
	"inventory_assets/pkg/utils"
)

// Inventory Asset Repository
type Repository interface {
	GetInventoryAssetsById(ctx context.Context, pq *utils.PaginationQuery, Id int) (*models.InventoryAssetList, error)
	GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error)
}
