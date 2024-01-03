package feeburner_test

import (
	"testing"

	"github.com/MonikaCat/neutron/v2/app"

	"github.com/stretchr/testify/require"

	"github.com/MonikaCat/neutron/v2/testutil/feeburner/keeper"
	"github.com/MonikaCat/neutron/v2/testutil/feeburner/nullify"
	"github.com/MonikaCat/neutron/v2/x/feeburner"
	"github.com/MonikaCat/neutron/v2/x/feeburner/types"
)

func TestGenesis(t *testing.T) {
	_ = app.GetDefaultConfig()

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k, ctx := keeper.FeeburnerKeeper(t)
	feeburner.InitGenesis(ctx, *k, genesisState)
	got := feeburner.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
