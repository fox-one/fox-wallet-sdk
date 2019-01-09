package sdk

import (
	"context"
	"encoding/json"

	"github.com/fox-one/mixin-sdk/mixin"
)

// WithdrawInput withdraw input
type WithdrawInput struct {
	mixin.WithdrawAddress

	Amount  string
	TraceID string
	Memo    string
}

// Transfer transfer to account
func (broker *Broker) Transfer(ctx context.Context, userID, pin, nonce string, input *mixin.TransferInput) (*Snapshot, error) {
	paras := map[string]interface{}{
		"asset_id":    input.AssetID,
		"opponent_id": input.OpponentID,
		"trace_id":    input.TraceID,
		"amount":      input.Amount,
		"memo":        input.Memo,
		"pin":         pin,
	}
	b, err := broker.RequestWithPIN(ctx, userID, pin, nonce, "POST", "/api/transfer", paras)
	if err != nil {
		return nil, requestError(err)
	}

	var data struct {
		Error
		Snapshot *Snapshot `json:"data,omitempty"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, requestError(err)
	}

	if data.Code == 0 {
		return data.Snapshot, nil
	}
	return nil, &data.Error
}

// Withdraw withdraw to address
//	address_id, opponent_id, amount, traceID, memo
func (broker *Broker) Withdraw(ctx context.Context, userID, pin, nonce string, input *WithdrawInput) (*Snapshot, error) {
	paras := map[string]interface{}{
		"asset_id": input.AssetID,
		"trace_id": input.TraceID,
		"amount":   input.Amount,
		"memo":     input.Memo,
		"pin":      pin,
	}
	if len(input.PublicKey) > 0 {
		paras["public_key"] = input.PublicKey
	} else {
		paras["account_name"] = input.AccountName
		if len(input.AccountTag) > 0 {
			paras["account_tag"] = input.AccountTag
		}
	}
	b, err := broker.RequestWithPIN(ctx, userID, pin, nonce, "POST", "/api/withdraw", paras)
	if err != nil {
		return nil, requestError(err)
	}

	var data struct {
		Error
		Snapshot *Snapshot `json:"data,omitempty"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, requestError(err)
	}

	if data.Code == 0 {
		return data.Snapshot, nil
	}
	return nil, &data.Error
}

// FetchWithdrawFee fetch withdraw fee
func (broker *Broker) FetchWithdrawFee(ctx context.Context, userID, pin, nonce string, input *mixin.WithdrawAddress) (string, error) {
	paras := map[string]interface{}{
		"asset_id": input.AssetID,
		"pin":      pin,
	}
	if len(input.PublicKey) > 0 {
		paras["public_key"] = input.PublicKey
	} else {
		paras["account_name"] = input.AccountName
		if len(input.AccountTag) > 0 {
			paras["account_tag"] = input.AccountTag
		}
	}
	b, err := broker.RequestWithPIN(ctx, userID, pin, nonce, "POST", "/api/withdraw-fee", paras)
	if err != nil {
		return "0", requestError(err)
	}

	var data struct {
		Error
		Data *struct {
			Fee string `json:"fee"`
		} `json:"data,omitempty"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return "0", requestError(err)
	}

	if data.Code == 0 {
		return data.Data.Fee, nil
	}
	return "0", &data.Error
}
