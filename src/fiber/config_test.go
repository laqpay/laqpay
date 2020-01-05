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
			GenesisSignatureStr: "571d33433aeb76327f70abc32ab4ef132e01fd1f86f8ce3946765a0c2acc82a1171440c82d42a9b48dec3348776f9b3ec9e497e4b4afe81d3bc04cee6c357cad01",
			GenesisAddressStr:   "zQ5Y7eZ5CJj749Ltj9MQsHyUx4NQtMWfYh",
			BlockchainPubkeyStr: "033f59f8cc6cec5d30613d4b7d2ef28e478dd74b6879661447f1cdb8649749f8c0",
			BlockchainSeckeyStr: "",
			GenesisTimestamp:    1578207105,
			GenesisCoinVolume:   100e12,
			DefaultConnections: []string{
				"138.201.196.174:6000",
			},
			Port:                           6000,
			PeerListURL:                    "https://downloads.laqpay.net/blockchain/peers.txt",
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
