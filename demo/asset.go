package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/fox-wallet-sdk"
)

func doChains(ctx context.Context, b *sdk.Broker) []*sdk.Chain {
	chains, err := b.FetchChains(ctx)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch chains", chains)

	return chains
}

func doNetworkAssets(ctx context.Context, b *sdk.Broker) []*sdk.Asset {
	assets, err := b.FetchNetworkAssets(ctx)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch network assets", assets)

	return assets
}

func doValidateBalance(ctx context.Context, b *sdk.Broker, userID string) {
	resp, err := b.ValidateBalances(ctx, userID)
	log.Println(err)
	printJSON("validate balances", resp)
}

func doExchangeRates(ctx context.Context, b *sdk.Broker) *sdk.ExRateResp {
	exRates, err := b.FetchExRates(ctx)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch exchange rates", exRates)

	return exRates
}

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

func doScanAssets(ctx context.Context, b *sdk.Broker, assetID string, timestamp int64) {
	balances, err := b.ScanAssets(ctx, assetID, timestamp)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("scan assets", balances)
}
