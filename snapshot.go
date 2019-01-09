package sdk

import (
	"context"
	"encoding/json"
	"time"
)

// Snapshot snapshot
type Snapshot struct {
	SnapshotID string `json:"snapshot_id"`
	TraceID    string `json:"trace_id"`

	UserID     string `json:"user_id"`
	AssetID    string `json:"asset_id"`
	OpponentID string `json:"opponent_id,omitempty"`

	Source string `json:"source"` // Source DEPOSIT_CONFIRMED, TRANSFER_INITIALIZED, WITHDRAWAL_INITIALIZED, WITHDRAWAL_FEE_CHARGED, WITHDRAWAL_FAILED
	Amount string `json:"amount"`
	Memo   string `json:"memo,omitempty"`

	Sender          string `json:"sender,omitempty"`
	Receiver        string `json:"receiver,omitempty"`
	TransactionHash string `json:"transaction_hash,omitempty"`

	CreatedAt int64 `json:"created_at"`

	Asset *Asset `json:"asset,omitempty"`
}

// FetchSnapshot fetch user snapshot
func (broker *Broker) FetchSnapshot(ctx context.Context, userID, traceID, snapshotID string) (*Snapshot, error) {
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
		Snapshot *Snapshot `json:"data"`
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
func (broker *Broker) FetchSnapshots(ctx context.Context, userID, assetID string, offset time.Time, order string, limit int) ([]*Snapshot, error) {
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
		Snapshots []*Snapshot `json:"data"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, requestError(err)
	}

	if data.Code == 0 {
		return data.Snapshots, nil
	}
	return nil, &data.Error
}
