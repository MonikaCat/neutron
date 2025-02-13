//nolint:unused,unparam // Lots of useful test helper fns that we don't want to delete, also extra params we need to keep
package keeper_test

import (
	"math"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/neutron-org/neutron/v2/testutil/apptesting"
	math_utils "github.com/neutron-org/neutron/v2/utils/math"
	. "github.com/neutron-org/neutron/v2/x/dex/keeper"
	. "github.com/neutron-org/neutron/v2/x/dex/keeper/internal/testutils"
	"github.com/neutron-org/neutron/v2/x/dex/types"
)

// Test suite
type DexTestSuite struct {
	apptesting.KeeperTestHelper
	msgServer types.MsgServer
	alice     sdk.AccAddress
	bob       sdk.AccAddress
	carol     sdk.AccAddress
	dan       sdk.AccAddress
}

var defaultPairID *types.PairID = &types.PairID{Token0: "TokenA", Token1: "TokenB"}

var denomMultiple = sdk.NewInt(1000000)

var defaultTradePairID0To1 *types.TradePairID = &types.TradePairID{
	TakerDenom: "TokenA",
	MakerDenom: "TokenB",
}

var defaultTradePairID1To0 *types.TradePairID = &types.TradePairID{
	TakerDenom: "TokenB",
	MakerDenom: "TokenA",
}

func TestDexTestSuite(t *testing.T) {
	suite.Run(t, new(DexTestSuite))
}

func (s *DexTestSuite) SetupTest() {
	s.Setup()

	s.alice = sdk.AccAddress([]byte("alice"))
	s.bob = sdk.AccAddress([]byte("bob"))
	s.carol = sdk.AccAddress([]byte("carol"))
	s.dan = sdk.AccAddress([]byte("dan"))

	s.msgServer = NewMsgServerImpl(s.App.DexKeeper)
}

// NOTE: In order to simulate more realistic trade volume and avoid inadvertent failures due to ErrInvalidPositionSpread
// all of the basic user operations (fundXXXBalance, assertXXXBalance, XXXLimitsSells, etc.) treat TokenA and TokenB
// as BIG tokens with an exponent of 6. Ie. fundAliceBalance(10, 10) funds alices account with 10,000,000 small TokenA and TokenB.
// For tests requiring more accuracy methods that take Ints (ie. assertXXXAccountBalancesInt, NewWithdrawlInt) are used
// and assume that amount are being provided in terms of small tokens.

// Example:
// s.fundAliceBalances(10, 10)
// s.assertAliceBalances(10, 10) ==> True
// s.assertAliceBalancesInt(sdkmath.NewInt(10_000_000), sdkmath.NewInt(10_000_000)) ==> true

// Fund accounts

func (s *DexTestSuite) fundAccountBalances(account sdk.AccAddress, aBalance, bBalance int64) {
	aBalanceInt := sdkmath.NewInt(aBalance).Mul(denomMultiple)
	bBalanceInt := sdkmath.NewInt(bBalance).Mul(denomMultiple)
	balances := sdk.NewCoins(NewACoin(aBalanceInt), NewBCoin(bBalanceInt))

	FundAccount(s.App.BankKeeper, s.Ctx, account, balances)
	s.assertAccountBalances(account, aBalance, bBalance)
}

func (s *DexTestSuite) fundAccountBalancesWithDenom(
	addr sdk.AccAddress,
	amounts sdk.Coins,
) {
	FundAccount(s.App.BankKeeper, s.Ctx, addr, amounts)
}

func (s *DexTestSuite) fundAliceBalances(a, b int64) {
	s.fundAccountBalances(s.alice, a, b)
}

func (s *DexTestSuite) fundBobBalances(a, b int64) {
	s.fundAccountBalances(s.bob, a, b)
}

func (s *DexTestSuite) fundCarolBalances(a, b int64) {
	s.fundAccountBalances(s.carol, a, b)
}

func (s *DexTestSuite) fundDanBalances(a, b int64) {
	s.fundAccountBalances(s.dan, a, b)
}

/// Assert balances

func (s *DexTestSuite) assertAccountBalancesInt(
	account sdk.AccAddress,
	aBalance sdkmath.Int,
	bBalance sdkmath.Int,
) {
	aActual := s.App.BankKeeper.GetBalance(s.Ctx, account, "TokenA").Amount
	s.Assert().True(aBalance.Equal(aActual), "expected %s != actual %s", aBalance, aActual)

	bActual := s.App.BankKeeper.GetBalance(s.Ctx, account, "TokenB").Amount
	s.Assert().True(bBalance.Equal(bActual), "expected %s != actual %s", bBalance, bActual)
}

func (s *DexTestSuite) assertAccountBalances(
	account sdk.AccAddress,
	aBalance int64,
	bBalance int64,
) {
	s.assertAccountBalancesInt(account, sdkmath.NewInt(aBalance).Mul(denomMultiple), sdkmath.NewInt(bBalance).Mul(denomMultiple))
}

func (s *DexTestSuite) assertAccountBalanceWithDenomInt(
	account sdk.AccAddress,
	denom string,
	expBalance sdkmath.Int,
) {
	actualBalance := s.App.BankKeeper.GetBalance(s.Ctx, account, denom).Amount
	s.Assert().
		True(expBalance.Equal(actualBalance), "expected %s != actual %s", expBalance, actualBalance)
}

func (s *DexTestSuite) assertAccountBalanceWithDenom(
	account sdk.AccAddress,
	denom string,
	expBalance int64,
) {
	expBalanceInt := sdkmath.NewInt(expBalance).Mul(denomMultiple)
	s.assertAccountBalanceWithDenomInt(account, denom, expBalanceInt)
}

func (s *DexTestSuite) assertAliceBalances(a, b int64) {
	s.assertAccountBalances(s.alice, a, b)
}

func (s *DexTestSuite) assertAliceBalancesInt(a, b sdkmath.Int) {
	s.assertAccountBalancesInt(s.alice, a, b)
}

func (s *DexTestSuite) assertBobBalances(a, b int64) {
	s.assertAccountBalances(s.bob, a, b)
}

func (s *DexTestSuite) assertBobBalancesInt(a, b sdkmath.Int) {
	s.assertAccountBalancesInt(s.bob, a, b)
}

func (s *DexTestSuite) assertCarolBalances(a, b int64) {
	s.assertAccountBalances(s.carol, a, b)
}

func (s *DexTestSuite) assertCarolBalancesInt(a, b sdkmath.Int) {
	s.assertAccountBalancesInt(s.carol, a, b)
}

func (s *DexTestSuite) assertDanBalances(a, b int64) {
	s.assertAccountBalances(s.dan, a, b)
}

func (s *DexTestSuite) assertDanBalancesInt(a, b sdkmath.Int) {
	s.assertAccountBalancesInt(s.dan, a, b)
}

func (s *DexTestSuite) assertDexBalances(a, b int64) {
	s.assertAccountBalances(s.App.AccountKeeper.GetModuleAddress("dex"), a, b)
}

func (s *DexTestSuite) assertDexBalancesInt(a, b sdkmath.Int) {
	s.assertAccountBalancesInt(s.App.AccountKeeper.GetModuleAddress("dex"), a, b)
}

func (s *DexTestSuite) assertDexBalanceWithDenom(denom string, expectedAmount int64) {
	s.assertAccountBalanceWithDenom(
		s.App.AccountKeeper.GetModuleAddress("dex"),
		denom,
		expectedAmount,
	)
}

