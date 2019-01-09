package main

import (
	"context"
	"log"

	sdk "github.com/fox-one/fox-wallet-sdk"
)

func doSnapshot(ctx context.Context, b *sdk.Broker, userID, traceID, snapshotID string) *sdk.Snapshot {
	snapshot, err := b.FetchSnapshot(ctx, userID, traceID, snapshotID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch snapshot", snapshot)

	return snapshot
}

func doSnapshots(ctx context.Context, b *sdk.Broker, userID string, assetID string) ([]*sdk.Snapshot, int64) {
	snapshots, nextOffset, err := b.FetchSnapshots(ctx, userID, assetID, 0, "DESC", 100)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch snapshots", snapshots)

	return snapshots, nextOffset
}
