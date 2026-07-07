package go_match2pay

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/listenfengyang/go-match2pay/utils"
)

// Deposit 发起入金请求（加密货币）
// POST /api/payment/request
//
// 支持的币种：
//   - DOGE：paymentCurrency="DOG", paymentGatewayName="DOGECOIN"
//   - SOL：paymentCurrency="SOL", paymentGatewayName="SOL"  （仅入金）
//   - ETH：paymentCurrency="ETH", paymentGatewayName="ETH"
func (cli *Client) Deposit(req Match2PayDepositReq) (*Match2PayDepositRsp, error) {
	// 自动填充固定参数
	req.APIToken = cli.Params.APIToken
	req.PaymentMethod = PaymentMethod
	if req.CallbackURL == "" {
		req.CallbackURL = cli.Params.CallbackURL
	}
	if req.CallbackURL == "" {
		return nil, fmt.Errorf("match2pay deposit: callbackUrl is required (set CallbackURL in params)")
	}
	if req.SuccessURL == "" {
		req.SuccessURL = cli.Params.SuccessURL
	}
	if req.FailureURL == "" {
		req.FailureURL = cli.Params.FailureURL
	}
	req.Timestamp = time.Now().Unix()

	// 格式化 customer 用于签名
	customerFields := utils.CustomerFields{
		FirstName:           req.Customer.FirstName,
		LastName:            req.Customer.LastName,
		AddressAddress:      req.Customer.Address.Address,
		AddressCity:         req.Customer.Address.City,
		AddressCountry:      req.Customer.Address.Country,
		AddressZipCode:      req.Customer.Address.ZipCode,
		AddressState:        req.Customer.Address.State,
		Email:               req.Customer.ContactInformation.Email,
		PhoneNumber:         req.Customer.ContactInformation.PhoneNumber,
		Locale:              req.Customer.Locale,
		DateOfBirth:         req.Customer.DateOfBirth,
		TradingAccountLogin: req.Customer.TradingAccountLogin,
		TradingAccountUUID:  req.Customer.TradingAccountUUID,
	}
	customerStr := utils.FormatCustomer(customerFields)

	// 构建签名参数 map
	params := map[string]string{
		"amount":             utils.FormatAmount(req.Amount),
		"apiToken":           req.APIToken,
		"callbackUrl":        req.CallbackURL,
		"currency":           req.Currency,
		"failureUrl":         req.FailureURL,
		"paymentCurrency":    req.PaymentCurrency,
		"paymentGatewayName": req.PaymentGatewayName,
		"paymentMethod":      req.PaymentMethod,
		"successUrl":         req.SuccessURL,
		"timestamp":          fmt.Sprintf("%d", req.Timestamp),
	}

	// 生成签名
	req.Signature = utils.SignDeposit(params, customerStr, cli.Params.APISecret)

	rawURL := cli.Params.BaseURL + cli.Params.DepositPath
	cli.logger.Infof("[Match2Pay] deposit url: %s, currency: %s, gateway: %s, amount: %v",
		rawURL, req.PaymentCurrency, req.PaymentGatewayName, req.Amount)

	resp, err := cli.ryClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetDebug(cli.debugMode).
		Post(rawURL)

	if resp != nil {
		restLog, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(utils.GetRestyLog(resp))
		cli.logger.Infof("PSPResty#match2pay#deposit->%s", string(restLog))
	}

	if err != nil {
		return nil, fmt.Errorf("match2pay deposit: request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("match2pay deposit: http status %d body: %s", resp.StatusCode(), resp.Body())
	}

	var result Match2PayDepositRsp
	if err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("match2pay deposit: parse response failed: %w", err)
	}
	if result.Status == "FAIL" || result.Status == "FAILED" || result.Status == StatusFailed {
		errMsg := "match2pay deposit: business error"
		if len(result.ErrorList) > 0 {
			if s, ok := result.ErrorList[0].(string); ok && s != "" {
				errMsg = fmt.Sprintf("match2pay deposit: %s", s)
			}
		}
		return nil, fmt.Errorf("%s (paymentId=%s)", errMsg, result.PaymentID)
	}
	return &result, nil
}
