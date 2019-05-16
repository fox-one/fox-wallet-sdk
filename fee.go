package sdk

import (
	"context"
	"errors"

	jsoniter "github.com/json-iterator/go"
)

// FetchWithdrawFee fetch withdraw fee
func (broker *BrokerHandler) FetchWithdrawFee(ctx context.Context, input *WithdrawAddress) (*Asset, string, error) {
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
	b, err := broker.Request(ctx, "POST", "/api/withdraw-fee", paras, "")
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
