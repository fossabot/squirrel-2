package tx

import (
	"fmt"
	"math/big"
	"squirrel/asset"
	"squirrel/rpc"
	"squirrel/smartcontract"
	"strings"
)

const (
	// RegisterTransaction represents registration transaction.
	RegisterTransaction = iota
	// MinerTransaction represents miner transaction.
	MinerTransaction
	// IssueTransaction represents issue transaction.
	IssueTransaction
	// InvocationTransaction represents invocation transaction.
	InvocationTransaction
	// ContractTransaction represents contract transaction.
	ContractTransaction
	// ClaimTransaction represents claim transaction.
	ClaimTransaction
	// PublishTransaction represents publish transaction.
	PublishTransaction
	// EnrollmentTransaction represents enrollment transaction.
	EnrollmentTransaction
)

// Bulk stores innner content of parsed raw block data.
type Bulk struct {
	TXs       []*Transaction
	TXAttrs   []*TransactionAttribute
	TXVins    []*TransactionVin
	TXVouts   []*TransactionVout
	TXScripts []*TransactionScripts
	Assets    []*asset.Asset
	Claims    []*TransactionClaims
}

// Transaction db model.
type Transaction struct {
	ID         uint
	BlockIndex uint
	BlockTime  uint64
	TxID       string
	Size       uint
	Type       string
	Version    uint
	// Attribute List
	// Vin List
	// Vout List
	SysFee *big.Float
	NetFee *big.Float
	// Scripts
	Nonce  int64
	Script string
	Gas    *big.Float
}

// TransactionAttribute of transactions.
type TransactionAttribute struct {
	ID    uint
	TxID  string
	Usage string
	Data  string
}

// TransactionVin of transacitons.
type TransactionVin struct {
	// ID   uint
	From string
	TxID string
	Vout uint16
}

// TransactionVout of transaction.
type TransactionVout struct {
	// ID      uint
	TxID    string
	N       uint16
	AssetID string
	Value   *big.Float
	Address string
}

// TransactionScripts of transaction.
type TransactionScripts struct {
	ID           uint
	TxID         string
	Invocation   string
	Verification string
}

// TransactionClaims of transaction.
type TransactionClaims struct {
	ID   uint
	TxID string
	Vout uint16
}

// AddrAssetIDTx is the bundle of address, asset_id and txid.
type AddrAssetIDTx struct {
	Address string
	AssetID string
	TxID    string
}

func typeName(txType int) string {
	switch txType {
	case RegisterTransaction:
		return "RegisterTransaction"
	case MinerTransaction:
		return "MinerTransaction"
	case IssueTransaction:
		return "IssueTransaction"
	case InvocationTransaction:
		return "InvocationTransaction"
	case ContractTransaction:
		return "ContractTransaction"
	case ClaimTransaction:
		return "ClaimTransaction"
	case PublishTransaction:
		return "PublishTransaction"
	case EnrollmentTransaction:
		return "EnrollmentTransaction"
	default:
		err := fmt.Errorf("unknown transaction type: %d", txType)
		panic(err)
	}
}

// ParseTxs parses all raw transactions in raw blocks to Bulk.
func ParseTxs(rawBlocks []*rpc.RawBlock) *Bulk {
	txs := Bulk{}

	for _, rawBlock := range rawBlocks {
		for _, rawTx := range rawBlock.Tx {
			txs.TXs = appendTx(txs.TXs, rawBlock.Index, rawBlock.Time, &rawTx)
			txs.TXAttrs = appendTxAttrs(txs.TXAttrs, &rawTx)
			txs.TXVins = appendTxVin(txs.TXVins, &rawTx)
			txs.TXVouts = appendTxVout(txs.TXVouts, &rawTx)
			txs.TXScripts = appendTxScripts(txs.TXScripts, &rawTx)
			txs.Assets = appendAsset(rawBlock, txs.Assets, &rawTx)
			txs.Claims = appendClaims(txs.Claims, &rawTx)
		}
	}

	return &txs
}

func appendTx(txs []*Transaction, blockIndex uint, blockTime uint64, rawTx *rpc.RawTx) []*Transaction {
	trans := Transaction{
		BlockIndex: blockIndex,
		BlockTime:  blockTime,
		TxID:       rawTx.TxID,
		Size:       rawTx.Size,
		Type:       rawTx.Type,
		Version:    rawTx.Version,
		SysFee:     rawTx.SysFee,
		NetFee:     rawTx.NetFee,
		Nonce:      rawTx.Nonce,
		Script:     rawTx.Script,
		Gas:        rawTx.Gas,
	}
	if rawTx.Gas == nil {
		trans.Gas = big.NewFloat(0)
	}
	txs = append(txs, &trans)

	return txs
}

