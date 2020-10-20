package sdk

import (
	"context"
	"errors"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// MixinToken sign token for mixin request
func (broker *Broker) MixinToken(ctx context.Context, userID, method, uri, body string, expire time.Duration) (string, error) {
	token, err := broker.SignToken(userID, 60)
	if err != nil {
		return "", err
	}

	return broker.BrokerHandler.MixinToken(ctx, method, uri, body, token, expire)
}

// MixinToken sign token for mixin request
func (broker *BrokerHandler) MixinToken(ctx context.Context, method, uri, body, token string, expire time.Duration) (string, error) {
	paras := map[string]interface{}{
		"method": method,
		"uri":    uri,
		"body":   body,
		"expire": expire,
	}

	b, err := broker.Request(ctx, "POST", "/api/token", paras, token)
	if err != nil {
		return "", err
	}

	var data struct {
		Error
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return "", errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Data.Token, nil
	}
	return "", errorWithWalletError(&data.Error)
}