func (s *DexTestSuite) assertDexBalanceWithDenomInt(denom string, expectedAmount sdkmath.Int) {
	s.assertAccountBalanceWithDenomInt(
		s.App.AccountKeeper.GetModuleAddress("dex"),
		denom,
		expectedAmount,
	)
}

func (s *DexTestSuite) traceBalances() {
	aliceA := s.App.BankKeeper.GetBalance(s.Ctx, s.alice, "TokenA")
	aliceB := s.App.BankKeeper.GetBalance(s.Ctx, s.alice, "TokenB")
	bobA := s.App.BankKeeper.GetBalance(s.Ctx, s.bob, "TokenA")
	bobB := s.App.BankKeeper.GetBalance(s.Ctx, s.bob, "TokenB")
	carolA := s.App.BankKeeper.GetBalance(s.Ctx, s.carol, "TokenA")
	carolB := s.App.BankKeeper.GetBalance(s.Ctx, s.carol, "TokenB")
	danA := s.App.BankKeeper.GetBalance(s.Ctx, s.dan, "TokenA")
	danB := s.App.BankKeeper.GetBalance(s.Ctx, s.dan, "TokenB")
	s.T().Logf(
		"Alice: %+v %+v\nBob: %+v %+v\nCarol: %+v %+v\nDan: %+v %+v",
		aliceA, aliceB,
		bobA, bobB,
		carolA, carolB,
		danA, danB,
	)
}

/// Place limit order

func (s *DexTestSuite) aliceLimitSells(
	selling string,
	tick, amountIn int,
	orderTypeOpt ...types.LimitOrderType,
) string {
	return s.limitSellsSuccess(s.alice, selling, tick, amountIn, orderTypeOpt...)
}

func (s *DexTestSuite) bobLimitSells(
	selling string,
	tick, amountIn int,
	orderTypeOpt ...types.LimitOrderType,
) string {
	return s.limitSellsSuccess(s.bob, selling, tick, amountIn, orderTypeOpt...)
}

func (s *DexTestSuite) carolLimitSells(
	selling string,
	tick, amountIn int,
	orderTypeOpt ...types.LimitOrderType,
) string {
	return s.limitSellsSuccess(s.carol, selling, tick, amountIn, orderTypeOpt...)
}

func (s *DexTestSuite) danLimitSells(
	selling string,
	tick, amountIn int,
	orderTypeOpt ...types.LimitOrderType,
) string {
	return s.limitSellsSuccess(s.dan, selling, tick, amountIn, orderTypeOpt...)
}

func (s *DexTestSuite) limitSellsSuccess(
	account sdk.AccAddress,
	tokenIn string,
	tick, amountIn int,
	orderTypeOpt ...types.LimitOrderType,
) string {
	trancheKey, err := s.limitSells(account, tokenIn, tick, amountIn, orderTypeOpt...)
	s.Assert().Nil(err)
	return trancheKey
}

func (s *DexTestSuite) aliceLimitSellsGoodTil(
	selling string,
	tick, amountIn int,
	goodTil time.Time,
) string {
	return s.limitSellsGoodTil(s.alice, selling, tick, amountIn, goodTil)
}

func (s *DexTestSuite) bobLimitSellsGoodTil(
	selling string,
	tick, amountIn int,
	goodTil time.Time,
) string {
	return s.limitSellsGoodTil(s.bob, selling, tick, amountIn, goodTil)
}

func (s *DexTestSuite) carolLimitSellsGoodTil(
	selling string,
	tick, amountIn int,
	goodTil time.Time,
) string {
	return s.limitSellsGoodTil(s.carol, selling, tick, amountIn, goodTil)
}

func (s *DexTestSuite) danLimitSellsGoodTil(
	selling string,
	tick, amountIn int,
	goodTil time.Time,
) string {
	return s.limitSellsGoodTil(s.dan, selling, tick, amountIn, goodTil)
}

func (s *DexTestSuite) assertAliceLimitSellFails(
	err error,
	selling string,
	tickIndexNormalized, amountIn int,
	orderTypeOpt ...types.LimitOrderType,
) {
	s.assertLimitSellFails(s.alice, err, selling, tickIndexNormalized, amountIn, orderTypeOpt...)
}

func (s *DexTestSuite) assertBobLimitSellFails(
	err error,
	selling string,
	tickIndexNormalized, amountIn int,
	orderTypeOpt ...types.LimitOrderType,
) {
	s.assertLimitSellFails(s.bob, err, selling, tickIndexNormalized, amountIn, orderTypeOpt...)
}

func (s *DexTestSuite) assertCarolLimitSellFails(
	err error,
	selling string,
	tickIndexNormalized, amountIn int,
	orderTypeOpt ...types.LimitOrderType,
) {
	s.assertLimitSellFails(s.carol, err, selling, tickIndexNormalized, amountIn, orderTypeOpt...)
}

func (s *DexTestSuite) assertDanLimitSellFails(
	err error,
	selling string,
	tickIndexNormalized, amountIn int,
	orderTypeOpt ...types.LimitOrderType,
) {
	s.assertLimitSellFails(s.dan, err, selling, tickIndexNormalized, amountIn, orderTypeOpt...)
}

func (s *DexTestSuite) assertLimitSellFails(
	account sdk.AccAddress,
	expectedErr error,
	tokenIn string,
	tickIndexNormalized, amountIn int,
	orderTypeOpt ...types.LimitOrderType,
) {
	_, err := s.limitSells(account, tokenIn, tickIndexNormalized, amountIn, orderTypeOpt...)
	s.Assert().ErrorIs(err, expectedErr)
}

func (s *DexTestSuite) aliceLimitSellsWithMaxOut(
	selling string,
	tick, amountIn, maxAmountOut int,
) string {
	return s.limitSellsWithMaxOut(s.alice, selling, tick, amountIn, maxAmountOut)
}

func (s *DexTestSuite) bobLimitSellsWithMaxOut(
	selling string,
	tick, amountIn, maxAmountOut int,
) string {
	return s.limitSellsWithMaxOut(s.bob, selling, tick, amountIn, maxAmountOut)
}

func (s *DexTestSuite) carolLimitSellsWithMaxOut(
	selling string,
	tick, amountIn, maxAmountOut int,
) string {
	return s.limitSellsWithMaxOut(s.carol, selling, tick, amountIn, maxAmountOut)
}

func (s *DexTestSuite) danLimitSellsWithMaxOut(
	selling string,
	tick, amountIn, maxAmountOut int,
) string {
	return s.limitSellsWithMaxOut(s.dan, selling, tick, amountIn, maxAmountOut)
}

func (s *DexTestSuite) limitSellsWithMaxOut(
	account sdk.AccAddress,
	tokenIn string,
	tick, amountIn int,
	maxAmoutOut int,
) string {
	tokenIn, tokenOut := GetInOutTokens(tokenIn, "TokenA", "TokenB")
	maxAmountOutInt := sdkmath.NewInt(int64(maxAmoutOut)).Mul(denomMultiple)

	msg, err := s.msgServer.PlaceLimitOrder(s.GoCtx, &types.MsgPlaceLimitOrder{
		Creator:          account.String(),
		Receiver:         account.String(),
		TokenIn:          tokenIn,
		TokenOut:         tokenOut,
		TickIndexInToOut: int64(tick),
		AmountIn:         sdkmath.NewInt(int64(amountIn)).Mul(denomMultiple),
		OrderType:        types.LimitOrderType_FILL_OR_KILL,
		MaxAmountOut:     &maxAmountOutInt,
	})

	s.Assert().NoError(err)

	return msg.TrancheKey
}

