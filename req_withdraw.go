package go_match2pay

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/listenfengyang/go-match2pay/utils"
)

// Withdraw 发起出金请求（加密货币）
// POST /api/payment/withdrawal
//
// 支持的币种（出金）：
//   - DOGE：withdrawCurrency="DOG", paymentGatewayName="DOGECOIN"
//   - ETH：withdrawCurrency="ETH", paymentGatewayName="ETH"
//
// 注意：Solana (SOL) 仅支持入金，不支持出金
// 如需 memo，address 格式为: "walletAddress;memo"
func (cli *Client) Withdraw(req Match2PayWithdrawReq) (*Match2PayWithdrawRsp, error) {
	// 自动填充固定参数
	req.APIToken = cli.Params.APIToken
	req.PaymentMethod = PaymentMethod
	if req.CallbackURL == "" {
		if cli.Params.WithdrawCallbackURL != "" {
			req.CallbackURL = cli.Params.WithdrawCallbackURL
		} else {
			req.CallbackURL = cli.Params.CallbackURL
		}
	}
	if req.CallbackURL == "" {
		return nil, fmt.Errorf("match2pay withdraw: callbackUrl is required (set WithdrawCallbackURL or CallbackURL in params)")
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
		"address":            req.Address,
		"amount":             utils.FormatAmount(req.Amount),
		"apiToken":           req.APIToken,
		"callbackUrl":        req.CallbackURL,
		"currency":           req.Currency,
		"paymentGatewayName": req.PaymentGatewayName,
		"paymentMethod":      req.PaymentMethod,
		"timestamp":          fmt.Sprintf("%d", req.Timestamp),
		"withdrawCurrency":   req.WithdrawCurrency,
	}

	// 生成签名
	req.Signature = utils.SignWithdraw(params, customerStr, cli.Params.APISecret)

	rawURL := cli.Params.BaseURL + cli.Params.WithdrawPath
	cli.logger.Infof("[Match2Pay] withdraw url: %s, currency: %s, gateway: %s, amount: %v, address: %s",
		rawURL, req.WithdrawCurrency, req.PaymentGatewayName, req.Amount, req.Address)

	resp, err := cli.ryClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetDebug(cli.debugMode).
		Post(rawURL)

	if resp != nil {
		restLog, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(utils.GetRestyLog(resp))
		cli.logger.Infof("PSPResty#match2pay#withdraw->%s", string(restLog))
	}

	if err != nil {
		return nil, fmt.Errorf("match2pay withdraw: request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("match2pay withdraw: http status %d body: %s", resp.StatusCode(), resp.Body())
	}

	var result Match2PayWithdrawRsp
	if err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("match2pay withdraw: parse response failed: %w", err)
	}
	if result.Status == "FAIL" || result.Status == "FAILED" || result.Status == StatusFailed {
		errMsg := "match2pay withdraw: business error"
		if len(result.ErrorList) > 0 {
			if s, ok := result.ErrorList[0].(string); ok && s != "" {
				errMsg = fmt.Sprintf("match2pay withdraw: %s", s)
			}
		}
		return nil, fmt.Errorf("%s (paymentId=%s)", errMsg, result.PaymentID)
	}
	return &result, nil
}
