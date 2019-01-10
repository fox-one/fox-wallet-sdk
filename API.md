# Fox Wallet API

## 鉴权

- Fox生成broker id, broker secret, broker pin secret
- 使用jwt生成token，SigningMethodHS256算法签名
- jwt map: {"i":"broker-id","u":"user-id","n":"nonce","e":123,"nr":2,"pt":"pin token"}
- e为expire, unix timestamp。
- n为nonce。nonce可以每次随机一个uuid。
- nr为nonce repeat，默认为1.代表一个nonce在10分钟内可以使用的次数。

```go
// 签名算法demo
func Sign(brokerSecret string, method, uri string, body []byte) (string, error) {
    secret, err := base64.StdEncoding.DecodeString(brokerSecret)
    if err != nil {
        return "", err
    }

    plain := append([]byte(method+uri), body...)
    h := hmac.New(sha256.New, secret)
    h.Write(plain)
    return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
```

## PIN 加密

- 以broker pin secret为aes key, broker id为aes iv, 对("data:"+pin+";"+nonce)进行aes加密
- nonce需要是请求中nonce参数的前10位

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
