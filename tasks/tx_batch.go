package tasks

import (
	"cosmossdk.io/store/multiversion"
	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DeliverTxEntry represents an individual transaction's request within a batch.
// This can be extended to include tx-level tracing or metadata
type DeliverTxEntry struct {
	Request            abci.RequestFinalizeBlock
	SdkTx              sdk.Tx
	Checksum           [32]byte
	AbsoluteIndex      int
	EstimatedWritesets MappedWritesets
}

// EstimatedWritesets represents an estimated writeset for a transaction mapped by storekey to the writeset estimate.
type MappedWritesets map[storetypes.StoreKey]multiversion.WriteSet

// DeliverTxBatchRequest represents a request object for a batch of transactions.
// This can be extended to include request-level tracing or metadata
type DeliverTxBatchRequest struct {
	TxEntries []*DeliverTxEntry
}

// DeliverTxResult represents an individual transaction's response within a batch.
// This can be extended to include tx-level tracing or metadata
type DeliverTxResult struct {
	Response abci.ResponseFinalizeBlock
}

// DeliverTxBatchResponse represents a response object for a batch of transactions.
// This can be extended to include response-level tracing or metadata
type DeliverTxBatchResponse struct {
	Results []*DeliverTxResult
}
