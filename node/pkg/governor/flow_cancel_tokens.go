package governor

// FlowCancelTokenList returns a list of `tokenConfigEntry`s representing tokens that can 'Flow Cancel'. This means that incoming transfers
// that use these tokens can reduce the 'daily usage' of the Governor configured for the destination chain.
// The list of tokens was generated by grepping the file `generated_mainnet_tokens.go` for "USDC", "USDT", and "DAI".
//
// Only tokens that are configured in the mainnet token list should be able to flow cancel. That is, if a token is
// present in this list but not in the mainnet token lists, it should not flow cancel.
//
// Note that the field `symbol` is unused. It is retained in this file only for convenience.
func FlowCancelTokenList() []tokenConfigEntry {
	return []tokenConfigEntry{
		// USDC variants
		{chain: 1, addr: "c6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d61", symbol: "USDC"},
		{chain: 2, addr: "000000000000000000000000a0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", symbol: "USDC"},

		// USDT variants
		{chain: 1, addr: "ce010e60afedb22717bd63192f54145a3f965a33bb82d2c7029eb2ce1e208264", symbol: "USDT"},
		{chain: 2, addr: "000000000000000000000000dac17f958d2ee523a2206206994597c13d831ec7", symbol: "USDT"},

		// DAI variants
		{chain: 2, addr: "0000000000000000000000006b175474e89094c44da98b954eedeac495271d0f", symbol: "DAI"},
	}
}