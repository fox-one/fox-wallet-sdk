package sdk

import (
	"context"
	"errors"

	jsoniter "github.com/json-iterator/go"
)

// ExRate exchange rate
type ExRate struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Price  string `json:"price"`
	Change string `json:"change"`
}

// ExRateResp exchange rate response
type ExRateResp struct {
	ExRates    []*ExRate         `json:"cnyTickers"`
	Currencies map[string]string `json:"currencies"`
}

// FetchExRatesRaw fetch exchange rates
func (broker *BrokerHandler) FetchExRatesRaw(ctx context.Context) ([]byte, error) {
	return broker.Request(ctx, "GET", "/api/exchange-rates", nil, "")
}

// FetchExRates fetch exchange rates
func (broker *BrokerHandler) FetchExRates(ctx context.Context) (*ExRateResp, error) {
	b, err := broker.FetchExRatesRaw(ctx)
	if err != nil {
		return nil, err
	}

	var data struct {
		Error
		Data *ExRateResp `json:"data"`
	}
	if err := jsoniter.Unmarshal(b, &data); err != nil {
		return nil, errors.New(string(b))
	}

	if data.Code == 0 {
		return data.Data, nil
	}
	return nil, errorWithWalletError(&data.Error)
}
