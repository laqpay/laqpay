package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"../../src/cipher"
	"../../src/cipher/bip39"
	"../../src/coin"
)

const (
	trustedPeerPort = 20000
	daemonPort      = 20100
	rpcPort         = 20200
	guiPort         = 20300
)

type Config struct {
	Secret SecretConfig `json:"secret"`
	Public PublicConfig `json:"public"`
}

type SecretConfig struct {
	MasterSecKey     string `json:"masterPrivateKey"`
	GenesisSignature string `json:"genesisSignature"`
}

type PublicConfig struct {
	MasterPubKey string `json:"masterPublicKey"`

	GenesisBlock  GenesisBlockConfig  `json:"genesisBlock"`
	GenesisWallet GenesisWalletConfig `json:"genesisWallet"`

	CoinCode string `json:"coinCode"`

	Port    int `json:"port"`
	RPCPort int `json:"rpcPort"`
	GUIPort int `json:"guiPort"`

	TrustedPeers []string `json:"trustedPeers"`
}

type GenesisBlockConfig struct {
	Address    string `json:"address"`
	CoinVolume uint64 `json:"coins"`
	Timestamp  uint64 `json:"timestamp"`
	BodyHash   string `json:"bodyHash"`
	HeaderHash string `json:"headerHash"`
}

type GenesisWalletConfig struct {
	Seed            string `json:"seed"`
	Addresses       uint64 `json:"addresses"`
	CoinsPerAddress uint64 `json:"coinsPerAddress"`
}

func main() {
	var (
		file = flag.String("file", "laqpay-genesis.json", "file to save configuration of new coin")
		coin = flag.String("code", "LAQ", "code of new coin")
		addrCount = flag.Int("addr", 1, "number of distribution addresses")
		coinVol = flag.Int("vol", 80000000000000, "coin volume to send to each of disribution addresses")
		peerCount = flag.Int("peers", 0, "number of trusted peers running on localhost")
	)

	flag.Parse()

	cfg := createCoin(*coin, *addrCount, *coinVol, *peerCount)

	// Print config
	out, err := json.MarshalIndent(&cfg, "", "    ")
	if err != nil {
		log.Fatalf("failed to marshal JSON - %s", err)
	}

	fmt.Println(string(out))

	// Save config to disk if required
	if *file != "" {
		if err := ioutil.WriteFile(*file, out, os.ModePerm); err != nil {
			log.Fatalf("failed to save coin configuration to file - %s", err)
		}
	}
}

func createCoin(coinCode string, addrCount, coinVol, peerCount int) Config {
	rb := cipher.RandByte(32)
	sk, _ := cipher.NewSecKey(rb)
	pk, _ := cipher.PubKeyFromSecKey(sk)

	// Geneate genesis block
	var (
		gbAddr, _  = cipher.AddressFromSecKey(sk)
		gbCoins = uint64(addrCount * coinVol)
		gbTs    = uint64(time.Now().Unix())
	)
	gb, err := coin.NewGenesisBlock(gbAddr, gbCoins, gbTs)
	if err != nil {
		log.Fatalf("failed to create genesis block - %s", err)
	}

	// Genesis block wallet
	gwSeed, err := bip39.NewDefaultMnemonic()
	if err != nil {
		log.Fatalf("failed to generate genesis wallet seed")
	}

	// Trusted peers of coin networks (default connections)
	peers := make([]string, peerCount)
	for i := 0; i < peerCount; i++ {
		peers[i] = fmt.Sprintf("127.0.0.1:%d", trustedPeerPort+i)
	}

	signhash := cipher.MustSignHash(gb.HashHeader(), sk).Hex()
	// Coin configuration
	cfg := Config{
		Secret: SecretConfig{
			MasterSecKey:     sk.Hex(),
			GenesisSignature: signhash,
		},

		Public: PublicConfig{
			MasterPubKey: pk.Hex(),

			GenesisBlock: GenesisBlockConfig{
				Address:    gbAddr.String(),
				CoinVolume: gbCoins,
				Timestamp:  gbTs,
				BodyHash:   gb.Body.Hash().Hex(),
				HeaderHash: gb.HashHeader().Hex(),
			},

			GenesisWallet: GenesisWalletConfig{
				Seed:            gwSeed,
				CoinsPerAddress: uint64(coinVol),
				Addresses:       uint64(addrCount),
			},

			CoinCode: coinCode,

			Port:    daemonPort,
			RPCPort: rpcPort,
			GUIPort: guiPort,

			TrustedPeers: peers,
		},
	}

	return cfg
}
