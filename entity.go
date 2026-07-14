package go_match2pay

// ============================================================
// Customer 客户信息
// ============================================================

// CustomerAddress 客户地址
type CustomerAddress struct {
	Address string `json:"address"` // 街道地址
	City    string `json:"city"`    // 城市
	Country string `json:"country"` // 国家（ISO-3166-1 alpha-2）
	ZipCode string `json:"zipCode"` // 邮政编码
	State   string `json:"state"`   // 州/省
}

// CustomerContactInformation 客户联系信息
type CustomerContactInformation struct {
	Email       string `json:"email"`       // 邮箱
	PhoneNumber string `json:"phoneNumber"` // 电话号码
}

// Customer 客户信息
type Customer struct {
	FirstName           string                     `json:"firstName"`           // 名
	LastName            string                     `json:"lastName"`            // 姓
	Address             CustomerAddress            `json:"address"`             // 地址
	ContactInformation  CustomerContactInformation `json:"contactInformation"`  // 联系信息
	Locale              string                     `json:"locale"`              // 语言（e.g. en_US）
	DateOfBirth         string                     `json:"dateOfBirth"`         // 出生日期
	TradingAccountLogin string                     `json:"tradingAccountLogin"` // 交易账户登录名（clientId）
	TradingAccountUUID  string                     `json:"tradingAccountUuid"`  // 交易账户 UUID（clientUid）
}

// ============================================================
// Deposit (入金)
// POST /api/payment/request
// ============================================================

// Match2PayDepositReq 入金请求参数
type Match2PayDepositReq struct {
	Amount             float64  `json:"amount"`             // 金额
	APIToken           string   `json:"apiToken"`           // API Token（自动填充）
	CallbackURL        string   `json:"callbackUrl"`        // 异步回调地址（自动填充）
	Currency           string   `json:"currency"`           // 法币种类，如 USD
	Customer           Customer `json:"customer"`           // 客户信息
	FailureURL         string   `json:"failureUrl"`         // 失败跳转地址（自动填充）
	PaymentCurrency    string   `json:"paymentCurrency"`    // 加密货币代码，如 DOG / SOL / ETH
	PaymentGatewayName string   `json:"paymentGatewayName"` // 支付网关名，如 DOGECOIN / SOL / ETH
	PaymentMethod      string   `json:"paymentMethod"`      // 支付方式，固定 CRYPTO_AGENT（自动填充）
	SuccessURL         string   `json:"successUrl"`         // 成功跳转地址（自动填充）
	Timestamp          int64    `json:"timestamp"`          // 当前时间秒（自动填充）
	Signature          string   `json:"signature"`          // 签名（自动填充）
}

// Match2PayDepositRsp 入金响应
type Match2PayDepositRsp struct {
	PaymentID           string        `json:"paymentId"`           // 支付 ID
	Status              string        `json:"status"`              // 状态：PENDING / DONE / FAILED / EXPIRED
	TransactionAmount   float64       `json:"transactionAmount"`   // 实际加密货币金额
	TransactionCurrency string        `json:"transactionCurrency"` // 实际加密货币种类
	NetAmount           float64       `json:"netAmount"`           // 净额
	FinalAmount         float64       `json:"finalAmount"`         // 最终法币金额
	FinalCurrency       string        `json:"finalCurrency"`       // 最终法币种类
	ProcessingFee       float64       `json:"processingFee"`       // 手续费
	ConvertRatio        float64       `json:"convertRatio"`        // 汇率
	ErrorList           []interface{} `json:"errorList"`           // 错误列表
	Address             string        `json:"address"`             // 充值地址
	TempTransactionID   string        `json:"tempTransactionId"`   // 临时交易 ID
	// Error response fields (API may return numeric status on error)
	ErrStatus   int    `json:"status_code,omitempty"` // HTTP 错误状态码（错误时）
	Title       string `json:"title,omitempty"`       // 错误标题
	Detail      string `json:"detail,omitempty"`      // 错误详情
	CheckoutUrl string `json:"checkoutUrl"`           // 下单地址
}

// ============================================================
// Withdrawal (出金)
// POST /api/payment/withdrawal
// ============================================================

