package sdk

import (
	"context"
	"encoding/json"

	"github.com/fox-one/fox-wallet/models"
)

// CreateUser create user
func (broker *Broker) CreateUser(ctx context.Context, fullname, pin string) (*models.UserExported, *Error) {
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
		User *models.UserExported `json:"data"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, requestError(err)
	}

	if data.Code == 0 {
		return data.User, nil
	}
	return nil, &data.Error
}
