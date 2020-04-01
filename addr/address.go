package addr

import (
	"math/big"
)

// Address db model.
type Address struct {
	ID                  uint
	Address             string
	CreatedAt           uint64
	LastTransactionTime uint64
	TransAsset          uint64
	TransNep5           uint64
}

// Asset db model.
type Asset struct {
	ID                  uint
	Address             string
	AssetID             string
	Balance             *big.Float
	Transactions        uint64
	LastTransactionTime uint64
}

// AssetInfo model.
type AssetInfo struct {
	Address             string
	CreatedAt           uint64
	LastTransactionTime uint64
	AssetID             string
	Balance             *big.Float
}

// Tx model.
type Tx struct {
	ID        uint
	TxID      string
	Address   string
	BlockTime uint64
	AssetType string
}
