package go_match2pay

import (
	"github.com/go-resty/resty/v2"
	"github.com/listenfengyang/go-match2pay/utils"
)

// Client Match2Pay SDK 客户端
type Client struct {
	Params    *Match2PayInitParams
	ryClient  *resty.Client
	debugMode bool
	logger    utils.Logger
}

// NewClient 创建 Match2Pay 客户端
func NewClient(logger utils.Logger, params *Match2PayInitParams) *Client {
	if params.BaseURL == "" {
		params.BaseURL = DefaultBaseURL
	}
	if params.DepositPath == "" {
		params.DepositPath = DepositPath
	}
	if params.WithdrawPath == "" {
		params.WithdrawPath = WithdrawPath
	}
	return &Client{
		Params:    params,
		ryClient:  resty.New(),
		debugMode: false,
		logger:    logger,
	}
}

// SetDebugMode 设置调试模式
func (cli *Client) SetDebugMode(debug bool) {
	cli.debugMode = debug
}
