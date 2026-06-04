package hypermid

import "encoding/json"

// apiResponse is the standard API response envelope.
// All API endpoints return this shape.
type apiResponse struct {
	Data  json.RawMessage `json:"data"`
	Error *ApiError       `json:"error"`
	Meta  ApiMeta         `json:"meta"`
}

// ─── Chains ──────────────────────────────────────────────────────────────

// NativeToken represents a chain's native token.
type NativeToken struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
}

// Chain represents a supported blockchain.
type Chain struct {
	ID          int         `json:"id"`
	Key         string      `json:"key"`
	Name        string      `json:"name"`
	ChainType   string      `json:"chainType"`
	NativeToken NativeToken `json:"nativeToken"`
	Provider    string      `json:"provider,omitempty"`
}

// ChainsResponse is the response from GetChains.
type ChainsResponse struct {
	Chains []Chain `json:"chains"`
}

// ─── Tokens ──────────────────────────────────────────────────────────────

// Token represents a token on a chain.
type Token struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	ChainID  int    `json:"chainId"`
	LogoURI  string `json:"logoURI,omitempty"`
	PriceUSD string `json:"priceUSD,omitempty"`
}

// TokensResponse is the response from GetTokens.
type TokensResponse struct {
	Tokens map[string][]Token `json:"tokens"`
}

// TokensParams are optional parameters for GetTokens.
type TokensParams struct {
	Chains   string `json:"chains,omitempty"`
	Keywords string `json:"keywords,omitempty"`
}

// ─── Connections ─────────────────────────────────────────────────────────

// ConnectionsParams are parameters for GetConnections.
type ConnectionsParams struct {
	FromChain string `json:"fromChain"`
	FromToken string `json:"fromToken"`
	ToChain   string `json:"toChain,omitempty"`
}

// ─── Gas ─────────────────────────────────────────────────────────────────

// GasPricesParams are parameters for GetGasPrices.
type GasPricesParams struct {
	Chains string `json:"chains"`
}

// ─── Quote ───────────────────────────────────────────────────────────────

// QuoteParams are parameters for GetQuote.
type QuoteParams struct {
	FromChain   string `json:"fromChain"`
	FromToken   string `json:"fromToken"`
	FromAmount  string `json:"fromAmount"`
	ToChain     string `json:"toChain"`
	ToToken     string `json:"toToken"`
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress,omitempty"`
	Slippage    string `json:"slippage,omitempty"`
	Order       string `json:"order,omitempty"`
}

// QuoteResponse is the response from GetQuote.
type QuoteResponse struct {
	Quote      json.RawMessage `json:"quote"`
	Provider   string          `json:"provider"`
	FeeBps     int             `json:"feeBps"`
	IsDryQuote bool            `json:"isDryQuote"`
}

// ─── Routes ──────────────────────────────────────────────────────────────

// RoutesParams are parameters for GetRoutes.
type RoutesParams struct {
	FromChain   string      `json:"fromChain"`
	FromToken   string      `json:"fromToken"`
	FromAmount  string      `json:"fromAmount"`
	ToChain     string      `json:"toChain"`
	ToToken     string      `json:"toToken"`
	FromAddress string      `json:"fromAddress"`
	ToAddress   string      `json:"toAddress,omitempty"`
	Slippage    interface{} `json:"slippage,omitempty"`
	Order       string      `json:"order,omitempty"`
}

// ─── Status ──────────────────────────────────────────────────────────────

// StatusParams are parameters for GetStatus.
// For LI.FI status: set TxHash (and optionally Bridge, FromChain, ToChain).
// For Near Intents status: set Provider to "near-intents" and CorrelationID.
type StatusParams struct {
	// LI.FI fields
	TxHash    string `json:"txHash,omitempty"`
	Bridge    string `json:"bridge,omitempty"`
	FromChain string `json:"fromChain,omitempty"`
	ToChain   string `json:"toChain,omitempty"`

	// Near Intents fields
	Provider      string `json:"provider,omitempty"`
	CorrelationID string `json:"correlationId,omitempty"`
}