func (s *DexTestSuite) limitSellsInt(
	account sdk.AccAddress,
	tokenIn string,
	tickIndexNormalized int, amountIn sdkmath.Int,
	orderTypeOpt ...types.LimitOrderType,
) (string, error) {
	var orderType types.LimitOrderType
	if len(orderTypeOpt) == 0 {
		orderType = types.LimitOrderType_GOOD_TIL_CANCELLED
	} else {
		orderType = orderTypeOpt[0]
	}

	tradePairID := types.NewTradePairIDFromTaker(defaultPairID, tokenIn)
	tickIndexTakerToMaker := tradePairID.TickIndexTakerToMaker(int64(tickIndexNormalized))
	msg, err := s.msgServer.PlaceLimitOrder(s.GoCtx, &types.MsgPlaceLimitOrder{
		Creator:          account.String(),
		Receiver:         account.String(),
		TokenIn:          tradePairID.TakerDenom,
		TokenOut:         tradePairID.MakerDenom,
		TickIndexInToOut: tickIndexTakerToMaker,
		AmountIn:         amountIn,
		OrderType:        orderType,
	})

	return msg.TrancheKey, err
}

func (s *DexTestSuite) limitSells(
	account sdk.AccAddress,
	tokenIn string,
	tickIndexNormalized, amountIn int,
	orderTypeOpt ...types.LimitOrderType,
) (string, error) {
	return s.limitSellsInt(account, tokenIn, tickIndexNormalized, sdkmath.NewInt(int64(amountIn)).Mul(denomMultiple), orderTypeOpt...)
}

func (s *DexTestSuite) limitSellsGoodTil(
	account sdk.AccAddress,
	tokenIn string,
	tick, amountIn int,
	goodTil time.Time,
) string {
	tradePairID := types.NewTradePairIDFromTaker(defaultPairID, tokenIn)
	tickIndexTakerToMaker := tradePairID.TickIndexTakerToMaker(int64(tick))

	msg, err := s.msgServer.PlaceLimitOrder(s.GoCtx, &types.MsgPlaceLimitOrder{
		Creator:          account.String(),
		Receiver:         account.String(),
		TokenIn:          tradePairID.TakerDenom,
		TokenOut:         tradePairID.MakerDenom,
		TickIndexInToOut: tickIndexTakerToMaker,
		AmountIn:         sdkmath.NewInt(int64(amountIn)).Mul(denomMultiple),
		OrderType:        types.LimitOrderType_GOOD_TIL_TIME,
		ExpirationTime:   &goodTil,
	})

	s.Assert().NoError(err)

	return msg.TrancheKey
}

// / Deposit
type Deposit struct {
	AmountA   sdkmath.Int
	AmountB   sdkmath.Int
	TickIndex int64
	Fee       uint64
}

type DepositOptions struct {
	DisableAutoswap bool
}

type DepositWithOptions struct {
	AmountA   sdkmath.Int
	AmountB   sdkmath.Int
	TickIndex int64
	Fee       uint64
	Options   DepositOptions
}

func NewDeposit(amountA, amountB, tickIndex, fee int) *Deposit {
	return &Deposit{
		AmountA:   sdkmath.NewInt(int64(amountA)).Mul(denomMultiple),
		AmountB:   sdkmath.NewInt(int64(amountB)).Mul(denomMultiple),
		TickIndex: int64(tickIndex),
		Fee:       uint64(fee),
	}
}

func NewDepositWithOptions(
	amountA, amountB, tickIndex, fee int,
	options DepositOptions,
) *DepositWithOptions {
	return &DepositWithOptions{
		AmountA:   sdkmath.NewInt(int64(amountA)).Mul(denomMultiple),
		AmountB:   sdkmath.NewInt(int64(amountB)).Mul(denomMultiple),
		TickIndex: int64(tickIndex),
		Fee:       uint64(fee),
		Options:   options,
	}
}

func (s *DexTestSuite) aliceDeposits(deposits ...*Deposit) {
	s.deposits(s.alice, deposits)
}

func (s *DexTestSuite) aliceDepositsWithOptions(deposits ...*DepositWithOptions) {
	s.depositsWithOptions(s.alice, deposits...)
}

func (s *DexTestSuite) bobDeposits(deposits ...*Deposit) {
	s.deposits(s.bob, deposits)
}

func (s *DexTestSuite) bobDepositsWithOptions(deposits ...*DepositWithOptions) {
	s.depositsWithOptions(s.bob, deposits...)
}

func (s *DexTestSuite) carolDeposits(deposits ...*Deposit) {
	s.deposits(s.carol, deposits)
}

func (s *DexTestSuite) danDeposits(deposits ...*Deposit) {
	s.deposits(s.dan, deposits)
}

func (s *DexTestSuite) deposits(
	account sdk.AccAddress,
	deposits []*Deposit,
	pairID ...types.PairID,
) {
	amountsA := make([]sdkmath.Int, len(deposits))
	amountsB := make([]sdkmath.Int, len(deposits))
	tickIndexes := make([]int64, len(deposits))
	fees := make([]uint64, len(deposits))
	options := make([]*types.DepositOptions, len(deposits))
	for i, e := range deposits {
		amountsA[i] = e.AmountA
		amountsB[i] = e.AmountB
		tickIndexes[i] = e.TickIndex
		fees[i] = e.Fee
		options[i] = &types.DepositOptions{DisableAutoswap: false}
	}

	var tokenA, tokenB string
	switch {
	case len(pairID) == 0:
		tokenA = "TokenA"
		tokenB = "TokenB"
	case len(pairID) == 1:
		tokenA = pairID[0].Token0
		tokenB = pairID[0].Token1
	case len(pairID) > 1:
		s.Assert().Fail("Only 1 pairID can be provided")
	}

	msg := &types.MsgDeposit{
		Creator:         account.String(),
		Receiver:        account.String(),
		TokenA:          tokenA,
		TokenB:          tokenB,
		AmountsA:        amountsA,
		AmountsB:        amountsB,
		TickIndexesAToB: tickIndexes,
		Fees:            fees,
		Options:         options,
	}
	err := msg.ValidateBasic()
	s.Assert().NoError(err)
	_, err = s.msgServer.Deposit(s.GoCtx, msg)
	s.Assert().Nil(err)
}

func (s *DexTestSuite) depositsWithOptions(
	account sdk.AccAddress,
	deposits ...*DepositWithOptions,
) {
	amountsA := make([]sdkmath.Int, len(deposits))
	amountsB := make([]sdkmath.Int, len(deposits))
	tickIndexes := make([]int64, len(deposits))
	fees := make([]uint64, len(deposits))
	options := make([]*types.DepositOptions, len(deposits))
	for i, e := range deposits {
		amountsA[i] = e.AmountA
		amountsB[i] = e.AmountB
		tickIndexes[i] = e.TickIndex
		fees[i] = e.Fee
		options[i] = &types.DepositOptions{
			DisableAutoswap: e.Options.DisableAutoswap,
		}
	}

	_, err := s.msgServer.Deposit(s.GoCtx, &types.MsgDeposit{
		Creator:         account.String(),
		Receiver:        account.String(),
		TokenA:          "TokenA",
		TokenB:          "TokenB",
		AmountsA:        amountsA,
		AmountsB:        amountsB,
		TickIndexesAToB: tickIndexes,
		Fees:            fees,
		Options:         options,
	})
	s.Assert().Nil(err)
}

