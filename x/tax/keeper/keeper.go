package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/tax/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper of the tax store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramSpace paramtypes.Subspace

	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
	distrKeeper   types.DistributionKeeper

	blockedAddrs map[string]bool
}

// NewKeeper returns a tax keeper. It handles:
// - creating new ModuleAccounts for each pool ReserveAccount
// - sending to and from ModuleAccounts
// - minting, burning PoolCoins
func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, distrKeeper types.DistributionKeeper,
	blockedAddrs map[string]bool,
) Keeper {
	// TODO: TBD module account for tax
	//// ensure tax module account is set
	//if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
	//	panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	//}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramSpace:    paramSpace,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrKeeper:   distrKeeper,
		blockedAddrs:  blockedAddrs,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetParams gets the parameters for the tax module.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the parameters for the tax module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetNextTaxIDWithUpdate returns and increments the global Tax ID counter.
// If the global tax number is not set, it initializes it with value 1.
func (k Keeper) GetNextTaxIDWithUpdate(ctx sdk.Context) uint64 {
	var id uint64
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GlobalTaxTaxIDKey)
	if bz == nil {
		// initialize the TaxId
		id = 1
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}

		id = val.GetValue()
	}
	bz = k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id + 1})
	store.Set(types.GlobalTaxTaxIDKey, bz)
	return id
}

func (k Keeper) decodeTax(bz []byte) types.TaxI {
	acc, err := k.UnmarshalTax(bz)
	if err != nil {
		panic(err)
	}

	return acc
}

// MarshalTax protobuf serializes an Tax interface
func (k Keeper) MarshalTax(tax types.TaxI) ([]byte, error) { // nolint:interfacer
	return k.cdc.MarshalInterface(tax)
}

// UnmarshalTax returns an Tax interface from raw encoded tax
// bytes of a Proto-based Tax type
func (k Keeper) UnmarshalTax(bz []byte) (types.TaxI, error) {
	var acc types.TaxI
	return acc, k.cdc.UnmarshalInterface(bz, &acc)
}

// GetCodec return codec.Codec object used by the keeper
func (k Keeper) GetCodec() codec.BinaryCodec { return k.cdc }
