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

func doSnapshots(ctx context.Context, b *sdk.Broker, userID string, assetID string) ([]*sdk.Snapshot, string) {
	snapshots, nextOffset, err := b.FetchSnapshots(ctx, userID, assetID, "", "DESC", 100)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch snapshots", snapshots)

	return snapshots, nextOffset
}

func doPendingSnapshots(ctx context.Context, b *sdk.Broker, userIDs []string, chainID, assetID string) []*sdk.PendingDeposit {
	snapshots, err := b.FetchPendingDeposits(ctx, userIDs, chainID, assetID)
	if err != nil {
		log.Panicln(err)
	}
	printJSON("fetch pending deposits", snapshots)

	return snapshots
}
