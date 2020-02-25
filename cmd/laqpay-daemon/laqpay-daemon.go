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
	GenesisSignatureStr = "7fad4367ed3632c4f45008b483af06657fcfaafe53ec2e8f342af1d7122e2fb7726bfeae654d775b78ba263a7b3a0087cc6439f1f358af503c32c0d08c93597000"
	// GenesisAddressStr genesis address string
	GenesisAddressStr = "2P6fPf4vd3mAvQeyyg4X5jdc78JmCotT7Ae"
	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = "0341ba8589b51981f0c9c16b51709def29979b5a81d4d6fb59c8f17c5303655fb4"
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""

	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64 = 1582289770
	// GenesisCoinVolume represents the coin capacity
	GenesisCoinVolume uint64 = 80000000000000

	// DefaultConnections the default trust node addresses
	DefaultConnections = []string{}

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