// Match2PayWithdrawReq 出金请求参数
type Match2PayWithdrawReq struct {
	Amount             float64  `json:"amount"`             // 金额
	APIToken           string   `json:"apiToken"`           // API Token（自动填充）
	CallbackURL        string   `json:"callbackUrl"`        // 异步回调地址（自动填充）
	Currency           string   `json:"currency"`           // 法币种类，如 USD
	Customer           Customer `json:"customer"`           // 客户信息
	WithdrawCurrency   string   `json:"withdrawCurrency"`   // 加密货币代码，如 DOG / ETH（注意：Solana 仅支持入金）
	PaymentGatewayName string   `json:"paymentGatewayName"` // 支付网关名，如 DOGECOIN / ETH
	PaymentMethod      string   `json:"paymentMethod"`      // 支付方式，固定 CRYPTO_AGENT（自动填充）
	Address            string   `json:"address"`            // 收款钱包地址（如有 memo 格式：address;memo）
	Timestamp          int64    `json:"timestamp"`          // 当前时间秒（自动填充）
	Signature          string   `json:"signature"`          // 签名（自动填充）
}

// Match2PayWithdrawRsp 出金响应
type Match2PayWithdrawRsp struct {
	PaymentID           string        `json:"paymentId"`           // 支付 ID
	Status              string        `json:"status"`              // 状态：PENDING / DONE / FAILED
	TransactionAmount   float64       `json:"transactionAmount"`   // 实际加密货币金额
	TransactionCurrency string        `json:"transactionCurrency"` // 实际加密货币种类
	NetAmount           float64       `json:"netAmount"`           // 净额
	FinalAmount         float64       `json:"finalAmount"`         // 最终法币金额
	FinalCurrency       string        `json:"finalCurrency"`       // 最终法币种类
	ProcessingFee       float64       `json:"processingFee"`       // 手续费
	ConvertRatio        float64       `json:"convertRatio"`        // 汇率
	ErrorList           []interface{} `json:"errorList"`           // 错误列表
	// Error response fields (API may return numeric status on error)
	ErrStatus int    `json:"status_code,omitempty"` // HTTP 错误状态码（错误时）
	Title     string `json:"title,omitempty"`       // 错误标题
	Detail    string `json:"detail,omitempty"`      // 错误详情
}

// ============================================================
// Callback（入金 & 出金回调）
// ============================================================

// CryptoTransactionInfo 加密货币交易详情
type CryptoTransactionInfo struct {
	TxID                string  `json:"txid"`                // 交易哈希
	Confirmations       int     `json:"confirmations"`       // 确认数
	Amount              float64 `json:"amount"`              // 交易金额
	TransactionCurrency string  `json:"transactionCurrency"` // 加密货币种类
	ConfirmedTime       string  `json:"confirmedTime"`       // 确认时间
	Status              string  `json:"status"`              // 状态：DONE
	ProcessingFee       float64 `json:"processingFee"`       // 手续费
	ConversionRate      float64 `json:"conversionRate"`      // 汇率
	BlockchainType      string  `json:"blockchainType"`      // 区块链类型（如 ETHEREUM）
}

// Match2PayCallbackReq 回调通知参数（入金 & 出金通用）
// 签名在 HTTP Header 中，字段名为 signature
type Match2PayCallbackReq struct {
	DepositAddress        string                  `json:"depositAddress"`        // 充值地址（入金）
	CryptoTransactionInfo []CryptoTransactionInfo `json:"cryptoTransactionInfo"` // 加密货币交易详情列表
	PaymentID             string                  `json:"paymentId"`             // 支付 ID
	Status                string                  `json:"status"`                // 状态：DONE / FAILED / EXPIRED
	TransactionAmount     float64                 `json:"transactionAmount"`     // 加密货币金额（用于签名验证）
	NetAmount             float64                 `json:"netAmount"`             // 净额
	TransactionCurrency   string                  `json:"transactionCurrency"`   // 加密货币种类（用于签名验证）
	ProcessingFee         float64                 `json:"processingFee"`         // 手续费
	FinalAmount           float64                 `json:"finalAmount"`           // 最终法币金额
	FinalCurrency         string                  `json:"finalCurrency"`         // 最终法币种类
	ConversionRate        float64                 `json:"conversionRate"`        // 汇率
	SettlementCurrency    string                  `json:"settlementCurrency"`    // 结算货币种类
	SettlementAmount      float64                 `json:"settlementAmount"`      // 结算金额
}
