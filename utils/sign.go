package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// SignDeposit 生成入金请求签名
//
// Match2Pay 签名算法（请求签名）：
//  1. 取所有请求字段的 key，按 A-Z 排序
//  2. 按排序后顺序拼接 value（customer 等嵌套对象需特殊格式化）
//  3. 末尾追加 API Secret
//  4. 对整个字符串进行 SHA-384 哈希，结果小写
//
// 注意：amount 格式化为浮点数字符串（去除尾零，但保留整数时的小数点）
// 例如：10 -> "10", 10.50 -> "10.5"
func SignDeposit(params map[string]string, customerStr string, apiSecret string) string {
	// 固定 A-Z 排序的 key 列表（入金）
	orderedKeys := []string{
		"amount",
		"apiToken",
		"callbackUrl",
		"currency",
		"customer",
		"failureUrl",
		"paymentCurrency",
		"paymentGatewayName",
		"paymentMethod",
		"successUrl",
		"timestamp",
	}

	var sb strings.Builder
	for _, k := range orderedKeys {
		if k == "customer" {
			sb.WriteString(customerStr)
		} else if v, ok := params[k]; ok {
			sb.WriteString(v)
		}
	}
	sb.WriteString(apiSecret)

	raw := sb.String()
	hash := sha512.New384()
	hash.Write([]byte(raw))
	return hex.EncodeToString(hash.Sum(nil))
}

// SignWithdraw 生成出金请求签名
//
// 出金 key 列表（A-Z 排序）与入金略有不同：
// address, amount, apiToken, callbackUrl, currency, customer,
// paymentGatewayName, paymentMethod, timestamp, withdrawCurrency
func SignWithdraw(params map[string]string, customerStr string, apiSecret string) string {
	orderedKeys := []string{
		"address",
		"amount",
		"apiToken",
		"callbackUrl",
		"currency",
		"customer",
		"paymentGatewayName",
		"paymentMethod",
		"timestamp",
		"withdrawCurrency",
	}

	var sb strings.Builder
	for _, k := range orderedKeys {
		if k == "customer" {
			sb.WriteString(customerStr)
		} else if v, ok := params[k]; ok {
			sb.WriteString(v)
		}
	}
	sb.WriteString(apiSecret)

	raw := sb.String()
	hash := sha512.New384()
	hash.Write([]byte(raw))
	return hex.EncodeToString(hash.Sum(nil))
}

// FormatCustomer 将 Customer 格式化为签名所需的字符串格式
// 格式: {firstName=xx, lastName=xx, address={address=xx, city=xx, country=xx, zipCode=xx, state=xx}, contactInformation={email=xx, phoneNumber=xx}, locale=xx, dateOfBirth=xx, tradingAccountLogin=xx, tradingAccountUuid=xx}
func FormatCustomer(c CustomerFields) string {
	addrStr := fmt.Sprintf("{address=%s, city=%s, country=%s, zipCode=%s, state=%s}",
		c.AddressAddress, c.AddressCity, c.AddressCountry, c.AddressZipCode, c.AddressState)
	contactStr := fmt.Sprintf("{email=%s, phoneNumber=%s}",
		c.Email, c.PhoneNumber)

	return fmt.Sprintf("{firstName=%s, lastName=%s, address=%s, contactInformation=%s, locale=%s, dateOfBirth=%s, tradingAccountLogin=%s, tradingAccountUuid=%s}",
		c.FirstName, c.LastName, addrStr, contactStr, c.Locale, c.DateOfBirth, c.TradingAccountLogin, c.TradingAccountUUID)
}

// CustomerFields 客户字段（用于签名格式化）
type CustomerFields struct {
	FirstName           string
	LastName            string
	AddressAddress      string
	AddressCity         string
	AddressCountry      string
	AddressZipCode      string
	AddressState        string
	Email               string
	PhoneNumber         string
	Locale              string
	DateOfBirth         string
	TradingAccountLogin string
	TradingAccountUUID  string
}

// FormatAmount 格式化金额（去除尾零，保留足够精度）
// 10 -> "10", 10.50 -> "10.5", 0.00011873 -> "0.00011873"
func FormatAmount(amount float64) string {
	// 使用 strconv 保留最短准确表示
	s := strconv.FormatFloat(amount, 'f', -1, 64)
	return s
}

// FormatAmountWith8Decimals 格式化金额为 8 位小数（用于回调签名）
// 1 -> "1.00000000", 0.00011873 -> "0.00011873"
func FormatAmountWith8Decimals(amount float64) string {
	return fmt.Sprintf("%.8f", amount)
}

// VerifyCallbackSignature 验证回调签名
//
// 回调签名算法：
//
//	SHA-384("transactionAmount" + "transactionCurrency" + "status" + "apiToken" + "apiSecret")
//	其中 transactionAmount 需格式化为 8 位小数字符串
//	签名在 HTTP Header 中
//	注意：仅在 status == "DONE" 时需要验证签名
//
// 参数：
//
//	transactionAmount - 交易加密货币金额
//	transactionCurrency - 交易加密货币种类
//	status - 交易状态
//	headerSignature - HTTP Header 中的 signature 值
//	apiToken - API Token
//	apiSecret - API Secret
func VerifyCallbackSignature(transactionAmount float64, transactionCurrency, status, headerSignature, apiToken, apiSecret string) bool {
	if status != "DONE" {
		// 仅 DONE 状态需验签
		return true
	}
	if headerSignature == "" {
		return false
	}

	amountStr := FormatAmountWith8Decimals(transactionAmount)
	raw := amountStr + transactionCurrency + status + apiToken + apiSecret

	hash := sha512.New384()
	hash.Write([]byte(raw))
	expected := hex.EncodeToString(hash.Sum(nil))

	return expected == strings.ToLower(headerSignature)
}

// SignGeneric 通用签名（按 key A-Z 排序，拼接 value + secret，SHA-384）
// 用于自定义字段排列的场景
func SignGeneric(params map[string]string, apiSecret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(params[k])
	}
	sb.WriteString(apiSecret)

	hash := sha512.New384()
	hash.Write([]byte(sb.String()))
	return hex.EncodeToString(hash.Sum(nil))
}
