# Fox Wallet API

## 鉴权

- Fox生成broker id, broker secret, broker pin secret
- 使用jwt生成token，SigningMethodHS256算法签名
- jwt map: {"i":"broker-id","u":"user-id","n":"nonce","e":123,"nr":2,"pt":"pin token"}
- u为可选项，当查询，操作具体钱包时，指明要操作的钱包地址。创建钱包时可以不传。
- e为expire, unix timestamp。
- n为nonce。nonce可以每次随机一个uuid。
- nr为nonce repeat，可选项，不传则默认为1.代表一个nonce在10分钟内可以使用的次数。

```go
func (b *Broker) SignTokenWithPIN(userID string, expire int64, pin string, nonceRepeats ...int) (string, error) {
    params := map[string]interface{}{
        "i": b.brokerID,
        "e": expire,
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
```

## PIN 加密

- 使用PKCS7 padding， AES CBC进行加密
- 以broker pin secret为aes key, broker id为aes iv, 对("data:"+pin+";"+nonce)进行aes加密
- nonce需要是请求中nonce参数的前10位

```go
AESCBC([]byte(encryptPrefix+pin+";"+nonce), pinSecret, []byte(brokerID))
```

## Create User

```go
// doCreateUser create user
//  Broker Sign Required
//  Method: POST
//  URL:    /api/users
//  Params: {"full_name":"xxxx","pin":"0123456"}
//  Return: {"data":{USER}}
```

## Modify PIN

```go
// doModinfyPIN create user
//  Broker Sign Required, PIN Required
//  Method: PUT
//  URL:    /api/pin
//  Params: {"pin":"xxx"}
//  Return: {}
```

## Verify PIN

```go
// doVerifyPIN verify pin
//  Broker Sign Required, PIN Required
//  Method: POST
//  URL:    /api/pin/verify
//  Return: {}
```

## Fetch Assets

```go
// doAssets create user
//  Broker Sign Required
//  Method: GET
//  URL:    /api/assets
//  Return: {"data":[{Asset}]}
```

## Fetch Asset

```go
// doAsset create user
//  Broker Sign Required
//  Method: GET
//  URL:    /api/asset
//  Return: {"data":{Asset}}
```

## Fetch Snapshot

```go
// doSnapshot fetch a snapshot
//  Broker Sign Required
//  Method: GET
//  URL:    /api/snapshot?snapshot_id=xxx&trace_id=xxx
//  Return: {"data":{Snapshot}}
```

## Fetch Snapshots

```go
// doSnapshots query snapshot
//  Broker Sign Required
//  Method: POST
//  URL:    /api/snapshots?user
//  Params: {"user_id":"xxx","asset_id":"xxx","offset":"xxxxx","limit":5,"order":"ASC"}
//  Return: {"data":[{Snapshot}],"next_offset":"xxxxx"}
```

## Transfer

```go
// doTransfer do a transfer
//  Broker Sign Required, PIN Required
//  Method: POST
//  URL:    /api/transfer
//  Params: {"assetID":"xx","opponent_id":"xxx","amount":"0.1","trace_id":"xxx","memo":"xxx"}
//  Return: {"data":{Snapshot}}
```

## Withdraw

```go
// doWithdraw do a withdraw
//  Broker Sign Required, PIN Required
//  Method: POST
//  URL:    /api/withdraw
//  Params: {"asset_id":"xxx","public_key":"xxx","label":"xxx","account_name":"xxx","account_tag":"xxx","amount":"0.1","trace_id":"xxx","memo":"xxx"}
//  Return: {"data":{Snapshot}}
```

## Fetch Withdraw Fee

```go
// doWithdrawFee verify pin
//  Broker Sign Required, PIN Required
//  Method: POST
//  URL:    /api/withdraw-fee
//  Params: {"asset_id":"xxx","public_key":"xxx","label":"xxx","account_name":"xxx","account_tag":"xxx"}
//  Return: {"data":{"fee":"0.1"}}
```
