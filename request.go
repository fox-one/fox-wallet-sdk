package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fox-one/mixin-sdk/utils"
)

// Request send request
//  jwt token, {"i":"broker-id","u":"user-id","n":"nonce","e":123,"nr":2,"pt":"pin token"}
func (broker *BrokerHandler) Request(ctx context.Context, method, uri string, params map[string]interface{}, token string, headers ...string) ([]byte, error) {
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

	headers = append(headers, "Content-Type", "application/json", "Authorization", "Bearer "+token)

	url := broker.apiBase + uri
	req, err := utils.NewRequest(url, method, string(body), headers...)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	resp, _ := utils.DoRequest(req)
	return utils.ReadResponse(resp)
}