func (s *DexTestSuite) getLiquidityAtTick(tickIndex int64, fee uint64) (sdkmath.Int, sdkmath.Int) {
	pool, err := s.App.DexKeeper.GetOrInitPool(s.Ctx, defaultPairID, tickIndex, fee)
	s.Assert().NoError(err)

	liquidityA := pool.LowerTick0.ReservesMakerDenom
	liquidityB := pool.UpperTick1.ReservesMakerDenom

	return liquidityA, liquidityB
}

func (s *DexTestSuite) getLiquidityAtTickWithDenom(
	pairID *types.PairID,
	tickIndex int64,
	fee uint64,
) (sdkmath.Int, sdkmath.Int) {
	pool, err := s.App.DexKeeper.GetOrInitPool(s.Ctx, pairID, tickIndex, fee)
	s.Assert().NoError(err)

	liquidityA := pool.LowerTick0.ReservesMakerDenom
	liquidityB := pool.UpperTick1.ReservesMakerDenom

	return liquidityA, liquidityB
}

func (s *DexTestSuite) assertAliceDepositFails(err error, deposits ...*Deposit) {
	s.assertDepositFails(s.alice, err, deposits...)
}

func (s *DexTestSuite) assertBobDepositFails(err error, deposits ...*Deposit) {
	s.assertDepositFails(s.bob, err, deposits...)
}

func (s *DexTestSuite) assertCarolDepositFails(err error, deposits ...*Deposit) {
	s.assertDepositFails(s.carol, err, deposits...)
}

func (s *DexTestSuite) assertDanDepositFails(err error, deposits ...*Deposit) {
	s.assertDepositFails(s.dan, err, deposits...)
}

func (s *DexTestSuite) assertDepositFails(
	account sdk.AccAddress,
	expectedErr error,
	deposits ...*Deposit,
) {
	amountsA := make([]sdkmath.Int, len(deposits))
	amountsB := make([]sdkmath.Int, len(deposits))
	tickIndexes := make([]int64, len(deposits))
	fees := make([]uint64, len(deposits))
	options := make([]*types.DepositOptions, len(deposits))
	for i, e := range deposits {
		amountsA[i] = e.AmountA
		amountsB[i] = e.AmountB
		tickIndexes[i] = e.TickIndex
		fees[i] = e.Fee
		options[i] = &types.DepositOptions{DisableAutoswap: true}
	}

	_, err := s.msgServer.Deposit(s.GoCtx, &types.MsgDeposit{
		Creator:         account.String(),
		Receiver:        account.String(),
		TokenA:          "TokenA",
		TokenB:          "TokenB",
		AmountsA:        amountsA,
		AmountsB:        amountsB,
		TickIndexesAToB: tickIndexes,
		Fees:            fees,
		Options:         options,
	})
	s.Assert().NotNil(err)
	s.Assert().ErrorIs(err, expectedErr)
}

func (s *DexTestSuite) assertDepositReponse(
	depositResponse, expectedDepositResponse DepositReponse,
) {
	for i := range expectedDepositResponse.amountsA {
		s.Assert().Equal(
			depositResponse.amountsA[i],
			expectedDepositResponse.amountsA[i],
			"Assertion failed for response.amountsA[%d]", i,
		)
		s.Assert().Equal(
			depositResponse.amountsB[i],
			expectedDepositResponse.amountsB[i],
			"Assertion failed for response.amountsB[%d]", i,
		)
	}
}

type DepositReponse struct {
	amountsA []sdkmath.Int
	amountsB []sdkmath.Int
}

// Withdraw
type Withdrawal struct {
	TickIndex int64
	Fee       uint64
	Shares    sdkmath.Int
}

func NewWithdrawalInt(shares sdkmath.Int, tick int64, fee uint64) *Withdrawal {
	return &Withdrawal{
		Shares:    shares,
		Fee:       fee,
		TickIndex: tick,
	}
}

// Multiples amount of shares to represent BIGtoken with exponent 6
func NewWithdrawal(shares, tick int64, fee uint64) *Withdrawal {
	return NewWithdrawalInt(sdkmath.NewInt(shares).Mul(denomMultiple), tick, fee)
}

func (s *DexTestSuite) aliceWithdraws(withdrawals ...*Withdrawal) {
	s.withdraws(s.alice, withdrawals...)
}

func (s *DexTestSuite) bobWithdraws(withdrawals ...*Withdrawal) {
	s.withdraws(s.bob, withdrawals...)
}

func (s *DexTestSuite) carolWithdraws(withdrawals ...*Withdrawal) {
	s.withdraws(s.carol, withdrawals...)
}

func (s *DexTestSuite) danWithdraws(withdrawals ...*Withdrawal) {
	s.withdraws(s.dan, withdrawals...)
}

func (s *DexTestSuite) withdraws(account sdk.AccAddress, withdrawals ...*Withdrawal) {
	tickIndexes := make([]int64, len(withdrawals))
	fee := make([]uint64, len(withdrawals))
	sharesToRemove := make([]sdkmath.Int, len(withdrawals))
	for i, e := range withdrawals {
		tickIndexes[i] = e.TickIndex
		fee[i] = e.Fee
		sharesToRemove[i] = e.Shares
	}

	_, err := s.msgServer.Withdrawal(s.GoCtx, &types.MsgWithdrawal{
		Creator:         account.String(),
		Receiver:        account.String(),
		TokenA:          "TokenA",
		TokenB:          "TokenB",
		SharesToRemove:  sharesToRemove,
		TickIndexesAToB: tickIndexes,
		Fees:            fee,
	})
	s.Assert().Nil(err)
}

func (s *DexTestSuite) aliceWithdrawFails(expectedErr error, withdrawals ...*Withdrawal) {
	s.withdrawFails(s.alice, expectedErr, withdrawals...)
}

func (s *DexTestSuite) bobWithdrawFails(expectedErr error, withdrawals ...*Withdrawal) {
	s.withdrawFails(s.bob, expectedErr, withdrawals...)
}

func (s *DexTestSuite) carolWithdrawFails(expectedErr error, withdrawals ...*Withdrawal) {
	s.withdrawFails(s.carol, expectedErr, withdrawals...)
}

func (s *DexTestSuite) danWithdrawFails(expectedErr error, withdrawals ...*Withdrawal) {
	s.withdrawFails(s.dan, expectedErr, withdrawals...)
}

func (s *DexTestSuite) withdrawFails(
	account sdk.AccAddress,
	expectedErr error,
	withdrawals ...*Withdrawal,
) {
	tickIndexes := make([]int64, len(withdrawals))
	fee := make([]uint64, len(withdrawals))
	sharesToRemove := make([]sdkmath.Int, len(withdrawals))
	for i, e := range withdrawals {
		tickIndexes[i] = e.TickIndex
		fee[i] = e.Fee
		sharesToRemove[i] = e.Shares
	}

	_, err := s.msgServer.Withdrawal(s.GoCtx, &types.MsgWithdrawal{
		Creator:         account.String(),
		Receiver:        account.String(),
		TokenA:          "TokenA",
		TokenB:          "TokenB",
		SharesToRemove:  sharesToRemove,
		TickIndexesAToB: tickIndexes,
		Fees:            fee,
	})
	s.Assert().NotNil(err)
	s.Assert().ErrorIs(err, expectedErr)
}

/// Cancel limit order

func (s *DexTestSuite) aliceCancelsLimitSell(trancheKey string) {
	s.cancelsLimitSell(s.alice, trancheKey)
}

