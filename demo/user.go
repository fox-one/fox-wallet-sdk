package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/fox-wallet-sdk"
	"github.com/fox-one/fox-wallet/models"
)

func doCreateUser(ctx context.Context, b *sdk.Broker, pin string) *models.UserExported {
	u, e := b.CreateUser(ctx, "Charlie", pin)
	if e != nil {
		log.Panicln(e)
	}
	printJSON("created user", u)

	return u
}

func doModifyPIN(ctx context.Context, b *sdk.Broker, userID string, pin, nonce, newPIN string) {
	err := b.ModifyPIN(ctx, userID, pin, nonce, newPIN)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("modify PIN succ")
}

func doVerifyPIN(ctx context.Context, b *sdk.Broker, userID string, pin, nonce string) {
	err := b.VerifyPIN(ctx, userID, pin, nonce)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("verify PIN succ")
}
