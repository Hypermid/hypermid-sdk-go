package hypermid

import (
	"context"
	"fmt"
	"time"
)

const (
	defaultPollInterval = 5 * time.Second
	defaultMaxWait      = 10 * time.Minute
)

func pollConfig(cfg *PollConfig) (interval time.Duration, maxWait time.Duration, maxPolls int) {
	interval = defaultPollInterval
	maxWait = defaultMaxWait
	maxPolls = 0 // unlimited

	if cfg != nil {
		if cfg.PollIntervalMs > 0 {
			interval = time.Duration(cfg.PollIntervalMs) * time.Millisecond
		}
		if cfg.MaxWaitMs > 0 {
			maxWait = time.Duration(cfg.MaxWaitMs) * time.Millisecond
		}
		if cfg.MaxPolls > 0 {
			maxPolls = cfg.MaxPolls
		}
	}
	return
}

// WaitForDepositCompletion polls Near Intents deposit/swap status until a terminal
// state is reached (SUCCESS, REFUNDED, or FAILED).
func (c *Client) WaitForDepositCompletion(ctx context.Context, params DepositStatusParams, cfg *PollConfig) (DepositStatusResponse, error) {
	interval, maxWait, maxPolls := pollConfig(cfg)
	start := time.Now()
	polls := 0

	for {
		status, err := c.GetDepositStatus(ctx, params)
		if err != nil {
			return DepositStatusResponse{}, err
		}
		polls++

		if IsNIStatusTerminal(status.Status) {
			return status, nil
		}

		if time.Since(start) >= maxWait {
			return DepositStatusResponse{}, &HypermidPollTimeoutError{
				Msg: fmt.Sprintf("deposit status polling timed out after %v (last status: %s)", maxWait, status.Status),
			}
		}
		if maxPolls > 0 && polls >= maxPolls {
			return DepositStatusResponse{}, &HypermidPollTimeoutError{
				Msg: fmt.Sprintf("deposit status polling exceeded %d attempts (last status: %s)", maxPolls, status.Status),
			}
		}

		select {
		case <-ctx.Done():
			return DepositStatusResponse{}, ctx.Err()
		case <-time.After(interval):
		}
	}
}

// WaitForLiFiCompletion polls LI.FI swap status until a terminal state is reached
// (DONE or FAILED).
func (c *Client) WaitForLiFiCompletion(ctx context.Context, params LiFiStatusParams, cfg *PollConfig) (StatusResponse, error) {
	interval, maxWait, maxPolls := pollConfig(cfg)
	start := time.Now()
	polls := 0

	for {
		status, err := c.GetStatus(ctx, StatusParams{
			TxHash:    params.TxHash,
			Bridge:    params.Bridge,
			FromChain: params.FromChain,
			ToChain:   params.ToChain,
		})
		if err != nil {
			return StatusResponse{}, err
		}
		polls++

		if status.Status != "" && IsLiFiStatusTerminal(status.Status) {
			return status, nil
		}

		if time.Since(start) >= maxWait {
			return StatusResponse{}, &HypermidPollTimeoutError{
				Msg: fmt.Sprintf("LI.FI status polling timed out after %v (last status: %s)", maxWait, status.Status),
			}
		}
		if maxPolls > 0 && polls >= maxPolls {
			return StatusResponse{}, &HypermidPollTimeoutError{
				Msg: fmt.Sprintf("LI.FI status polling exceeded %d attempts (last status: %s)", maxPolls, status.Status),
			}
		}

		select {
		case <-ctx.Done():
			return StatusResponse{}, ctx.Err()
		case <-time.After(interval):
		}
	}
}

// ExecuteSwapResult contains the final state of an ExecuteSwap call.
type ExecuteSwapResult struct {
	// Provider is "lifi" or "near-intents".
	Provider string
	// ExecuteResponse is the initial execute response.
	ExecuteResponse *ExecuteResponse
	// DepositStatus is the final Near Intents deposit status (if applicable).
	DepositStatus *DepositStatusResponse
	// LiFiStatus is the final LI.FI status (if applicable).
	LiFiStatus *StatusResponse
	// Error message if the swap failed.
	Error string
}

