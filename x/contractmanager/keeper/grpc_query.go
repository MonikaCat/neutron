package keeper

import (
	"github.com/MonikaCat/neutron/v2/x/contractmanager/types"
)

var _ types.QueryServer = Keeper{}
