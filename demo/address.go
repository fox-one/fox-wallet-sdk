package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/fox-wallet-sdk"
)

func doAddresses(ctx context.Context, b *sdk.Broker, userID, assetID string) []*sdk.WithdrawAddress {
	addresses, err := b.FetchAddresses(ctx, userID, assetID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch addresses", addresses)

	return addresses
}

func doUpsertAddress(ctx context.Context, b *sdk.Broker, userID, pin string, addr *sdk.WithdrawAddress) *sdk.WithdrawAddress {
	a, err := b.UpsertAddress(ctx, userID, pin, addr)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("upsert address", a)

	return a
}

func doDeleteAddress(ctx context.Context, b *sdk.Broker, userID, addressID, pin string) {
	err := b.DeleteAddress(ctx, userID, addressID, pin)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("address deleted", addressID)
}
