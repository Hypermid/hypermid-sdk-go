// Package hypermid provides a Go client for the Hypermid Partner API.
//
// The SDK supports all Hypermid API endpoints including cross-chain swaps
// (LI.FI + Near Intents), fiat on-ramp, partner analytics, and webhooks.
//
// Usage:
//
//	client := hypermid.New(nil) // anonymous, 100 req/min
//
//	// With API key for higher rate limits (2000 req/min)
//	client := hypermid.New(&hypermid.Config{
//	    APIKey: "your-api-key",
//	})
//
//	chains, err := client.GetChains(context.Background())
package hypermid

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.hypermid.io"
	defaultTimeout = 30 * time.Second
	apiVersion     = "/v1"
)

// Config configures the Hypermid client.
type Config struct {
	// APIKey for authenticated access (2000 req/min, partner fee tier).
	// Optional — anonymous access allows 100 req/min.
	APIKey string
	// BaseURL override (default: https://api.hypermid.io).
	BaseURL string
	// Timeout for HTTP requests (default: 30s).
	Timeout time.Duration
	// HTTPClient is a custom http.Client to use. If nil, a default client is created.
	HTTPClient *http.Client
}

// Client is the Hypermid API client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// New creates a new Hypermid client. Pass nil for default configuration.
func New(cfg *Config) *Client {
	c := &Client{
		baseURL: defaultBaseURL,
	}

	if cfg != nil {
		if cfg.BaseURL != "" {
			c.baseURL = strings.TrimRight(cfg.BaseURL, "/")
		}
		c.apiKey = cfg.APIKey

		if cfg.HTTPClient != nil {
			c.httpClient = cfg.HTTPClient
		} else {
			timeout := cfg.Timeout
			if timeout == 0 {
				timeout = defaultTimeout
			}
			c.httpClient = &http.Client{Timeout: timeout}
		}
	} else {
		c.httpClient = &http.Client{Timeout: defaultTimeout}
	}

	return c
}

// ─── Internal HTTP helpers ───────────────────────────────────────────────

func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values, body interface{}) (json.RawMessage, error) {
	u := c.baseURL + apiVersion + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, &HypermidNetworkError{Msg: "failed to marshal request body", Cause: err}
		}
		bodyReader = strings.NewReader(string(b))
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, &HypermidNetworkError{Msg: "failed to create request", Cause: err}
	}

	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}
	if method == http.MethodPost || method == http.MethodDelete {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, &HypermidTimeoutError{TimeoutMs: int(c.httpClient.Timeout.Milliseconds())}
		}
		return nil, &HypermidNetworkError{Msg: "request failed", Cause: err}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &HypermidNetworkError{Msg: "failed to read response body", Cause: err}
	}

	var envelope apiResponse
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, &HypermidNetworkError{Msg: fmt.Sprintf("invalid JSON response (HTTP %d)", resp.StatusCode)}
	}

	if envelope.Error != nil {
		return nil, &HypermidError{
			Code:    envelope.Error.Code,
			Msg:     envelope.Error.Message,
			Status:  resp.StatusCode,
			Meta:    envelope.Meta,
			Details: envelope.Error.Details,
		}
	}

	return envelope.Data, nil
}

func (c *Client) doGet(ctx context.Context, path string, query url.Values) (json.RawMessage, error) {
	return c.doRequest(ctx, http.MethodGet, path, query, nil)
}

func (c *Client) doPost(ctx context.Context, path string, body interface{}) (json.RawMessage, error) {
	return c.doRequest(ctx, http.MethodPost, path, nil, body)
}

func (c *Client) doDelete(ctx context.Context, path string) (json.RawMessage, error) {
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// unmarshal is a helper to decode json.RawMessage into a typed value.
func unmarshal[T any](data json.RawMessage, err error) (T, error) {
	var zero T
	if err != nil {
		return zero, err
	}
	if err := json.Unmarshal(data, &zero); err != nil {
		return zero, &HypermidNetworkError{Msg: "failed to decode response data"}
	}
	return zero, nil
}

// ─── Core Swap Endpoints ─────────────────────────────────────────────────

// GetChains returns all supported chains (LI.FI + Near Intents).
func (c *Client) GetChains(ctx context.Context) (ChainsResponse, error) {
	return unmarshal[ChainsResponse](c.doGet(ctx, "/chains", nil))
}

// GetTokens returns available tokens, optionally filtered by chains and keywords.
func (c *Client) GetTokens(ctx context.Context, params *TokensParams) (TokensResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.Chains != "" {
			q.Set("chains", params.Chains)
		}
		if params.Keywords != "" {
			q.Set("keywords", params.Keywords)
		}
	}
	return unmarshal[TokensResponse](c.doGet(ctx, "/tokens", q))
}

