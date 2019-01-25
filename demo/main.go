package main

import (
	"context"
	"crypto/x509"
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
	b := sdk.NewBroker(apiBase, brokerID, brokerSecret, brokerPINSecret)

	assetID := "965e5c6e-434c-3fa9-b780-c50f43cd955c"
	publicKey := "0xe20FE5C04Fa6b044b720F8CA019Cd896881ED13B"

	tmpPIN := "123456"
	userID := doCreateUser(ctx, b, tmpPIN).UserID

	doVerifyPIN(ctx, b, userID, tmpPIN)

	doModifyPIN(ctx, b, userID, tmpPIN, default_pin)

	doAssets(ctx, b, userID)
	doAsset(ctx, b, userID, assetID)
	time.Sleep(1 * time.Second)
	if len(doAssets(ctx, b, userID)) == 0 {
		log.Panicln("should have at least one asset")
	}

	doWithdrawFee(ctx, b, userID, assetID, publicKey, default_pin)

	snapshot := doTransfer(ctx, b, dapp, userID, assetID, "0.1", default_pin)

	doWithdraw(ctx, b, dapp, userID, assetID, publicKey, "0.1", default_pin)

	time.Sleep(2 * time.Second)
	doSnapshot(ctx, b, userID, "", snapshot.SnapshotID)
	_, offset := doSnapshots(ctx, b, userID, "")
	log.Println("next offset", offset)

	doAsset(ctx, b, userID, assetID)
	doAssets(ctx, b, userID)

	doChains(ctx, b)
	doNetworkAssets(ctx, b)
}
