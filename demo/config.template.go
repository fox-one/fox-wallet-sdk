// +build template

package main

const (
	// ClientID mixin user id
	ClientID = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	// PIN pin
	PIN = "000000"
	// SessionID session id
	SessionID = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	// PINToken pin token
	PINToken = "xxxxx"
	// SessionKey private key in pem
	SessionKey = `-----BEGIN RSA PRIVATE KEY-----
xxx
-----END RSA PRIVATE KEY-----`

	apiBase         = "http://localhost:8081"
	brokerID        = "xxx"
	brokerSecret    = "xxx"
	brokerPINSecret = "xxx"

	default_pin = "111111"
)
