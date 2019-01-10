package sdk

import (
	"context"
	"encoding/json"
)

// ModifyPIN modify pin
func (broker *Broker) ModifyPIN(ctx context.Context, userID, pin, newPIN string) error {
	pinToken, nonce, err := broker.PINToken(pin)
	if err != nil {
		return requestError(err)
	}

	newPINToken, _, err := broker.PINToken(newPIN)
	if err != nil {
		return requestError(err)
	}

	paras := map[string]interface{}{
		"pin": newPINToken,
	}
	b, err := broker.RequestWithPIN(ctx, userID, pinToken, nonce, "PUT", "/api/pin", paras)
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
	pinToken, nonce, err := broker.PINToken(pin)
	if err != nil {
		return requestError(err)
	}

	b, err := broker.RequestWithPIN(ctx, userID, pinToken, nonce, "POST", "/api/pin/verify", nil)
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
