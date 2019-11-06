package sdk

import (
	"context"
	"errors"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// Snapshot snapshot
type Snapshot struct {
	ID         int64  `json:"id"`
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

// PendingDeposit pending deposit
type PendingDeposit struct {
	Type string `json:"type"`

	TransactionID   string    `json:"transaction_id"`
	TransactionHash string    `json:"transaction_hash"`
	CreatedAt       time.Time `json:"created_at"`

	AssetID       string `json:"asset_id,omitempty"`
	ChainID       string `json:"chain_id,omitempty"`
	Amount        string `json:"amount"`
	Confirmations int    `json:"confirmations"`
	Threshold     int    `json:"threshold"`

	BrokerID    string `json:"broker_id"`
	UserID      string `json:"user_id"`
	Sender      string `json:"sender"`
	Destination string `json:"destination"`
	Tag         string `json:"tag"`

	// TODO Deprecated
	PublicKey   string `json:"public_key"`
	AccountName string `json:"account_name"`
	AccountTag  string `json:"account_tag"`
}

// FetchPendingDeposits fetch pending deposits
func (broker *Broker) FetchPendingDeposits(ctx context.Context, userIDs []string, chainID, assetID string) ([]*PendingDeposit, error) {
	token, err := broker.SignToken("", 60)
	if err != nil {
		return nil, err
	}

	return broker.BrokerHandler.FetchPendingDeposits(ctx, userIDs, chainID, assetID, token)
}

// FetchPendingDeposits fetch pending deposits
func (broker *BrokerHandler) FetchPendingDeposits(ctx context.Context, userIDs []string, chainID, assetID string, token string) ([]*PendingDeposit, error) {
	paras := map[string]interface{}{
		"user_ids": strings.Join(userIDs, ","),
		"asset_id": assetID,
		"chain_id": chainID,
	}

	b, err := broker.Request(ctx, "GET", "/api/snapshots/pending-deposits", paras, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Snapshots []*PendingDeposit `json:"data"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Snapshots, nil
	}
	return nil, errorWithWalletError(&data.Error)
}

// ExternalSnapshot external snapshot
type ExternalSnapshot struct {
	SnapshotID string `json:"snapshot_id"`
	Source     string `json:"source"`
	UserID     string `json:"user_id"`
	AssetID    string `json:"asset_id"`
	Amount     string `json:"amount"`
	CreatedAt  int64  `json:"created_at"`
}

// FetchExternalSnapshots fetch external snapshots
func (broker *Broker) FetchExternalSnapshots(ctx context.Context, userID string, from, to int64, limit int) ([]*ExternalSnapshot, error) {
	token, err := broker.SignToken("", 60)
	if err != nil {
		return nil, err
	}

	return broker.BrokerHandler.FetchExternalSnapshots(ctx, userID, from, to, limit, token)
}

// FetchExternalSnapshots fetch external snapshots
func (broker *BrokerHandler) FetchExternalSnapshots(ctx context.Context, userID string, from, to int64, limit int, token string) ([]*ExternalSnapshot, error) {
	paras := map[string]interface{}{
		"from":  from,
		"limit": limit,
	}
	if userID != "" {
		paras["user_id"] = userID
	}
	if to > from {
		paras["to"] = to
	}

	b, err := broker.Request(ctx, "GET", "/api/snapshots/external", paras, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Snapshots []*ExternalSnapshot `json:"data"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Snapshots, nil
	}
	return nil, errorWithWalletError(&data.Error)
}
