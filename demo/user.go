package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/fox-wallet-sdk"
)

func doCreateUser(ctx context.Context, b *sdk.Broker, pin string) *sdk.User {
	u, e := b.CreateUser(ctx, "", pin, "")
	if e != nil {
		log.Panicln(e)
	}
	printJSON("created user", u)

	return u
}

func doModifyUser(ctx context.Context, b *sdk.Broker, userID string) *sdk.User {
	fullname := "test"
	avatar := "https://images.mixin.one/7y_5w1rpGkFs_65XTvvJ37ZqsOy0t2D-qdSUVfw9LykUyE4gAWH3OMU5kxGVsJSdq3oVXx6Qz1KuJYRiOi6_5fQ=s256"
	u, e := b.ModifyUser(ctx, userID, fullname, avatar)
	if e != nil {
		log.Panicln(e)
	}

	printJSON("fetch user", u)

	return u
}

func doFetchUser(ctx context.Context, b *sdk.Broker, userID string) *sdk.User {
	u, e := b.FetchUser(ctx, userID)
	if e != nil {
		log.Panicln(e)
	}
	printJSON("fetch user", u)

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
