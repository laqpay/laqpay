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

	"../../src/fiber"
	"../../src/laqpay"
	"../../src/readable"
	"../../src/util/logging"
)

var (
	// Version of the node. Can be set by -ldflags
	Version = "0.1.1"
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
	GenesisSignatureStr = "a38878da1d8594929489a84ddd07201f9c7dc72d302dacbd4e892c2c1bbcfcc72d9bfd711a1a20084c7ecd2fc7bbf14f410767587f11592aa92d15dd644a9cea01"
	// GenesisAddressStr genesis address string
	GenesisAddressStr = "2E8KVvGvoMsC9Cohj7MBKgHLmcmdWBCdRkv"
	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = "02ca3b946a100b02ec3f9e92a5093c562e53a17fbad441c5dfc1be067fec987b8c"
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""

	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64 = 1583062083
	// GenesisCoinVolume represents the coin capacity
	GenesisCoinVolume uint64 = 80000000000000

	// DefaultConnections the default trust node addresses
	DefaultConnections = []string{
		"138.201.196.174:6000",
		"193.47.33.235:6000",
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
		PeerListURL:         "https://api.laqpay.com/network/peers",
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

		DisplayName:           "LAQ",
		Ticker:                "LAQ",
		CoinHoursName:         "LAQH",
		CoinHoursNameSingular: "LAQH",
		CoinHoursTicker:       "LAQH",
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
