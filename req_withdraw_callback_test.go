package go_match2pay

import (
	"fmt"
	"testing"
)

// TestWithdrawCallback_MockSign 使用本地生成签名测试出金回调完整流程
func TestWithdrawCallback_MockSign(t *testing.T) {
	cli := newTestClient()

	amount := 0.5
	currency := "ETH"
	status := StatusDone

	sig := makeCallbackSig(amount, currency, status, TestAPIToken, TestAPISecret)
	t.Logf("generated withdraw callback sig: %s", sig)

	req := Match2PayCallbackReq{
		PaymentID:           "pay-withdraw-001",
		Status:              status,
		TransactionAmount:   amount,
		TransactionCurrency: currency,
		FinalAmount:         1500.0,
		FinalCurrency:       "USD",
	}

	err := cli.WithdrawCallback(req, sig, func(cb Match2PayCallbackReq) error {
		t.Logf("withdraw callback processed: paymentId=%s status=%s amount=%v currency=%s",
			cb.PaymentID, cb.Status, cb.TransactionAmount, cb.TransactionCurrency)
		return nil
	})
	if err != nil {
		t.Fatalf("withdraw callback failed: %v", err)
	}
}

// TestWithdrawCallback_DOGE 测试 DOGE 出金回调
func TestWithdrawCallback_DOGE(t *testing.T) {
	cli := newTestClient()

	amount := 500.0
	currency := "DOG"
	status := StatusDone

	sig := makeCallbackSig(amount, currency, status, TestAPIToken, TestAPISecret)

	req := Match2PayCallbackReq{
		PaymentID:           "pay-withdraw-doge-001",
		Status:              status,
		TransactionAmount:   amount,
		TransactionCurrency: currency,
		FinalAmount:         50.0,
		FinalCurrency:       "USD",
	}

	err := cli.WithdrawCallback(req, sig, func(cb Match2PayCallbackReq) error {
		t.Logf("DOGE withdraw callback: paymentId=%s status=%s amount=%v",
			cb.PaymentID, cb.Status, cb.TransactionAmount)
		return nil
	})
	if err != nil {
		t.Fatalf("DOGE withdraw callback failed: %v", err)
	}
}

// TestWithdrawCallback_InvalidSig 测试签名错误被拒绝
func TestWithdrawCallback_InvalidSig(t *testing.T) {
	cli := newTestClient()

	req := Match2PayCallbackReq{
		PaymentID:           "pay-withdraw-002",
		Status:              StatusDone,
		TransactionAmount:   100.0,
		TransactionCurrency: "DOG",
	}

	err := cli.WithdrawCallback(req, "invalidsignature", func(cb Match2PayCallbackReq) error {
		return nil
	})
	if err == nil {
		t.Fatal("expected sign verify error, got nil")
	}
	t.Logf("correctly rejected invalid sig: %v", err)
}

// TestWithdrawCallback_FailedStatus 测试 FAILED 状态不验签，直接通过
func TestWithdrawCallback_FailedStatus(t *testing.T) {
	cli := newTestClient()

	req := Match2PayCallbackReq{
		PaymentID:           "pay-withdraw-003",
		Status:              StatusFailed,
		TransactionAmount:   0,
		TransactionCurrency: "ETH",
	}

	err := cli.WithdrawCallback(req, "anysignature", func(cb Match2PayCallbackReq) error {
		t.Logf("failed withdraw callback: paymentId=%s status=%s", cb.PaymentID, cb.Status)
		return nil
	})
	if err != nil {
		t.Fatalf("failed status callback should pass: %v", err)
	}
}

// TestWithdrawCallback_ProcessorError 测试 processor 返回错误时正确传递
func TestWithdrawCallback_ProcessorError(t *testing.T) {
	cli := newTestClient()

	amount := 200.0
	currency := "DOG"
	status := StatusDone
	sig := makeCallbackSig(amount, currency, status, TestAPIToken, TestAPISecret)

	req := Match2PayCallbackReq{
		PaymentID:           "pay-withdraw-004",
		Status:              status,
		TransactionAmount:   amount,
		TransactionCurrency: currency,
	}

	err := cli.WithdrawCallback(req, sig, func(cb Match2PayCallbackReq) error {
		return fmt.Errorf("withdraw business logic error")
	})
	if err == nil || err.Error() != "withdraw business logic error" {
		t.Fatalf("expected business logic error, got: %v", err)
	}
	t.Logf("correctly propagated processor error: %v", err)
}
