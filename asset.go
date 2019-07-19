package sdk

import (
	"context"
	"errors"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
)

// Asset asset
type Asset struct {
	AssetID  string `json:"asset_id"`
	AssetKey string `json:"asset_key,omitempty"`
	ChainID  string `json:"chain_id"`

	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	IconURL string `json:"icon_url"`

	Price     decimal.Decimal `json:"price"`
	PriceUSD  decimal.Decimal `json:"price_usd"`
	PriceBTC  decimal.Decimal `json:"price_btc"`
	Change    decimal.Decimal `json:"change"`
	ChangeUSD decimal.Decimal `json:"change_usd"`
	ChangeBTC decimal.Decimal `json:"change_btc"`
}

// Chain chain
type Chain struct {
	Asset

	Fee           decimal.Decimal `json:"fee"`
	Confirmations int             `json:"confirmations"`
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
	AssetID           string          `json:"asset_id"`
	Balance           decimal.Decimal `json:"balance"`
	TransactionAmount decimal.Decimal `json:"transaction_amount"`
	TransactionCount  int64           `json:"transaction_count"`
	Asset             *Asset          `json:"asset,omitempty"`
	Chain             *Chain          `json:"chain,omitempty"`
	Address           *UserAddress    `json:"address,omitempty"`

	LastTransactionTime time.Time `json:"last_transaction_time,omitempty"`
}

// UserBalance user asset balance
type UserBalance struct {
	AssetID    string          `json:"asset_id"`
	Matched    bool            `json:"matched"`
	Balance    decimal.Decimal `json:"balance"`
	FoxBalance decimal.Decimal `json:"fox_balance"`
}

// ValidateUserBalanceResp validate user assets response
type ValidateUserBalanceResp struct {
	UserID  string                  `json:"user_id"`
	Matched bool                    `json:"matched"`
	Assets  map[string]*UserBalance `json:"assets"`
}

// FetchChains fetch chains
func (broker *BrokerHandler) FetchChains(ctx context.Context) ([]*Chain, error) {
	b, err := broker.Request(ctx, "GET", "/api/chains", nil, "")
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Chains []*Chain `json:"data"`
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
func (broker *BrokerHandler) FetchNetworkAssets(ctx context.Context) ([]*Asset, error) {
	b, err := broker.Request(ctx, "GET", "/api/network-assets", nil, "")
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

// ValidateBalances validate balances
func (broker *Broker) ValidateBalances(ctx context.Context, userIDs ...string) ([]*ValidateUserBalanceResp, error) {
	token, err := broker.SignToken("", 60)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"users": userIDs,
	}
	b, err := broker.Request(ctx, "POST", "/api/balance-validate", params, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Data []*ValidateUserBalanceResp `json:"data"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Data, nil
	}
	return nil, errorWithWalletError(&data.Error)
}
