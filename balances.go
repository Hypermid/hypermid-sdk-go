package hypermid

import (
	"context"
	"net/url"
	"strconv"
	"strings"
)

// GetBalances returns multi-ecosystem token balances and the total USD value
// for an address. The backend auto-detects the address ecosystem
// (EVM / Sui / Tron / NEAR / Solana / Bitcoin); pass ChainIDs to restrict
// EVM coverage to specific chains.
func (c *Client) GetBalances(ctx context.Context, params BalancesParams) (BalancesResponse, error) {
	q := url.Values{}
	q.Set("address", params.Address)
	if len(params.ChainIDs) > 0 {
		ids := make([]string, len(params.ChainIDs))
		for i, id := range params.ChainIDs {
			ids[i] = strconv.Itoa(id)
		}
		q.Set("chainIds", strings.Join(ids, ","))
	}
	return unmarshal[BalancesResponse](c.doGet(ctx, "/balances", q))
}

// RegisterInboundReceiver registers a SuperSwap V2 inbound deposit so the
// backend executes the PulseChain-side output. The deposit must already be
// on-chain, and an EIP-712 signature over the registration is required.
func (c *Client) RegisterInboundReceiver(ctx context.Context, params InboundReceiverParams) (InboundReceiverResponse, error) {
	return unmarshal[InboundReceiverResponse](c.doPost(ctx, "/inbound-receiver/register", params))
}
