package sdk

import (
	"context"
	"encoding/json"
)

// ModifyPIN modify pin
func (broker *Broker) ModifyPIN(ctx context.Context, userID, pin, nonce, newPIN string) *Error {
	paras := map[string]interface{}{
		"pin": newPIN,
	}
	b, err := broker.RequestWithPIN(ctx, userID, pin, nonce, "PUT", "/api/pin", paras)
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
func (broker *Broker) VerifyPIN(ctx context.Context, userID, pin, nonce string) *Error {
	b, err := broker.RequestWithPIN(ctx, userID, pin, nonce, "POST", "/api/pin/verify", nil)
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
