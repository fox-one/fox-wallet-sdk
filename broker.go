package sdk

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fox-one/mixin-sdk/utils"
	"github.com/satori/go.uuid"
)

const (
	encryptPrefix = "data:"
)

// Broker broker
type Broker struct {
	brokerID  string
	secret    []byte
	pinSecret []byte
	apiBase   string
}

// NewBroker create broker
func NewBroker(apiBase, brokerID string, secret []byte, pinSecret ...[]byte) *Broker {
	b := &Broker{
		brokerID: brokerID,
		apiBase:  apiBase,
		secret:   secret,
	}
	if len(pinSecret) > 0 && len(pinSecret[0]) > 0 {
		b.pinSecret = pinSecret[0]
	}
	return b
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
	if len(pin) != 6 {
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

// Sign sign a request token
//  jwt token, {"i":"broker-id","u":"user-id","n":"nonce","e":123,"nr":2,"pt":"pin token"}
func (b *Broker) Sign(params map[string]interface{}, expire int64, nonce ...string) (string, error) {
	if len(b.secret) == 0 {
		return "", fmt.Errorf("empty secret")
	}

	var n = uuid.Must(uuid.NewV4()).String()
	if len(nonce) > 0 && len(nonce[0]) > 0 {
		n = nonce[0]
	}

	params["e"] = expire
	params["n"] = n
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(params))

	return token.SignedString(b.secret)
}
