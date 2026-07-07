package go_match2pay

// Match2PayInitParams Match2Pay 初始化参数
type Match2PayInitParams struct {
	BaseURL             string `json:"baseUrl"             mapstructure:"baseUrl"             config:"baseUrl"             yaml:"baseUrl"`             // API 根地址，默认 https://wallet.match2pay.com
	APIToken            string `json:"apiToken"            mapstructure:"apiToken"            config:"apiToken"            yaml:"apiToken"`            // API Token，由 Match-Trade 提供
	APISecret           string `json:"apiSecret"           mapstructure:"apiSecret"           config:"apiSecret"           yaml:"apiSecret"`           // API Secret，由 Match-Trade 提供
	CallbackURL         string `json:"callbackUrl"         mapstructure:"callbackUrl"         config:"callbackUrl"         yaml:"callbackUrl"`         // 入金异步回调地址
	WithdrawCallbackURL string `json:"withdrawCallbackUrl" mapstructure:"withdrawCallbackUrl" config:"withdrawCallbackUrl" yaml:"withdrawCallbackUrl"` // 出金异步回调地址（不填时使用 callbackUrl）
	SuccessURL          string `json:"successUrl"          mapstructure:"successUrl"          config:"successUrl"          yaml:"successUrl"`          // 入金成功跳转地址
	FailureURL          string `json:"failureUrl"          mapstructure:"failureUrl"          config:"failureUrl"          yaml:"failureUrl"`          // 入金失败跳转地址
	DepositPath         string `json:"depositPath"         mapstructure:"depositPath"         config:"depositPath"         yaml:"depositPath"`         // 入金路径，默认 /api/payment/request
	WithdrawPath        string `json:"withdrawPath"        mapstructure:"withdrawPath"        config:"withdrawPath"        yaml:"withdrawPath"`        // 出金路径，默认 /api/payment/withdrawal
}
