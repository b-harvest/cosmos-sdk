package tasks

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"testing"

	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"cosmossdk.io/store/cachekv"
	"cosmossdk.io/store/cachemulti"
	"cosmossdk.io/store/dbadapter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/utils/tracing"
)

type mockDeliverTxFunc func(ctx sdk.Context, req abci.RequestFinalizeBlock, tx sdk.Tx, checksum [32]byte) (res abci.ResponseFinalizeBlock)

var testStoreKey = storetypes.NewKVStoreKey("mock")
var itemKey = []byte("key")

func requestList(n int) []*DeliverTxEntry {
	tasks := make([]*DeliverTxEntry, n)
	for i := 0; i < n; i++ {
		tasks[i] = &DeliverTxEntry{
			Request: abci.RequestFinalizeBlock{
				Tx: []byte(fmt.Sprintf("%d", i)),
			},
			AbsoluteIndex: i,
			// TODO: maybe we need to add dummy sdkTx message types and handler routers too
		}

	}
	return tasks
}

func initTestCtx(injectStores bool) sdk.Context {
	ctx := sdk.Context{}.WithContext(context.Background())
	keys := make(map[string]storetypes.StoreKey)
	stores := make(map[storetypes.StoreKey]storetypes.CacheWrapper)
	db := dbm.NewMemDB()
	if injectStores {
		mem := dbadapter.Store{DB: db}
		stores[testStoreKey] = cachekv.NewStore(mem)
		keys[testStoreKey.Name()] = testStoreKey
	}
	store := cachemulti.NewStore(db, stores, keys, nil, nil, nil)
	ctx = ctx.WithMultiStore(&store)
	return ctx
}

func generateTasks(count int) []*deliverTxTask {
	var res []*deliverTxTask
	for i := 0; i < count; i++ {
		res = append(res, &deliverTxTask{Index: i})
	}
	return res
}