// StatusResponse is the response from GetStatus.
type StatusResponse struct {
	Provider string `json:"provider"`
	Status   string `json:"status,omitempty"`
	// Additional fields vary by provider; decode with json.RawMessage if needed.
	Extra map[string]interface{} `json:"-"`
}

// LiFiStatusParams are parameters for polling LI.FI swap status.
type LiFiStatusParams struct {
	TxHash    string `json:"txHash"`
	Bridge    string `json:"bridge,omitempty"`
	FromChain string `json:"fromChain,omitempty"`
	ToChain   string `json:"toChain,omitempty"`
}

// ─── Execute ─────────────────────────────────────────────────────────────

// ExecuteParams are parameters for Execute.
type ExecuteParams struct {
	FromChain    string `json:"fromChain"`
	FromToken    string `json:"fromToken"`
	FromAmount   string `json:"fromAmount"`
	ToChain      string `json:"toChain"`
	ToToken      string `json:"toToken"`
	FromAddress  string `json:"fromAddress"`
	ToAddress    string `json:"toAddress"`
	DepositMode  string `json:"depositMode,omitempty"`
	Slippage     string `json:"slippage,omitempty"`
	Order        string `json:"order,omitempty"`
	RefundAddress string `json:"refundAddress,omitempty"`
}

// ExecuteResponse is the response from Execute.
// Check Provider to determine if it's a LI.FI or Near Intents response.
type ExecuteResponse struct {
	Provider    string `json:"provider"`
	DepositMode string `json:"depositMode,omitempty"`
	FeeBps      int    `json:"feeBps"`

	// LI.FI fields
	TransactionRequest *TransactionRequest `json:"transactionRequest,omitempty"`
	Quote              *ExecuteQuote       `json:"quote,omitempty"`

	// Near Intents fields
	DepositAddress  string   `json:"depositAddress,omitempty"`
	DepositMemo     string   `json:"depositMemo,omitempty"`
	ExpectedOutput  string   `json:"expectedOutput,omitempty"`
	ExpectedOutputUsd *float64 `json:"expectedOutputUsd,omitempty"`
	MinAmountOut    string   `json:"minAmountOut,omitempty"`
	TimeEstimate    *int     `json:"timeEstimate,omitempty"`
	CorrelationID   string   `json:"correlationId,omitempty"`

	// Instructions (present for both providers)
	Instructions map[string]string `json:"instructions,omitempty"`
}

