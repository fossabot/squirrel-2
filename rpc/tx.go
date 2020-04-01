package rpc

import "math/big"

// RawTx is the transaction part of block data.
type RawTx struct {
	TxID       string `json:"txid"`
	Size       uint   `json:"size"`
	Type       string `json:"type"`
	Version    uint   `json:"version"`
	Attributes []struct {
		Usage string `json:"usage"`
		Data  string `json:"data"`
	}
	Vin []struct {
		TxID string `json:"txid"`
		Vout uint16 `json:"vout"`
	}

	Vout []struct {
		N       uint16     `json:"n"`
		Asset   string     `json:"asset"`
		Value   *big.Float `json:"value"`
		Address string     `json:"address"`
	}
	SysFee  *big.Float `json:"sys_fee"`
	NetFee  *big.Float `json:"net_fee"`
	Scripts []struct {
		Invocation   string `json:"invocation"`
		Verification string `json:"verification"`
	}
	Asset struct {
		Type string `json:"type"`
		Name []struct {
			Lang string `json:"lang"`
			Name string `json:"name"`
		}
		Amount    *big.Float `json:"amount"`
		Precision uint8      `json:"precision"`
		Owner     string     `json:"owner"`
		Admin     string     `json:"admin"`
	}
	Claims []struct {
		TxID string `json:"txid"`
		Vout uint16 `json:"vout"`
	}
	Script string     `json:"script"`
	Nonce  int64      `json:"nonce"`
	Gas    *big.Float `json:"gas"`
}
