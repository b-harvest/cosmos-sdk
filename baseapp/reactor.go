package baseapp

import (
	"fmt"
	"sync"

	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (app *BaseApp) startReactors() {
	go app.checkTxAsyncReactor()
}

type RequestCheckTxAsync struct {
	txBytes  []byte
	txType   abci.CheckTxType
	callback abci.CheckTxCallback
	prepare  *sync.WaitGroup
	tx       sdk.Tx
	err      error
}

func (app *BaseApp) checkTxAsyncReactor() {
	for req := range app.chCheckTx {
		req.prepare.Wait()
		if req.err != nil {
			req.callback(sdkerrors.ResponseCheckTxWithEvents(req.err, 0, 0, nil, false))
			continue
		}

		waits, signals := app.checkAccountWGs.Register(app.cdc, req.tx)

		go app.checkTxAsync(req, waits, signals)
	}
}

func (app *BaseApp) prepareCheckTx(req *RequestCheckTxAsync) {
	defer req.prepare.Done()
	req.tx, req.err = app.preCheckTx(req.txBytes)
}

func (app *BaseApp) checkTxAsync(req *RequestCheckTxAsync, waits []*sync.WaitGroup, signals []*AccountWG) {
	app.checkAccountWGs.Wait(waits)
	defer app.checkAccountWGs.Done(signals)

	var mode execMode
	if req.txType == abci.CheckTxType_New {
		mode = execModeCheck
	} else if req.txType == abci.CheckTxType_Recheck {
		mode = execModeReCheck
	} else {
		panic(fmt.Sprintf("unknown RequestCheckTx type: %s", req.txType))
	}

	gInfo, err := app.checkTx(mode, req.txBytes, req.tx)

	if err != nil {
		req.callback(sdkerrors.ResponseCheckTxWithEvents(err, gInfo.GasWanted, gInfo.GasUsed, nil, false))
		return
	}

	req.callback(&abci.ResponseCheckTx{
		GasWanted: int64(gInfo.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
	})
}