func appendTxAttrs(txAttrs []*TransactionAttribute, rawTx *rpc.RawTx) []*TransactionAttribute {
	for _, rawAttr := range rawTx.Attributes {
		attr := TransactionAttribute{
			TxID:  rawTx.TxID,
			Usage: rawAttr.Usage,
			Data:  rawAttr.Data,
		}
		txAttrs = append(txAttrs, &attr)
	}
	return txAttrs
}

func appendTxVin(txVin []*TransactionVin, rawTx *rpc.RawTx) []*TransactionVin {
	for _, rawVin := range rawTx.Vin {
		vin := TransactionVin{
			From: rawTx.TxID,
			TxID: rawVin.TxID,
			Vout: rawVin.Vout,
		}
		txVin = append(txVin, &vin)
	}
	return txVin
}

func appendTxVout(txVout []*TransactionVout, rawTx *rpc.RawTx) []*TransactionVout {
	for _, rawVout := range rawTx.Vout {
		vout := TransactionVout{
			TxID:    rawTx.TxID,
			N:       rawVout.N,
			AssetID: rawVout.Asset,
			Value:   rawVout.Value,
			Address: rawVout.Address,
		}
		txVout = append(txVout, &vout)
	}
	return txVout
}

func appendTxScripts(txScripts []*TransactionScripts, rawTx *rpc.RawTx) []*TransactionScripts {
	for _, rawScript := range rawTx.Scripts {
		script := TransactionScripts{
			TxID:         rawTx.TxID,
			Invocation:   rawScript.Invocation,
			Verification: rawScript.Verification,
		}
		txScripts = append(txScripts, &script)
	}
	return txScripts
}

func appendAsset(rawBlock *rpc.RawBlock, assets []*asset.Asset, rawTx *rpc.RawTx) []*asset.Asset {
	var asset *asset.Asset
	if rawTx.Type == typeName(RegisterTransaction) {
		asset = parseAssetFromRegisterTransaction(rawBlock.Index, rawTx)
	} else if rawTx.Type == typeName(InvocationTransaction) {
		// Example: 0x4a629db0af0d9c7ee0e11f4f4894765f5ab2579bcc8b4a203e4c6814a9784f00(testnet).
		if strings.HasSuffix(rawTx.Script, smartcontract.AssetFingerPrint) {
			asset = parseAssetFromInvocationTransaction(rawTx.Script)
			if asset == nil {
				return assets
			}

			// Supplement the rest fields.
			asset.Version = 0
			asset.AssetID = rawTx.TxID
			asset.Expiration = uint64(rawBlock.Index + 2000000)
		}
	}

	if asset == nil {
		return assets
	}

	asset.BlockIndex = rawBlock.Index
	asset.BlockTime = rawBlock.Time
	asset.Addresses = 0
	asset.Transactions = 0

	assets = append(assets, asset)
	return assets
}

func appendClaims(claims []*TransactionClaims, rawTx *rpc.RawTx) []*TransactionClaims {
	for _, rawClaim := range rawTx.Claims {
		claim := TransactionClaims{
			TxID: rawClaim.TxID,
			Vout: rawClaim.Vout,
		}
		claims = append(claims, &claim)
	}
	return claims
}

func parseAssetFromRegisterTransaction(blockIndex uint, rawTx *rpc.RawTx) *asset.Asset {
	newAsset := asset.Asset{
		Version:    rawTx.Version,
		AssetID:    rawTx.TxID,
		Type:       rawTx.Asset.Type,
		Name:       rawTx.Asset.Name[0].Name,
		Amount:     rawTx.Asset.Amount,
		Available:  big.NewFloat(0),
		Precision:  rawTx.Asset.Precision,
		Owner:      rawTx.Asset.Owner,
		Admin:      rawTx.Asset.Admin,
		Issuer:     rawTx.Asset.Owner,
		Expiration: uint64(blockIndex + 2*2000000),
		Frozen:     false,
	}

	if newAsset.AssetID == asset.NEOAssetID {
		newAsset.Name = asset.NEO
	} else if newAsset.AssetID == asset.GASAssetID {
		newAsset.Name = asset.GAS
	}

	return &newAsset
}

func parseAssetFromInvocationTransaction(sc string) *asset.Asset {
	asset := smartcontract.GetAssetInfo(sc)
	return asset
}
