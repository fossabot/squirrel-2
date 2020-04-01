package rpc

// BlockCountRespponse returns block height of chain.
type BlockCountRespponse struct {
	jsonRPCResponse
	Result int `json:"result"`
}

// BlockResponse returns full block data of a specific index.
type BlockResponse struct {
	jsonRPCResponse
	Result *RawBlock `json:"result"`
}

// RawBlock is the raw block structure used in rpc response.
type RawBlock struct {
	Hash              string `json:"hash"`
	Size              int    `json:"size"`
	Version           uint   `json:"version"`
	PreviousBlockHash string `json:"previousblockhash"`
	MerkleRoot        string `json:"merkleroot"`
	Time              uint64 `json:"time"`
	Index             uint   `json:"index"`
	Nonce             string `json:"nonce"`
	NextConsensus     string `json:"nextconsensus"`
	Script            struct {
		Invocation   string `json:"invocation"`
		Verification string `json:"verification"`
	}
	Tx            []RawTx
	NextBlockHash string `json:"nextblockhash"`
}

// DownloadBlock from rpc server.
func DownloadBlock(index int) *RawBlock {
	params := []interface{}{index, 1}
	args := getRPCRequestBody("getblock", params)

	respData := BlockResponse{}
	rpcCall(index, args, &respData)

	return respData.Result
}
