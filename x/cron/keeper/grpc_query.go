package keeper

import (
	"github.com/MonikaCat/neutron/v2/x/cron/types"
)

var _ types.QueryServer = Keeper{}
