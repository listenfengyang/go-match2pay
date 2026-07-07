package utils

import (
	"testing"
)

// TestFormatAmount 测试金额格式化
func TestFormatAmount(t *testing.T) {
	cases := []struct {
		input    float64
		expected string
	}{
		{10, "10"},
		{10.5, "10.5"},
		{10.50, "10.5"},
		{0.00011873, "0.00011873"},
		{1.0, "1"},
	}
	for _, c := range cases {
		got := FormatAmount(c.input)
		if got != c.expected {
			t.Errorf("FormatAmount(%v) = %q, want %q", c.input, got, c.expected)
		}
	}
}

// TestFormatAmountWith8Decimals 测试回调签名金额格式化
func TestFormatAmountWith8Decimals(t *testing.T) {
	cases := []struct {
		input    float64
		expected string
	}{
		{1, "1.00000000"},
		{0.00011873, "0.00011873"},
		{10, "10.00000000"},
	}
	for _, c := range cases {
		got := FormatAmountWith8Decimals(c.input)
		if got != c.expected {
			t.Errorf("FormatAmountWith8Decimals(%v) = %q, want %q", c.input, got, c.expected)
		}
	}
}

// TestFormatCustomer 测试客户信息格式化
func TestFormatCustomer(t *testing.T) {
	c := CustomerFields{
		FirstName:           "firstName_4da0af01617c",
		LastName:            "lastName_801eb285edd1",
		AddressAddress:      "address_52c10ed842fb",
		AddressCity:         "city_62da6faaeb17",
		AddressCountry:      "country_a6be7ed127cc",
		AddressZipCode:      "zipCode_3e168862ef49",
		AddressState:        "state_b8d531055c90",
		Email:               "email_e4ac63093536",
		PhoneNumber:         "phoneNumber_8fcb0237f7ee",
		Locale:              "en_US",
		DateOfBirth:         "dateOfBirth_338c08d95dd6",
		TradingAccountLogin: "clientId_d6811bff2963",
		TradingAccountUUID:  "clientUid_56b798ba2ae2",
	}

	expected := "{firstName=firstName_4da0af01617c, lastName=lastName_801eb285edd1, address={address=address_52c10ed842fb, city=city_62da6faaeb17, country=country_a6be7ed127cc, zipCode=zipCode_3e168862ef49, state=state_b8d531055c90}, contactInformation={email=email_e4ac63093536, phoneNumber=phoneNumber_8fcb0237f7ee}, locale=en_US, dateOfBirth=dateOfBirth_338c08d95dd6, tradingAccountLogin=clientId_d6811bff2963, tradingAccountUuid=clientUid_56b798ba2ae2}"

	got := FormatCustomer(c)
	if got != expected {
		t.Errorf("FormatCustomer mismatch:\ngot:  %s\nwant: %s", got, expected)
	}
}

// TestSignDeposit 使用官方文档示例验证入金签名
// 官方文档：https://app.theneo.io/match-trade/match2pay-v2/signature
// 预期签名：e058a6fb0d310335985f69d686ae47bdf1c34ad33bcafd83852b6034d8eb3254e7d300a63590f815baf2fe3d8eff11ec
func TestSignDeposit(t *testing.T) {
	params := map[string]string{
		"amount":             "10",
		"apiToken":           "ApiTokenProvidedBySupport",
		"callbackUrl":        "http://test/deposit/callback",
		"currency":           "USD",
		"failureUrl":         "http://test/failed-payment",
		"paymentCurrency":    "USX",
		"paymentGatewayName": "USDT TRC20",
		"paymentMethod":      "CRYPTO_AGENT",
		"successUrl":         "http://test/thanku",
		"timestamp":          "1764149779000",
	}

	customerFields := CustomerFields{
		FirstName:           "firstName_4da0af01617c",
		LastName:            "lastName_801eb285edd1",
		AddressAddress:      "address_52c10ed842fb",
		AddressCity:         "city_62da6faaeb17",
		AddressCountry:      "country_a6be7ed127cc",
		AddressZipCode:      "zipCode_3e168862ef49",
		AddressState:        "state_b8d531055c90",
		Email:               "email_e4ac63093536",
		PhoneNumber:         "phoneNumber_8fcb0237f7ee",
		Locale:              "en_US",
		DateOfBirth:         "dateOfBirth_338c08d95dd6",
		TradingAccountLogin: "clientId_d6811bff2963",
		TradingAccountUUID:  "clientUid_56b798ba2ae2",
	}
	customerStr := FormatCustomer(customerFields)
	apiSecret := "ApiSecretProvidedBySupport"

	got := SignDeposit(params, customerStr, apiSecret)
	expected := "e058a6fb0d310335985f69d686ae47bdf1c34ad33bcafd83852b6034d8eb3254e7d300a63590f815baf2fe3d8eff11ec"

	if got != expected {
		t.Errorf("SignDeposit mismatch:\ngot:  %s\nwant: %s", got, expected)
	}
}

// TestVerifyCallbackSignature 验证回调签名
func TestVerifyCallbackSignature(t *testing.T) {
	// 使用文档示例构造回调签名
	// amount = 0.00011873, currency = BTC, status = DONE
	// raw = "0.00011873BTCDONE" + apiToken + apiSecret
	apiToken := "ApiTokenProvidedBySupport"
	apiSecret := "ApiSecretProvidedBySupport"

	// 验证非 DONE 状态直接返回 true
	ok := VerifyCallbackSignature(0.00011873, "BTC", "PENDING", "anysig", apiToken, apiSecret)
	if !ok {
		t.Error("VerifyCallbackSignature: non-DONE status should return true")
	}

	// 验证 DONE 状态但签名为空
	ok = VerifyCallbackSignature(0.00011873, "BTC", "DONE", "", apiToken, apiSecret)
	if ok {
		t.Error("VerifyCallbackSignature: empty signature should return false")
	}
}
