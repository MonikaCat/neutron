package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/neutron-org/neutron/x/interchaintxs/types"
)

const (
	LabelSubmitTx                  = "submit_tx"
	LabelHandleAcknowledgment      = "handle_ack"
	LabelLabelHandleChanOpenAck    = "handle_chan_open_ack"
	LabelRegisterInterchainAccount = "register_interchain_account"
	LabelHandleTimeout             = "handle_timeout"
)

type (
	Keeper struct {
		Codec               codec.BinaryCodec
		storeKey            storetypes.StoreKey
		memKey              storetypes.StoreKey
		channelKeeper       types.ChannelKeeper
		feeKeeper           types.FeeRefunderKeeper
		icaControllerKeeper types.ICAControllerKeeper
		sudoKeeper          types.WasmKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	channelKeeper types.ChannelKeeper,
	icaControllerKeeper types.ICAControllerKeeper,
	contractManagerKeeper types.WasmKeeper,
	feeKeeper types.FeeRefunderKeeper,
) *Keeper {
	return &Keeper{
		Codec:               cdc,
		storeKey:            storeKey,
		memKey:              memKey,
		channelKeeper:       channelKeeper,
		icaControllerKeeper: icaControllerKeeper,
		sudoKeeper:          contractManagerKeeper,
		feeKeeper:           feeKeeper,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
