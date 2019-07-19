package sdk

import (
	"context"
	"errors"

	jsoniter "github.com/json-iterator/go"
)

// FetchAddresses fetch addresses
func (broker *Broker) FetchAddresses(ctx context.Context, userID, assetID string) ([]*WithdrawAddress, error) {
	token, err := broker.SignToken(userID, 60)
	if err != nil {
		return nil, err
	}
	return broker.BrokerHandler.FetchAddresses(ctx, assetID, token)
}

// FetchAddresses fetch user addresses
func (broker *BrokerHandler) FetchAddresses(ctx context.Context, assetID, token string) ([]*WithdrawAddress, error) {
	paras := map[string]interface{}{
		"asset_id": assetID,
	}
	b, err := broker.Request(ctx, "GET", "/api/wallet/addresses", paras, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Addresses []*WithdrawAddress `json:"data"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Addresses, nil
	}
	return nil, errorWithWalletError(&data.Error)
}

// UpsertAddress add/update user address
func (broker *Broker) UpsertAddress(ctx context.Context, userID, pin string, addr *WithdrawAddress) (*WithdrawAddress, error) {
	token, err := broker.SignTokenWithPIN(userID, 60, pin)
	if err != nil {
		return nil, err
	}
	return broker.BrokerHandler.UpsertAddress(ctx, addr, token)
}

// UpsertAddress add/update user address
func (broker *BrokerHandler) UpsertAddress(ctx context.Context, addr *WithdrawAddress, token string) (*WithdrawAddress, error) {
	paras := map[string]interface{}{
		"asset_id":     addr.AssetID,
		"public_key":   addr.PublicKey,
		"label":        addr.Label,
		"account_name": addr.AccountName,
		"account_tag":  addr.AccountTag,
	}
	b, err := broker.Request(ctx, "POST", "/api/wallet/address", paras, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Address *WithdrawAddress `json:"data"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Address, nil
	}
	return nil, errorWithWalletError(&data.Error)
}

// DeleteAddress delete user address
func (broker *Broker) DeleteAddress(ctx context.Context, userID, addressID, pin string) error {
	token, err := broker.SignTokenWithPIN(userID, 60, pin)
	if err != nil {
		return err
	}
	return broker.BrokerHandler.DeleteAddress(ctx, addressID, token)
}

// DeleteAddress delete user address
func (broker *BrokerHandler) DeleteAddress(ctx context.Context, addressID, token string) error {
	b, err := broker.Request(ctx, "DELETE", "/api/wallet/address/"+addressID, nil, token)
	if err != nil {
		return err
	}

	var data struct {
		Error
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return errors.New(string(b))
	}

	if data.Code == 0 {
		return nil
	}
	return errorWithWalletError(&data.Error)
}