// GetConnections returns available connections (which token pairs can be swapped).
func (c *Client) GetConnections(ctx context.Context, params ConnectionsParams) (json.RawMessage, error) {
	q := url.Values{}
	q.Set("fromChain", params.FromChain)
	q.Set("fromToken", params.FromToken)
	if params.ToChain != "" {
		q.Set("toChain", params.ToChain)
	}
	return c.doGet(ctx, "/connections", q)
}

// GetTools returns available bridge/swap tools.
func (c *Client) GetTools(ctx context.Context) (json.RawMessage, error) {
	return c.doGet(ctx, "/tools", nil)
}

// GetGasPrices returns gas prices for specified chains.
func (c *Client) GetGasPrices(ctx context.Context, params GasPricesParams) (json.RawMessage, error) {
	q := url.Values{}
	q.Set("chains", params.Chains)
	return c.doGet(ctx, "/gas-prices", q)
}

// GetQuote returns the best swap quote for a token pair.
func (c *Client) GetQuote(ctx context.Context, params QuoteParams) (QuoteResponse, error) {
	q := url.Values{}
	q.Set("fromChain", params.FromChain)
	q.Set("fromToken", params.FromToken)
	q.Set("fromAmount", params.FromAmount)
	q.Set("toChain", params.ToChain)
	q.Set("toToken", params.ToToken)
	q.Set("fromAddress", params.FromAddress)
	if params.ToAddress != "" {
		q.Set("toAddress", params.ToAddress)
	}
	if params.Slippage != "" {
		q.Set("slippage", params.Slippage)
	}
	if params.Order != "" {
		q.Set("order", params.Order)
	}
	return unmarshal[QuoteResponse](c.doGet(ctx, "/quote", q))
}

// GetRoutes returns available routes for a token pair (multi-route comparison).
func (c *Client) GetRoutes(ctx context.Context, params RoutesParams) (json.RawMessage, error) {
	body := map[string]interface{}{
		"fromChain":   params.FromChain,
		"fromToken":   params.FromToken,
		"fromAmount":  params.FromAmount,
		"toChain":     params.ToChain,
		"toToken":     params.ToToken,
		"fromAddress": params.FromAddress,
	}
	if params.ToAddress != "" {
		body["toAddress"] = params.ToAddress
	}
	if params.Slippage != nil {
		body["slippage"] = params.Slippage
	}
	if params.Order != "" {
		body["order"] = params.Order
	}
	return c.doPost(ctx, "/routes", body)
}

// GetStatus checks the status of a cross-chain swap.
func (c *Client) GetStatus(ctx context.Context, params StatusParams) (StatusResponse, error) {
	q := url.Values{}
	if params.Provider == "near-intents" {
		q.Set("provider", "near-intents")
		q.Set("correlationId", params.CorrelationID)
	} else {
		q.Set("txHash", params.TxHash)
		if params.Bridge != "" {
			q.Set("bridge", params.Bridge)
		}
		if params.FromChain != "" {
			q.Set("fromChain", params.FromChain)
		}
		if params.ToChain != "" {
			q.Set("toChain", params.ToChain)
		}
	}
	return unmarshal[StatusResponse](c.doGet(ctx, "/status", q))
}

// ─── Execute ─────────────────────────────────────────────────────────────

// Execute returns full transaction data for execution.
func (c *Client) Execute(ctx context.Context, params ExecuteParams) (ExecuteResponse, error) {
	body := map[string]interface{}{
		"fromChain":   params.FromChain,
		"fromToken":   params.FromToken,
		"fromAmount":  params.FromAmount,
		"toChain":     params.ToChain,
		"toToken":     params.ToToken,
		"fromAddress": params.FromAddress,
		"toAddress":   params.ToAddress,
	}
	if params.DepositMode != "" {
		body["depositMode"] = params.DepositMode
	}
	if params.Slippage != "" {
		body["slippage"] = params.Slippage
	}
	if params.Order != "" {
		body["order"] = params.Order
	}
	if params.RefundAddress != "" {
		body["refundAddress"] = params.RefundAddress
	}
	return unmarshal[ExecuteResponse](c.doPost(ctx, "/execute", body))
}

// SubmitDeposit submits a deposit transaction hash after sending tokens
// to a Near Intents deposit address.
func (c *Client) SubmitDeposit(ctx context.Context, params DepositSubmitParams) (DepositSubmitResponse, error) {
	return unmarshal[DepositSubmitResponse](c.doPost(ctx, "/execute/deposit/submit", params))
}

// GetDepositStatus checks the status of a Near Intents deposit/swap.
func (c *Client) GetDepositStatus(ctx context.Context, params DepositStatusParams) (DepositStatusResponse, error) {
	q := url.Values{}
	q.Set("depositAddress", params.DepositAddress)
	if params.DepositMemo != "" {
		q.Set("depositMemo", params.DepositMemo)
	}
	return unmarshal[DepositStatusResponse](c.doGet(ctx, "/execute/deposit/status", q))
}

