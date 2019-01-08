package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fox-one/mixin-sdk/utils"
)

// GenerateToken generate jwt token
func (broker *Broker) GenerateToken(userID, pinToken, nonce string) (string, error) {
	jwtMap := map[string]interface{}{
		"i": broker.BrokerID,
	}
	if len(userID) > 0 {
		jwtMap["u"] = userID
	}

	if len(pinToken) > 0 {
		jwtMap["pt"] = pinToken
	}

	token, err := broker.Sign(jwtMap, time.Now().Unix()+60, nonce)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Request sign and request
//  jwt token, {"i":"broker-id","u":"user-id","n":"nonce","e":123,"nr":2,"pt":"pin token"}
func (broker *Broker) Request(ctx context.Context, userID string, method, uri string, params map[string]interface{}, headers ...string) ([]byte, error) {
	return broker.RequestWithPIN(ctx, userID, "", "", method, uri, params, headers...)
}

// RequestWithPIN sign and request
//  jwt token, {"i":"broker-id","u":"user-id","n":"nonce","e":123,"nr":2,"pt":"pin token"}
func (broker *Broker) RequestWithPIN(ctx context.Context, userID, pinToken, nonce string, method, uri string, params map[string]interface{}, headers ...string) ([]byte, error) {
	body := []byte{}
	switch method {
	case "GET", "DELETE":
		arr := make([]string, 0, len(params))
		for k, v := range params {
			arr = append(arr, fmt.Sprintf("%s=%v", k, v))
		}
		uri = uri + "?" + strings.Join(arr, "&")

	case "POST", "PUT":
		b, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		body = b
	}

	token, err := broker.GenerateToken(userID, pinToken, nonce)
	if err != nil {
		return nil, err
	}

	headers = append(headers, "Content-Type", "application/json",
		"Authorization", "Bearer "+token)

	url := broker.apiBase + uri
	req, err := utils.NewRequest(url, method, string(body), headers...)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	resp, _ := utils.DoRequest(req)
	return utils.ReadResponse(resp)
}
