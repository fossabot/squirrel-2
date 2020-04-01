package tasks

import (
	"os"
	"squirrel/config"
	"squirrel/log"
	"squirrel/rpc"
	"squirrel/util"
	"testing"
)

func TestQueryNep5AssetBalance(t *testing.T) {
	log.Init()
	defer func() {
		os.Remove("error.log")
	}()

	config.Load(false)

	rpc.RefreshServers()

	assetID := "af7c7328eee5a275a3bcaee2bf0cf662b5e739be"
	scriptHash := util.GetScriptHashFromAssetID(assetID)

	addrBytesList := [][]byte{
		util.GetScriptHashFromAddress("AKQjaQ7Hor11BfRnXUBvYYiY1CwUkLywyc"),
		util.GetScriptHashFromAddress(""),
	}

	nep5AssetDecimals = map[string]uint8{
		"af7c7328eee5a275a3bcaee2bf0cf662b5e739be": 8,
	}

	_, ok := queryBalances(99999999, scriptHash, assetID, addrBytesList)
	if !ok {
		t.Errorf("Failed to get nep5 balance from rpc")
		return
	}
}
