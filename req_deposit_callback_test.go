package go_match2pay

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/listenfengyang/go-match2pay/utils"
)

// makeCallbackSig 构造回调签名（用于测试）
func makeCallbackSig(amount float64, currency, status, apiToken, apiSecret string) string {
	amountStr := utils.FormatAmountWith8Decimals(amount)
	raw := amountStr + currency + status + apiToken + apiSecret
	h := sha512.New384()
	h.Write([]byte(raw))
	return hex.EncodeToString(h.Sum(nil))
}

// TestDepositCallback_MockSign 使用本地生成签名测试入金回调完整流程
func TestDepositCallback_MockSign(t *testing.T) {
	cli := newTestClient()

	amount := 0.00011873
	currency := "BTC"
	status := StatusDone

	sig := makeCallbackSig(amount, currency, status, TestAPIToken, TestAPISecret)
	t.Logf("generated deposit callback sig: %s", sig)

	req := Match2PayCallbackReq{
		PaymentID:           "pay-deposit-001",
		Status:              status,
		TransactionAmount:   amount,
		TransactionCurrency: currency,
		FinalAmount:         100.0,
		FinalCurrency:       "USD",
		DepositAddress:      "DTestWalletAddress123",
	}

	err := cli.DepositCallback(req, sig, func(cb Match2PayCallbackReq) error {
		t.Logf("deposit callback processed: paymentId=%s status=%s amount=%v currency=%s",
			cb.PaymentID, cb.Status, cb.TransactionAmount, cb.TransactionCurrency)
		return nil
	})
	if err != nil {
		t.Fatalf("deposit callback failed: %v", err)
	}
}

// TestDepositCallback_InvalidSig 测试签名错误被拒绝
func TestDepositCallback_InvalidSig(t *testing.T) {
	cli := newTestClient()

	req := Match2PayCallbackReq{
		PaymentID:           "pay-deposit-002",
		Status:              StatusDone,
		TransactionAmount:   100.0,
		TransactionCurrency: "ETH",
	}

	err := cli.DepositCallback(req, "invalidsignature", func(cb Match2PayCallbackReq) error {
		return nil
	})
	if err == nil {
		t.Fatal("expected sign verify error, got nil")
	}
	t.Logf("correctly rejected invalid sig: %v", err)
}

// TestDepositCallback_PendingStatus 测试非 DONE 状态不验签，直接通过
func TestDepositCallback_PendingStatus(t *testing.T) {
	cli := newTestClient()

	req := Match2PayCallbackReq{
		PaymentID:           "pay-deposit-003",
		Status:              StatusPending,
		TransactionAmount:   100.0,
		TransactionCurrency: "SOL",
	}

	processed := false
	err := cli.DepositCallback(req, "anysignature", func(cb Match2PayCallbackReq) error {
		processed = true
		t.Logf("pending deposit callback processed: paymentId=%s status=%s", cb.PaymentID, cb.Status)
		return nil
	})
	if err != nil {
		t.Fatalf("pending callback should pass: %v", err)
	}
	if !processed {
		t.Fatal("processor was not called")
	}
}

// TestDepositCallback_ExpiredStatus 测试 EXPIRED 状态不验签，直接通过
func TestDepositCallback_ExpiredStatus(t *testing.T) {
	cli := newTestClient()

	req := Match2PayCallbackReq{
		PaymentID:           "pay-deposit-004",
		Status:              StatusExpired,
		TransactionAmount:   0,
		TransactionCurrency: "DOG",
	}

	err := cli.DepositCallback(req, "", func(cb Match2PayCallbackReq) error {
		t.Logf("expired deposit callback: paymentId=%s status=%s", cb.PaymentID, cb.Status)
		return nil
	})
	if err != nil {
		t.Fatalf("expired callback should pass: %v", err)
	}
}

// TestDepositCallback_ProcessorError 测试 processor 返回错误时正确传递
func TestDepositCallback_ProcessorError(t *testing.T) {
	cli := newTestClient()

	amount := 50.0
	currency := "ETH"
	status := StatusDone
	sig := makeCallbackSig(amount, currency, status, TestAPIToken, TestAPISecret)

	req := Match2PayCallbackReq{
		PaymentID:           "pay-deposit-005",
		Status:              status,
		TransactionAmount:   amount,
		TransactionCurrency: currency,
	}

	err := cli.DepositCallback(req, sig, func(cb Match2PayCallbackReq) error {
		return fmt.Errorf("business logic error")
	})
	if err == nil || err.Error() != "business logic error" {
		t.Fatalf("expected business logic error, got: %v", err)
	}
	t.Logf("correctly propagated processor error: %v", err)
}
