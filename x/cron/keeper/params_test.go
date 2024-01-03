package keeper_test

import (
	"testing"

	"github.com/MonikaCat/neutron/v2/testutil"

	"github.com/MonikaCat/neutron/v2/app"

	testkeeper "github.com/MonikaCat/neutron/v2/testutil/cron/keeper"

	"github.com/stretchr/testify/require"

	"github.com/MonikaCat/neutron/v2/x/cron/types"
)

func TestGetParams(t *testing.T) {
	_ = app.GetDefaultConfig()

	k, ctx := testkeeper.CronKeeper(t, nil, nil)
	params := types.Params{
		SecurityAddress: testutil.TestOwnerAddress,
		Limit:           5,
	}

	err := k.SetParams(ctx, params)
	require.NoError(t, err)

	require.EqualValues(t, params, k.GetParams(ctx))
}
