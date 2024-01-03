package keeper_test

import (
	"testing"

	testkeeper "github.com/MonikaCat/neutron/v2/testutil/feerefunder/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/MonikaCat/neutron/v2/x/feerefunder/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.FeeKeeper(t, nil, nil)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