func (s *DexTestSuite) bobCancelsLimitSell(trancheKey string) {
	s.cancelsLimitSell(s.bob, trancheKey)
}

func (s *DexTestSuite) carolCancelsLimitSell(trancheKey string) {
	s.cancelsLimitSell(s.carol, trancheKey)
}

func (s *DexTestSuite) danCancelsLimitSell(trancheKey string) {
	s.cancelsLimitSell(s.dan, trancheKey)
}

func (s *DexTestSuite) cancelsLimitSell(account sdk.AccAddress, trancheKey string) {
	_, err := s.msgServer.CancelLimitOrder(s.GoCtx, &types.MsgCancelLimitOrder{
		Creator:    account.String(),
		TrancheKey: trancheKey,
	})
	s.Assert().Nil(err)
}

func (s *DexTestSuite) aliceCancelsLimitSellFails(trancheKey string, expectedErr error) {
	s.cancelsLimitSellFails(s.alice, trancheKey, expectedErr)
}

func (s *DexTestSuite) bobCancelsLimitSellFails(trancheKey string, expectedErr error) {
	s.cancelsLimitSellFails(s.bob, trancheKey, expectedErr)
}

func (s *DexTestSuite) carolCancelsLimitSellFails(trancheKey string, expectedErr error) {
	s.cancelsLimitSellFails(s.carol, trancheKey, expectedErr)
}

func (s *DexTestSuite) danCancelsLimitSellFails(trancheKey string, expectedErr error) {
	s.cancelsLimitSellFails(s.dan, trancheKey, expectedErr)
}

func (s *DexTestSuite) cancelsLimitSellFails(
	account sdk.AccAddress,
	trancheKey string,
	expectedErr error,
) {
	_, err := s.msgServer.CancelLimitOrder(s.GoCtx, &types.MsgCancelLimitOrder{
		Creator:    account.String(),
		TrancheKey: trancheKey,
	})
	s.Assert().ErrorIs(err, expectedErr)
}

/// MultiHopSwap

func (s *DexTestSuite) aliceMultiHopSwaps(
	routes [][]string,
	amountIn int,
	exitLimitPrice math_utils.PrecDec,
	pickBest bool,
) {
	s.multiHopSwaps(s.alice, routes, amountIn, exitLimitPrice, pickBest)
}

func (s *DexTestSuite) bobMultiHopSwaps(
	routes [][]string,
	amountIn int,
	exitLimitPrice math_utils.PrecDec,
	pickBest bool,
) {
	s.multiHopSwaps(s.bob, routes, amountIn, exitLimitPrice, pickBest)
}

func (s *DexTestSuite) carolMultiHopSwaps(
	routes [][]string,
	amountIn int,
	exitLimitPrice math_utils.PrecDec,
	pickBest bool,
) {
	s.multiHopSwaps(s.carol, routes, amountIn, exitLimitPrice, pickBest)
}

func (s *DexTestSuite) danMultiHopSwaps(
	routes [][]string,
	amountIn int,
	exitLimitPrice math_utils.PrecDec,
	pickBest bool,
) {
	s.multiHopSwaps(s.dan, routes, amountIn, exitLimitPrice, pickBest)
}

func (s *DexTestSuite) multiHopSwaps(
	account sdk.AccAddress,
	routes [][]string,
	amountIn int,
	exitLimitPrice math_utils.PrecDec,
	pickBest bool,
) {
	msg := types.NewMsgMultiHopSwap(
		account.String(),
		account.String(),
		routes,
		sdkmath.NewInt(int64(amountIn)).Mul(denomMultiple),
		exitLimitPrice,
		pickBest,
	)
	_, err := s.msgServer.MultiHopSwap(s.GoCtx, msg)
	s.Assert().Nil(err)
}

func (s *DexTestSuite) aliceEstimatesMultiHopSwap(
	routes [][]string,
	amountIn int,
	exitLimitPrice math_utils.PrecDec,
	pickBest bool,
) (coinOut sdk.Coin) {
	multiHopRoutes := make([]*types.MultiHopRoute, len(routes))
	for i, hops := range routes {
		multiHopRoutes[i] = &types.MultiHopRoute{Hops: hops}
	}
	msg := &types.QueryEstimateMultiHopSwapRequest{
		Creator:        s.alice.String(),
		Receiver:       s.alice.String(),
		Routes:         multiHopRoutes,
		AmountIn:       sdkmath.NewInt(int64(amountIn)).Mul(denomMultiple),
		ExitLimitPrice: exitLimitPrice,
		PickBestRoute:  pickBest,
	}
	res, err := s.App.DexKeeper.EstimateMultiHopSwap(s.GoCtx, msg)
	s.Require().Nil(err)
	return res.CoinOut
}

func (s *DexTestSuite) aliceEstimatesMultiHopSwapFails(
	expectedErr error,
	routes [][]string,
	amountIn int,
	exitLimitPrice math_utils.PrecDec,
	pickBest bool,
) {
	multiHopRoutes := make([]*types.MultiHopRoute, len(routes))
	for i, hops := range routes {
		multiHopRoutes[i] = &types.MultiHopRoute{Hops: hops}
	}
	msg := &types.QueryEstimateMultiHopSwapRequest{
		Creator:        s.alice.String(),
		Receiver:       s.alice.String(),
		Routes:         multiHopRoutes,
		AmountIn:       sdkmath.NewInt(int64(amountIn)).Mul(denomMultiple),
		ExitLimitPrice: exitLimitPrice,
		PickBestRoute:  pickBest,
	}
	_, err := s.App.DexKeeper.EstimateMultiHopSwap(s.GoCtx, msg)
	s.Assert().ErrorIs(err, expectedErr)
}

func (s *DexTestSuite) aliceMultiHopSwapFails(
	err error,
	routes [][]string,
	amountIn int,
	exitLimitPrice math_utils.PrecDec,
	pickBest bool,
) {
	s.multiHopSwapFails(s.alice, err, routes, amountIn, exitLimitPrice, pickBest)
}

func (s *DexTestSuite) bobMultiHopSwapFails(
	err error,
	routes [][]string,
	amountIn int,
	exitLimitPrice math_utils.PrecDec,
	pickBest bool,
) {
	s.multiHopSwapFails(s.bob, err, routes, amountIn, exitLimitPrice, pickBest)
}

func (s *DexTestSuite) carolMultiHopSwapFails(
	err error,
	routes [][]string,
	amountIn int,
	exitLimitPrice math_utils.PrecDec,
	pickBest bool,
) {
	s.multiHopSwapFails(s.carol, err, routes, amountIn, exitLimitPrice, pickBest)
}

func (s *DexTestSuite) danMultiHopSwapFails(
	err error,
	routes [][]string,
	amountIn int,
	exitLimitPrice math_utils.PrecDec,
	pickBest bool,
) {
	s.multiHopSwapFails(s.dan, err, routes, amountIn, exitLimitPrice, pickBest)
}

func (s *DexTestSuite) multiHopSwapFails(
	account sdk.AccAddress,
	expectedErr error,
	routes [][]string,
	amountIn int,
	exitLimitPrice math_utils.PrecDec,
	pickBest bool,
) {
	msg := types.NewMsgMultiHopSwap(
		account.String(),
		account.String(),
		routes,
		sdkmath.NewInt(int64(amountIn)).Mul(denomMultiple),
		exitLimitPrice,
		pickBest,
	)
	_, err := s.msgServer.MultiHopSwap(s.GoCtx, msg)
	s.Assert().ErrorIs(err, expectedErr)
}

/// Withdraw filled limit order

