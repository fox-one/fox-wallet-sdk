package sdk

import (
	"context"
	"encoding/json"
	"time"
)

// ModifyPIN modify pin
func (broker *Broker) ModifyPIN(ctx context.Context, userID, pin, newPIN string) error {
	token, err := broker.SignTokenWithPIN(userID, time.Now().Unix()+60, pin)
	if err != nil {
		return requestError(err)
	}

	newPINToken, _, err := broker.PINToken(newPIN)
	if err != nil {
		return requestError(err)
	}

	return broker.BrokerHandler.ModifyPIN(ctx, newPINToken, token)
}

// ModifyPIN modify pin
func (broker *BrokerHandler) ModifyPIN(ctx context.Context, newPINToken, token string) error {
	paras := map[string]interface{}{
		"pin": newPINToken,
	}
	b, err := broker.Request(ctx, "PUT", "/api/pin", paras, token)
	if err != nil {
		return requestError(err)
	}

	var data struct {
		Error
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return requestError(err)
	}

	if data.Code == 0 {
		return nil
	}
	return &data.Error
}

// VerifyPIN verify pin
func (broker *Broker) VerifyPIN(ctx context.Context, userID, pin string) error {
	token, err := broker.SignTokenWithPIN(userID, time.Now().Unix()+60, pin)
	if err != nil {
		return requestError(err)
	}

	return broker.BrokerHandler.VerifyPIN(ctx, token)
}

// VerifyPIN verify pin
func (broker *BrokerHandler) VerifyPIN(ctx context.Context, token string) error {
	b, err := broker.Request(ctx, "POST", "/api/pin/verify", nil, token)
	if err != nil {
		return requestError(err)
	}

	var data struct {
		Error
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return requestError(err)
	}

	if data.Code == 0 {
		return nil
	}
	return &data.Error
}
