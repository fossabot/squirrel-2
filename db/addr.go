package db

import (
	"database/sql"
	"fmt"
	"squirrel/addr"
	"squirrel/asset"
	"squirrel/cache"
	"squirrel/log"
	"squirrel/util"
)

// GetAddrAssetInfo returns all addresses with it's assets.
func GetAddrAssetInfo() []*addr.AssetInfo {
	const query = "SELECT `address`.`address`, `address`.`created_at`, `address`.`last_transaction_time`, `addr_asset`.`asset_id`, `addr_asset`.`balance` FROM `addr_asset` LEFT JOIN `address` ON `address`.`address`=`addr_asset`.`address`"

	result := []*addr.AssetInfo{}

	rows, err := wrappedQuery(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		m := &addr.AssetInfo{}
		var balanceStr string

		err := rows.Scan(
			&m.Address,
			&m.CreatedAt,
			&m.LastTransactionTime,
			&m.AssetID,
			&balanceStr,
		)

		if err != nil {
			panic(err)
		}

		m.Balance = util.StrToBigFloat(balanceStr)

		result = append(result, m)
	}

	return result
}

func updateAddrInfo(tx *sql.Tx, blockTime uint64, txID string, addr string, assetType string) error {
	var incrAsset, incrNep5 = 0, 0
	switch assetType {
	case asset.ASSET:
		incrAsset = 1
	case asset.NEP5:
		incrNep5 = 1
	default:
		panic("Unsupported asset Type: " + assetType)
	}

	addrCache, created := cache.GetAddrOrCreate(addr, blockTime)

	if created {
		const createAddrQuery = "INSERT INTO `address` (`address`, `created_at`, `last_transaction_time`, `trans_asset`, `trans_nep5`) VALUES (?, ?, ?, ?, ?)"
		_, err := tx.Exec(createAddrQuery, addr, blockTime, blockTime, incrAsset, incrNep5)
		if err != nil {
			log.Error.Printf("TxID: %s, addr=%s, assetType=%s\n", txID, addr, assetType)
			return err
		}
	} else {
		query := fmt.Sprintf("UPDATE `address` SET `trans_asset` = `trans_asset` + %d, `trans_nep5` = `trans_nep5` + %d", incrAsset, incrNep5)
		// Because task tx and task nep5 run in parallel,
		// maybe one task executes before the other one with a bigger blockTime.
		if addrCache.UpdateCreatedTime(blockTime) {
			query += fmt.Sprintf(", `created_at` = %d", blockTime)
		}
		if addrCache.UpdateLastTxTime(blockTime) {
			query += fmt.Sprintf(", `last_transaction_time` = %d", blockTime)
		}
		query += fmt.Sprintf(" WHERE `address` = '%s' LIMIT 1", addr)

		_, err := tx.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func createAddrInfoIfNotExist(tx *sql.Tx, blockTime uint64, addr string) error {
	_, created := cache.GetAddrOrCreate(addr, blockTime)
	if created {
		const createAddrQuery = "INSERT INTO `address` (`address`, `created_at`, `last_transaction_time`, `trans_asset`, `trans_nep5`) VALUES (?, ?, ?, ?, ?)"
		_, err := tx.Exec(createAddrQuery, addr, blockTime, blockTime, 0, 0)
		if err != nil {
			log.Error.Printf("addr=%s\n", addr)
			return err
		}
	}

	return nil
}