// TransactionRequest contains the transaction data to sign and broadcast (LI.FI).
type TransactionRequest struct {
	To       string `json:"to"`
	Data     string `json:"data"`
	Value    string `json:"value"`
	From     string `json:"from"`
	ChainID  int    `json:"chainId"`
	GasLimit string `json:"gasLimit,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
}

// ExecuteQuote contains quote details returned with a LI.FI execute response.
type ExecuteQuote struct {
	FromToken     Token           `json:"fromToken"`
	ToToken       Token           `json:"toToken"`
	FromAmount    string          `json:"fromAmount"`
	ToAmount      string          `json:"toAmount"`
	ToAmountMin   string          `json:"toAmountMin"`
	EstimatedTime int             `json:"estimatedTime"`
	GasCosts      json.RawMessage `json:"gasCosts"`
	FeeCosts      json.RawMessage `json:"feeCosts"`
}

// ─── Deposit (Near Intents) ──────────────────────────────────────────────

// DepositSubmitParams are parameters for SubmitDeposit.
type DepositSubmitParams struct {
	TxHash         string `json:"txHash"`
	DepositAddress string `json:"depositAddress"`
}

// DepositSubmitResponse is the response from SubmitDeposit.
type DepositSubmitResponse struct {
	Submitted      bool   `json:"submitted"`
	TxHash         string `json:"txHash"`
	DepositAddress string `json:"depositAddress"`
	NextStep       string `json:"nextStep"`
}

// DepositStatusParams are parameters for GetDepositStatus.
type DepositStatusParams struct {
	DepositAddress string `json:"depositAddress"`
	DepositMemo    string `json:"depositMemo,omitempty"`
}

// SwapDetails contains swap completion details for Near Intents deposits.
type SwapDetails struct {
	AmountOut              string   `json:"amountOut,omitempty"`
	AmountOutFormatted     string   `json:"amountOutFormatted,omitempty"`
	AmountOutUsd           *float64 `json:"amountOutUsd,omitempty"`
	DestinationChainTxHashes []string `json:"destinationChainTxHashes,omitempty"`
	RefundedAmount         string   `json:"refundedAmount,omitempty"`
	RefundReason           string   `json:"refundReason,omitempty"`
}

// DepositStatusResponse is the response from GetDepositStatus.
type DepositStatusResponse struct {
	Provider       string       `json:"provider"`
	Status         string       `json:"status"`
	DepositAddress string       `json:"depositAddress"`
	SwapDetails    *SwapDetails `json:"swapDetails,omitempty"`
}

// ─── On-Ramp ─────────────────────────────────────────────────────────────

// OnrampQuoteParams are parameters for GetOnrampQuote.
type OnrampQuoteParams struct {
	FiatAmount    interface{} `json:"fiatAmount"`
	FiatCurrency  string      `json:"fiatCurrency"`
	CryptoToken   string      `json:"cryptoToken"`
	CryptoChain   string      `json:"cryptoChain"`
	WalletAddress string      `json:"walletAddress,omitempty"`
	PaymentMode   string      `json:"paymentMode,omitempty"`
	UserCountry   string      `json:"userCountry,omitempty"`
}

// OnrampCheckoutParams are parameters for CreateOnrampCheckout.
type OnrampCheckoutParams struct {
	WalletAddress string      `json:"walletAddress"`
	CryptoToken   string      `json:"cryptoToken"`
	CryptoChain   string      `json:"cryptoChain"`
	FiatCurrency  string      `json:"fiatCurrency"`
	FiatAmount    interface{} `json:"fiatAmount"`
	Email         string      `json:"email,omitempty"`
	ReturnURL     string      `json:"returnUrl,omitempty"`
	PaymentMode   string      `json:"paymentMode,omitempty"`
}

// OnrampCheckoutResponse is the response from CreateOnrampCheckout.
type OnrampCheckoutResponse struct {
	RedirectURL      string `json:"redirectUrl"`
	OrderUID         string `json:"orderUid"`
	ExternalOrderUID string `json:"externalOrderUid"`
}

// OnrampStatusResponse is the response from GetOnrampStatus.
type OnrampStatusResponse struct {
	Status   string `json:"status"`
	OrderUID string `json:"orderUid"`
	DstAmount string `json:"dstAmount,omitempty"`
	TxHash   string `json:"txHash,omitempty"`
	Message  string `json:"message,omitempty"`
}

// OnrampConfigResponse is the response from GetOnrampConfig.
type OnrampConfigResponse struct {
	Chains map[string][]string `json:"chains"`
}

// OnrampAssetsParams are parameters for GetOnrampAssets.
type OnrampAssetsParams struct {
	Currency      string `json:"currency"`
	Chain         string `json:"chain"`
	OrderCurrency string `json:"orderCurrency,omitempty"`
}

// ─── Swap Event ──────────────────────────────────────────────────────────

// SwapEventParams are parameters for RecordSwapEvent.
type SwapEventParams struct {
	Provider        string   `json:"provider,omitempty"`
	FromChain       string   `json:"from_chain,omitempty"`
	FromToken       string   `json:"from_token,omitempty"`
	ToChain         string   `json:"to_chain,omitempty"`
	ToToken         string   `json:"to_token,omitempty"`
	AmountUsd       *float64 `json:"amount_usd,omitempty"`
	FeeUsd          *float64 `json:"fee_usd,omitempty"`
	TxHash          string   `json:"tx_hash,omitempty"`
	WalletHash      string   `json:"wallet_hash,omitempty"`
	Status          string   `json:"status"`
	FromAmount      string   `json:"from_amount,omitempty"`
	ToAmount        string   `json:"to_amount,omitempty"`
	DurationSeconds *int     `json:"duration_seconds,omitempty"`
	ErrorMessage    string   `json:"error_message,omitempty"`
}

// SwapEventResponse is the response from RecordSwapEvent.
type SwapEventResponse struct {
	Updated bool  `json:"updated"`
	ID      int64 `json:"id"`
}

// ─── Partner ─────────────────────────────────────────────────────────────

// PartnerInfo contains partner account details.
type PartnerInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Status      string  `json:"status"`
	Tier        string  `json:"tier"`
	FeeBps      int     `json:"fee_bps"`
	VolumeTotal float64 `json:"volume_total"`
	TxCount     int     `json:"tx_count"`
	CreatedAt   string  `json:"created_at"`
}

// PartnerStats contains aggregated partner statistics.
type PartnerStats struct {
	TxCount            int                `json:"tx_count"`
	CompletedCount     int                `json:"completed_count"`
	FailedCount        int                `json:"failed_count"`
	VolumeUsd          float64            `json:"volume_usd"`
	FeesEarnedUsd      float64            `json:"fees_earned_usd"`
	AvgDurationSeconds float64            `json:"avg_duration_seconds"`
	ByChain            []ChainStat        `json:"by_chain"`
	ByProvider         []ProviderStat     `json:"by_provider"`
}

// ChainStat contains per-chain statistics.
type ChainStat struct {
	Chain  string  `json:"chain"`
	Count  int     `json:"count"`
	Volume float64 `json:"volume"`
}

// ProviderStat contains per-provider statistics.
type ProviderStat struct {
	Provider string `json:"provider"`
	Count    int    `json:"count"`
}

// PartnerStatsParams are optional parameters for GetPartnerStats.
type PartnerStatsParams struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// Transaction represents a single partner transaction.
type Transaction struct {
	ID              int     `json:"id"`
	Provider        string  `json:"provider"`
	FromChain       string  `json:"from_chain"`
	FromToken       string  `json:"from_token"`
	ToChain         string  `json:"to_chain"`
	ToToken         string  `json:"to_token"`
	AmountUsd       float64 `json:"amount_usd"`
	FeeUsd          float64 `json:"fee_usd"`
	TxHash          string  `json:"tx_hash"`
	WalletHash      string  `json:"wallet_hash"`
	Status          string  `json:"status"`
	FromAmount      string  `json:"from_amount"`
	ToAmount        string  `json:"to_amount"`
	DurationSeconds int     `json:"duration_seconds"`
	CreatedAt       string  `json:"created_at"`
}

// PaginatedTransactions is a paginated list of transactions.
type PaginatedTransactions struct {
	Items      []Transaction `json:"items"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"totalPages"`
}

// PaginationParams are optional pagination parameters.
type PaginationParams struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

// ─── Webhooks ────────────────────────────────────────────────────────────

// CreateWebhookParams are parameters for CreateWebhook.
type CreateWebhookParams struct {
	URL    string   `json:"url"`
	Events []string `json:"events,omitempty"`
}

// Webhook represents a registered webhook.
type Webhook struct {
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Events    []string `json:"events"`
	Status    string   `json:"status"`
	CreatedAt string   `json:"created_at"`
}

// WebhookCreated is returned from CreateWebhook and includes the signing secret.
type WebhookCreated struct {
	Webhook
	// Secret is the webhook signing secret. Only returned on creation.
	Secret string `json:"secret"`
}

// WebhooksListResponse is the response from ListWebhooks.
type WebhooksListResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

// DeleteWebhookResponse is the response from DeleteWebhook.
type DeleteWebhookResponse struct {
	Deleted bool   `json:"deleted"`
	ID      string `json:"id"`
}

// ─── Ping ────────────────────────────────────────────────────────────────

// PingProviders contains provider health statuses.
type PingProviders struct {
	LiFi        string `json:"lifi"`
	NearIntents string `json:"nearIntents"`
	RampNow     string `json:"rampnow"`
}

// PingResponse is the response from Ping.
type PingResponse struct {
	Status    string        `json:"status"`
	Version   string        `json:"version"`
	Uptime    float64       `json:"uptime"`
	Timestamp int64         `json:"timestamp"`
	Providers PingProviders `json:"providers"`
}

// ─── Poll Config ─────────────────────────────────────────────────────────

// PollConfig configures polling behavior for WaitForDepositCompletion and WaitForLiFiCompletion.
type PollConfig struct {
	// PollInterval is the time between polls (default: 5s).
	PollIntervalMs int
	// MaxWaitMs is the maximum total time to wait (default: 600000 = 10 min).
	MaxWaitMs int
	// MaxPolls is the maximum number of poll attempts (default: 0 = unlimited).
	MaxPolls int
}

// ─── Balances ────────────────────────────────────────────────────────────

// BalancesParams is the input for GetBalances. ChainIDs is optional and is
// sent as a comma-separated `chainIds` query param (not a JSON body).
type BalancesParams struct {
	Address  string `json:"-"`
	ChainIDs []int  `json:"-"`
}

// TokenBalance is a single token holding for an address.
type TokenBalance struct {
	ChainID    int      `json:"chainId"`
	Address    string   `json:"address"`
	Symbol     string   `json:"symbol"`
	Name       string   `json:"name"`
	Decimals   int      `json:"decimals"`
	Balance    string   `json:"balance"`
	PriceUSD   float64  `json:"priceUSD"`
	BalanceUSD float64  `json:"balanceUSD"`
	LogoURI    string   `json:"logoURI"`
	Providers  []string `json:"providers"`
}

// BalanceChainMeta is the per-chain fetch status (use to render retry chips
// for failing chains rather than hiding them).
type BalanceChainMeta struct {
	OK         bool   `json:"ok"`
	Error      string `json:"error,omitempty"`
	Source     string `json:"source,omitempty"`
	DurationMs int    `json:"durationMs"`
	Stale      bool   `json:"stale,omitempty"`
}

// BalancesResponse is the GetBalances result.
type BalancesResponse struct {
	Address         string                      `json:"address"`
	TotalBalanceUSD string                      `json:"totalBalanceUSD"`
	Balances        map[string][]TokenBalance   `json:"balances"`
	ChainMeta       map[string]BalanceChainMeta `json:"chainMeta,omitempty"`
	CachedAt        string                      `json:"cachedAt,omitempty"`
	CacheHit        bool                        `json:"cacheHit,omitempty"`
}

// ─── Inbound receiver (SuperSwap V2) ─────────────────────────────────────

// InboundReceiverParams registers a SuperSwap V2 inbound deposit.
type InboundReceiverParams struct {
	TxHash            string `json:"txHash"`
	FromAddress       string `json:"fromAddress"`
	ToAddress         string `json:"toAddress"`
	OutputToken       string `json:"outputToken"`
	DestinationDomain int    `json:"destinationDomain"`
	Signature         string `json:"signature"`
}

// InboundReceiverResponse is the RegisterInboundReceiver result.
type InboundReceiverResponse struct {
	Registered bool   `json:"registered"`
	RecordID   string `json:"recordId"`
	USDCAmount string `json:"usdcAmount"`
	Status     string `json:"status"`
}
