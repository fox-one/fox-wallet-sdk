# fox-wallet-sdk

demo 中包含sdk的初始化以及各个接口的使用范例。

## 初始化Broker

```go
sdk.NewBroker(apiBase, brokerID, brokerSecret, brokerPINSecret)
```

## 生成pin token

```go
pinToken, _, err := b.PINToken(pin)
if err != nil {
    log.Panicln(err)
}
```

## 通过broker访问Fox Wallet

```go
// 创建用户，name, pinToken 可选
b.CreateUser(ctx, name, pinToken)

// 修改PIN
b.ModifyPIN(ctx, userID, pinToken, nonce, newPINToken)

// 验证PIN
b.VerifyPIN(ctx, userID, pinToken, nonce)

// 用户的所有持仓，地址信息
b.FetchAssets(ctx, userID)

// 用户的某个币的持仓，地址信息
b.FetchAsset(ctx, userID, assetID)

// 转账
b.Transfer(ctx, userID, pinToken, nonce, &mixin.TransferInput{
    AssetID:    assetID,
    Amount:     amount,
    OpponentID: receiverID,
    TraceID:    uuid.Must(uuid.NewV4()).String(),
    Memo:       "ping",
})

// 提现
b.Withdraw(ctx, userID, pinToken, nonce, &&sdk.WithdrawInput{
    WithdrawAddress: mixin.WithdrawAddress{
        AssetID:     assetID,
        PublicKey:   publicKey,
        AccountName: "",
        AccountTag:  "",
    },

    Amount:  amount,
    TraceID: uuid.Must(uuid.NewV4()).String(),
    Memo:    "pong",
})

// 提现手续费查询
b.FetchWithdrawFee(ctx, userID, pinToken, nonce, &mixin.WithdrawAddress{
    AssetID:     assetID,
    PublicKey:   publicKey,
    AccountName: "",
    AccountTag:  "",
})

// 查询用户的单条snapshot
b.FetchSnapshot(ctx, userID, traceID, snapshotID)

// 查询转账记录
//  userID, assetID为可选;
//  order为ASC/DESC;
//  limit最大为500
b.FetchSnapshots(ctx, userID, assetID, 0, "DESC", 100)
```