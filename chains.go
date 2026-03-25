package hypermid

const niBase = 900_000_000

// ChainID constants for all supported chains.
const (
	// EVM Chains
	ChainIDEthereum  = 1
	ChainIDOptimism  = 10
	ChainIDBSC       = 56
	ChainIDGnosis    = 100
	ChainIDPolygon   = 137
	ChainIDXLayer    = 196
	ChainIDArbitrum  = 42161
	ChainIDAvalanche = 43114
	ChainIDBase      = 8453
	ChainIDPlasma    = 1012
	ChainIDBerachain = 80094
	ChainIDMonad     = 10143

	// Non-EVM (LI.FI supported)
	ChainIDSolana  = 1151111081099710
	ChainIDBitcoin = 20000000000001
	ChainIDSui     = 9270000000000000

	// Near Intents-only chains
	ChainIDNear        = niBase + 1
	ChainIDTon         = niBase + 2
	ChainIDTron        = niBase + 3
	ChainIDXRP         = niBase + 4
	ChainIDDogecoin    = niBase + 5
	ChainIDLitecoin    = niBase + 6
	ChainIDBitcoinCash = niBase + 7
	ChainIDStellar     = niBase + 8
	ChainIDCardano     = niBase + 9
	ChainIDAptos       = niBase + 10
	ChainIDStarknet    = niBase + 11
	ChainIDDash        = niBase + 12
	ChainIDZcash       = niBase + 13
	ChainIDAleo        = niBase + 14
	ChainIDAdi         = niBase + 15
)

// IsNearIntentsChain returns true if the chain ID belongs to a Near Intents-only chain.
func IsNearIntentsChain(chainID int) bool {
	return chainID >= niBase && chainID < niBase+1000
}

// SupportsWalletDeposit returns true if the chain supports wallet-connected deposit mode.
// Chains with wallet connectors: EVM, Solana, Bitcoin, Sui, TON, Tron.
// Other Near Intents chains (NEAR, XRP, DOGE, etc.) require manual deposit.
func SupportsWalletDeposit(chainID int) bool {
	if chainID > 0 && chainID < niBase {
		return true // EVM
	}
	switch chainID {
	case ChainIDSolana, ChainIDBitcoin, ChainIDSui, ChainIDTon, ChainIDTron:
		return true
	}
	return false
}
