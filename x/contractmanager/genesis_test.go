package contractmanager_test

import (
	"testing"

	"github.com/MonikaCat/neutron/v2/x/contractmanager/keeper"

	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/stretchr/testify/require"

	keepertest "github.com/MonikaCat/neutron/v2/testutil/contractmanager/keeper"
	"github.com/MonikaCat/neutron/v2/testutil/contractmanager/nullify"
	"github.com/MonikaCat/neutron/v2/x/contractmanager"
	"github.com/MonikaCat/neutron/v2/x/contractmanager/types"
)

func TestGenesis(t *testing.T) {
	payload1, err := keeper.PrepareSudoCallbackMessage(
		channeltypes.Packet{
			Sequence: 1,
		},
		&channeltypes.Acknowledgement{
			Response: &channeltypes.Acknowledgement_Result{
				Result: []byte("Result"),
			},
		})
	require.NoError(t, err)
	payload2, err := keeper.PrepareSudoCallbackMessage(
		channeltypes.Packet{
			Sequence: 2,
		}, nil)
	require.NoError(t, err)
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		FailuresList: []types.Failure{
			{
				Address:     "address1",
				Id:          1,
				SudoPayload: payload1,
			},
			{
				Address:     "address1",
				Id:          2,
				SudoPayload: payload2,
			},
		},
	}

	k, ctx := keepertest.ContractManagerKeeper(t, nil)
	contractmanager.InitGenesis(ctx, *k, genesisState)
	got := contractmanager.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.FailuresList, got.FailuresList)
}
