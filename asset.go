package sdk

import (
	"context"
	"encoding/json"

	"github.com/fox-one/fox-wallet/models"
)

// FetchAssets fetch user assets
func (broker *Broker) FetchAssets(ctx context.Context, userID string) ([]*models.UserAssetExported, *Error) {
	b, err := broker.Request(ctx, userID, "GET", "/api/assets", nil)
	if err != nil {
		return nil, requestError(err)
	}

	var data struct {
		Error
		Assets []*models.UserAssetExported `json:"data"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, requestError(err)
	}

	if data.Code == 0 {
		return data.Assets, nil
	}
	return nil, &data.Error
}

// FetchAsset fetch user asset
func (broker *Broker) FetchAsset(ctx context.Context, userID, assetID string) (*models.UserAssetExported, *Error) {
	b, err := broker.Request(ctx, userID, "GET", "/api/asset/"+assetID, nil)
	if err != nil {
		return nil, requestError(err)
	}

	var data struct {
		Error
		Asset *models.UserAssetExported `json:"data"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, requestError(err)
	}

	if data.Code == 0 {
		return data.Asset, nil
	}
	return nil, &data.Error
}
