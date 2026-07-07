package go_match2pay

import (
	"fmt"
	"testing"
	"time"
)

// TestWithdraw_DOGE 测试 DOGE 出金
func TestWithdraw_DOGE(t *testing.T) {
	cli := newTestClient()
	cli.SetDebugMode(false)

	req := Match2PayWithdrawReq{
		Currency:           "USD",
		Amount:             50.00,
		WithdrawCurrency:   CurrencyDOGE,
		PaymentGatewayName: GatewayDOGECOIN,
		Address:            "DH5yaieqoZN36fDVciNyRueRGvGLR3mr7L", // Dogecoin Foundation public address
		Customer:           newTestCustomer(),
	}

	rsp, err := cli.Withdraw(req)
	if err != nil {
		t.Logf("DOGE withdraw err (expected without real credentials): %v", err)
		return
	}
	t.Logf("DOGE withdraw rsp: paymentId=%s status=%s", rsp.PaymentID, rsp.Status)
}

// TestWithdraw_ETH 测试 ETH 出金
func TestWithdraw_ETH(t *testing.T) {
	cli := newTestClient()
	cli.SetDebugMode(false)

	req := Match2PayWithdrawReq{
		Currency:           "USD",
		Amount:             50.00,
		WithdrawCurrency:   CurrencyETH,
		PaymentGatewayName: GatewayETH,
		Address:            "0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe", // Ethereum Foundation public address
		Customer:           newTestCustomer(),
	}

	rsp, err := cli.Withdraw(req)
	if err != nil {
		t.Logf("ETH withdraw err (expected without real credentials): %v", err)
		return
	}
	t.Logf("ETH withdraw rsp: paymentId=%s status=%s", rsp.PaymentID, rsp.Status)
}

// TestWithdraw_WithMemo 测试带 memo 的出金（address;memo 格式）
func TestWithdraw_WithMemo(t *testing.T) {
	cli := newTestClient()
	cli.SetDebugMode(false)

	walletAddress := "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh" // XRP genesis address (public)
	memo := "123456789"
	addressWithMemo := fmt.Sprintf("%s;%s", walletAddress, memo)

	req := Match2PayWithdrawReq{
		Currency:           "USD",
		Amount:             50.00,
		WithdrawCurrency:   CurrencyXRP,
		PaymentGatewayName: GatewayXRP,
		Address:            addressWithMemo,
		Customer:           newTestCustomer(),
	}

	rsp, err := cli.Withdraw(req)
	if err != nil {
		t.Logf("XRP withdraw with memo err (expected without real credentials): %v", err)
		return
	}
	t.Logf("XRP withdraw with memo rsp: paymentId=%s status=%s", rsp.PaymentID, rsp.Status)
}

// TestWithdraw_CustomCallbackURL 测试自定义回调地址的出金
func TestWithdraw_CustomCallbackURL(t *testing.T) {
	cli := newTestClient()
	cli.SetDebugMode(false)

	req := Match2PayWithdrawReq{
		Currency:           "USD",
		Amount:             100.00,
		WithdrawCurrency:   CurrencyDOGE,
		PaymentGatewayName: GatewayDOGECOIN,
		Address:            "DH5yaieqoZN36fDVciNyRueRGvGLR3mr7L", // Dogecoin Foundation public address
		CallbackURL:        fmt.Sprintf("http://test/withdraw/callback?orderId=%d", time.Now().Unix()),
		Customer:           newTestCustomer(),
	}

	rsp, err := cli.Withdraw(req)
	if err != nil {
		t.Logf("DOGE withdraw custom callback err (expected without real credentials): %v", err)
		return
	}
	t.Logf("DOGE withdraw custom callback rsp: paymentId=%s status=%s", rsp.PaymentID, rsp.Status)
}
