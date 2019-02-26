package sdk

import (
	"context"
	"errors"

	jsoniter "github.com/json-iterator/go"
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
	token, err := broker.SignToken(userID, 60)
	if err != nil {
		return nil, err
	}

	return broker.BrokerHandler.FetchSnapshot(ctx, traceID, snapshotID, token)
}

// FetchSnapshot fetch a snapshot
func (broker *BrokerHandler) FetchSnapshot(ctx context.Context, traceID, snapshotID, token string) (*Snapshot, error) {
	paras := map[string]interface{}{}
	if len(traceID) > 0 {
		paras["trace_id"] = traceID
	} else {
		paras["snapshot_id"] = snapshotID
	}
	b, err := broker.Request(ctx, "GET", "/api/snapshot", paras, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Snapshot *Snapshot `json:"data"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Snapshot, nil
	}
	return nil, errorWithWalletError(&data.Error)
}

// FetchSnapshots fetch snapshots
func (broker *Broker) FetchSnapshots(ctx context.Context, userID, assetID, offset, order string, limit int) ([]*Snapshot, string, error) {
	token, err := broker.SignToken(userID, 60)
	if err != nil {
		return nil, offset, err
	}

	return broker.BrokerHandler.FetchSnapshots(ctx, assetID, offset, order, limit, token)
}

// FetchSnapshots fetch snapshots
func (broker *BrokerHandler) FetchSnapshots(ctx context.Context, assetID, offset, order string, limit int, token string) ([]*Snapshot, string, error) {
	paras := map[string]interface{}{
		"order":    order,
		"asset_id": assetID,
		"offset":   offset,
	}

	if limit > 0 {
		paras["limit"] = limit
	}

	b, err := broker.Request(ctx, "POST", "/api/snapshots", paras, token)
	if err != nil {
		return nil, offset, err
	}

	var data struct {
		Error
		Snapshots  []*Snapshot `json:"data"`
		NextOffset string      `json:"next_offset"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, offset, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Snapshots, data.NextOffset, nil
	}
	return nil, offset, errorWithWalletError(&data.Error)
}
