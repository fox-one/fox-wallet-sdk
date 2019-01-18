package sdk

import (
	"context"
	"encoding/json"
)

// Asset asset
type Asset struct {
	AssetID  string `json:"asset_id"`
	AssetKey string `json:"asset_key,omitempty"`
	ChainID  string `json:"chain_id"`

	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	IconURL string `json:"icon_url"`
}

// UserAddress user address
type UserAddress struct {
	UserID  string `json:"user_id"`
	ChainID string `json:"chain_id"`

	PublicKey   string `json:"public_key"`
	AccountName string `json:"account_name"`
	AccountTag  string `json:"account_tag"`

	Confirmations  int     `json:"confirmations"`
	Capitalization float64 `json:"capitalization"`
}

// UserAsset user asset
type UserAsset struct {
	AssetID           string       `json:"asset_id"`
	Balance           string       `json:"balance"`
	TransactionAmount string       `json:"transaction_amount"`
	TransactionCount  int64        `json:"transaction_count"`
	Asset             *Asset       `json:"asset,omitempty"`
	Address           *UserAddress `json:"address,omitempty"`
}

// FetchAssets fetch assets
func (broker *Broker) FetchAssets(ctx context.Context, userID string) ([]*UserAsset, error) {
	token, err := broker.SignToken(userID, 60)
	if err != nil {
		return nil, err
	}
	return broker.BrokerHandler.FetchAssets(ctx, token)
}

// FetchAssets fetch user assets
func (broker *BrokerHandler) FetchAssets(ctx context.Context, token string) ([]*UserAsset, error) {
	b, err := broker.Request(ctx, "GET", "/api/assets", nil, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Assets []*UserAsset `json:"data"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	if data.Code == 0 {
		return data.Assets, nil
	}
	return nil, errorWithWalletError(&data.Error)
}

// FetchAsset fetch asset
func (broker *Broker) FetchAsset(ctx context.Context, userID, assetID string) (*UserAsset, error) {
	token, err := broker.SignToken(userID, 60)
	if err != nil {
		return nil, err
	}
	return broker.BrokerHandler.FetchAsset(ctx, assetID, token)
}

// FetchAsset fetch user asset
func (broker *BrokerHandler) FetchAsset(ctx context.Context, assetID, token string) (*UserAsset, error) {
	b, err := broker.Request(ctx, "GET", "/api/asset/"+assetID, nil, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Asset *UserAsset `json:"data"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	if data.Code == 0 {
		return data.Asset, nil
	}
	return nil, errorWithWalletError(&data.Error)
}
