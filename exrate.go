package sdk

import (
	"context"
	"encoding/json"
)

// ExRate exchange rate
type ExRate struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Price  string `json:"price"`
	Change string `json:"change"`
}

// FetchExRates fetch exchange rates, if currencies is empty, will return all
func (broker *BrokerHandler) FetchExRates(ctx context.Context, currencies ...string) ([]*ExRate, error) {
	b, err := broker.Request(ctx, "GET", "/api/exchange-rates", nil, "")
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Data *struct {
			ExRates []*ExRate `json:"cnyTickers"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	if data.Code == 0 {
		return data.Data.ExRates, nil
	}
	return nil, errorWithWalletError(&data.Error)
}
