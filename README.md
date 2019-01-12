# fox-wallet-sdk

demo 中包含sdk的初始化以及各个接口的使用范例。

## 初始化Broker

```go
// 若需要自行创建钱包，管理钱包，则使用broker id, secret, pin secret初始化sdk.Broker
sdk.NewBroker(apiBase, brokerID, brokerSecret, brokerPINSecret)
// 若只接收会员层下发的token，则使用api base初始化BrokerHandler
sdk.NewBrokerHandler(apiBase)
```

## 生成pin token

```go
// pin 必须为6位数字
pinToken, _, err := b.PINToken(pin)
if err != nil {
    log.Panicln(err)
}
```

## 生成jwt token

```go
// SignToken 生成jwt token
func (b *Broker) SignToken(userID string, expire int64, nonceRepeat ...int) (string, error)
// SignTokenWithPIN 生成带有pin信息的jwt token
func (b *Broker) SignTokenWithPIN(userID string, expire int64, pin string, nonceRepeats ...int) (string, error)
```

## 错误码

```go
var (
    // ErrInternalServerError code 2, internal server error
    ErrInternalServerError = errors.New("internal server error")
    // ErrInvalidBroker code 1500, invalid broker
    ErrInvalidBroker = errors.New("invalid broker")
    // ErrReplayRequest code 1501, replay request
    ErrReplayRequest = errors.New("replay request")
    // ErrInvalidUser code 1502, invalid user
    ErrInvalidUser = errors.New("invalid user")
    // ErrForbidden code 1503, forbiddenn
    ErrForbidden = errors.New("forbidden")
    // ErrInvalidPINToken code 1504, invalid pin token
    ErrInvalidPINToken = errors.New("invalid pin token")
    // ErrIncorrectPIN code 1505, incorrect pin
    ErrIncorrectPIN = errors.New("incorrect pin")
    // ErrIncorrectPINNonce code 1506, invalid pin nonce
    ErrIncorrectPINNonce = errors.New("invalid pin nonce")
    // ErrNotFound code 1507, not found
    ErrNotFound = errors.New("not found")
    // ErrInvalidRequest code 1508, invalid request
    ErrInvalidRequest = errors.New("invalid request")
    // ErrTooManyRequests code 1509, too many requests
    ErrTooManyRequests = errors.New("too many requests")
    // ErrInsufficientBalance code 1510, insufficient balance
    ErrInsufficientBalance = errors.New("insufficient balance")
    // ErrInsufficientFee code 1511, insufficient fee
    ErrInsufficientFee = errors.New("insufficient fee")
    // ErrAmountTooSmall code 1512, amount too small
    ErrAmountTooSmall = errors.New("amount too small")
    // ErrInvalidTrace code 1513, invalid trace id
    ErrInvalidTrace = errors.New("invalid trace id")
)
```

## 通过broker访问Fox Wallet

```go
// 创建用户，name, pin 可选。
b.CreateUser(ctx, name, pin)

// 修改PIN
b.ModifyPIN(ctx, userID, pin, newPIN)

// 验证PIN
b.VerifyPIN(ctx, userID, pin)

// 用户的所有持仓，地址信息
b.FetchAssets(ctx, userID)

// 用户的某个币的持仓，地址信息
b.FetchAsset(ctx, userID, assetID)

// 转账
b.Transfer(ctx, userID, pin, &sdk.TransferInput{
    AssetID:    assetID,
    Amount:     amount,
    OpponentID: receiverID,
    TraceID:    uuid.Must(uuid.NewV4()).String(),
    Memo:       "ping",
})

// 提现
//  对于EOS及其链上的币种，public key留空,只需填入用户的account name, account tag; memo不可填。
//  对于其他币种，只需填入public key。留空account name, account tag.
b.Withdraw(ctx, userID, pin, &sdk.WithdrawInput{
    WithdrawAddress: sdk.WithdrawAddress{
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
b.FetchWithdrawFee(ctx, userID, pin, &sdk.WithdrawAddress{
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
b.FetchSnapshots(ctx, userID, assetID, "xxx", "DESC", 100)
```