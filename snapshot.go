package sdk

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fox-one/fox-wallet/models"
)

// FetchSnapshot fetch user snapshot
func (broker *Broker) FetchSnapshot(ctx context.Context, userID, traceID, snapshotID string) (*models.Snapshot, *Error) {
	paras := map[string]interface{}{}
	if len(traceID) > 0 {
		paras["trace_id"] = traceID
	} else {
		paras["snapshot_id"] = snapshotID
	}
	b, err := broker.Request(ctx, userID, "GET", "/api/snapshot", paras)
	if err != nil {
		return nil, requestError(err)
	}

	var data struct {
		Error
		Snapshot *models.Snapshot `json:"data"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, requestError(err)
	}

	if data.Code == 0 {
		return data.Snapshot, nil
	}
	return nil, &data.Error
}

// FetchSnapshots fetch snapshots
func (broker *Broker) FetchSnapshots(ctx context.Context, userID, assetID string, offset time.Time, order string, limit int) ([]*models.Snapshot, *Error) {
	paras := map[string]interface{}{
		"order":    order,
		"asset_id": assetID,
	}

	if !offset.IsZero() {
		paras["offset"] = offset.UnixNano()
	}

	if limit > 0 {
		paras["limit"] = limit
	}

	b, err := broker.Request(ctx, userID, "POST", "/api/snapshots", paras)
	if err != nil {
		return nil, requestError(err)
	}

	var data struct {
		Error
		Snapshots []*models.Snapshot `json:"data"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, requestError(err)
	}

	if data.Code == 0 {
		return data.Snapshots, nil
	}
	return nil, &data.Error
}
