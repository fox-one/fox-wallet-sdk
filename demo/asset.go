package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/fox-wallet-sdk"
)

func doAssets(ctx context.Context, b *sdk.Broker, userID string) []*sdk.UserAsset {
	assets, err := b.FetchAssets(ctx, userID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch assets", assets)

	return assets
}

func doAsset(ctx context.Context, b *sdk.Broker, userID string, assetID string) *sdk.UserAsset {
	asset, err := b.FetchAsset(ctx, userID, assetID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch asset", asset)

	return asset
}
