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

	// OpponentKey used for raw transaction
	OpponentKey string `json:"opponent_key,omitempty"`

	OpponentMultisig struct {
		Receivers []string `json:"receivers,omitempty"`
		Threshold uint8    `json:"threshold,omitempty"`
	} `json:"opponent_multisig,omitempty"`
}

// WithdrawAddress withdraw address
type WithdrawAddress struct {
	AddressID string `json:"address_id"`
	AssetID   string `json:"asset_id"`
	Label     string `json:"label,omitempty"`

	Destination string `json:"destination,omitempty"`
	Tag         string `json:"tag,omitempty"`

	// TODO Deprecated
	PublicKey   string `json:"public_key,omitempty"`
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
		"asset_id": input.AssetID,
		"trace_id": input.TraceID,
		"amount":   input.Amount,
		"memo":     input.Memo,
	}

	if input.OpponentID != "" {
		paras["opponent_id"] = input.OpponentID
	} else if input.OpponentKey != "" {
		paras["opponent_key"] = input.OpponentKey
	} else {
		paras["opponent_multisig"] = map[string]interface{}{
			"receivers": input.OpponentMultisig.Receivers,
			"threshold": input.OpponentMultisig.Threshold,
		}
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
	if input.Destination == "" {
		if input.PublicKey != "" {
			input.Destination = input.PublicKey
			input.Tag = ""
		} else {
			input.Destination = input.AccountName
			input.Tag = input.AccountTag
		}
	}

	paras := map[string]interface{}{
		"asset_id": input.AssetID,
		"trace_id": input.TraceID,
		"amount":   input.Amount,
		"memo":     input.Memo,
		"label":    input.Label,
	}
	if input.AddressID != "" {
		paras["address_id"] = input.AddressID
	} else {
		paras["destination"] = input.Destination
		paras["tag"] = input.Tag
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
