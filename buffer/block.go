package buffer

import (
	"squirrel/rpc"
	"sync"
)

// BlockBuffer is used to temporarily buffer fetched blocks.
type BlockBuffer struct {
	mu sync.Mutex
	// maxHeight indicates the highest existing height.
	maxHeight int
	// nextHeight indicates the next block height to fetch,
	// used before blockchain fully synchronized.
	nextHeight int
	buffer     map[int]*rpc.RawBlock
}

// NewBuffer inits a new block buffer.
func NewBuffer(height int) BlockBuffer {
	return BlockBuffer{
		maxHeight:  height,
		nextHeight: height,
		buffer:     make(map[int]*rpc.RawBlock),
	}
}

// Pop pops the specific block by id.
func (b *BlockBuffer) Pop(index int) (*rpc.RawBlock, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if block, ok := b.buffer[index]; ok {
		delete(b.buffer, index)
		return block, true
	}
	return nil, false
}

// GetHighest returns the highest existing block height.
func (b *BlockBuffer) GetHighest() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.maxHeight
}

// GetNextPending returns the next fetching block index.
func (b *BlockBuffer) GetNextPending() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.nextHeight++
	return b.nextHeight
}

// Put adds the given block into buffer and update maxHeight.
func (b *BlockBuffer) Put(block *rpc.RawBlock) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buffer[int(block.Index)] = block
	if b.maxHeight < int(block.Index) {
		b.maxHeight = int(block.Index)
	}
}

// Size returns size of current buffer.
func (b *BlockBuffer) Size() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return len(b.buffer)
}
