package go_match2pay

import (
	"fmt"
	"testing"
	"time"

	"github.com/listenfengyang/go-match2pay/utils"
)

// testLogger 测试用 Logger 实现
type testLogger struct{}

func (l *testLogger) Debugf(f string, a ...interface{}) { fmt.Printf("[DEBUG] "+f+"\n", a...) }
func (l *testLogger) Infof(f string, a ...interface{})  { fmt.Printf("[INFO]  "+f+"\n", a...) }
func (l *testLogger) Warnf(f string, a ...interface{})  { fmt.Printf("[WARN]  "+f+"\n", a...) }
func (l *testLogger) Errorf(f string, a ...interface{}) { fmt.Printf("[ERROR] "+f+"\n", a...) }

var _ utils.Logger = (*testLogger)(nil)

// newTestClient 创建测试客户端（使用官方文档示例 token/secret）
func newTestClient() *Client {
	return NewClient(&testLogger{}, &Match2PayInitParams{
		APIToken:            TestAPIToken,
		APISecret:           TestAPISecret,
		BaseURL:             DefaultBaseURL, //StagingBaseURL,
		CallbackURL:         TestCallbackURL,
		WithdrawCallbackURL: TestWithdrawCallbackURL,
		SuccessURL:          TestSuccessURL,
		FailureURL:          TestFailureURL,
	})
}

// newTestCustomer 创建测试客户信息
func newTestCustomer() Customer {
	return Customer{
		FirstName: "John",
		LastName:  "Doe",
		Address: CustomerAddress{
			Address: "123 Main St",
			City:    "New York",
			Country: "US",
			ZipCode: "10001",
			State:   "NY",
		},
		ContactInformation: CustomerContactInformation{
			Email:       "john.doe@example.com",
			PhoneNumber: "+1234567890",
		},
		Locale:              "en_US",
		DateOfBirth:         "1990-01-01",
		TradingAccountLogin: fmt.Sprintf("client%d", time.Now().Unix()),
		TradingAccountUUID:  "uuid-test-xxxx",
	}
}

// TestDeposit_DOGE 测试 DOGE 入金
func TestDeposit_DOGE(t *testing.T) {
	cli := newTestClient()
	cli.SetDebugMode(false)

	req := Match2PayDepositReq{
		Currency:           "USD",
		Amount:             100.00,
		PaymentCurrency:    CurrencyDOGE,
		PaymentGatewayName: GatewayDOGECOIN,
		Customer:           newTestCustomer(),
	}

	rsp, err := cli.Deposit(req)
	if err != nil {
		t.Logf("DOGE deposit err (expected without real credentials): %v", err)
		return
	}
	t.Logf("DOGE deposit rsp: paymentId=%s status=%s address=%s checkoutUrl=%s",
		rsp.PaymentID, rsp.Status, rsp.Address, rsp.CheckoutUrl)
}

// TestDeposit_SOL 测试 Solana 入金
func TestDeposit_SOL(t *testing.T) {
	cli := newTestClient()
	cli.SetDebugMode(false)

	req := Match2PayDepositReq{
		Currency:           "USD",
		Amount:             100.00,
		PaymentCurrency:    CurrencySOL,
		PaymentGatewayName: GatewaySOL,
		Customer:           newTestCustomer(),
	}

	rsp, err := cli.Deposit(req)
	if err != nil {
		t.Logf("SOL deposit err (expected without real credentials): %v", err)
		return
	}
	t.Logf("SOL deposit rsp: paymentId=%s status=%s address=%s checkoutUrl=%s",
		rsp.PaymentID, rsp.Status, rsp.Address, rsp.CheckoutUrl)
}

// TestDeposit_ETH 测试 ETH 入金
func TestDeposit_ETH(t *testing.T) {
	cli := newTestClient()
	cli.SetDebugMode(false)

	req := Match2PayDepositReq{
		Currency:           "USD",
		Amount:             100.00,
		PaymentCurrency:    CurrencyETH,
		PaymentGatewayName: GatewayETH,
		Customer:           newTestCustomer(),
	}

	rsp, err := cli.Deposit(req)
	if err != nil {
		t.Logf("ETH deposit err (expected without real credentials): %v", err)
		return
	}
	t.Logf("ETH deposit rsp: paymentId=%s status=%s address=%s checkoutUrl=%s",
		rsp.PaymentID, rsp.Status, rsp.Address, rsp.CheckoutUrl)
}
