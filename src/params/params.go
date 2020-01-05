package params

/*
CODE GENERATED AUTOMATICALLY WITH FIBER COIN CREATOR
AVOID EDITING THIS MANUALLY
*/

var (
	// MainNetDistribution Laqpay mainnet coin distribution parameters
	MainNetDistribution = Distribution{
		MaxCoinSupply:        100000000,
		InitialUnlockedCount: 25,
		UnlockAddressRate:    5,
		UnlockTimeInterval:   31536000,
		Addresses: []string{
			"2EuyufQijRLMNZExtXwdvSzezRdgefh2NwD"
		},
	}

	// UserVerifyTxn transaction verification parameters for user-created transactions
	UserVerifyTxn = VerifyTxn{
		// BurnFactor can be overriden with `USER_BURN_FACTOR` env var
		BurnFactor: 10,
		// MaxTransactionSize can be overriden with `USER_MAX_TXN_SIZE` env var
		MaxTransactionSize: 32768, // in bytes
		// MaxDropletPrecision can be overriden with `USER_MAX_DECIMALS` env var
		MaxDropletPrecision: 3,
	}
)
