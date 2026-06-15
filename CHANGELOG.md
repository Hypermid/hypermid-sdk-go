# Changelog

## v1.1.0

### Added

- **SuperSwap V2 execute fields** on `ExecuteResponse` (Provider == "superswap"):
  `Source`, `ApprovalAddress`, `EstimatedOutput`, `MinOutput`, `V2`. Approve
  `ApprovalAddress` (== `TransactionRequest.To`), then sign & broadcast
  `TransactionRequest`.
- **`StatusResponse.Extra` now populated.** Added a custom `UnmarshalJSON` that
  decodes `provider`/`status` and collects every other key into `Extra` — so
  SuperSwap V2 status fields (`hyperlaneMessageId`, `subStatus`, `sending`,
  `receiving`, `destinationTxHash`) are no longer dropped. (Previously `Extra`
  had `json:"-"` and was never filled.)
- `IsSuperSwapRoute(*ExecuteResponse) bool`.
- `IsSuperSwapStatusTerminal(string) bool` (terminal: `DONE`, `FAILED`).
- `IsWalletDeposit` now returns true for SuperSwap routes.
