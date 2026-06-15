package hypermid

// IsLiFiRoute returns true if the execute response is a LI.FI route
// (has transactionRequest).
func IsLiFiRoute(resp *ExecuteResponse) bool {
	return resp.Provider == "lifi"
}

// IsNearIntentsRoute returns true if the execute response is a Near Intents route
// (has depositAddress).
func IsNearIntentsRoute(resp *ExecuteResponse) bool {
	return resp.Provider == "near-intents"
}

// IsSuperSwapRoute returns true if the execute response is a SuperSwap V2 route
// (wallet, has transactionRequest).
func IsSuperSwapRoute(resp *ExecuteResponse) bool {
	return resp.Provider == "superswap"
}

// IsManualDeposit returns true if a Near Intents deposit requires manual user action
// (QR code / copy address).
func IsManualDeposit(resp *ExecuteResponse) bool {
	return resp.Provider == "near-intents" && resp.DepositMode == "manual"
}

// IsWalletDeposit returns true if the deposit can be done programmatically via wallet.
func IsWalletDeposit(resp *ExecuteResponse) bool {
	return resp.Provider == "lifi" ||
		resp.Provider == "superswap" ||
		(resp.Provider == "near-intents" && resp.DepositMode == "wallet")
}

// IsNIStatusTerminal returns true if a Near Intents deposit status is terminal
// (no more polling needed).
func IsNIStatusTerminal(status string) bool {
	switch status {
	case "SUCCESS", "REFUNDED", "FAILED":
		return true
	}
	return false
}

// IsLiFiStatusTerminal returns true if a LI.FI status is terminal.
func IsLiFiStatusTerminal(status string) bool {
	switch status {
	case "DONE", "FAILED":
		return true
	}
	return false
}

// IsSuperSwapStatusTerminal returns true if a SuperSwap V2 status is terminal
// (vocabulary: PENDING | DONE | FAILED | NOT_FOUND | INVALID).
func IsSuperSwapStatusTerminal(status string) bool {
	switch status {
	case "DONE", "FAILED":
		return true
	}
	return false
}

// IsDepositSuccess returns true if a Near Intents swap completed successfully.
func IsDepositSuccess(resp *DepositStatusResponse) bool {
	return resp.Status == "SUCCESS"
}

// IsDepositRefunded returns true if a Near Intents swap was refunded.
func IsDepositRefunded(resp *DepositStatusResponse) bool {
	return resp.Status == "REFUNDED"
}

// IsDepositFailed returns true if a Near Intents swap failed.
func IsDepositFailed(resp *DepositStatusResponse) bool {
	return resp.Status == "FAILED"
}