func (s *DexTestSuite) aliceWithdrawsLimitSell(trancheKey string) {
	s.withdrawsLimitSell(s.alice, trancheKey)
}

func (s *DexTestSuite) bobWithdrawsLimitSell(trancheKey string) {
	s.withdrawsLimitSell(s.bob, trancheKey)
}

func (s *DexTestSuite) carolWithdrawsLimitSell(trancheKey string) {
	s.withdrawsLimitSell(s.carol, trancheKey)
}

func (s *DexTestSuite) danWithdrawsLimitSell(trancheKey string) {
	s.withdrawsLimitSell(s.dan, trancheKey)
}

func (s *DexTestSuite) withdrawsLimitSell(account sdk.AccAddress, trancheKey string) {
	_, err := s.msgServer.WithdrawFilledLimitOrder(s.GoCtx, &types.MsgWithdrawFilledLimitOrder{
		Creator:    account.String(),
		TrancheKey: trancheKey,
	})
	s.Assert().Nil(err)
}

func (s *DexTestSuite) aliceWithdrawLimitSellFails(expectedErr error, trancheKey string) {
	s.withdrawLimitSellFails(s.alice, expectedErr, trancheKey)
}

func (s *DexTestSuite) bobWithdrawLimitSellFails(expectedErr error, trancheKey string) {
	s.withdrawLimitSellFails(s.bob, expectedErr, trancheKey)
}

func (s *DexTestSuite) carolWithdrawLimitSellFails(expectedErr error, trancheKey string) {
	s.withdrawLimitSellFails(s.carol, expectedErr, trancheKey)
}

func (s *DexTestSuite) danWithdrawLimitSellFails(expectedErr error, trancheKey string) {
	s.withdrawLimitSellFails(s.dan, expectedErr, trancheKey)
}

func (s *DexTestSuite) withdrawLimitSellFails(
	account sdk.AccAddress,
	expectedErr error,
	trancheKey string,
) {
	_, err := s.msgServer.WithdrawFilledLimitOrder(s.GoCtx, &types.MsgWithdrawFilledLimitOrder{
		Creator:    account.String(),
		TrancheKey: trancheKey,
	})
	s.Assert().ErrorIs(err, expectedErr)
}

// Shares
func (s *DexTestSuite) getPoolShares(
	token0 string,
	token1 string,
	tick int64,
	fee uint64,
) (shares sdkmath.Int) {
	poolID, found := s.App.DexKeeper.GetPoolIDByParams(s.Ctx, &types.PairID{Token0: token0, Token1: token1}, tick, fee)
	if !found {
		return sdkmath.ZeroInt()
	}
	poolDenom := types.NewPoolDenom(poolID)
	return s.App.BankKeeper.GetSupply(s.Ctx, poolDenom).Amount
}

func (s *DexTestSuite) assertPoolShares(
	tick int64,
	fee uint64,
	sharesExpected uint64,
) {
	sharesExpectedInt := sdkmath.NewIntFromUint64(sharesExpected).Mul(denomMultiple)
	sharesOwned := s.getPoolShares("TokenA", "TokenB", tick, fee)
	s.Assert().Equal(sharesExpectedInt, sharesOwned)
}

func (s *DexTestSuite) getAccountShares(
	account sdk.AccAddress,
	token0 string,
	token1 string,
	tick int64,
	fee uint64,
) (shares sdkmath.Int) {
	id, found := s.App.DexKeeper.GetPoolIDByParams(s.Ctx, types.MustNewPairID(token0, token1), tick, fee)
	if !found {
		return sdkmath.ZeroInt()
	}

	poolDenom := types.NewPoolDenom(id)
	return s.App.BankKeeper.GetBalance(s.Ctx, account, poolDenom).Amount
}

func (s *DexTestSuite) assertAccountSharesInt(
	account sdk.AccAddress,
	tick int64,
	fee uint64,
	sharesExpected sdkmath.Int,
) {
	sharesOwned := s.getAccountShares(account, "TokenA", "TokenB", tick, fee)
	s.Assert().
		Equal(sharesExpected, sharesOwned, "expected %s != actual %s", sharesExpected, sharesOwned)
}

func (s *DexTestSuite) assertAccountShares(
	account sdk.AccAddress,
	tick int64,
	fee uint64,
	sharesExpected uint64,
) {
	sharesExpectedInt := sdkmath.NewIntFromUint64(sharesExpected).Mul(denomMultiple)
	s.assertAccountSharesInt(account, tick, fee, sharesExpectedInt)
}

func (s *DexTestSuite) assertAliceShares(tick int64, fee, sharesExpected uint64) {
	s.assertAccountShares(s.alice, tick, fee, sharesExpected)
}

func (s *DexTestSuite) assertBobShares(tick int64, fee, sharesExpected uint64) {
	s.assertAccountShares(s.bob, tick, fee, sharesExpected)
}

func (s *DexTestSuite) assertCarolShares(tick int64, fee, sharesExpected uint64) {
	s.assertAccountShares(s.carol, tick, fee, sharesExpected)
}

func (s *DexTestSuite) assertDanShares(tick int64, fee, sharesExpected uint64) {
	s.assertAccountShares(s.dan, tick, fee, sharesExpected)
}

// Ticks
func (s *DexTestSuite) assertCurrentTicks(
	expected1To0 int64,
	expected0To1 int64,
) {
	s.assertCurr0To1(expected0To1)
	s.assertCurr1To0(expected1To0)
}

func (s *DexTestSuite) assertCurr0To1(curr0To1Expected int64) {
	curr0To1Actual, found := s.App.DexKeeper.GetCurrTickIndexTakerToMakerNormalized(
		s.Ctx,
		defaultTradePairID0To1,
	)
	if curr0To1Expected == math.MaxInt64 {
		s.Assert().False(found)
	} else {
		s.Assert().Equal(curr0To1Expected, curr0To1Actual)
	}
}

func (s *DexTestSuite) assertCurr1To0(curr1To0Expected int64) {
	curr1to0Actual, found := s.App.DexKeeper.GetCurrTickIndexTakerToMakerNormalized(
		s.Ctx,
		defaultTradePairID1To0,
	)
	if curr1To0Expected == math.MinInt64 {
		s.Assert().False(found)
	} else {
		s.Assert().Equal(curr1To0Expected, curr1to0Actual)
	}
}

// Pool liquidity (i.e. deposited rather than LO)
func (s *DexTestSuite) assertLiquidityAtTickInt(
	amountA, amountB sdkmath.Int,
	tickIndex int64,
	fee uint64,
) {
	liquidityA, liquidityB := s.getLiquidityAtTick(tickIndex, fee)
	s.Assert().
		True(amountA.Equal(liquidityA), "liquidity A: actual %s, expected %s", liquidityA, amountA)
	s.Assert().
		True(amountB.Equal(liquidityB), "liquidity B: actual %s, expected %s", liquidityB, amountB)
}

func (s *DexTestSuite) assertLiquidityAtTick(
	amountA, amountB int64,
	tickIndex int64,
	fee uint64,
) {
	amountAInt := sdkmath.NewInt(amountA).Mul(denomMultiple)
	amountBInt := sdkmath.NewInt(amountB).Mul(denomMultiple)
	s.assertLiquidityAtTickInt(amountAInt, amountBInt, tickIndex, fee)
}

