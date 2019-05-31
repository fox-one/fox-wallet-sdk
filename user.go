package sdk

import (
	"context"
	"errors"

	jsoniter "github.com/json-iterator/go"
)

// User user
type User struct {
	UserID   string `json:"user_id"`
	BrokerID string `json:"broker_id,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Avatar   string `json:"avatar_url,omitempty"`
	Inside   bool   `json:"inside"`
}

// CreateUser create user
func (broker *Broker) CreateUser(ctx context.Context, fullname, avatar, pin string) (*User, error) {
	paras := map[string]interface{}{}
	if fullname != "" {
		paras["full_name"] = fullname
	}
	if avatar != "" {
		paras["avatar"] = avatar
	}
	if pin != "" {
		pinToken, _, err := broker.PINToken(pin)
		if err != nil {
			return nil, err
		}
		paras["pin"] = pinToken
	}

	token, err := broker.SignToken("", 60, 1)
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
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.User, nil
	}
	return nil, errorWithWalletError(&data.Error)
}

// ModifyUser modify user
func (broker *Broker) ModifyUser(ctx context.Context, userID, fullname, avatar string) (*User, error) {
	token, err := broker.SignToken(userID, 60)
	if err != nil {
		return nil, err
	}
	return broker.BrokerHandler.ModifyUser(ctx, fullname, avatar, token)
}

// ModifyUser modify user
func (broker *BrokerHandler) ModifyUser(ctx context.Context, fullname, avatar, token string) (*User, error) {
	paras := map[string]interface{}{}
	if fullname != "" {
		paras["full_name"] = fullname
	}
	if avatar != "" {
		paras["avatar"] = avatar
	}
	b, err := broker.Request(ctx, "PUT", "/api/users", paras, token)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		User *User `json:"data"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.User, nil
	}
	return nil, errorWithWalletError(&data.Error)
}

// FetchUser fetch user
func (broker *BrokerHandler) FetchUser(ctx context.Context, userID string) (*User, error) {
	b, err := broker.Request(ctx, "GET", "/api/users/"+userID, nil, "")
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		User *User `json:"data"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.User, nil
	}
	return nil, errorWithWalletError(&data.Error)
}
