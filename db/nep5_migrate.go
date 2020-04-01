package db

import (
	"database/sql"
	"squirrel/cache"
)

// HandleNEP5Migrate handles nep5 contract migration.
func HandleNEP5Migrate(newAssetAdmin, oldAssetID, newAssetID string, txPK uint, txID string) error {
	return transact(func(tx *sql.Tx) error {
		query := "UPDATE `nep5` SET `visible` = FALSE WHERE `asset_id` = ? LIMIT 1"
		if _, err := tx.Exec(query, oldAssetID); err != nil {
			return err
		}

		query = "DELETE FROM `addr_asset` WHERE `asset_id` = ? AND `address` IN ("
		query += "SELECT `address` FROM (SELECT `address` FROM `addr_asset` WHERE asset_id=? AND `address` IN ("
		query += "SELECT `address` FROM `addr_asset` WHERE `asset_id` IN (?, ?) GROUP BY `address` HAVING COUNT(`asset_id`) = 2))a)"
		if _, err := tx.Exec(query, newAssetID, newAssetID, oldAssetID, newAssetID); err != nil {
			return err
		}

		query = "UPDATE `addr_asset` SET `asset_id` = ? WHERE `asset_id` = ?"
		if _, err := tx.Exec(query, newAssetID, oldAssetID); err != nil {
			return err
		}

		addrs, holdingAddrs := cache.MigrateNEP5(newAssetAdmin, oldAssetID, newAssetID)
		query = "UPDATE `nep5` SET `addresses` = ?, `holding_addresses` = ? WHERE `asset_id` = ? LIMIT 1"
		if _, err := tx.Exec(query, addrs, holdingAddrs, newAssetID); err != nil {
			return err
		}

		query = "INSERT INTO `nep5_migrate`(`old_asset_id`, `new_asset_id`, `migrate_txid`) VALUES (?, ?, ?)"
		if _, err := tx.Exec(query, oldAssetID, newAssetID, txID); err != nil {
			return err
		}

		err := updateNep5Counter(tx, txPK, -1)
		return err
	})
}
