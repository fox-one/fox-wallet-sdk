package main

import (
	"context"
	"log"
	"time"

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

func doSnapshots(ctx context.Context, b *sdk.Broker, userID string, assetID string) []*sdk.Snapshot {
	snapshots, err := b.FetchSnapshots(ctx, userID, assetID, time.Time{}, "DESC", 100)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch snapshots", snapshots)

	return snapshots
}
