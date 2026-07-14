package go_match2pay

// Callback status constants
const (
	StatusPending    = "PENDING"
	StatusDone       = "DONE"
	StatusFailed     = "FAILED"
	StatusError      = "ERROR"
	StatusExpired    = "EXPIRED"
	StatusNew        = "NEW"
	StatusInProgress = "IN_PROGRESS"
)

// URL and path constants
const (
	DefaultBaseURL = "https://wallet.match2pay.com"
	StagingBaseURL = "https://wallet-staging.match2pay.com"
	DepositPath    = "/api/v2/payment/deposit"    // 入金
	WithdrawPath   = "/api/v2/payment/withdrawal" // 出金
	PaymentMethod  = "CRYPTO_AGENT"               // 加密货币固定方法
)

const (

	// Supported coins - paymentCurrency codes
	CurrencyDOGE = "DOG" // Dogecoin
	CurrencySOL  = "SOL" // Solana (Deposits only)
	CurrencyETH  = "ETH" // Ethereum

	// Additional currencies
	CurrencyBTC = "BTC" // Bitcoin
	CurrencyBNB = "BNB" // Binance Coin (BSC)
	CurrencyLTC = "LTC" // Litecoin
	CurrencyXRP = "XRP" // XRP (Deposits only)
	CurrencyUST = "UST" // USDT ERC20
	CurrencyUSX = "USX" // USDT TRC20
	CurrencyUSB = "USB" // USDT BEP20
	CurrencyUSS = "USS" // USDT Solana
	CurrencyUCS = "UCS" // USDC Solana
	CurrencyUCC = "UCC" // USDC ERC20
	CurrencyUCB = "UCB" // USDC BEP20
	CurrencyUSP = "USP" // USDT Polygon
	CurrencyDAE = "DAE" // DAI ERC20
	CurrencyDAB = "DAB" // DAI BEP20
	CurrencyDAP = "DAP" // DAI Polygon
	CurrencyDAS = "DAS" // DAI Solanas

	// paymentGatewayName values
	GatewayDOGECOIN   = "DOGECOIN"
	GatewaySOL        = "SOL"
	GatewayETH        = "ETH"
	GatewayBTC        = "BTC"
	GatewayBNB        = "BNB"
	GatewayLTC        = "LTC"
	GatewayXRP        = "XRP"
	GatewayUSDTTRC20  = "USDT TRC20"
	GatewayUSDTERC20  = "USDT ERC20"
	GatewayUSDTBEP20  = "USDT BEP20"
	GatewayUSDTSOL    = "USDT SOL"
	GatewayUSDCSOL    = "USDC SOL"
	GatewayUSDCERC20  = "USDC ERC20"
	GatewayUSDCBEP20  = "USDC BEP20"
	GatewayBINANCEPAY = "BINANCEPAY" // 币安支付（Binance Pay）：入金后 checkoutUrl 跳转币安 App/Web
)
