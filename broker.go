package sdk

import (
	"encoding/base64"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fox-one/mixin-sdk/utils"
	uuid "github.com/gofrs/uuid"
)

const (
	encryptPrefix = "data:"
)

// BrokerHandler broker request handler
type BrokerHandler struct {
	apiBase string
}

// Broker broker sign jwt token
type Broker struct {
	*BrokerHandler
	brokerID  string
	secret    []byte
	pinSecret []byte
}

// NewBroker new broker
func NewBroker(apiBase string, brokerID, secret, pinSecret string) *Broker {
	b := &Broker{
		BrokerHandler: NewBrokerHandler(apiBase),
		brokerID:      brokerID,
	}

	s, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		panic("secret is not base64")
	}
	b.secret = s

	if len(pinSecret) > 0 {
		ps, err := base64.StdEncoding.DecodeString(pinSecret)
		if err != nil {
			panic("pin secret is not base64")
		}
		b.pinSecret = ps
	}

	return b
}

// NewBrokerHandler create broker
func NewBrokerHandler(apiBase string) *BrokerHandler {
	return &BrokerHandler{
		apiBase: apiBase,
	}
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

// EncryptPIN encrypt pin
func (b *Broker) EncryptPIN(pin, nonce string) (string, error) {
	if len(pin) == 0 {
		return "", fmt.Errorf("invalid pin")
	}

	if len(b.pinSecret) == 0 {
		return "", fmt.Errorf("broker pin secret not set")
	}

	if len(nonce) > 10 {
		nonce = nonce[:10]
	}

	return utils.Encrypt([]byte(encryptPrefix+pin+";"+nonce), b.pinSecret, []byte(b.brokerID))
}

// SignTokenWithPIN sign a request token with pin
//  jwt token, {"i":"broker-id","u":"user-id","n":"nonce","e":123,"nr":2,"pt":"pin token"}
func (b *Broker) SignTokenWithPIN(userID string, expire int64, pin string, nonceRepeats ...int) (string, error) {
	if len(b.secret) == 0 {
		return "", fmt.Errorf("empty secret")
	}

	params := map[string]interface{}{
		"i": b.brokerID,
		"e": time.Now().Unix() + expire,
		"n": uuid.Must(uuid.NewV4()).String(),
	}

	if len(userID) > 0 {
		params["u"] = userID
	}

	if len(pin) > 0 {
		pinToken, nonce, err := b.PINToken(pin)
		if err != nil {
			return "", err
		}
		params["n"] = nonce
		params["pt"] = pinToken
	}

	if len(nonceRepeats) > 0 && nonceRepeats[0] > 0 {
		params["nr"] = nonceRepeats[0]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(params))
	return token.SignedString(b.secret)
}

// SignToken sign a request token
//  jwt token, {"i":"broker-id","u":"user-id","n":"nonce","e":123,"nr":2,"pt":"pin token"}
func (b *Broker) SignToken(userID string, expire int64, nonceRepeat ...int) (string, error) {
	return b.SignTokenWithPIN(userID, expire, "", nonceRepeat...)
}
