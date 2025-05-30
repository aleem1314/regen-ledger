package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	marketplacev1 "github.com/regen-network/regen-ledger/api/v2/regen/ecocredit/marketplace/v1"
	api "github.com/regen-network/regen-ledger/api/v2/regen/ecocredit/v1"
	regentypes "github.com/regen-network/regen-ledger/types/v2"
	regenerrors "github.com/regen-network/regen-ledger/types/v2/errors"
	types "github.com/regen-network/regen-ledger/x/ecocredit/v3/base/types/v1"
)

// Params queries the ecocredit module parameters.
// Deprecated: This rpc method is deprecated and will be removed in next version.
// Use individual param query instead.
func (k Keeper) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {

	allowlistEnabled, err := k.stateStore.ClassCreatorAllowlistTable().Get(ctx)
	if err != nil {
		return nil, regenerrors.ErrInternal.Wrapf("unable to get allowlist param: %s", err.Error())
	}

	itr, err := k.stateStore.AllowedClassCreatorTable().List(ctx, api.AllowedClassCreatorPrimaryKey{})
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	creators := []string{}
	for itr.Next() {
		val, err := itr.Value()
		if err != nil {
			return nil, err
		}

		creators = append(creators, sdk.AccAddress(val.Address).String())
	}

	classFee, err := k.stateStore.ClassFeeTable().Get(ctx)
	if err != nil {
		return nil, regenerrors.ErrInternal.Wrapf("unable to get class fee param: %s", err.Error())
	}
	basketFee, err := k.basketStore.BasketFeeTable().Get(ctx)
	if err != nil {
		return nil, regenerrors.ErrInternal.Wrapf("unable to get basket fee: %s", err.Error())
	}

	allowedDenomsItr, err := k.marketStore.AllowedDenomTable().List(ctx, marketplacev1.AllowedDenomPrimaryKey{})
	if err != nil {
		return nil, err
	}
	defer allowedDenomsItr.Close()

	var allowedDenoms []*types.AllowedDenom
	for allowedDenomsItr.Next() {
		val, err := allowedDenomsItr.Value()
		if err != nil {
			return nil, err
		}

		allowedDenoms = append(allowedDenoms, &types.AllowedDenom{
			BankDenom:    val.BankDenom,
			DisplayDenom: val.DisplayDenom,
			Exponent:     val.Exponent,
		})
	}

	var allowedBridgeChains []string
	it, err := k.stateStore.AllowedBridgeChainTable().List(ctx, api.AllowedBridgeChainPrimaryKey{})
	if err != nil {
		return nil, err
	}
	defer it.Close()

	for it.Next() {
		entry, err := it.Value()
		if err != nil {
			return nil, err
		}
		allowedBridgeChains = append(allowedBridgeChains, entry.ChainName)
	}

	return &types.QueryParamsResponse{
		Params: &types.Params{
			AllowedClassCreators: creators,
			AllowlistEnabled:     allowlistEnabled.Enabled,
			CreditClassFee:       sdk.Coins{regentypes.CoinFromCosmosAPILegacy(classFee.Fee)},
			BasketFee:            sdk.Coins{regentypes.CoinFromCosmosAPILegacy(basketFee.Fee)},
			AllowedDenoms:        allowedDenoms,
			AllowedBridgeChains:  allowedBridgeChains,
		},
	}, nil
}
