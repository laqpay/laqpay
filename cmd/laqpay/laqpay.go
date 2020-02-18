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
	"github.com/laqpay/laqpay/src/laqpay"
	"github.com/laqpay/laqpay/src/readable"
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
	GenesisSignatureStr = "db00864790c5ebfd7c702f8305ae4eff499282ae7177f2a06b5ee6bc65ad3ca8735f3db3636eed4e40fe9bc3ebace6850d08abe59e27c95e885d0a0c330c9c7200"
	// GenesisAddressStr genesis address string
	GenesisAddressStr = "7x7oZJhvR6n9QP7hsZCrtAVP7Apg42TGm4"
	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = "03de90898df039c28c984a29823537491e1aa0dd61f21ecd18d983c9f5f5244afa"
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""

	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64 = 1578231479
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
