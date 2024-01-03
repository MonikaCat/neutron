package wasmbinding

import (
	contractmanagerkeeper "github.com/MonikaCat/neutron/v2/x/contractmanager/keeper"
	feeburnerkeeper "github.com/MonikaCat/neutron/v2/x/feeburner/keeper"
	feerefunderkeeper "github.com/MonikaCat/neutron/v2/x/feerefunder/keeper"
	icqkeeper "github.com/MonikaCat/neutron/v2/x/interchainqueries/keeper"
	icacontrollerkeeper "github.com/MonikaCat/neutron/v2/x/interchaintxs/keeper"

	tokenfactorykeeper "github.com/MonikaCat/neutron/v2/x/tokenfactory/keeper"
)

type QueryPlugin struct {
	icaControllerKeeper   *icacontrollerkeeper.Keeper
	icqKeeper             *icqkeeper.Keeper
	feeBurnerKeeper       *feeburnerkeeper.Keeper
	feeRefunderKeeper     *feerefunderkeeper.Keeper
	tokenFactoryKeeper    *tokenfactorykeeper.Keeper
	contractmanagerKeeper *contractmanagerkeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(icaControllerKeeper *icacontrollerkeeper.Keeper, icqKeeper *icqkeeper.Keeper, feeBurnerKeeper *feeburnerkeeper.Keeper, feeRefunderKeeper *feerefunderkeeper.Keeper, tfk *tokenfactorykeeper.Keeper, contractmanagerKeeper *contractmanagerkeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		icaControllerKeeper:   icaControllerKeeper,
		icqKeeper:             icqKeeper,
		feeBurnerKeeper:       feeBurnerKeeper,
		feeRefunderKeeper:     feeRefunderKeeper,
		tokenFactoryKeeper:    tfk,
		contractmanagerKeeper: contractmanagerKeeper,
	}
}