func TestProcessAll(t *testing.T) {
	runtime.SetBlockProfileRate(1)

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	tests := []struct {
		name          string
		workers       int
		runs          int
		before        func(ctx sdk.Context)
		requests      []*DeliverTxEntry
		deliverTxFunc mockDeliverTxFunc
		addStores     bool
		expectedErr   error
		assertions    func(t *testing.T, ctx sdk.Context, res []abci.ResponseFinalizeBlock)
	}{
		{
			name:      "Test zero txs does not hang",
			workers:   20,
			runs:      10,
			addStores: true,
			requests:  requestList(0),
			deliverTxFunc: func(ctx sdk.Context, req abci.RequestFinalizeBlock, tx sdk.Tx, checksum [32]byte) (res abci.ResponseFinalizeBlock) {
				panic("should not deliver")
			},
			assertions: func(t *testing.T, ctx sdk.Context, res []abci.ResponseFinalizeBlock) {
				require.Len(t, res, 0)
			},
			expectedErr: nil,
		},
		{
			name:      "Test tx writing to a store that another tx is iterating",
			workers:   50,
			runs:      1,
			requests:  requestList(500),
			addStores: true,
			before: func(ctx sdk.Context) {
				kv := ctx.MultiStore().GetKVStore(testStoreKey)
				// initialize 100 test values in the base kv store so iterating isn't too fast
				for i := 0; i < 10; i++ {
					kv.Set([]byte(fmt.Sprintf("%d", i)), []byte(fmt.Sprintf("%d", i)))
				}
			},
			deliverTxFunc: func(ctx sdk.Context, req abci.RequestFinalizeBlock, tx sdk.Tx, checksum [32]byte) (res abci.ResponseFinalizeBlock) {
				kv := ctx.MultiStore().GetKVStore(testStoreKey)
				if ctx.TxIndex()%2 == 0 {
					// For even-indexed transactions, write to the store
					kv.Set(req.Tx, req.Tx)
					return abci.ResponseFinalizeBlock{
						Info: "write",
					}
				} else {
					// For odd-indexed transactions, iterate over the store

					// just write so we have more writes going on
					kv.Set(req.Tx, req.Tx)
					iterator := kv.Iterator(nil, nil)
					defer iterator.Close()
					for ; iterator.Valid(); iterator.Next() {
						// Do nothing, just iterate
					}
					return abci.ResponseFinalizeBlock{
						Info: "iterate",
					}
				}
			},
			assertions: func(t *testing.T, ctx sdk.Context, res []abci.ResponseFinalizeBlock) {
				for idx, response := range res {
					if idx%2 == 0 {
						require.Equal(t, "write", response.Info)
					} else {
						require.Equal(t, "iterate", response.Info)
					}
				}
			},
			expectedErr: nil,
		},
		{
			name:      "Test no overlap txs",
			workers:   20,
			runs:      10,
			addStores: true,
			requests:  requestList(1000),
			deliverTxFunc: func(ctx sdk.Context, req abci.RequestFinalizeBlock, tx sdk.Tx, checksum [32]byte) (res abci.ResponseFinalizeBlock) {
				// all txs read and write to the same key to maximize conflicts
				kv := ctx.MultiStore().GetKVStore(testStoreKey)

				// write to the store with this tx's index
				kv.Set(req.Tx, req.Tx)
				val := string(kv.Get(req.Tx))

				// return what was read from the store (final attempt should be index-1)
				return abci.ResponseFinalizeBlock{
					Info: val,
				}
			},
			assertions: func(t *testing.T, ctx sdk.Context, res []abci.ResponseFinalizeBlock) {
				for idx, response := range res {
					require.Equal(t, fmt.Sprintf("%d", idx), response.Info)
				}
				store := ctx.MultiStore().GetKVStore(testStoreKey)
				for i := 0; i < len(res); i++ {
					val := store.Get([]byte(fmt.Sprintf("%d", i)))
					require.Equal(t, []byte(fmt.Sprintf("%d", i)), val)
				}
			},
			expectedErr: nil,
		},
		{
			name:      "Test every tx accesses same key",
			workers:   50,
			runs:      1,
			addStores: true,
			requests:  requestList(1000),
			deliverTxFunc: func(ctx sdk.Context, req abci.RequestFinalizeBlock, tx sdk.Tx, checksum [32]byte) (res abci.ResponseFinalizeBlock) {
				// all txs read and write to the same key to maximize conflicts
				kv := ctx.MultiStore().GetKVStore(testStoreKey)
				val := string(kv.Get(itemKey))

				// write to the store with this tx's index
				kv.Set(itemKey, req.Tx)

				// return what was read from the store (final attempt should be index-1)
				return abci.ResponseFinalizeBlock{
					Info: val,
				}
			},
			assertions: func(t *testing.T, ctx sdk.Context, res []abci.ResponseFinalizeBlock) {
				for idx, response := range res {
					if idx == 0 {
						require.Equal(t, "", response.Info)
					} else {
						// the info is what was read from the kv store by the tx
						// each tx writes its own index, so the info should be the index of the previous tx
						require.Equal(t, fmt.Sprintf("%d", idx-1), response.Info)
					}
				}
				// confirm last write made it to the parent store
				latest := ctx.MultiStore().GetKVStore(testStoreKey).Get(itemKey)
				require.Equal(t, []byte(fmt.Sprintf("%d", len(res)-1)), latest)
			},
			expectedErr: nil,
		},
		{
			name:      "Test no stores on context should not panic",
			workers:   50,
			runs:      10,
			addStores: false,
			requests:  requestList(10),
			deliverTxFunc: func(ctx sdk.Context, req abci.RequestFinalizeBlock, tx sdk.Tx, checksum [32]byte) (res abci.ResponseFinalizeBlock) {
				return abci.ResponseFinalizeBlock{
					Info: fmt.Sprintf("%d", ctx.TxIndex()),
				}
			},
			assertions: func(t *testing.T, ctx sdk.Context, res []abci.ResponseFinalizeBlock) {
				for idx, response := range res {
					require.Equal(t, fmt.Sprintf("%d", idx), response.Info)
				}
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.runs; i++ {
				// set a tracer provider
				tp := trace.NewNoopTracerProvider()
				otel.SetTracerProvider(trace.NewNoopTracerProvider())
				tr := tp.Tracer("scheduler-test")
				ti := &tracing.Info{
					Tracer: &tr,
				}

				s := NewScheduler(tt.workers, ti, tt.deliverTxFunc)
				ctx := initTestCtx(tt.addStores)

				if tt.before != nil {
					tt.before(ctx)
				}

				res, err := s.ProcessAll(ctx, tt.requests)
				require.Len(t, res, len(tt.requests))

				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
				} else {
					tt.assertions(t, ctx, res)
				}
			}
		})
	}
}