func (s *DexTestSuite) assertLiquidityAtTickWithDenomInt(
	pairID *types.PairID,
	expected0, expected1 sdkmath.Int,
	tickIndex int64,
	fee uint64,
) {
	liquidity0, liquidity1 := s.getLiquidityAtTickWithDenom(pairID, tickIndex, fee)
	s.Assert().
		True(expected0.Equal(liquidity0), "liquidity 0: actual %s, expected %s", liquidity0, expected0)
	s.Assert().
		True(expected1.Equal(liquidity1), "liquidity 1: actual %s, expected %s", liquidity1, expected1)
}

func (s *DexTestSuite) assertLiquidityAtTickWithDenom(
	pairID *types.PairID,
	expected0,
	expected1,
	tickIndex int64,
	fee uint64,
) {
	expected0Int := sdkmath.NewInt(expected0).Mul(denomMultiple)
	expected1Int := sdkmath.NewInt(expected1).Mul(denomMultiple)
	s.assertLiquidityAtTickWithDenomInt(pairID, expected0Int, expected1Int, tickIndex, fee)
}

func (s *DexTestSuite) assertPoolLiquidity(
	amountA, amountB int64,
	tickIndex int64,
	fee uint64,
) {
	s.assertLiquidityAtTick(amountA, amountB, tickIndex, fee)
}

func (s *DexTestSuite) assertNoLiquidityAtTick(tickIndex int64, fee uint64) {
	s.assertLiquidityAtTick(0, 0, tickIndex, fee)
}

// Filled limit liquidity
func (s *DexTestSuite) assertAliceLimitFilledAtTickAtIndex(
	selling string,
	amount int,
	tickIndex int64,
	trancheKey string,
) {
	s.assertLimitFilledAtTickAtIndex(s.alice, selling, amount, tickIndex, trancheKey)
}

func (s *DexTestSuite) assertBobLimitFilledAtTickAtIndex(
	selling string,
	amount int,
	tickIndex int64,
	trancheKey string,
) {
	s.assertLimitFilledAtTickAtIndex(s.bob, selling, amount, tickIndex, trancheKey)
}

func (s *DexTestSuite) assertCarolLimitFilledAtTickAtIndex(
	selling string,
	amount int,
	tickIndex int64,
	trancheKey string,
) {
	s.assertLimitFilledAtTickAtIndex(s.carol, selling, amount, tickIndex, trancheKey)
}

func (s *DexTestSuite) assertDanLimitFilledAtTickAtIndex(
	selling string,
	amount int,
	tickIndex int64,
	trancheKey string,
) {
	s.assertLimitFilledAtTickAtIndex(s.dan, selling, amount, tickIndex, trancheKey)
}

func (s *DexTestSuite) assertLimitFilledAtTickAtIndex(
	account sdk.AccAddress,
	selling string,
	amount int,
	tickIndex int64,
	trancheKey string,
) {
	userShares, totalShares := s.getLimitUserSharesAtTick(
		account,
		selling,
		tickIndex,
	), s.getLimitTotalSharesAtTick(
		selling,
		tickIndex,
	)
	userRatio := math_utils.NewPrecDecFromInt(userShares).QuoInt(totalShares)
	filled := s.getLimitFilledLiquidityAtTickAtIndex(selling, tickIndex, trancheKey)
	amt := sdkmath.NewInt(int64(amount)).Mul(denomMultiple)
	userFilled := userRatio.MulInt(filled).RoundInt()
	s.Assert().True(amt.Equal(userFilled))
}

// Limit liquidity
func (s *DexTestSuite) assertAliceLimitLiquidityAtTick(
	selling string,
	amount int,
	tickIndex int64,
) {
	s.assertAccountLimitLiquidityAtTick(s.alice, selling, amount, tickIndex)
}

func (s *DexTestSuite) assertBobLimitLiquidityAtTick(
	selling string,
	amount int,
	tickIndex int64,
) {
	s.assertAccountLimitLiquidityAtTick(s.bob, selling, amount, tickIndex)
}

func (s *DexTestSuite) assertCarolLimitLiquidityAtTick(
	selling string,
	amount int,
	tickIndex int64,
) {
	s.assertAccountLimitLiquidityAtTick(s.carol, selling, amount, tickIndex)
}

func (s *DexTestSuite) assertDanLimitLiquidityAtTick(
	selling string,
	amount int,
	tickIndex int64,
) {
	s.assertAccountLimitLiquidityAtTick(s.dan, selling, amount, tickIndex)
}

func (s *DexTestSuite) assertAccountLimitLiquidityAtTick(
	account sdk.AccAddress,
	selling string,
	amount int,
	tickIndexNormalized int64,
) {
	userShares := s.getLimitUserSharesAtTick(account, selling, tickIndexNormalized)
	totalShares := s.getLimitTotalSharesAtTick(selling, tickIndexNormalized)
	userRatio := math_utils.NewPrecDecFromInt(userShares).QuoInt(totalShares)
	userLiquidity := userRatio.MulInt64(int64(amount)).TruncateInt()

	s.assertLimitLiquidityAtTick(selling, tickIndexNormalized, userLiquidity.Int64())
}

func (s *DexTestSuite) assertLimitLiquidityAtTick(
	selling string,
	tickIndexNormalized, amount int64,
) {
	s.assertLimitLiquidityAtTickInt(selling, tickIndexNormalized, sdkmath.NewInt(amount))
}

func (s *DexTestSuite) assertLimitLiquidityAtTickInt(
	selling string,
	tickIndexNormalized int64,
	amount sdkmath.Int,
) {
	amount = amount.Mul(denomMultiple)
	tradePairID := defaultPairID.MustTradePairIDFromMaker(selling)
	tickIndexTakerToMaker := tradePairID.TickIndexTakerToMaker(tickIndexNormalized)
	tranches := s.App.DexKeeper.GetAllLimitOrderTrancheAtIndex(
		s.Ctx,
		tradePairID,
		tickIndexTakerToMaker,
	)
	liquidity := sdkmath.ZeroInt()
	for _, t := range tranches {
		if !t.IsExpired(s.Ctx) {
			liquidity = liquidity.Add(t.ReservesMakerDenom)
		}
	}

	s.Assert().
		True(amount.Equal(liquidity), "Incorrect liquidity: expected %s, have %s", amount.String(), liquidity.String())
}

func (s *DexTestSuite) assertFillAndPlaceTrancheKeys(
	selling string,
	tickIndexNormalized int64,
	expectedFill, expectedPlace string,
) {
	tradePairID := defaultPairID.MustTradePairIDFromMaker(selling)
	tickIndexTakerToMaker := tradePairID.TickIndexTakerToMaker(tickIndexNormalized)
	placeTranche := s.App.DexKeeper.GetPlaceTranche(s.Ctx, tradePairID, tickIndexTakerToMaker)
	fillTranche, foundFill := s.App.DexKeeper.GetFillTranche(
		s.Ctx,
		tradePairID,
		tickIndexTakerToMaker,
	)
	placeKey, fillKey := "", ""
	if placeTranche != nil {
		placeKey = placeTranche.Key.TrancheKey
	}

	if foundFill {
		fillKey = fillTranche.Key.TrancheKey
	}
	s.Assert().Equal(expectedFill, fillKey)
	s.Assert().Equal(expectedPlace, placeKey)
}

