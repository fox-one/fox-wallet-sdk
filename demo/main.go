package main

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"time"

	sdk "github.com/fox-one/fox-wallet-sdk"
	"github.com/fox-one/mixin-sdk/mixin"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

func printJSON(prefix string, item interface{}) {
	msg, err := jsoniter.MarshalToString(item)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(prefix, msg)
}

func ensurePinToken(b *sdk.Broker, pin string, nonce ...string) (string, string) {
	pinToken, n, err := b.PINToken(pin, nonce...)
	if err != nil {
		log.Panicln(err)
	}

	return pinToken, n
}

func main() {
	log.SetLevel(log.DebugLevel)

	dapp := &mixin.User{
		UserID:    ClientID,
		SessionID: SessionID,
		PINToken:  PINToken,
	}

	block, _ := pem.Decode([]byte(SessionKey))
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Panicln(err)
	}
	dapp.SetPrivateKey(privateKey)

	ctx := context.TODO()
	var secret []byte
	var pinSecret []byte
	if s, err := base64.StdEncoding.DecodeString(brokerSecret); err != nil {
		log.Panicln(err)
	} else {
		secret = s
	}

	if s, err := base64.StdEncoding.DecodeString(brokerPINSecret); err != nil {
		log.Panicln(err)
	} else {
		pinSecret = s
	}

	b := sdk.NewBroker(apiBase, brokerID, secret, pinSecret)

	assetID := "965e5c6e-434c-3fa9-b780-c50f43cd955c"
	publicKey := "0xe20FE5C04Fa6b044b720F8CA019Cd896881ED13B"
	var pinToken, nonce string

	// userID := "cedf2dee-a142-3893-86ac-66127df36b55"
	tmpPIN := "123456"
	pinToken, nonce = ensurePinToken(b, tmpPIN)
	userID := doCreateUser(ctx, b, pinToken).UserID

	pinToken, nonce = ensurePinToken(b, tmpPIN)
	doVerifyPIN(ctx, b, userID, pinToken, nonce)

	pinToken, nonce = ensurePinToken(b, tmpPIN)
	newPIN, _ := ensurePinToken(b, default_pin)
	doModifyPIN(ctx, b, userID, pinToken, nonce, newPIN)

	doAssets(ctx, b, userID)
	doAsset(ctx, b, userID, assetID)
	time.Sleep(1 * time.Second)
	if len(doAssets(ctx, b, userID)) == 0 {
		log.Panicln("should have at least one asset")
	}

	pinToken, nonce = ensurePinToken(b, default_pin)
	doWithdrawFee(ctx, b, userID, assetID, publicKey, pinToken, nonce)

	pinToken, nonce = ensurePinToken(b, default_pin)
	snapshot := doTransfer(ctx, b, dapp, userID, assetID, "0.1", pinToken, nonce)

	pinToken, nonce = ensurePinToken(b, default_pin)
	doWithdraw(ctx, b, dapp, userID, assetID, publicKey, "0.1", pinToken, nonce)

	time.Sleep(10 * time.Second)
	doSnapshot(ctx, b, userID, "", snapshot.SnapshotID)
	doSnapshots(ctx, b, userID, "")

	doAsset(ctx, b, userID, assetID)
	doAssets(ctx, b, userID)
}