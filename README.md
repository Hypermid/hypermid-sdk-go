# hypermid-sdk-go

Go SDK for the [Hypermid](https://hypermid.io) Partner API — swap,
bridge, and on-ramp across 90+ chains (EVM, Solana, Bitcoin, Sui, NEAR,
Tron, TON, XRP, Doge).

```bash
go get github.com/Hypermid/hypermid-sdk-go
```

## Quick start

No API key required. The SDK works anonymously out of the box at the
default fee tier — pass an API key only if you're a partner with
custom fee terms.

```go
package main

import (
    "context"
    "fmt"
    "log"

    hypermid "github.com/Hypermid/hypermid-sdk-go"
)

func main() {
    // Anonymous — works immediately, no signup
    hm := hypermid.New(nil)

    // Or, partner with custom fees:
    // hm := hypermid.New(&hypermid.Config{APIKey: os.Getenv("HYPERMID_API_KEY")})

    chains, err := hm.GetChains(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Supported chains: %d\n", len(chains.Chains))
}
```

## Features

- `GetQuote` / `Execute` / `GetStatus` — the swap pipeline
- `GetChains` / `GetTokens` — supported chains and tokens
- `GetBalances` — multi-ecosystem wallet balances + USD totals
- `RegisterInboundReceiver` — SuperSwap V2 inbound deposits
- On-ramp helpers — `GetOnrampQuote`, `GetOnrampCheckout`, `GetOnrampStatus`
- `VerifyWebhookSignature` — HMAC webhook signature verification

## Authentication

The API is open by default — every endpoint works without
authentication, so you can integrate, test, and ship without a signup.

An **API key is only needed if you're a partner** with negotiated terms
(custom fee splits, fee discounts, volume tiers, higher rate limits,
webhook events scoped to your traffic). When set, the SDK sends it as
the `X-API-Key` header.

Apply for a partner account at [partner.hypermid.io](https://partner.hypermid.io).

## Documentation

Full reference: <https://docs.hypermid.io>

## License

MIT
