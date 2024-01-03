package keeper

import (
	"github.com/MonikaCat/neutron/v2/x/interchaintxs/types"
)

var _ types.QueryServer = Keeper{}
