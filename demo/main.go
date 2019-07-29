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

	{
		token, err := b.MixinToken(ctx, userID, "GET", "/assets", "", time.Hour*24*365)
		log.Println("token", token, err)
	}

	doFetchUser(ctx, b, userID)
	doFetchUsers(ctx, b, userID)

	fullname := "test"
	avatar := "https://images.mixin.one/7y_5w1rpGkFs_65XTvvJ37ZqsOy0t2D-qdSUVfw9LykUyE4gAWH3OMU5kxGVsJSdq3oVXx6Qz1KuJYRiOi6_5fQ=s256"
	doModifyUser(ctx, b, userID, fullname, avatar)

	doVerifyPIN(ctx, b, userID, tmpPIN)

	doModifyPIN(ctx, b, userID, tmpPIN, default_pin)

	doAssets(ctx, b, userID)
	doAsset(ctx, b, userID, assetID)
	time.Sleep(1 * time.Second)
	if len(doAssets(ctx, b, userID)) == 0 {
		log.Panicln("should have at least one asset")
	}

	doWithdrawFee(ctx, b, assetID, publicKey)

	snapshot := doTransfer(ctx, b, dapp, userID, assetID, "0.1", default_pin)
	log.Println(snapshot)

	doWithdraw(ctx, b, dapp, userID, assetID, publicKey, "0.1", default_pin)

	time.Sleep(2 * time.Second)
	doSnapshot(ctx, b, userID, "", snapshot.SnapshotID)
	_, offset := doSnapshots(ctx, b, userID, "")
	log.Println("next offset", offset)

	doAsset(ctx, b, userID, assetID)
	doAssets(ctx, b, userID)

	doAddresses(ctx, b, userID, assetID)
	addr := doUpsertAddress(ctx, b, userID, default_pin, &sdk.WithdrawAddress{
		AssetID:   assetID,
		PublicKey: publicKey,
		Label:     "test",
	})
	doAddresses(ctx, b, userID, assetID)
	doDeleteAddress(ctx, b, userID, addr.AddressID, default_pin)

	doChains(ctx, b)
	doNetworkAssets(ctx, b)
	doExchangeRates(ctx, b)

	doValidateBalance(ctx, b, userID)

	doFetchUserSession(ctx, b, userID)
}
