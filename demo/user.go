package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/fox-wallet-sdk"
)

func doCreateUser(ctx context.Context, b *sdk.Broker, pin string) *sdk.User {
	u, e := b.CreateUser(ctx, "", pin)
	if e != nil {
		log.Panicln(e)
	}
	printJSON("created user", u)

	return u
}

func doModifyPIN(ctx context.Context, b *sdk.Broker, userID string, pin, newPIN string) {
	err := b.ModifyPIN(ctx, userID, pin, newPIN)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("modify PIN succ")
}

func doVerifyPIN(ctx context.Context, b *sdk.Broker, userID string, pin string) {
	err := b.VerifyPIN(ctx, userID, pin)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("verify PIN succ")
}
