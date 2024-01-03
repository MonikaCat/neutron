package types_test

import (
	"testing"

	"github.com/MonikaCat/neutron/v2/app"

	"github.com/stretchr/testify/require"

	"github.com/MonikaCat/neutron/v2/x/cron/types"
)

func TestGenesisState_Validate(t *testing.T) {
	app.GetDefaultConfig()

	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params: types.Params{
					SecurityAddress: "neutron17dtl0mjt3t77kpuhg2edqzjpszulwhgzcdvagh",
					Limit:           1,
				},
			},
			valid: true,
		},
		{
			desc: "invalid genesis state - params are invalid",
			genState: &types.GenesisState{
				Params: types.Params{
					SecurityAddress: "",
					Limit:           0,
				},
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
