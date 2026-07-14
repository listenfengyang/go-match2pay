# go-match2pay

Match2Pay v2 Go SDK，支持加密货币入金（Deposit）和出金（Withdrawal）。

## 支持的币种

| 操作 | 币种 | `paymentCurrency` / `withdrawCurrency` | `paymentGatewayName` |
|------|------|----------------------------------------|----------------------|
| 入金/出金 | DOGE | `DOG` | `DOGECOIN` |
| 入金/出金 | ETH | `ETH` | `ETH` |
| **仅入金** | Solana | `SOL` | `SOL` |
| **仅入金** | 币安支付 | `BNB`（或其他币安支持币种） | `BINANCEPAY` |

> ⚠️ Solana (SOL) 仅支持入金，不支持出金。
> 💡 Binance Pay 入金后，响应中的 `checkoutUrl` 即为跳转币安 App/Web 的支付链接。

## 安装

```bash
go get github.com/listenfengyang/go-match2pay
```

## 快速开始

```go
import (
    match2pay "github.com/listenfengyang/go-match2pay"
)

// 初始化客户端
client := match2pay.NewClient(logger, &match2pay.Match2PayInitParams{
    APIToken:    "your-api-token",
    APISecret:   "your-api-secret",
    BaseURL:     "https://platform.match-trade.com",   // 可选，默认值
    CallbackURL: "https://your-domain.com/callback",
    SuccessURL:  "https://your-domain.com/success",
    FailureURL:  "https://your-domain.com/failure",
})
```

### 入金（Deposit）

```go
// DOGE 入金
rsp, err := client.Deposit(match2pay.Match2PayDepositReq{
    Currency:           "USD",    // 法币种类
    Amount:             100.00,   // 法币金额
    PaymentCurrency:    "DOG",    // 加密货币种类
    PaymentGatewayName: "DOGECOIN",
    Customer: match2pay.Customer{
        FirstName: "John",
        LastName:  "Doe",
        Address: match2pay.Address{
            Address: "123 Main St",
            City:    "New York",
            Country: "US",
            ZipCode: "10001",
            State:   "NY",
        },
        ContactInformation: match2pay.ContactInformation{
            Email:       "john.doe@example.com",
            PhoneNumber: "+1234567890",
        },
        Locale:              "en_US",
        DateOfBirth:         "1990-01-01",
        TradingAccountLogin: "client123",
        TradingAccountUUID:  "uuid-xxxx",
    },
})
// rsp.RedirectUrl 为支付页面跳转地址

// SOL 入金
rsp, err := client.Deposit(match2pay.Match2PayDepositReq{
    Currency:           "USD",
    Amount:             100.00,
    PaymentCurrency:    "SOL",
    PaymentGatewayName: "SOL",
    Customer:           customer,
})

// ETH 入金
rsp, err := client.Deposit(match2pay.Match2PayDepositReq{
    Currency:           "USD",
    Amount:             100.00,
    PaymentCurrency:    "ETH",
    PaymentGatewayName: "ETH",
    Customer:           customer,
})
```

### 币安支付（Binance Pay）入金

Match2Pay 支持通过 Binance Pay 完成入金。调用后，响应中的 `checkoutUrl` 即为跳转币安 App 或 Web 的支付链接，将用户重定向到该地址即可完成支付。

```go
// 方式一：使用便捷方法 DepositWithBinancePay
rsp, err := client.DepositWithBinancePay(
    "USD",                    // 法币种类
    100.00,                   // 法币金额
    match2pay.CurrencyBNB,    // 加密货币代码（BNB 或其他币安支持的币种）
    customer,
)
if err != nil {
    // 处理错误
}
// 将用户重定向到 checkoutUrl，即可在币安 App/Web 完成支付
redirectURL := rsp.CheckoutUrl

// 方式二：使用通用 Deposit 方法
rsp, err := client.Deposit(match2pay.Match2PayDepositReq{
    Currency:           "USD",
    Amount:             100.00,
    PaymentCurrency:    match2pay.CurrencyBNB,
    PaymentGatewayName: match2pay.GatewayBINANCEPAY,
    Customer:           customer,
})
```

### 出金（Withdrawal）

```go
// DOGE 出金
rsp, err := client.Withdraw(match2pay.Match2PayWithdrawReq{
    Currency:           "USD",
    Amount:             50.00,
    WithdrawCurrency:   "DOG",
    PaymentGatewayName: "DOGECOIN",
    Address:            "DRWallet123Address",  // 钱包地址（如有 memo，格式: "address;memo"）
    Customer:           customer,
})
// rsp.Status 为处理状态

// ETH 出金
rsp, err := client.Withdraw(match2pay.Match2PayWithdrawReq{
    Currency:           "USD",
    Amount:             50.00,
    WithdrawCurrency:   "ETH",
    PaymentGatewayName: "ETH",
    Address:            "0xYourEthWalletAddress",
    Customer:           customer,
})
```

### 回调验签

```go
// 在回调 HTTP Handler 中
func callbackHandler(w http.ResponseWriter, r *http.Request) {
    headerSig := r.Header.Get("signature")

    var cb match2pay.Match2PayCallbackReq
    json.NewDecoder(r.Body).Decode(&cb)

    if !client.VerifyCallbackSignature(&cb, headerSig) {
        http.Error(w, "invalid signature", http.StatusUnauthorized)
        return
    }

    switch cb.Status {
    case match2pay.StatusDone:
        // 交易成功
    case match2pay.StatusFailed:
        // 交易失败
    case match2pay.StatusPending:
        // 交易处理中
    }
    w.WriteHeader(http.StatusOK)
}
```

## 回调状态码

| 常量 | 值 | 说明 |
|------|-----|------|
| `StatusPending` | `PENDING` | 处理中 |
| `StatusDone` | `DONE` | 成功 |
| `StatusFailed` | `FAILED` | 失败 |
| `StatusError` | `ERROR` | 错误 |
| `StatusExpired` | `EXPIRED` | 已过期 |
| `StatusNew` | `NEW` | 新建 |
| `StatusInProgress` | `IN_PROGRESS` | 进行中 |

## 配置项

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `APIToken` | string | ✅ | Match2Pay API Token |
| `APISecret` | string | ✅ | Match2Pay API Secret |
| `BaseURL` | string | ❌ | API 基础地址，默认 `https://platform.match-trade.com` |
| `CallbackURL` | string | ✅ | 支付回调地址 |
| `SuccessURL` | string | ✅（入金） | 支付成功跳转地址 |
| `FailureURL` | string | ✅（入金） | 支付失败跳转地址 |
| `DepositPath` | string | ❌ | 入金接口路径，默认 `/api/payment/request` |
| `WithdrawPath` | string | ❌ | 出金接口路径，默认 `/api/payment/withdrawal` |

## 文件结构

```
go-match2pay/
├── client.go          # 客户端初始化
├── conf.go            # 配置参数结构体
├── default.go         # 默认值常量
├── entity.go          # 请求/响应/回调结构体
├── req_deposit.go     # 入金实现
├── req_withdraw.go    # 出金实现
├── callback.go        # 回调验签
└── utils/
    ├── logger.go      # Logger 接口
    ├── sign.go        # 签名算法
    ├── sign_test.go   # 签名单元测试
    └── request.go     # HTTP 请求日志工具
```

## API 文档

- 官方文档：https://app.theneo.io/match-trade/match2pay-v2