// ─── On-Ramp ─────────────────────────────────────────────────────────────

// GetOnrampQuote returns a fiat-to-crypto price quote.
func (c *Client) GetOnrampQuote(ctx context.Context, params OnrampQuoteParams) (json.RawMessage, error) {
	return c.doPost(ctx, "/onramp/quote", params)
}

// CreateOnrampCheckout creates a fiat-to-crypto purchase session.
func (c *Client) CreateOnrampCheckout(ctx context.Context, params OnrampCheckoutParams) (OnrampCheckoutResponse, error) {
	return unmarshal[OnrampCheckoutResponse](c.doPost(ctx, "/onramp/checkout", params))
}

// GetOnrampStatus checks on-ramp order status.
func (c *Client) GetOnrampStatus(ctx context.Context, orderUID string) (OnrampStatusResponse, error) {
	q := url.Values{}
	q.Set("orderUid", orderUID)
	return unmarshal[OnrampStatusResponse](c.doGet(ctx, "/onramp/status", q))
}

// GetOnrampConfig returns supported chains and tokens for on-ramp.
func (c *Client) GetOnrampConfig(ctx context.Context) (OnrampConfigResponse, error) {
	return unmarshal[OnrampConfigResponse](c.doGet(ctx, "/onramp/config", nil))
}

// GetOnrampAssets returns asset configuration (min/max amounts, precision, payment methods).
func (c *Client) GetOnrampAssets(ctx context.Context, params OnrampAssetsParams) (json.RawMessage, error) {
	q := url.Values{}
	q.Set("currency", params.Currency)
	q.Set("chain", params.Chain)
	if params.OrderCurrency != "" {
		q.Set("orderCurrency", params.OrderCurrency)
	}
	return c.doGet(ctx, "/onramp/assets", q)
}

// ─── Swap Event ──────────────────────────────────────────────────────────

// RecordSwapEvent records a swap event for analytics.
func (c *Client) RecordSwapEvent(ctx context.Context, params SwapEventParams) (SwapEventResponse, error) {
	return unmarshal[SwapEventResponse](c.doPost(ctx, "/swap-event", params))
}

// ─── Partner (requires API key) ──────────────────────────────────────────

// GetPartnerInfo returns partner account details.
func (c *Client) GetPartnerInfo(ctx context.Context) (PartnerInfo, error) {
	return unmarshal[PartnerInfo](c.doGet(ctx, "/partner/me", nil))
}

// GetPartnerStats returns aggregated partner statistics.
func (c *Client) GetPartnerStats(ctx context.Context, params *PartnerStatsParams) (PartnerStats, error) {
	q := url.Values{}
	if params != nil {
		if params.From != "" {
			q.Set("from", params.From)
		}
		if params.To != "" {
			q.Set("to", params.To)
		}
	}
	return unmarshal[PartnerStats](c.doGet(ctx, "/partner/stats", q))
}

// GetPartnerTransactions returns paginated transaction history.
func (c *Client) GetPartnerTransactions(ctx context.Context, params *PaginationParams) (PaginatedTransactions, error) {
	q := url.Values{}
	if params != nil {
		if params.Page > 0 {
			q.Set("page", fmt.Sprintf("%d", params.Page))
		}
		if params.Limit > 0 {
			q.Set("limit", fmt.Sprintf("%d", params.Limit))
		}
	}
	return unmarshal[PaginatedTransactions](c.doGet(ctx, "/partner/transactions", q))
}

// ─── Webhooks (requires API key) ─────────────────────────────────────────

// CreateWebhook registers a webhook endpoint. The returned WebhookCreated
// includes the signing secret, which is only returned on creation.
func (c *Client) CreateWebhook(ctx context.Context, params CreateWebhookParams) (WebhookCreated, error) {
	return unmarshal[WebhookCreated](c.doPost(ctx, "/partner/webhooks", params))
}

// ListWebhooks returns all registered webhooks.
func (c *Client) ListWebhooks(ctx context.Context) (WebhooksListResponse, error) {
	return unmarshal[WebhooksListResponse](c.doGet(ctx, "/partner/webhooks", nil))
}

// DeleteWebhook deletes a webhook by ID.
func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) (DeleteWebhookResponse, error) {
	return unmarshal[DeleteWebhookResponse](c.doDelete(ctx, "/partner/webhooks/"+webhookID))
}

// ─── Health Check ────────────────────────────────────────────────────────

// Ping returns API health status, version, uptime, and provider statuses.
func (c *Client) Ping(ctx context.Context) (PingResponse, error) {
	return unmarshal[PingResponse](c.doGet(ctx, "/ping", nil))
}