// ExecuteSwapHooks provides callback hooks for the swap execution lifecycle.
type ExecuteSwapHooks struct {
	// OnExecute is called when the execute response is received.
	OnExecute func(resp *ExecuteResponse)
	// OnTransactionRequest is called when a LI.FI transactionRequest is ready.
	// You MUST sign and broadcast it, then return the transaction hash.
	// If nil, ExecuteSwap returns after receiving the transactionRequest.
	OnTransactionRequest func(resp *ExecuteResponse) (txHash string, err error)
	// OnDepositRequired is called when a Near Intents deposit address is ready
	// and depositMode is "wallet". You MUST send tokens and return the tx hash.
	// If nil, ExecuteSwap returns after receiving the deposit address.
	OnDepositRequired func(resp *ExecuteResponse) (txHash string, err error)
}

// ExecuteSwap runs the full swap lifecycle: execute -> sign/deposit -> poll -> complete.
func (c *Client) ExecuteSwap(ctx context.Context, params ExecuteParams, hooks *ExecuteSwapHooks, cfg *PollConfig) (*ExecuteSwapResult, error) {
	execResp, err := c.Execute(ctx, params)
	if err != nil {
		return nil, err
	}

	if hooks != nil && hooks.OnExecute != nil {
		hooks.OnExecute(&execResp)
	}

	result := &ExecuteSwapResult{
		Provider:        execResp.Provider,
		ExecuteResponse: &execResp,
	}

	if IsLiFiRoute(&execResp) {
		if hooks == nil || hooks.OnTransactionRequest == nil {
			return result, nil
		}

		txHash, err := hooks.OnTransactionRequest(&execResp)
		if err != nil {
			return nil, err
		}

		finalStatus, err := c.WaitForLiFiCompletion(ctx, LiFiStatusParams{
			TxHash:    txHash,
			FromChain: params.FromChain,
			ToChain:   params.ToChain,
		}, cfg)
		if err != nil {
			return nil, err
		}

		result.LiFiStatus = &finalStatus
		if finalStatus.Status == "FAILED" {
			result.Error = "LI.FI swap failed"
		}
		return result, nil
	}

	if IsNearIntentsRoute(&execResp) {
		isManual := execResp.DepositMode == "manual"

		if isManual {
			finalStatus, err := c.WaitForDepositCompletion(ctx, DepositStatusParams{
				DepositAddress: execResp.DepositAddress,
				DepositMemo:    execResp.DepositMemo,
			}, cfg)
			if err != nil {
				return nil, err
			}
			result.DepositStatus = &finalStatus
			if finalStatus.Status == "FAILED" {
				result.Error = "Near Intents swap failed"
			}
			return result, nil
		}

		// Wallet deposit
		if hooks == nil || hooks.OnDepositRequired == nil {
			return result, nil
		}

		txHash, err := hooks.OnDepositRequired(&execResp)
		if err != nil {
			return nil, err
		}

		_, submitErr := c.SubmitDeposit(ctx, DepositSubmitParams{
			TxHash:         txHash,
			DepositAddress: execResp.DepositAddress,
		})
		if submitErr != nil {
			return nil, submitErr
		}

		finalStatus, err := c.WaitForDepositCompletion(ctx, DepositStatusParams{
			DepositAddress: execResp.DepositAddress,
			DepositMemo:    execResp.DepositMemo,
		}, cfg)
		if err != nil {
			return nil, err
		}
		result.DepositStatus = &finalStatus
		if finalStatus.Status == "FAILED" {
			result.Error = "Near Intents swap failed"
		}
		return result, nil
	}

	result.Error = "unknown provider in execute response"
	return result, nil
}
