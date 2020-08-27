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
	Version = "0.1.2"
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
	GenesisSignatureStr = "28d0e64fc177e589e1cc39a9a10ebe19a8a3f18bb3f9daf14d443c12d131346278935cbf9ae7097b4b4baa5d5259ea53fca76ac081f2f8f5eb9c0f9a7520429e00"
	// GenesisAddressStr genesis address string
	GenesisAddressStr = "2UtyPbZ6xyBMDEncV5u1ZZDZF6E9edwciZf"
	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = "028605c8ea5f05b238d590829f4597ed1e52c621a556d3e7b8ade0f1410742f632"
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""

	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64 = 1583070642
	// GenesisCoinVolume represents the coin capacity
	GenesisCoinVolume uint64 = 80000000000000

	// DefaultConnections the default trust node addresses
	DefaultConnections = []string{
		"138.201.196.174:6000",
		"193.47.33.235:6000",
		"91.188.222.22:6000",
		"91.188.222.23:6000",
		"91.188.222.24:6000",
		"91.188.222.33:6000",
		"193.47.33.204:6000",
		"193.47.33.206:6000",
		"193.47.33.242:6000",
		"193.47.33.247:6000",
		"193.47.33.249:6000",
		"193.47.33.250:6000",
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
