package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/tax/types"
)

// NewTax sets the next tax number to a given tax interface
func (k Keeper) NewTax(ctx sdk.Context, tax types.TaxI) types.TaxI {
	if err := tax.SetId(k.GetNextTaxIDWithUpdate(ctx)); err != nil {
		panic(err)
	}

	return tax
}

// GetTax implements TaxI.
func (k Keeper) GetTax(ctx sdk.Context, id uint64) (tax types.TaxI, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTaxKey(id))
	if bz == nil {
		return tax, false
	}

	return k.decodeTax(bz), true
}

// GetAllTaxes returns all taxes in the Keeper.
func (k Keeper) GetAllTaxes(ctx sdk.Context) (taxes []types.TaxI) {
	k.IterateAllTaxes(ctx, func(tax types.TaxI) (stop bool) {
		taxes = append(taxes, tax)
		return false
	})

	return taxes
}

// SetTax implements TaxI.
func (k Keeper) SetTax(ctx sdk.Context, tax types.TaxI) {
	id := tax.GetId()
	store := ctx.KVStore(k.storeKey)

	bz, err := k.MarshalTax(tax)
	if err != nil {
		panic(err)
	}

	store.Set(types.GetTaxKey(id), bz)
}

// RemoveTax removes an tax for the tax mapper store.
// NOTE: this will cause supply invariant violation if called
func (k Keeper) RemoveTax(ctx sdk.Context, tax types.TaxI) {
	id := tax.GetId()
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetTaxKey(id))
}

// IterateAllTaxes iterates over all the stored taxes and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateAllTaxes(ctx sdk.Context, cb func(tax types.TaxI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.TaxKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		tax := k.decodeTax(iterator.Value())

		if cb(tax) {
			break
		}
	}
}

// GetTaxesByFarmerAddrIndex reads from kvstore and return a specific Tax indexed by given farmer address
func (k Keeper) GetTaxesByFarmerAddrIndex(ctx sdk.Context, farmerAcc sdk.AccAddress) (taxes []types.TaxI) {
	k.IterateTaxesByFarmerAddr(ctx, farmerAcc, func(tax types.TaxI) bool {
		taxes = append(taxes, tax)
		return false
	})

	return taxes
}

// IterateTaxesByFarmerAddr iterates over all the stored taxes and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateTaxesByFarmerAddr(ctx sdk.Context, farmerAcc sdk.AccAddress, cb func(tax types.TaxI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetTaxesByFarmerAddrIndexKey(farmerAcc))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		taxID := gogotypes.UInt64Value{}

		err := k.cdc.Unmarshal(iterator.Value(), &taxID)
		if err != nil {
			panic(err)
		}
		tax, _ := k.GetTax(ctx, taxID.GetValue())
		if cb(tax) {
			break
		}
	}
}

// SetTaxIDByFarmerAddrIndex sets Index by FarmerAddr
// TODO: need to gas cost check for existing check or update everytime
func (k Keeper) SetTaxIDByFarmerAddrIndex(ctx sdk.Context, taxID uint64, farmerAcc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: taxID})
	store.Set(types.GetTaxByFarmerAddrIndexKey(farmerAcc, taxID), b)
}

