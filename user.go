package sdk

import (
	"context"
	"encoding/json"
)

// User user
type User struct {
	UserID   string `json:"user_id"`
	BrokerID string `json:"broker_id,omitempty"`
	FullName string `json:"full_name,omitempty"`
}

// CreateUser create user
func (broker *Broker) CreateUser(ctx context.Context, fullname, pin string) (*User, error) {
	paras := map[string]interface{}{
		"full_name": fullname,
		"pin":       pin,
	}
	b, err := broker.Request(ctx, "", "POST", "/api/users", paras)
	if err != nil {
		return nil, requestError(err)
	}

	var data struct {
		Error
		User *User `json:"data"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, requestError(err)
	}

	if data.Code == 0 {
		return data.User, nil
	}
	return nil, &data.Error
}
