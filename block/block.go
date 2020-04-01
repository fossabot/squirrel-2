package block

import "squirrel/rpc"

// Block db model.
type Block struct {
	ID                 uint
	Hash               string
	Size               int
	Version            uint
	PreviousBlockHash  string
	MerkleRoot         string
	Time               uint64
	Index              uint
	Nonce              string
	NextConsensus      string
	ScriptInvocation   string
	ScriptVerification string
	// Txs moited
	NextBlockhash string
}

// ParseBlocks parses struct RawBlock to struct Block.
func ParseBlocks(rawBlocks []*rpc.RawBlock) []*Block {
	blocks := []*Block{}

	for _, rawBlock := range rawBlocks {
		block := Block{
			Hash:               rawBlock.Hash,
			Size:               rawBlock.Size,
			Version:            rawBlock.Version,
			PreviousBlockHash:  rawBlock.PreviousBlockHash,
			MerkleRoot:         rawBlock.MerkleRoot,
			Time:               rawBlock.Time,
			Index:              rawBlock.Index,
			Nonce:              rawBlock.Nonce,
			NextConsensus:      rawBlock.NextConsensus,
			ScriptInvocation:   rawBlock.Script.Invocation,
			ScriptVerification: rawBlock.Script.Verification,
			NextBlockhash:      rawBlock.NextBlockHash,
		}
		blocks = append(blocks, &block)
	}
	return blocks
}
