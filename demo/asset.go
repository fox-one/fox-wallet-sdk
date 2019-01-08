package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/fox-wallet-sdk"
	"github.com/fox-one/fox-wallet/models"
)

func doAssets(ctx context.Context, b *sdk.Broker, userID string) []*models.UserAssetExported {
	assets, err := b.FetchAssets(ctx, userID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch assets", assets)

	return assets
}

func doAsset(ctx context.Context, b *sdk.Broker, userID string, assetID string) *models.UserAssetExported {
	asset, err := b.FetchAsset(ctx, userID, assetID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch asset", asset)

	return asset
}
