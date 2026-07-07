package go_match2pay

import (
	"encoding/json"
	"errors"

	"github.com/listenfengyang/go-match2pay/utils"
)

// WithdrawCallback 出金回调验签并处理
//
// Match2Pay 出金回调通过 HTTP POST JSON 发送，签名在 Header "signature" 字段中。
// 仅当 status == "DONE" 时需要验证签名。
// 回调成功后建议返回 HTTP 200。
//
// 用法：
//
//	err := client.WithdrawCallback(req, headerSig, func(cb Match2PayCallbackReq) error {
//	    // 处理出金成功/失败逻辑
//	    return nil
//	})
func (cli *Client) WithdrawCallback(req Match2PayCallbackReq, headerSignature string, processor func(Match2PayCallbackReq) error) error {
	if !utils.VerifyCallbackSignature(
		req.TransactionAmount,
		req.TransactionCurrency,
		req.Status,
		headerSignature,
		cli.Params.APIToken,
		cli.Params.APISecret,
	) {
		raw, _ := json.Marshal(req)
		cli.logger.Errorf("[Match2Pay] withdraw callback verify failed: sig=%s body=%s", headerSignature, string(raw))
		return errors.New("sign verify error")
	}
	return processor(req)
}
