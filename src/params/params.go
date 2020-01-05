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
			"knuD1GNFxVq8h489g2EkixwnMMXdzJxRv2",
			"JTYGfZf9RAtLMX8fX9HQGvgdR43vmyte3e",
			"2DxDynFwsB4hLeENzDfT9Qusn6MR8ovVjJ9",
			"MrFhCebjazo3DYfwHxgzwdaiv4VWQ9fMfs",
			"4u68xDgjzLhFUEmHUwkViffDpYe4t5r3Nh",
			"QmUfKLGUEU8zDVdD1uznZmwwepakBcakzs",
			"2MKgrbu4XWiyAwUhSFybEM3dBr5RyqemLmU",
			"268AtZHvX2jTExbNbfZot2FrmAa4vKfgQzw",
			"SH6kbkb3wT7zoauwFCJfvppcr6B9h74vTz",
			"51EYGzABTuKWCebRBhQ8HKSVrnM3KZKmiV",
			"VjwnFi95XoEcdHGgCg3ocA4EZYxM6rZzs9",
			"2ZpdGTB1oefbVgUrfiBZoiZA5zjJ4ee6aR",
			"HpyjJmr3v7gsnd8oX5smaQSquhs8NeX7ga",
			"2KPU8jPr8fRRJpwhzTRpa8Ea3r7LQSgxuHm",
			"29C7LksK2twK3x7Wf2e5r8PgyrUoi7GSGzP",
			"2RSVhj6avAX7Cwo758n9d1q1vyefCpUZG9L",
			"jHDeaSumg8BY3LodRDctQua5A7pcmHqXAg",
			"2ZrwfqZSLu3witrjkJyfbf2SgCisdJw14Ap",
			"EDoj3jWPs7CXYM86CfNRcX5MwNr5RoUvNM",
			"QdJCYpaDATnP95uiGtv16tAofbB2pgMuRy",
			"29tpXB2czPxF251RzkgH7ZC3bhZH9DQct2A",
			"VN39PB65pciZBqSwmSVHAHkVf2LrkajbkG",
			"b34BYnh9TyicNZhzsyjvVwgD6bSeLfD8H3",
			"21R4G8QvyWr5MbpHgF7MCH9tmmAdomBB1YS",
			"28bwz5JNeqKKJfmqHDzQT2KPSJPwZZ3aTgB",
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
