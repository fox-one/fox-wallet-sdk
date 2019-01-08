package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/fox-wallet-sdk"
	"github.com/fox-one/fox-wallet/models"
	"github.com/fox-one/mixin-sdk/mixin"
	"github.com/satori/go.uuid"
)

func doTransfer(ctx context.Context, b *sdk.Broker, dapp *mixin.User, userID, assetID, amount, pin, nonce string) *models.Snapshot {
	input := &mixin.TransferInput{
		AssetID:    assetID,
		Amount:     amount,
		OpponentID: userID,
		TraceID:    uuid.Must(uuid.NewV4()).String(),
		Memo:       "ping",
	}
	_, e := dapp.Transfer(ctx, input, PIN)
	if e != nil {
		log.Panicln(e)
	}

	log.Println("ping done")

	input.OpponentID = dapp.UserID
	input.Memo = "pong"
	input.TraceID = uuid.Must(uuid.NewV4()).String()

	snapshot, err := b.Transfer(ctx, userID, pin, nonce, input)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("do transfer", snapshot)

	return snapshot
}

func doWithdraw(ctx context.Context, b *sdk.Broker, dapp *mixin.User, userID, assetID, publicKey, amount, pin, nonce string) *models.Snapshot {
	input := &mixin.TransferInput{
		AssetID:    assetID,
		Amount:     amount,
		OpponentID: userID,
		TraceID:    uuid.Must(uuid.NewV4()).String(),
		Memo:       "ping",
	}
	_, e := dapp.Transfer(ctx, input, PIN)
	if e != nil {
		log.Panicln(e)
	}

	withdrawInput := &sdk.WithdrawInput{
		WithdrawAddress: mixin.WithdrawAddress{
			AssetID:   assetID,
			PublicKey: publicKey,
		},

		Amount:  amount,
		TraceID: uuid.Must(uuid.NewV4()).String(),
		Memo:    "pong",
	}

	snapshot, err := b.Withdraw(ctx, userID, pin, nonce, withdrawInput)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("do withdraw", snapshot)

	return snapshot
}

func doWithdrawFee(ctx context.Context, b *sdk.Broker, userID, assetID, publicKey, pin, nonce string) string {
	input := &mixin.WithdrawAddress{
		AssetID:   assetID,
		PublicKey: publicKey,
	}
	fee, err := b.FetchWithdrawFee(ctx, userID, pin, nonce, input)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch withdraw fee", fee)

	return fee
}
