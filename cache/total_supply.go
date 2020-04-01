package cache

import (
	"math/big"
	"sync"
)

// AssetTotalSupplyCacheItem caches assetID with its total supply.
type AssetTotalSupplyCacheItem struct {
	TotalSupply *big.Float
	BlockIndex  uint
}

var (
	totalSupplyCache = make(map[string]*AssetTotalSupplyCacheItem)
	assetCacheLock   sync.Mutex
)

// GetAssetTotalSupply returns assetID total supply cache record.
func GetAssetTotalSupply(assetID string) (*big.Float, uint, bool) {
	assetCacheLock.Lock()
	defer assetCacheLock.Unlock()

	rec, ok := totalSupplyCache[assetID]
	if !ok {
		return nil, 0, false
	}

	return rec.TotalSupply, rec.BlockIndex, true
}

// UpdateAssetTotalSupply updates or sets total supply for assetID.
func UpdateAssetTotalSupply(assetID string, totalSupply *big.Float, blockIndex uint) bool {
	assetCacheLock.Lock()
	defer assetCacheLock.Unlock()

	if rec, ok := totalSupplyCache[assetID]; ok {
		if rec.BlockIndex > blockIndex {
			return false
		}
	}

	totalSupplyCache[assetID] = &AssetTotalSupplyCacheItem{
		TotalSupply: totalSupply,
		BlockIndex:  blockIndex,
	}

	return true
}
