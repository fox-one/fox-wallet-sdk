package sdk

import (
	"context"
	"encoding/json"
	"time"
)

// User user
type User struct {
	UserID   string `json:"user_id"`
	BrokerID string `json:"broker_id,omitempty"`
	FullName string `json:"full_name,omitempty"`
}

// CreateUser create user
func (broker *Broker) CreateUser(ctx context.Context, fullname, pin string) (*User, error) {
	paras := map[string]interface{}{}
	if len(fullname) > 0 {
		paras["full_name"] = fullname
	}
	if len(pin) > 0 {
		pinToken, _, err := broker.PINToken(pin)
		if err != nil {
			return nil, err
		}
		paras["pin"] = pinToken
	}

	token, err := broker.SignToken("", time.Now().Unix()+60, 1)
	if err != nil {
		return nil, err
	}

	b, err := broker.Request(ctx, "POST", "/api/users", paras, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		User *User `json:"data"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	if data.Code == 0 {
		return data.User, nil
	}
	return nil, errorWithWalletError(&data.Error)
}
