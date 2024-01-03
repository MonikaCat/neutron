package keeper

import (
	"github.com/MonikaCat/neutron/v2/x/dex/types"
)

var _ types.QueryServer = Keeper{}