// Limit order map helpers
func (s *DexTestSuite) getLimitUserSharesAtTick(
	account sdk.AccAddress,
	selling string,
	tickIndexNormalized int64,
) sdkmath.Int {
	tradePairID := defaultPairID.MustTradePairIDFromMaker(selling)
	tickIndexTakerToMaker := tradePairID.TickIndexTakerToMaker(tickIndexNormalized)
	tranches := s.App.DexKeeper.GetAllLimitOrderTrancheAtIndex(
		s.Ctx,
		tradePairID,
		tickIndexTakerToMaker,
	)
	fillTranche := tranches[0]
	// get user shares and total shares
	userShares := s.getLimitUserSharesAtTickAtIndex(account, fillTranche.Key.TrancheKey)
	if len(tranches) >= 2 {
		userShares = userShares.Add(
			s.getLimitUserSharesAtTickAtIndex(account, tranches[1].Key.TrancheKey),
		)
	}

	return userShares
}

func (s *DexTestSuite) getLimitUserSharesAtTickAtIndex(
	account sdk.AccAddress,
	trancheKey string,
) sdkmath.Int {
	userShares, found := s.App.DexKeeper.GetLimitOrderTrancheUser(
		s.Ctx,
		account.String(),
		trancheKey,
	)
	s.Assert().True(found, "Failed to get limit order user shares for index %s", trancheKey)
	return userShares.SharesOwned
}

func (s *DexTestSuite) getLimitTotalSharesAtTick(
	selling string,
	tickIndexNormalized int64,
) sdkmath.Int {
	tradePairID := defaultPairID.MustTradePairIDFromMaker(selling)
	tickIndexTakerToMaker := tradePairID.TickIndexTakerToMaker(tickIndexNormalized)
	tranches := s.App.DexKeeper.GetAllLimitOrderTrancheAtIndex(
		s.Ctx,
		tradePairID,
		tickIndexTakerToMaker,
	)
	// get user shares and total shares
	totalShares := sdkmath.ZeroInt()
	for _, t := range tranches {
		totalShares = totalShares.Add(t.TotalMakerDenom)
	}

	return totalShares
}

func (s *DexTestSuite) getLimitFilledLiquidityAtTickAtIndex(
	selling string,
	tickIndex int64,
	trancheKey string,
) sdkmath.Int {
	// grab fill tranche reserves and shares
	tradePairID := defaultPairID.MustTradePairIDFromMaker(selling)
	tranche, _, found := s.App.DexKeeper.FindLimitOrderTranche(s.Ctx, &types.LimitOrderTrancheKey{
		TradePairId:           tradePairID,
		TickIndexTakerToMaker: tickIndex,
		TrancheKey:            trancheKey,
	})
	s.Assert().True(found, "Failed to get limit order filled reserves for index %s", trancheKey)

	return tranche.ReservesTakerDenom
}

func (s *DexTestSuite) getLimitReservesAtTickAtKey(
	selling string,
	tickIndex int64,
	trancheKey string,
) sdkmath.Int {
	// grab fill tranche reserves and shares
	tradePairID := defaultPairID.MustTradePairIDFromMaker(selling)
	tranche, _, found := s.App.DexKeeper.FindLimitOrderTranche(s.Ctx, &types.LimitOrderTrancheKey{
		TradePairId:           tradePairID,
		TickIndexTakerToMaker: tickIndex,
		TrancheKey:            trancheKey,
	})
	s.Assert().True(found, "Failed to get limit order reserves for index %s", trancheKey)

	return tranche.ReservesMakerDenom
}

func (s *DexTestSuite) assertNLimitOrderExpiration(expected int) {
	exps := s.App.DexKeeper.GetAllLimitOrderExpiration(s.Ctx)
	s.Assert().Equal(expected, len(exps))
}

func (s *DexTestSuite) calcAutoswapSharesMinted(
	centerTick int64,
	fee uint64,
	residual0, residual1, balanced0, balanced1, totalShares, valuePool int64,
) sdkmath.Int {
	residual0Int, residual1Int, balanced0Int, balanced1Int, totalSharesInt, valuePoolInt := sdkmath.NewInt(residual0),
		sdkmath.NewInt(residual1),
		sdkmath.NewInt(balanced0),
		sdkmath.NewInt(balanced1),
		sdkmath.NewInt(totalShares),
		sdkmath.NewInt(valuePool)

	// residualValue = 1.0001^-f * residualAmount0 + 1.0001^{i-f} * residualAmount1
	// balancedValue = balancedAmount0 + 1.0001^{i} * balancedAmount1
	// value = residualValue + balancedValue
	// shares minted = value * totalShares / valuePool

	centerPrice := types.MustCalcPrice(-1 * centerTick)
	leftPrice := types.MustCalcPrice(-1 * (centerTick - int64(fee)))
	discountPrice := types.MustCalcPrice(-1 * int64(fee))

	balancedValue := math_utils.NewPrecDecFromInt(balanced0Int).
		Add(centerPrice.MulInt(balanced1Int)).
		TruncateInt()
	residualValue := discountPrice.MulInt(residual0Int).
		Add(leftPrice.Mul(math_utils.NewPrecDecFromInt(residual1Int))).
		TruncateInt()
	valueMint := balancedValue.Add(residualValue)

	return valueMint.Mul(totalSharesInt).Quo(valuePoolInt)
}

func (s *DexTestSuite) calcSharesMinted(centerTick, amount0Int, amount1Int int64) sdkmath.Int {
	amount0, amount1 := sdkmath.NewInt(amount0Int), sdkmath.NewInt(amount1Int)
	centerPrice := types.MustCalcPrice(-1 * centerTick)

	return math_utils.NewPrecDecFromInt(amount0).Add(centerPrice.Mul(math_utils.NewPrecDecFromInt(amount1))).TruncateInt()
}

func (s *DexTestSuite) calcExpectedBalancesAfterWithdrawOnePool(
	sharesMinted sdkmath.Int,
	account sdk.AccAddress,
	tickIndex int64,
	fee uint64,
) (sdkmath.Int, sdkmath.Int, sdkmath.Int, sdkmath.Int) {
	dexCurrentBalance0 := s.App.BankKeeper.GetBalance(
		s.Ctx,
		s.App.AccountKeeper.GetModuleAddress("dex"),
		"TokenA",
	).Amount
	dexCurrentBalance1 := s.App.BankKeeper.GetBalance(
		s.Ctx,
		s.App.AccountKeeper.GetModuleAddress("dex"),
		"TokenB",
	).Amount
	currentBalance0 := s.App.BankKeeper.GetBalance(s.Ctx, account, "TokenA").Amount
	currentBalance1 := s.App.BankKeeper.GetBalance(s.Ctx, account, "TokenB").Amount
	amountPool0, amountPool1 := s.getLiquidityAtTick(tickIndex, fee)
	poolShares := s.getPoolShares("TokenA", "TokenB", tickIndex, fee)

	amountOut0 := amountPool0.Mul(sharesMinted).Quo(poolShares)
	amountOut1 := amountPool1.Mul(sharesMinted).Quo(poolShares)

	expectedBalance0 := currentBalance0.Add(amountOut0)
	expectedBalance1 := currentBalance1.Add(amountOut1)
	dexExpectedBalance0 := dexCurrentBalance0.Sub(amountOut0)
	dexExpectedBalance1 := dexCurrentBalance1.Sub(amountOut1)

	return expectedBalance0, expectedBalance1, dexExpectedBalance0, dexExpectedBalance1
}

func (s *DexTestSuite) nextBlockWithTime(blockTime time.Time) {
	newCtx := s.Ctx.WithBlockTime(blockTime)
	s.Ctx = newCtx
	s.GoCtx = sdk.WrapSDKContext(newCtx)
	s.App.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		Height: s.App.LastBlockHeight() + 1, AppHash: s.App.LastCommitID().Hash,
		Time: blockTime,
	}})
}
