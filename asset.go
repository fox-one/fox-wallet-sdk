package sdk

import (
	"context"
	"errors"

	jsoniter "github.com/json-iterator/go"
)

// Asset asset
type Asset struct {
	AssetID  string `json:"asset_id"`
	AssetKey string `json:"asset_key,omitempty"`
	ChainID  string `json:"chain_id"`

	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	IconURL string `json:"icon_url"`

	Price    string `json:"price"`
	PriceUSD string `json:"price_usd"`
	Change   string `json:"change"`
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

// FetchChains fetch chains
func (broker *Broker) FetchChains(ctx context.Context) ([]*Asset, error) {
	token, err := broker.SignToken("", 60)
	if err != nil {
		return nil, err
	}
	return broker.BrokerHandler.FetchChains(ctx, token)
}

// FetchChains fetch chains
func (broker *BrokerHandler) FetchChains(ctx context.Context, token string) ([]*Asset, error) {
	b, err := broker.Request(ctx, "GET", "/api/chains", nil, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Chains []*Asset `json:"data"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Chains, nil
	}
	return nil, errorWithWalletError(&data.Error)
}

// FetchNetworkAssets fetch network assets
func (broker *Broker) FetchNetworkAssets(ctx context.Context) ([]*Asset, error) {
	token, err := broker.SignToken("", 60)
	if err != nil {
		return nil, err
	}
	return broker.BrokerHandler.FetchNetworkAssets(ctx, token)
}

// FetchNetworkAssets fetch network assets
func (broker *BrokerHandler) FetchNetworkAssets(ctx context.Context, token string) ([]*Asset, error) {
	b, err := broker.Request(ctx, "GET", "/api/network-assets", nil, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Assets []*Asset `json:"data"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Assets, nil
	}
	return nil, errorWithWalletError(&data.Error)
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
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
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
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Asset, nil
	}
	return nil, errorWithWalletError(&data.Error)
}
