package sdk

import (
	"encoding/base64"

	"github.com/fox-one/fox-wallet/models"
	"github.com/satori/go.uuid"
)

// Broker broker
type Broker struct {
	*models.Broker
	apiBase string
}

// NewBroker create broker
func NewBroker(apiBase, brokerID, secret string, pinSecret ...string) (*Broker, error) {
	b := &models.Broker{}

	if len(secret) > 0 {
		s, err := base64.StdEncoding.DecodeString(secret)
		if err != nil {
			return nil, err
		}
		b.Secret = s
	}

	b.BrokerID = brokerID
	if len(pinSecret) > 0 {
		s, err := base64.StdEncoding.DecodeString(pinSecret[0])
		if err != nil {
			return nil, err
		}
		b.PINSecret = s
	}
	return &Broker{
		Broker:  b,
		apiBase: apiBase,
	}, nil
}

// PINToken generate pin token, with nonce
func (b *Broker) PINToken(pin string, nonce ...string) (string, string, error) {
	n := uuid.Must(uuid.NewV4()).String()
	if len(nonce) > 0 {
		n = nonce[0]
	}
	p, err := b.EncryptPIN(pin, n)
	if err != nil {
		return "", "", err
	}
	return p, n, nil
}
