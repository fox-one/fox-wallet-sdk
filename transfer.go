package sdk

import (
	"context"
	"errors"

	jsoniter "github.com/json-iterator/go"
)

// TransferInput input for transfer/verify payment request
type TransferInput struct {
	AddressID  string `json:"address_id,omitempty"`
	AssetID    string `json:"asset_id,omitempty"`
	OpponentID string `json:"opponent_id,omitempty"`
	Amount     string `json:"amount,omitempty"`
	TraceID    string `json:"trace_id,omitempty"`
	Memo       string `json:"memo,omitempty"`
}

// WithdrawAddress withdraw address
type WithdrawAddress struct {
	AssetID string `json:"asset_id"`

	PublicKey string `json:"public_key,omitempty"`
	Label     string `json:"label,omitempty"`

	AccountName string `json:"account_name,omitempty"`
	AccountTag  string `json:"account_tag,omitempty"`
}

// WithdrawInput withdraw input
type WithdrawInput struct {
	WithdrawAddress

	Amount  string
	TraceID string
	Memo    string
}

// Transfer transfer to account
func (broker *Broker) Transfer(ctx context.Context, userID, pin string, input *TransferInput) (*Snapshot, error) {
	token, err := broker.SignTokenWithPIN(userID, 60, pin)
	if err != nil {
		return nil, err
	}

	return broker.BrokerHandler.Transfer(ctx, input, token)
}

// Transfer transfer to account
func (broker *BrokerHandler) Transfer(ctx context.Context, input *TransferInput, token string) (*Snapshot, error) {
	paras := map[string]interface{}{
		"asset_id":    input.AssetID,
		"opponent_id": input.OpponentID,
		"trace_id":    input.TraceID,
		"amount":      input.Amount,
		"memo":        input.Memo,
	}
	b, err := broker.Request(ctx, "POST", "/api/transfer", paras, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Snapshot *Snapshot `json:"data,omitempty"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Snapshot, nil
	}
	return nil, errorWithWalletError(&data.Error)
}

// Withdraw withdraw to address
//	address_id, opponent_id, amount, traceID, memo
func (broker *Broker) Withdraw(ctx context.Context, userID, pin string, input *WithdrawInput) (*Snapshot, error) {
	token, err := broker.SignTokenWithPIN(userID, 60, pin)
	if err != nil {
		return nil, err
	}

	return broker.BrokerHandler.Withdraw(ctx, input, token)
}

// Withdraw withdraw to address
//	address_id, opponent_id, amount, traceID, memo
func (broker *BrokerHandler) Withdraw(ctx context.Context, input *WithdrawInput, token string) (*Snapshot, error) {
	paras := map[string]interface{}{
		"asset_id": input.AssetID,
		"trace_id": input.TraceID,
		"amount":   input.Amount,
		"memo":     input.Memo,
	}
	if len(input.PublicKey) > 0 {
		paras["public_key"] = input.PublicKey
	} else {
		paras["account_name"] = input.AccountName
		if len(input.AccountTag) > 0 {
			paras["account_tag"] = input.AccountTag
		}
	}
	b, err := broker.Request(ctx, "POST", "/api/withdraw", paras, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Snapshot *Snapshot `json:"data,omitempty"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Snapshot, nil
	}
	return nil, errorWithWalletError(&data.Error)
}

// FetchWithdrawFee fetch withdraw fee
func (broker *Broker) FetchWithdrawFee(ctx context.Context, userID, pin string, input *WithdrawAddress) (*Asset, string, error) {
	token, err := broker.SignTokenWithPIN(userID, 60, pin)
	if err != nil {
		return nil, "0", err
	}

	return broker.BrokerHandler.FetchWithdrawFee(ctx, input, token)
}

// FetchWithdrawFee fetch withdraw fee
func (broker *BrokerHandler) FetchWithdrawFee(ctx context.Context, input *WithdrawAddress, token string) (*Asset, string, error) {
	paras := map[string]interface{}{
		"asset_id": input.AssetID,
	}
	if len(input.PublicKey) > 0 {
		paras["public_key"] = input.PublicKey
	} else {
		paras["account_name"] = input.AccountName
		if len(input.AccountTag) > 0 {
			paras["account_tag"] = input.AccountTag
		}
	}
	b, err := broker.Request(ctx, "POST", "/api/withdraw-fee", paras, token)
	if err != nil {
		return nil, "0", err
	}

	var data struct {
		Error
		Data *struct {
			Fee      string `json:"fee"`
			FeeAsset *Asset `json:"fee_asset"`
		} `json:"data,omitempty"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, "0", errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Data.FeeAsset, data.Data.Fee, nil
	}
	return nil, "0", errorWithWalletError(&data.Error)
}
