package fiber

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/laqpay/laqpay/src/cipher/bip44"
)

// TODO(therealssj): write better tests
func TestNewConfig(t *testing.T) {
	coinConfig, err := NewConfig("test.fiber.toml", "./testdata")
	require.NoError(t, err)
	require.Equal(t, Config{
		Node: NodeConfig{
			GenesisSignatureStr: "db00864790c5ebfd7c702f8305ae4eff499282ae7177f2a06b5ee6bc65ad3ca8735f3db3636eed4e40fe9bc3ebace6850d08abe59e27c95e885d0a0c330c9c7200",
			GenesisAddressStr:   "7x7oZJhvR6n9QP7hsZCrtAVP7Apg42TGm4",
			BlockchainPubkeyStr: "03de90898df039c28c984a29823537491e1aa0dd61f21ecd18d983c9f5f5244afa",
			BlockchainSeckeyStr: "",
			GenesisTimestamp:    1578231479,
			GenesisCoinVolume:   100e12,
			DefaultConnections: []string{
				"138.201.196.174:6000",
			},
			Port:                           6000,
			PeerListURL:                    "https://downloads.laqpay.com/blockchain/peers.txt",
			WebInterfacePort:               6420,
			UnconfirmedBurnFactor:          10,
			UnconfirmedMaxTransactionSize:  777,
			UnconfirmedMaxDropletPrecision: 3,
			CreateBlockBurnFactor:          9,
			CreateBlockMaxTransactionSize:  1234,
			CreateBlockMaxDropletPrecision: 4,
			MaxBlockTransactionsSize:       1111,
			DisplayName:                    "Testcoin",
			Ticker:                         "TST",
			CoinHoursName:                  "Testcoin Hours",
			CoinHoursNameSingular:          "Testcoin Hour",
			CoinHoursTicker:                "TCH",
			ExplorerURL:                    "https://explorer.testcoin.com",
			Bip44Coin:                      bip44.CoinTypeLaqpay,
		},
		Params: ParamsConfig{
			MaxCoinSupply:           1e8,
			UnlockAddressRate:       5,
			InitialUnlockedCount:    25,
			UnlockTimeInterval:      60 * 60 * 24 * 365,
			UserBurnFactor:          3,
			UserMaxTransactionSize:  999,
			UserMaxDropletPrecision: 2,
		},
	}, coinConfig)
}
