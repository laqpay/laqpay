/*
laqpay daemon
*/
package main

/*
CODE GENERATED AUTOMATICALLY WITH FIBER COIN CREATOR
AVOID EDITING THIS MANUALLY
*/

import (
	"flag"
	_ "net/http/pprof"
	"os"

	"github.com/laqpay/laqpay/src/fiber"
	"github.com/laqpay/laqpay/src/readable"
	"github.com/laqpay/laqpay/src/laqpay"
	"github.com/laqpay/laqpay/src/util/logging"
)

var (
	// Version of the node. Can be set by -ldflags
	Version = "0.26.0"
	// Commit ID. Can be set by -ldflags
	Commit = ""
	// Branch name. Can be set by -ldflags
	Branch = ""
	// ConfigMode (possible values are "", "STANDALONE_CLIENT").
	// This is used to change the default configuration.
	// Can be set by -ldflags
	ConfigMode = ""

	logger = logging.MustGetLogger("main")

	// CoinName name of coin
	CoinName = "laqpay"

	// GenesisSignatureStr hex string of genesis signature
	GenesisSignatureStr = "571d33433aeb76327f70abc32ab4ef132e01fd1f86f8ce3946765a0c2acc82a1171440c82d42a9b48dec3348776f9b3ec9e497e4b4afe81d3bc04cee6c357cad01"
	// GenesisAddressStr genesis address string
	GenesisAddressStr = "zQ5Y7eZ5CJj749Ltj9MQsHyUx4NQtMWfYh"
	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = "033f59f8cc6cec5d30613d4b7d2ef28e478dd74b6879661447f1cdb8649749f8c0"
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""

	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64 = 1578207105
	// GenesisCoinVolume represents the coin capacity
	GenesisCoinVolume uint64 = 100000000000000

	// DefaultConnections the default trust node addresses
	DefaultConnections = []string{
		"138.201.196.174:6000",
	}

	nodeConfig = laqpay.NewNodeConfig(ConfigMode, fiber.NodeConfig{
		CoinName:            CoinName,
		GenesisSignatureStr: GenesisSignatureStr,
		GenesisAddressStr:   GenesisAddressStr,
		GenesisCoinVolume:   GenesisCoinVolume,
		GenesisTimestamp:    GenesisTimestamp,
		BlockchainPubkeyStr: BlockchainPubkeyStr,
		BlockchainSeckeyStr: BlockchainSeckeyStr,
		DefaultConnections:  DefaultConnections,
		PeerListURL:         "https://downloads.laqpay.com/blockchain/peers.txt",
		Port:                6000,
		WebInterfacePort:    6420,
		DataDirectory:       "$HOME/.laqpay",

		UnconfirmedBurnFactor:          10,
		UnconfirmedMaxTransactionSize:  32768,
		UnconfirmedMaxDropletPrecision: 3,
		CreateBlockBurnFactor:          10,
		CreateBlockMaxTransactionSize:  32768,
		CreateBlockMaxDropletPrecision: 3,
		MaxBlockTransactionsSize:       32768,

		DisplayName:           "Laqpay",
		Ticker:                "LAQ",
		CoinHoursName:         "Coin Hours",
		CoinHoursNameSingular: "Coin Hour",
		CoinHoursTicker:       "SCH",
		ExplorerURL:           "https://explorer.laqpay.com",
		Bip44Coin:             8000,
	})

	parseFlags = true
)

func init() {
	nodeConfig.RegisterFlags()
}

func main() {
	if parseFlags {
		flag.Parse()
	}

	// create a new fiber coin instance
	coin := laqpay.NewCoin(laqpay.Config{
		Node: nodeConfig,
		Build: readable.BuildInfo{
			Version: Version,
			Commit:  Commit,
			Branch:  Branch,
		},
	}, logger)

	// parse config values
	if err := coin.ParseConfig(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// run fiber coin node
	if err := coin.Run(); err != nil {
		os.Exit(1)
	}
}
