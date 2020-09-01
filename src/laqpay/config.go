/*
Package laqpay implements the main daemon cmd's configuration and setup
*/
package laqpay

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/blang/semver"
	"github.com/toqueteos/webbrowser"

	"../../src/api"
	"../../src/cipher"
	"../../src/coin"
	"../../src/daemon"
	"../../src/fiber"
	"../../src/kvstorage"
	"../../src/params"
	"../../src/readable"
	"../../src/util/apputil"
	"../../src/util/certutil"
	"../../src/util/droplet"
	"../../src/util/file"
	"../../src/util/logging"
	"../../src/util/useragent"
	"../../src/visor"
	"../../src/visor/dbutil"
	"../../src/wallet"
)

var (
	help = false
)

// Config records laqpay node and build config
type Config struct {
	Node  NodeConfig
	Build readable.BuildInfo
}

// NodeConfig records the node's configuration
type NodeConfig struct {
	// Name of the coin
	CoinName string

	// Disable peer exchange
	DisablePEX bool
	// Download peer list
	DownloadPeerList bool
	// Download the peers list from this URL
	PeerListURL string
	// Don't make any outgoing connections
	DisableOutgoingConnections bool
	// Don't allowing incoming connections
	DisableIncomingConnections bool
	// Disables networking altogether
	DisableNetworking bool
	// Enable GUI
	EnableGUI bool
	// Disable CSRF check in the wallet API
	DisableCSRF bool
	// Disable Host, Origin and Referer header check in the wallet API
	DisableHeaderCheck bool
	// Disable CSP disable content-security-policy in http response
	DisableCSP bool
	// Comma separated list of API sets enabled on the remote web interface
	EnabledAPISets string
	// Comma separated list of API sets disabled on the remote web interface
	DisabledAPISets string
	// Enable all of API sets. Applies before disabling individual sets
	EnableAllAPISets bool

	enabledAPISets map[string]struct{}
	// Comma separate list of hostnames to accept in the Host header, used to bypass the Host header check which only applies to localhost addresses
	HostWhitelist string
	hostWhitelist []string

	// Only run on localhost and only connect to others on localhost
	LocalhostOnly bool
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address string
	// gnet uses this for TCP incoming and outgoing
	Port int
	// MaxConnections is the maximum number of total connections allowed
	MaxConnections int
	// Maximum outgoing connections to maintain
	MaxOutgoingConnections int
	// Maximum default outgoing connections
	MaxDefaultPeerOutgoingConnections int
	// How often to make outgoing connections
	OutgoingConnectionsRate time.Duration
	// MaxOutgoingMessageLength maximum size of outgoing messages
	MaxOutgoingMessageLength int
	// MaxIncomingMessageLength maximum size of incoming messages
	MaxIncomingMessageLength int
	// PeerlistSize represents the maximum number of peers that the pex would maintain
	PeerlistSize int
	// Wallet Address Version
	// AddressVersion string
	// Remote web interface
	WebInterface bool
	// Remote web interface port
	WebInterfacePort int
	// Remote web interface address
	WebInterfaceAddr string
	// Remote web interface certificate
	WebInterfaceCert string
	// Remote web interface key
	WebInterfaceKey string
	// Remote web interface HTTPS support
	WebInterfaceHTTPS bool
	// Remote web interface username and password
	WebInterfaceUsername string
	WebInterfacePassword string
	// Allow web interface auth without HTTPS
	WebInterfacePlaintextAuth bool

	// Launch System Default Browser after client startup
	LaunchBrowser bool

	// Data directory holds app data -- defaults to ~/.laqpay
	DataDirectory string
	// GUI directory contains assets for the HTML interface
	GUIDirectory string

	// Timeouts for the HTTP listener
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration

	// Remark to include in user agent sent in the wire protocol introduction
	UserAgentRemark string
	userAgent       useragent.Data

	// Logging
	ColorLog bool
	// This is the value registered with flag, it is converted to LogLevel after parsing
	LogLevel string
	// Disable "Reply to ping", "Received pong" log messages
	DisablePingPong bool

	// Verify the database integrity after loading
	VerifyDB bool
	// Reset the database if integrity checks fail, and continue running
	ResetCorruptDB bool

	// Transaction verification parameters for unconfirmed transactions
	UnconfirmedVerifyTxn params.VerifyTxn
	// Transaction verification parameters for transactions when creating blocks
	CreateBlockVerifyTxn params.VerifyTxn
	// Maximum total size of transactions in a block
	MaxBlockTransactionsSize uint32

	unconfirmedBurnFactor          uint64
	maxUnconfirmedTransactionSize  uint64
	unconfirmedMaxDropletPrecision uint64
	createBlockBurnFactor          uint64
	createBlockMaxTransactionSize  uint64
	createBlockMaxDropletPrecision uint64
	maxBlockSize                   uint64

	// Wallets
	// Defaults to ${DataDirectory}/wallets/
	WalletDirectory string
	// Wallet crypto type
	WalletCryptoType string

	// Key-value storage
	// Default to ${DataDirectory}/data
	KVStorageDirectory  string
	EnabledStorageTypes []kvstorage.Type

	// Disable the hardcoded default peers
	DisableDefaultPeers bool
	// Load custom peers from disk
	CustomPeersFile string

	RunBlockPublisher bool

	/* Developer options */

	// Enable cpu profiling
	ProfileCPU bool
	// Where the file is written to
	ProfileCPUFile string
	// Enable HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf bool
	// Expose HTTP profiling on this interface
	HTTPProfHost string

	DBPath     string
	DBReadOnly bool
	LogToFile  bool
	Version    bool // show node version

	GenesisSignatureStr string
	GenesisAddressStr   string
	BlockchainPubkeyStr string
	BlockchainSeckeyStr string
	GenesisTimestamp    uint64
	GenesisCoinVolume   uint64
	DefaultConnections  []string

	genesisSignature cipher.Sig
	genesisAddress   cipher.Address
	genesisHash      cipher.SHA256

	blockchainPubkey cipher.PubKey
	blockchainSeckey cipher.SecKey

	Fiber readable.FiberConfig
}

// NewNodeConfig returns a new node config instance
func NewNodeConfig(mode string, node fiber.NodeConfig) NodeConfig {
	nodeConfig := NodeConfig{
		CoinName:            node.CoinName,
		GenesisSignatureStr: node.GenesisSignatureStr,
		GenesisAddressStr:   node.GenesisAddressStr,
		GenesisCoinVolume:   node.GenesisCoinVolume,
		GenesisTimestamp:    node.GenesisTimestamp,
		BlockchainPubkeyStr: node.BlockchainPubkeyStr,
		BlockchainSeckeyStr: node.BlockchainSeckeyStr,
		DefaultConnections:  node.DefaultConnections,
		// Disable peer exchange
		DisablePEX: false,
		// Don't make any outgoing connections
		DisableOutgoingConnections: false,
		// Don't allowing incoming connections
		DisableIncomingConnections: false,
		// Disables networking altogether
		DisableNetworking: false,
		// Enable GUI
		EnableGUI: false,
		// Disable CSRF check in the wallet API
		DisableCSRF: false,
		// Disable Host, Origin and Referer header check in the wallet API
		DisableHeaderCheck: false,
		// DisableCSP disable content-security-policy in http response
		DisableCSP: false,
		// Only run on localhost and only connect to others on localhost
		LocalhostOnly: false,
		// Which address to serve on. Leave blank to automatically assign to a
		// public interface
		Address: "",
		// gnet uses this for TCP incoming and outgoing
		Port: node.Port,
		// MaxConnections is the maximum number of total connections allowed
		MaxConnections: 128,
		// MaxOutgoingConnections is the maximum outgoing connections allowed
		MaxOutgoingConnections: 16,
		// MaxDefaultOutgoingConnections is the maximum default outgoing connections allowed
		MaxDefaultPeerOutgoingConnections: 16,
		DownloadPeerList:                  true,
		PeerListURL:                       node.PeerListURL,
		// How often to make outgoing connections, in seconds
		OutgoingConnectionsRate:  time.Second * 5,
		MaxOutgoingMessageLength: 256 * 1024,
		MaxIncomingMessageLength: 1024 * 1024,
		PeerlistSize:             65535,
		// Wallet Address Version
		// AddressVersion: "test",
		// Remote web interface
		WebInterface:      true,
		WebInterfacePort:  node.WebInterfacePort,
		WebInterfaceAddr:  "127.0.0.1",
		WebInterfaceCert:  "",
		WebInterfaceKey:   "",
		WebInterfaceHTTPS: false,
		EnabledAPISets: strings.Join([]string{
			api.EndpointsRead,
			api.EndpointsTransaction,
		}, ","),
		DisabledAPISets: "",
		EnableAllAPISets: true,

		LaunchBrowser: false, // GUI
		// Data directory holds app data
		DataDirectory: node.DataDirectory,
		// Web GUI static resources
		GUIDirectory: "./src/gui/static/",
		// Logging
		ColorLog:        true,
		LogLevel:        "info",
		LogToFile:       false,
		DisablePingPong: false,

		VerifyDB:       true,
		ResetCorruptDB: true,

		// Blockchain/transaction validation
		UnconfirmedVerifyTxn: params.VerifyTxn{
			BurnFactor:          node.UnconfirmedBurnFactor,
			MaxTransactionSize:  node.UnconfirmedMaxTransactionSize,
			MaxDropletPrecision: node.UnconfirmedMaxDropletPrecision,
		},
		CreateBlockVerifyTxn: params.VerifyTxn{
			BurnFactor:          node.CreateBlockBurnFactor,
			MaxTransactionSize:  node.CreateBlockMaxTransactionSize,
			MaxDropletPrecision: node.CreateBlockMaxDropletPrecision,
		},
		MaxBlockTransactionsSize: node.MaxBlockTransactionsSize,

		// Wallets
		WalletDirectory:  "",
		WalletCryptoType: string(wallet.DefaultCryptoType),

		// Key-value storage
		KVStorageDirectory: "",
		EnabledStorageTypes: []kvstorage.Type{
			kvstorage.TypeTxIDNotes,
			kvstorage.TypeGeneral,
		},

		// Timeout settings for http.Server
		// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		HTTPReadTimeout:  time.Second * 10,
		HTTPWriteTimeout: time.Second * 60,
		HTTPIdleTimeout:  time.Second * 120,

		RunBlockPublisher: false,

		// Enable cpu profiling
		ProfileCPU: false,
		// Where the file is written to
		ProfileCPUFile: "cpu.prof",
		// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
		HTTPProf:     false,
		HTTPProfHost: "localhost:6060",

		Fiber: readable.FiberConfig{
			Name:                  node.CoinName,
			DisplayName:           node.DisplayName,
			Ticker:                node.Ticker,
			CoinHoursName:         node.CoinHoursName,
			CoinHoursNameSingular: node.CoinHoursNameSingular,
			CoinHoursTicker:       node.CoinHoursTicker,
			ExplorerURL:           node.ExplorerURL,
			Bip44Coin:             node.Bip44Coin,
		},
	}

	nodeConfig.applyConfigMode(mode)

	return nodeConfig
}

func (c *Config) postProcess() error {
	if help {
		flag.Usage()
		os.Exit(0)
	}

	var err error
	if c.Node.GenesisSignatureStr != "" {
		c.Node.genesisSignature, err = cipher.SigFromHex(c.Node.GenesisSignatureStr)
		panicIfError(err, "Invalid Signature")
	}

	if c.Node.GenesisAddressStr != "" {
		c.Node.genesisAddress, err = cipher.DecodeBase58Address(c.Node.GenesisAddressStr)
		panicIfError(err, "Invalid Address")
	}

	// Compute genesis block hash
	gb, err := coin.NewGenesisBlock(c.Node.genesisAddress, c.Node.GenesisCoinVolume, c.Node.GenesisTimestamp)
	if err != nil {
		panicIfError(err, "Create genesis hash failed")
	}
	c.Node.genesisHash = gb.HashHeader()

	if c.Node.BlockchainPubkeyStr != "" {
		c.Node.blockchainPubkey, err = cipher.PubKeyFromHex(c.Node.BlockchainPubkeyStr)
		panicIfError(err, "Invalid Pubkey")
	}
	if c.Node.BlockchainSeckeyStr != "" {
		c.Node.blockchainSeckey, err = cipher.SecKeyFromHex(c.Node.BlockchainSeckeyStr)
		panicIfError(err, "Invalid Seckey")
		c.Node.BlockchainSeckeyStr = ""
	}
	if c.Node.BlockchainSeckeyStr != "" {
		c.Node.blockchainSeckey = cipher.SecKey{}
	}

	home := file.UserHome()
	c.Node.DataDirectory, err = file.InitDataDir(replaceHome(c.Node.DataDirectory, home))
	panicIfError(err, "Invalid DataDirectory")

	if c.Node.WebInterfaceCert == "" {
		c.Node.WebInterfaceCert = filepath.Join(c.Node.DataDirectory, "laqpayd.cert")
	} else {
		c.Node.WebInterfaceCert = replaceHome(c.Node.WebInterfaceCert, home)
	}

	if c.Node.WebInterfaceKey == "" {
		c.Node.WebInterfaceKey = filepath.Join(c.Node.DataDirectory, "laqpayd.key")
	} else {
		c.Node.WebInterfaceKey = replaceHome(c.Node.WebInterfaceKey, home)
	}

	if c.Node.WalletDirectory == "" {
		c.Node.WalletDirectory = filepath.Join(c.Node.DataDirectory, "wallets")
	} else {
		c.Node.WalletDirectory = replaceHome(c.Node.WalletDirectory, home)
	}

	if c.Node.KVStorageDirectory == "" {
		c.Node.KVStorageDirectory = filepath.Join(c.Node.DataDirectory, "data")
	} else {
		c.Node.KVStorageDirectory = replaceHome(c.Node.KVStorageDirectory, home)
	}
	if len(c.Node.EnabledStorageTypes) == 0 {
		c.Node.EnabledStorageTypes = []kvstorage.Type{
			kvstorage.TypeGeneral,
			kvstorage.TypeTxIDNotes,
		}
	}

	if c.Node.DBPath == "" {
		c.Node.DBPath = filepath.Join(c.Node.DataDirectory, "data.db")
	} else {
		c.Node.DBPath = replaceHome(c.Node.DBPath, home)
	}

	userAgentData := useragent.Data{
		Coin:    c.Node.CoinName,
		Version: c.Build.Version,
		Remark:  c.Node.UserAgentRemark,
	}

	if _, err := userAgentData.Build(); err != nil {
		return err
	}

	c.Node.userAgent = userAgentData

	apiSets, err := buildAPISets(c.Node)
	if err != nil {
		return err
	}

	// Don't open browser to load wallets if wallet apis are disabled.
	c.Node.enabledAPISets = apiSets
	if _, ok := c.Node.enabledAPISets[api.EndpointsWallet]; !ok {
		c.Node.EnableGUI = false
		c.Node.LaunchBrowser = false
	}

	if c.Node.EnableGUI {
		c.Node.GUIDirectory = file.ResolveResourceDirectory(c.Node.GUIDirectory)
	}

	if c.Node.DisableDefaultPeers {
		c.Node.DefaultConnections = nil
	}

	if c.Node.HostWhitelist != "" {
		if c.Node.DisableHeaderCheck {
			return errors.New("host whitelist should be empty when header check is disabled")
		}
		c.Node.hostWhitelist = strings.Split(c.Node.HostWhitelist, ",")
	}

	httpAuthEnabled := c.Node.WebInterfaceUsername != "" || c.Node.WebInterfacePassword != ""
	if httpAuthEnabled && !c.Node.WebInterfaceHTTPS && !c.Node.WebInterfacePlaintextAuth {
		return errors.New("Web interface auth enabled but HTTPS is not enabled. Use -web-interface-plaintext-auth=true if this is desired")
	}

	if c.Node.MaxConnections < c.Node.MaxOutgoingConnections+c.Node.MaxDefaultPeerOutgoingConnections {
		return errors.New("-max-connections must be >= -max-outgoing-connections + -max-default-peer-outgoing-connections")
	}

	if c.Node.MaxOutgoingConnections > c.Node.MaxConnections {
		return errors.New("-max-outgoing-connections cannot be higher than -max-connections")
	}

	if c.Node.maxBlockSize > math.MaxUint32 {
		return errors.New("-max-block-size exceeds MaxUint32")
	}
	if c.Node.maxUnconfirmedTransactionSize > math.MaxUint32 {
		return errors.New("-max-txn-size-unconfirmed exceeds MaxUint32")
	}
	if c.Node.unconfirmedBurnFactor > math.MaxUint32 {
		return errors.New("-burn-factor-unconfirmed exceeds MaxUint32")
	}
	if c.Node.createBlockBurnFactor > math.MaxUint32 {
		return errors.New("-burn-factor-create-block exceeds MaxUint32")
	}

	if c.Node.unconfirmedMaxDropletPrecision > math.MaxUint8 {
		return errors.New("-max-decimals-unconfirmed exceeds MaxUint8")
	}
	if c.Node.createBlockMaxDropletPrecision > math.MaxUint8 {
		return errors.New("-max-decimals-create-block exceeds MaxUint8")
	}

	c.Node.UnconfirmedVerifyTxn.BurnFactor = uint32(c.Node.unconfirmedBurnFactor)
	c.Node.UnconfirmedVerifyTxn.MaxTransactionSize = uint32(c.Node.maxUnconfirmedTransactionSize)
	c.Node.UnconfirmedVerifyTxn.MaxDropletPrecision = uint8(c.Node.unconfirmedMaxDropletPrecision)
	c.Node.CreateBlockVerifyTxn.BurnFactor = uint32(c.Node.createBlockBurnFactor)
	c.Node.CreateBlockVerifyTxn.MaxTransactionSize = uint32(c.Node.createBlockMaxTransactionSize)
	c.Node.CreateBlockVerifyTxn.MaxDropletPrecision = uint8(c.Node.createBlockMaxDropletPrecision)
	c.Node.MaxBlockTransactionsSize = uint32(c.Node.maxBlockSize)

	if c.Node.UnconfirmedVerifyTxn.MaxTransactionSize < params.MinTransactionSize {
		return fmt.Errorf("-max-txn-size-unconfirmed must be >= params.MinTransactionSize (%d)", params.MinTransactionSize)
	}
	if c.Node.UnconfirmedVerifyTxn.MaxTransactionSize < params.UserVerifyTxn.MaxTransactionSize {
		return fmt.Errorf("-max-txn-size-unconfirmed must be >= params.UserVerifyTxn.MaxTransactionSize (%d)", params.UserVerifyTxn.MaxTransactionSize)
	}
	if c.Node.CreateBlockVerifyTxn.MaxTransactionSize < params.MinTransactionSize {
		return fmt.Errorf("-max-txn-size-create-block must be >= params.MinTransactionSize (%d)", params.MinTransactionSize)
	}
	if c.Node.CreateBlockVerifyTxn.MaxTransactionSize < params.UserVerifyTxn.MaxTransactionSize {
		return fmt.Errorf("-max-txn-size-create-block must be >= params.UserVerifyTxn.MaxTransactionSize (%d)", params.UserVerifyTxn.MaxTransactionSize)
	}

	if c.Node.MaxBlockTransactionsSize < params.MinTransactionSize {
		return fmt.Errorf("-max-block-size must be >= params.MinTransactionSize (%d)", params.MinTransactionSize)
	}
	if c.Node.MaxBlockTransactionsSize < params.UserVerifyTxn.MaxTransactionSize {
		return fmt.Errorf("-max-block-size must be >= params.UserVerifyTxn.MaxTransactionSize (%d)", params.UserVerifyTxn.MaxTransactionSize)
	}
	if c.Node.MaxBlockTransactionsSize < c.Node.UnconfirmedVerifyTxn.MaxTransactionSize {
		return errors.New("-max-block-size must be >= -max-txn-size-unconfirmed")
	}
	if c.Node.MaxBlockTransactionsSize < c.Node.CreateBlockVerifyTxn.MaxTransactionSize {
		return errors.New("-max-block-size must be >= -max-txn-size-create-block")
	}

	if c.Node.UnconfirmedVerifyTxn.BurnFactor < params.MinBurnFactor {
		return fmt.Errorf("-burn-factor-unconfirmed must be >= params.MinBurnFactor (%d)", params.MinBurnFactor)
	}
	if c.Node.UnconfirmedVerifyTxn.BurnFactor < params.UserVerifyTxn.BurnFactor {
		return fmt.Errorf("-burn-factor-unconfirmed must be >= params.UserVerifyTxn.BurnFactor (%d)", params.UserVerifyTxn.BurnFactor)
	}

	if c.Node.CreateBlockVerifyTxn.BurnFactor < params.MinBurnFactor {
		return fmt.Errorf("-burn-factor-create-block must be >= params.MinBurnFactor (%d)", params.MinBurnFactor)
	}
	if c.Node.CreateBlockVerifyTxn.BurnFactor < params.UserVerifyTxn.BurnFactor {
		return fmt.Errorf("-burn-factor-create-block must be >= params.UserVerifyTxn.BurnFactor (%d)", params.UserVerifyTxn.BurnFactor)
	}

	if c.Node.UnconfirmedVerifyTxn.MaxDropletPrecision > droplet.Exponent {
		return fmt.Errorf("-max-decimals-unconfirmed must be <= droplet.Exponent (%d)", droplet.Exponent)
	}
	if c.Node.UnconfirmedVerifyTxn.MaxDropletPrecision < params.UserVerifyTxn.MaxDropletPrecision {
		return fmt.Errorf("-max-decimals-unconfirmed must be >= params.UserVerifyTxn.MaxDropletPrecision (%d)", params.UserVerifyTxn.MaxDropletPrecision)
	}

	if c.Node.CreateBlockVerifyTxn.MaxDropletPrecision > droplet.Exponent {
		return fmt.Errorf("-max-decimals-create-block must be <= droplet.Exponent (%d)", droplet.Exponent)
	}
	if c.Node.CreateBlockVerifyTxn.MaxDropletPrecision < params.UserVerifyTxn.MaxDropletPrecision {
		return fmt.Errorf("-max-decimals-create-block must be >= params.UserVerifyTxn.MaxDropletPrecision (%d)", params.UserVerifyTxn.MaxDropletPrecision)
	}

	return nil
}

// buildAPISets builds the set of enable APIs by the following rules:
// * If EnableAll, all API sets are added
// * For each api set in EnabledAPISets, add
// * For each api set in DisabledAPISets, remove
func buildAPISets(c NodeConfig) (map[string]struct{}, error) {
	enabledAPISets := strings.Split(c.EnabledAPISets, ",")
	if err := validateAPISets("-enable-api-sets", enabledAPISets); err != nil {
		return nil, err
	}

	disabledAPISets := strings.Split(c.DisabledAPISets, ",")
	if err := validateAPISets("-disable-api-sets", disabledAPISets); err != nil {
		return nil, err
	}

	apiSets := make(map[string]struct{})

	allAPISets := []string{
		api.EndpointsRead,
		api.EndpointsStatus,
		api.EndpointsWallet,
		api.EndpointsTransaction,
		api.EndpointsPrometheus,
		api.EndpointsNetCtrl,
		api.EndpointsStorage,
		// Do not include insecure or deprecated API sets, they must always
		// be explicitly enabled through -enable-api-sets
	}

	if c.EnableAllAPISets {
		for _, s := range allAPISets {
			apiSets[s] = struct{}{}
		}
	}

	// Add the enabled API sets
	for _, k := range enabledAPISets {
		apiSets[k] = struct{}{}
	}

	// Remove the disabled API sets
	for _, k := range disabledAPISets {
		delete(apiSets, k)
	}

	return apiSets, nil
}

func validateAPISets(opt string, apiSets []string) error {
	for _, k := range apiSets {
		k = strings.ToUpper(strings.TrimSpace(k))
		switch k {
		case api.EndpointsRead,
			api.EndpointsStatus,
			api.EndpointsTransaction,
			api.EndpointsWallet,
			api.EndpointsInsecureWalletSeed,
			api.EndpointsPrometheus,
			api.EndpointsNetCtrl,
			api.EndpointsStorage:
		case "":
			continue
		default:
			return fmt.Errorf("Invalid value in %s: %q", opt, k)
		}
	}
	return nil
}

// RegisterFlags binds CLI flags to config values
func (c *NodeConfig) RegisterFlags() {
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&c.DisablePEX, "disable-pex", c.DisablePEX, "disable PEX peer discovery")
	flag.BoolVar(&c.DownloadPeerList, "download-peerlist", c.DownloadPeerList, "download a peers.txt from -peerlist-url")
	flag.StringVar(&c.PeerListURL, "peerlist-url", c.PeerListURL, "with -download-peerlist=true, download a peers.txt file from this url")
	flag.BoolVar(&c.DisableOutgoingConnections, "disable-outgoing", c.DisableOutgoingConnections, "Don't make outgoing connections")
	flag.BoolVar(&c.DisableIncomingConnections, "disable-incoming", c.DisableIncomingConnections, "Don't allow incoming connections")
	flag.BoolVar(&c.DisableNetworking, "disable-networking", c.DisableNetworking, "Disable all network activity")
	flag.BoolVar(&c.EnableGUI, "enable-gui", c.EnableGUI, "Enable GUI")
	flag.BoolVar(&c.DisableCSRF, "disable-csrf", c.DisableCSRF, "disable CSRF check")
	flag.BoolVar(&c.DisableHeaderCheck, "disable-header-check", c.DisableHeaderCheck, "disables the host, origin and referer header checks.")
	flag.BoolVar(&c.DisableCSP, "disable-csp", c.DisableCSP, "disable content-security-policy in http response")
	flag.StringVar(&c.Address, "address", c.Address, "IP Address to run application on. Leave empty to default to a public interface")
	flag.IntVar(&c.Port, "port", c.Port, "Port to run application on")

	flag.BoolVar(&c.WebInterface, "web-interface", c.WebInterface, "enable the web interface")
	flag.IntVar(&c.WebInterfacePort, "web-interface-port", c.WebInterfacePort, "port to serve web interface on")
	flag.StringVar(&c.WebInterfaceAddr, "web-interface-addr", c.WebInterfaceAddr, "addr to serve web interface on")
	flag.StringVar(&c.WebInterfaceCert, "web-interface-cert", c.WebInterfaceCert, "laqpayd.cert file for web interface HTTPS. If not provided, will autogenerate or use laqpayd.cert in --data-dir")
	flag.StringVar(&c.WebInterfaceKey, "web-interface-key", c.WebInterfaceKey, "laqpayd.key file for web interface HTTPS. If not provided, will autogenerate or use laqpayd.key in --data-dir")
	flag.BoolVar(&c.WebInterfaceHTTPS, "web-interface-https", c.WebInterfaceHTTPS, "enable HTTPS for web interface")
	flag.StringVar(&c.HostWhitelist, "host-whitelist", c.HostWhitelist, "Hostnames to whitelist in the Host header check. Only applies when the web interface is bound to localhost.")

	allAPISets := []string{
		api.EndpointsRead,
		api.EndpointsStatus,
		api.EndpointsWallet,
		api.EndpointsTransaction,
		api.EndpointsPrometheus,
		api.EndpointsNetCtrl,
		api.EndpointsInsecureWalletSeed,
		api.EndpointsStorage,
	}
	flag.StringVar(&c.EnabledAPISets, "enable-api-sets", c.EnabledAPISets, fmt.Sprintf("enable API set. Options are %s. Multiple values should be separated by comma", strings.Join(allAPISets, ", ")))
	flag.StringVar(&c.DisabledAPISets, "disable-api-sets", c.DisabledAPISets, fmt.Sprintf("disable API set. Options are %s. Multiple values should be separated by comma", strings.Join(allAPISets, ", ")))
	flag.BoolVar(&c.EnableAllAPISets, "enable-all-api-sets", c.EnableAllAPISets, "enable all API sets, except for deprecated or insecure sets. This option is applied before -disable-api-sets.")

	flag.StringVar(&c.WebInterfaceUsername, "web-interface-username", c.WebInterfaceUsername, "username for the web interface")
	flag.StringVar(&c.WebInterfacePassword, "web-interface-password", c.WebInterfacePassword, "password for the web interface")
	flag.BoolVar(&c.WebInterfacePlaintextAuth, "web-interface-plaintext-auth", c.WebInterfacePlaintextAuth, "allow web interface auth without https")

	flag.BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser, "launch system default webbrowser at client startup")
	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory, "directory to store app data (defaults to ~/.laqpay)")
	flag.StringVar(&c.DBPath, "db-path", c.DBPath, "path of database file (defaults to ~/.laqpay/data.db)")
	flag.BoolVar(&c.DBReadOnly, "db-read-only", c.DBReadOnly, "open bolt db read-only")
	flag.BoolVar(&c.ProfileCPU, "profile-cpu", c.ProfileCPU, "enable cpu profiling")
	flag.StringVar(&c.ProfileCPUFile, "profile-cpu-file", c.ProfileCPUFile, "where to write the cpu profile file")
	flag.BoolVar(&c.HTTPProf, "http-prof", c.HTTPProf, "run the HTTP profiling interface")
	flag.StringVar(&c.HTTPProfHost, "http-prof-host", c.HTTPProfHost, "hostname to bind the HTTP profiling interface to")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Choices are: debug, info, warn, error, fatal, panic")
	flag.BoolVar(&c.ColorLog, "color-log", c.ColorLog, "Add terminal colors to log output")
	flag.BoolVar(&c.DisablePingPong, "no-ping-log", c.DisablePingPong, `disable "reply to ping" and "received pong" debug log messages`)
	flag.BoolVar(&c.LogToFile, "logtofile", c.LogToFile, "log to file")
	flag.StringVar(&c.GUIDirectory, "gui-dir", c.GUIDirectory, "static content directory for the HTML interface")

	flag.BoolVar(&c.VerifyDB, "verify-db", c.VerifyDB, "check the database for corruption")
	flag.BoolVar(&c.ResetCorruptDB, "reset-corrupt-db", c.ResetCorruptDB, "reset the database if corrupted, and continue running instead of exiting")

	flag.BoolVar(&c.DisableDefaultPeers, "disable-default-peers", c.DisableDefaultPeers, "disable the hardcoded default peers")
	flag.StringVar(&c.CustomPeersFile, "custom-peers-file", c.CustomPeersFile, "load custom peers from a newline separate list of ip:port in a file. Note that this is different from the peers.json file in the data directory")

	flag.StringVar(&c.UserAgentRemark, "user-agent-remark", c.UserAgentRemark, "additional remark to include in the user agent sent over the wire protocol")

	flag.Uint64Var(&c.maxUnconfirmedTransactionSize, "max-txn-size-unconfirmed", uint64(c.UnconfirmedVerifyTxn.MaxTransactionSize), "maximum size of an unconfirmed transaction")
	flag.Uint64Var(&c.unconfirmedBurnFactor, "burn-factor-unconfirmed", uint64(c.UnconfirmedVerifyTxn.BurnFactor), "coinhour burn factor applied to unconfirmed transactions")
	flag.Uint64Var(&c.unconfirmedMaxDropletPrecision, "max-decimals-unconfirmed", uint64(c.UnconfirmedVerifyTxn.MaxDropletPrecision), "max number of decimal places applied to unconfirmed transactions")
	flag.Uint64Var(&c.createBlockBurnFactor, "burn-factor-create-block", uint64(c.CreateBlockVerifyTxn.BurnFactor), "coinhour burn factor applied when creating blocks")
	flag.Uint64Var(&c.createBlockMaxTransactionSize, "max-txn-size-create-block", uint64(c.CreateBlockVerifyTxn.MaxTransactionSize), "maximum size of a transaction applied when creating blocks")
	flag.Uint64Var(&c.createBlockMaxDropletPrecision, "max-decimals-create-block", uint64(c.CreateBlockVerifyTxn.MaxDropletPrecision), "max number of decimal places applied when creating blocks")
	flag.Uint64Var(&c.maxBlockSize, "max-block-size", uint64(c.MaxBlockTransactionsSize), "maximum total size of transactions in a block")

	flag.BoolVar(&c.RunBlockPublisher, "block-publisher", c.RunBlockPublisher, "run the daemon as a block publisher")
	flag.StringVar(&c.BlockchainPubkeyStr, "blockchain-public-key", c.BlockchainPubkeyStr, "public key of the blockchain")
	flag.StringVar(&c.BlockchainSeckeyStr, "blockchain-secret-key", c.BlockchainSeckeyStr, "secret key of the blockchain")

	flag.StringVar(&c.GenesisAddressStr, "genesis-address", c.GenesisAddressStr, "genesis address")
	flag.StringVar(&c.GenesisSignatureStr, "genesis-signature", c.GenesisSignatureStr, "genesis block signature")
	flag.Uint64Var(&c.GenesisTimestamp, "genesis-timestamp", c.GenesisTimestamp, "genesis block timestamp")

	flag.StringVar(&c.WalletDirectory, "wallet-dir", c.WalletDirectory, "location of the wallet files. Defaults to ~/.laqpay/wallet/")
	flag.StringVar(&c.KVStorageDirectory, "storage-dir", c.KVStorageDirectory, "location of the storage data files. Defaults to ~/.laqpay/data/")
	flag.IntVar(&c.MaxConnections, "max-connections", c.MaxConnections, "Maximum number of total connections allowed")
	flag.IntVar(&c.MaxOutgoingConnections, "max-outgoing-connections", c.MaxOutgoingConnections, "Maximum number of outgoing connections allowed")
	flag.IntVar(&c.MaxDefaultPeerOutgoingConnections, "max-default-peer-outgoing-connections", c.MaxDefaultPeerOutgoingConnections, "The maximum default peer outgoing connections allowed")
	flag.IntVar(&c.PeerlistSize, "peerlist-size", c.PeerlistSize, "Max number of peers to track in peerlist")
	flag.DurationVar(&c.OutgoingConnectionsRate, "connection-rate", c.OutgoingConnectionsRate, "How often to make an outgoing connection")
	flag.IntVar(&c.MaxOutgoingMessageLength, "max-out-msg-len", c.MaxOutgoingMessageLength, "Maximum length of outgoing wire messages")
	flag.IntVar(&c.MaxIncomingMessageLength, "max-in-msg-len", c.MaxIncomingMessageLength, "Maximum length of incoming wire messages")
	flag.BoolVar(&c.LocalhostOnly, "localhost-only", c.LocalhostOnly, "Run on localhost and only connect to localhost peers")
	flag.StringVar(&c.WalletCryptoType, "wallet-crypto-type", c.WalletCryptoType, "wallet crypto type. Can be sha256-xor or scrypt-chacha20poly1305")
	flag.BoolVar(&c.Version, "version", false, "show node version")
}

func (c *NodeConfig) applyConfigMode(configMode string) {
	if runtime.GOOS == "windows" {
		c.ColorLog = false
	}
	switch configMode {
	case "":
	case "STANDALONE_CLIENT":
		c.EnableAllAPISets = true
		c.EnabledAPISets = api.EndpointsInsecureWalletSeed
		c.EnableGUI = false
		c.LaunchBrowser = false
		c.DisableCSRF = false
		c.DisableHeaderCheck = false
		c.DisableCSP = false
		c.DownloadPeerList = true
		c.WebInterface = true
		c.LogToFile = false
		c.ResetCorruptDB = false
		c.WebInterfacePort = 6420 // randomize web interface port
	default:
		panic("Invalid ConfigMode")
	}
}

func panicIfError(err error, msg string, args ...interface{}) { //nolint:unparam
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func replaceHome(path, home string) string {
	return strings.Replace(path, "$HOME", home, 1)
}

var (
	// DBVerifyCheckpointVersion is a checkpoint for determining if DB verification should be run.
	// Any DB upgrading from less than this version to equal or higher than this version will be forced to verify.
	// Update this version checkpoint if a newer version requires a new verification run.
	DBVerifyCheckpointVersion       = "0.1.4"
	dbVerifyCheckpointVersionParsed semver.Version
)

// Coin represents a fiber coin instance
type Coin struct {
	config Config
	logger *logging.Logger
}

// Run starts the node
func (c *Coin) Run() error {
	var db *dbutil.DB
	var w *wallet.Service
	var v *visor.Visor
	var d *daemon.Daemon
	var s *kvstorage.Manager
	var gw *api.Gateway
	var webInterface *api.Server
	var retErr error
	errC := make(chan error, 10)

	if c.config.Node.Version {
		fmt.Println(c.config.Build.Version)
		return nil
	}

	logLevel, err := logging.LevelFromString(c.config.Node.LogLevel)
	if err != nil {
		err = fmt.Errorf("Invalid -log-level: %v", err)
		c.logger.Error(err)
		return err
	}

	logging.SetLevel(logLevel)

	if c.config.Node.ColorLog {
		logging.EnableColors()
	} else {
		logging.DisableColors()
	}

	var logFile *os.File
	if c.config.Node.LogToFile {
		var err error
		logFile, err = c.initLogFile()
		if err != nil {
			c.logger.Error(err)
			return err
		}
	}

	var fullAddress string
	scheme := "http"
	if c.config.Node.WebInterfaceHTTPS {
		scheme = "https"
	}
	host := fmt.Sprintf("%s:%d", c.config.Node.WebInterfaceAddr, c.config.Node.WebInterfacePort)

	if c.config.Node.ProfileCPU {
		f, err := os.Create(c.config.Node.ProfileCPUFile)
		if err != nil {
			c.logger.Error(err)
			return err
		}

		if err := pprof.StartCPUProfile(f); err != nil {
			c.logger.Error(err)
			return err
		}
		defer pprof.StopCPUProfile()
	}

	if c.config.Node.HTTPProf {
		go func() {
			if err := http.ListenAndServe(c.config.Node.HTTPProfHost, nil); err != nil {
				c.logger.WithError(err).Errorf("Listen on HTTP profiling interface %s failed", c.config.Node.HTTPProfHost)
			}
		}()
	}

	var wg sync.WaitGroup

	quit := make(chan struct{})

	// Catch SIGINT (CTRL-C) (closes the quit channel)
	go apputil.CatchInterrupt(quit)

	// Catch SIGUSR1 (prints runtime stack to stdout)
	go apputil.CatchDebug()

	// Parse the current app version
	appVersion, err := c.config.Build.Semver()
	if err != nil {
		c.logger.WithError(err).Errorf("Version %s is not a valid semver", c.config.Build.Version)
		return err
	}

	c.logger.Infof("App version: %s", appVersion)
	c.logger.Infof("OS: %s", runtime.GOOS)
	c.logger.Infof("Arch: %s", runtime.GOARCH)

	wconf := c.ConfigureWallet()
	dconf := c.ConfigureDaemon()
	vconf := c.ConfigureVisor()
	sconf := c.ConfigureStorage()

	// Open the database
	c.logger.Infof("Opening database %s", c.config.Node.DBPath)
	db, err = visor.OpenDB(c.config.Node.DBPath, c.config.Node.DBReadOnly)
	if err != nil {
		c.logger.Errorf("Database failed to open: %v. Is another laqpay instance running?", err)
		return err
	}

	// Look for saved app version
	dbVersion, err := visor.GetDBVersion(db)
	if err != nil {
		c.logger.WithError(err).Error("visor.GetDBVersion failed")
		retErr = err
		goto earlyShutdown
	}

	if dbVersion == nil {
		c.logger.Info("DB version not found in DB")
	} else {
		c.logger.Infof("DB version: %s", dbVersion)
	}

	c.logger.Infof("DB verify checkpoint version: %s", DBVerifyCheckpointVersion)

	// If the saved DB version is higher than the app version, abort.
	// Otherwise DB corruption could occur.
	if dbVersion != nil && dbVersion.GT(*appVersion) {
		err = fmt.Errorf("Cannot use newer DB version=%v with older software version=%v", dbVersion, appVersion)
		c.logger.WithError(err).Error()
		retErr = err
		goto earlyShutdown
	}

	// Verify the DB if the version detection says to, or if it was requested on the command line
	if shouldVerifyDB(appVersion, dbVersion) || c.config.Node.VerifyDB {
		if c.config.Node.ResetCorruptDB {
			// Check the database integrity and recreate it if necessary
			c.logger.Info("Checking database and resetting if corrupted")
			if newDB, err := visor.ResetCorruptDB(db, c.config.Node.blockchainPubkey, quit); err != nil {
				if err != visor.ErrVerifyStopped {
					c.logger.WithError(err).Error("visor.ResetCorruptDB failed")
					retErr = err
				}
				goto earlyShutdown
			} else {
				db = newDB
			}
		} else {
			c.logger.Info("Checking database")
			if err := visor.CheckDatabase(db, c.config.Node.blockchainPubkey, quit); err != nil {
				if err != visor.ErrVerifyStopped {
					c.logger.WithError(err).Error("visor.CheckDatabase failed")
					retErr = err
				}
				goto earlyShutdown
			}
		}
	}

	// Update the DB version
	if !db.IsReadOnly() {
		if err := visor.SetDBVersion(db, *appVersion); err != nil {
			c.logger.WithError(err).Error("visor.SetDBVersion failed")
			retErr = err
			goto earlyShutdown
		}
	}

	c.logger.Infof("Coinhour burn factor for user transactions is %d", params.UserVerifyTxn.BurnFactor)
	c.logger.Infof("Max transaction size for user transactions is %d", params.UserVerifyTxn.MaxTransactionSize)
	c.logger.Infof("Max decimals for user transactions is %d", params.UserVerifyTxn.MaxDropletPrecision)

	c.logger.Info("wallet.NewService")
	w, err = wallet.NewService(wconf)
	if err != nil {
		c.logger.WithError(err).Error("wallet.NewService failed")
		retErr = err
		goto earlyShutdown
	}

	c.logger.Info("visor.New")
	v, err = visor.New(vconf, db, w)
	if err != nil {
		c.logger.WithError(err).Error("visor.New failed")
		retErr = err
		goto earlyShutdown
	}

	c.logger.Info("daemon.New")
	d, err = daemon.New(dconf, v)
	if err != nil {
		c.logger.WithError(err).Error("daemon.New failed")
		retErr = err
		goto earlyShutdown
	}

	c.logger.Info("kvstorage.NewManager")
	s, err = kvstorage.NewManager(sconf)
	if err != nil {
		c.logger.WithError(err).Error("kvstorage.NewManager failed")
		retErr = err
		goto earlyShutdown
	}

	c.logger.Info("api.NewGateway")
	gw = api.NewGateway(d, v, w, s)

	if c.config.Node.WebInterface {
		webInterface, err = c.createGUI(gw, host)
		if err != nil {
			c.logger.WithError(err).Error("c.createGUI failed")
			retErr = err
			goto earlyShutdown
		}

		fullAddress = fmt.Sprintf("%s://%s", scheme, webInterface.Addr())
		c.logger.Critical().Infof("Full address: %s", fullAddress)
	}

	c.logger.Info("visor.Init")
	if err := v.Init(); err != nil {
		c.logger.WithError(err).Error("visor.Init failed")
		retErr = err
		goto earlyShutdown
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		c.logger.Info("daemon.Run")
		if err := d.Run(); err != nil {
			c.logger.WithError(err).Error("daemon.Run failed")
			errC <- err
		}
	}()

	if c.config.Node.WebInterface {
		cancelLaunchBrowser := make(chan struct{})

		wg.Add(1)
		go func() {
			defer wg.Done()

			c.logger.Info("webInterface.Serve")
			if err := webInterface.Serve(); err != nil {
				close(cancelLaunchBrowser)
				c.logger.WithError(err).Error("webInterface.Serve failed")
				errC <- err
			}
		}()

		if c.config.Node.LaunchBrowser {
			go func() {
				select {
				case <-cancelLaunchBrowser:
					c.logger.Warning("Browser launching canceled")

					// Wait a moment just to make sure the http interface is up
				case <-time.After(time.Millisecond * 100):
					c.logger.Infof("Launching System Browser with %s", fullAddress)
					if err := webbrowser.Open(fullAddress); err != nil {
						c.logger.WithError(err).Error("webbrowser.Open failed")
					}
				}
			}()
		}
	}

	select {
	case <-quit:
	case retErr = <-errC:
		c.logger.WithError(err).Error("Received error from errC (something prior has failed)")
	}

	c.logger.Info("Shutting down...")

	if webInterface != nil {
		c.logger.Info("Closing web interface")
		webInterface.Shutdown()
	}

	c.logger.Info("Closing daemon")
	d.Shutdown()

	c.logger.Info("Waiting for goroutines to finish")
	wg.Wait()

earlyShutdown:
	if db != nil {
		c.logger.Info("Closing database")
		if err := db.Close(); err != nil {
			c.logger.WithError(err).Error("Failed to close DB")
		}
	}

	c.logger.Info("Goodbye")

	if logFile != nil {
		if err := logFile.Close(); err != nil {
			fmt.Println("Failed to close log file")
		}
	}

	return retErr
}

// NewCoin returns a new fiber coin instance
func NewCoin(config Config, logger *logging.Logger) *Coin {
	return &Coin{
		config: config,
		logger: logger,
	}
}

func (c *Coin) initLogFile() (*os.File, error) {
	logDir := filepath.Join(c.config.Node.DataDirectory, "logs")
	if err := createDirIfNotExist(logDir); err != nil {
		c.logger.WithError(err).Errorf("createDirIfNotExist(%s) failed", logDir)
		return nil, fmt.Errorf("createDirIfNotExist(%s) failed: %v", logDir, err)
	}

	// open log file
	tf := "2020-03-01-030405"
	logfile := filepath.Join(logDir, fmt.Sprintf("%s-v%s.log", time.Now().Format(tf), c.config.Build.Version))

	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		c.logger.WithError(err).Errorf("os.OpenFile(%s) failed", logfile)
		return nil, err
	}

	hook := logging.NewWriteHook(f)
	logging.AddHook(hook)

	return f, nil
}

// ConfigureVisor sets the visor config values
func (c *Coin) ConfigureVisor() visor.Config {
	vc := visor.NewConfig()

	vc.Distribution = params.MainNetDistribution

	vc.IsBlockPublisher = c.config.Node.RunBlockPublisher
	vc.Arbitrating = c.config.Node.RunBlockPublisher

	vc.BlockchainPubkey = c.config.Node.blockchainPubkey
	vc.BlockchainSeckey = c.config.Node.blockchainSeckey

	vc.UnconfirmedVerifyTxn = c.config.Node.UnconfirmedVerifyTxn
	vc.CreateBlockVerifyTxn = c.config.Node.CreateBlockVerifyTxn
	vc.MaxBlockTransactionsSize = c.config.Node.MaxBlockTransactionsSize

	vc.GenesisAddress = c.config.Node.genesisAddress
	vc.GenesisSignature = c.config.Node.genesisSignature
	vc.GenesisTimestamp = c.config.Node.GenesisTimestamp
	vc.GenesisCoinVolume = c.config.Node.GenesisCoinVolume

	return vc
}

// ConfigureWallet sets the wallet config values
func (c *Coin) ConfigureWallet() wallet.Config {
	wc := wallet.NewConfig()

	wc.WalletDir = c.config.Node.WalletDirectory
	_, wc.EnableWalletAPI = c.config.Node.enabledAPISets[api.EndpointsWallet]
	_, wc.EnableSeedAPI = c.config.Node.enabledAPISets[api.EndpointsInsecureWalletSeed]

	// Initialize wallet default crypto type
	cryptoType, err := wallet.CryptoTypeFromString(c.config.Node.WalletCryptoType)
	if err != nil {
		log.Panic(err)
	}

	wc.CryptoType = cryptoType

	bc := c.config.Node.Fiber.Bip44Coin
	wc.Bip44Coin = &bc

	return wc
}

// ConfigureStorage sets the key-value storage config values
func (c *Coin) ConfigureStorage() kvstorage.Config {
	sc := kvstorage.NewConfig()

	sc.StorageDir = c.config.Node.KVStorageDirectory
	_, sc.EnableStorageAPI = c.config.Node.enabledAPISets[api.EndpointsStorage]
	sc.EnabledStorages = c.config.Node.EnabledStorageTypes

	return sc
}

// ConfigureDaemon sets the daemon config values
func (c *Coin) ConfigureDaemon() daemon.Config {
	dc := daemon.NewConfig()

	dc.Pool.DefaultConnections = c.config.Node.DefaultConnections
	dc.Pool.MaxDefaultPeerOutgoingConnections = c.config.Node.MaxDefaultPeerOutgoingConnections
	dc.Pool.MaxIncomingMessageLength = c.config.Node.MaxIncomingMessageLength
	dc.Pool.MaxOutgoingMessageLength = c.config.Node.MaxOutgoingMessageLength

	dc.Pex.DataDirectory = c.config.Node.DataDirectory
	dc.Pex.Disabled = c.config.Node.DisablePEX
	dc.Pex.NetworkDisabled = c.config.Node.DisableNetworking
	dc.Pex.Max = c.config.Node.PeerlistSize
	dc.Pex.DownloadPeerList = c.config.Node.DownloadPeerList
	dc.Pex.PeerListURL = c.config.Node.PeerListURL
	dc.Pex.DisableTrustedPeers = c.config.Node.DisableDefaultPeers
	dc.Pex.CustomPeersFile = c.config.Node.CustomPeersFile
	dc.Pex.DefaultConnections = c.config.Node.DefaultConnections

	dc.Daemon.MaxOutgoingMessageLength = uint64(c.config.Node.MaxOutgoingMessageLength)
	dc.Daemon.MaxIncomingMessageLength = uint64(c.config.Node.MaxIncomingMessageLength)
	dc.Daemon.MaxBlockTransactionsSize = c.config.Node.MaxBlockTransactionsSize
	dc.Daemon.DefaultConnections = c.config.Node.DefaultConnections
	dc.Daemon.DisableOutgoingConnections = c.config.Node.DisableOutgoingConnections
	dc.Daemon.DisableIncomingConnections = c.config.Node.DisableIncomingConnections
	dc.Daemon.DisableNetworking = c.config.Node.DisableNetworking
	dc.Daemon.Port = c.config.Node.Port
	dc.Daemon.Address = c.config.Node.Address
	dc.Daemon.LocalhostOnly = c.config.Node.LocalhostOnly
	dc.Daemon.MaxConnections = c.config.Node.MaxConnections
	dc.Daemon.MaxOutgoingConnections = c.config.Node.MaxOutgoingConnections
	dc.Daemon.DataDirectory = c.config.Node.DataDirectory
	dc.Daemon.LogPings = !c.config.Node.DisablePingPong
	dc.Daemon.BlockchainPubkey = c.config.Node.blockchainPubkey
	dc.Daemon.GenesisHash = c.config.Node.genesisHash
	dc.Daemon.UserAgent = c.config.Node.userAgent
	dc.Daemon.UnconfirmedVerifyTxn = c.config.Node.UnconfirmedVerifyTxn

	if c.config.Node.OutgoingConnectionsRate == 0 {
		c.config.Node.OutgoingConnectionsRate = time.Millisecond
	}
	dc.Daemon.OutgoingRate = c.config.Node.OutgoingConnectionsRate

	return dc
}

func (c *Coin) createGUI(gw *api.Gateway, host string) (*api.Server, error) {
	config := api.Config{
		StaticDir:          c.config.Node.GUIDirectory,
		DisableCSRF:        c.config.Node.DisableCSRF,
		DisableHeaderCheck: c.config.Node.DisableHeaderCheck,
		DisableCSP:         c.config.Node.DisableCSP,
		EnableGUI:          c.config.Node.EnableGUI,
		ReadTimeout:        c.config.Node.HTTPReadTimeout,
		WriteTimeout:       c.config.Node.HTTPWriteTimeout,
		IdleTimeout:        c.config.Node.HTTPIdleTimeout,
		EnabledAPISets:     c.config.Node.enabledAPISets,
		HostWhitelist:      c.config.Node.hostWhitelist,
		Health: api.HealthConfig{
			BuildInfo: readable.BuildInfo{
				Version: c.config.Build.Version,
				Commit:  c.config.Build.Commit,
				Branch:  c.config.Build.Branch,
			},
			Fiber:           c.config.Node.Fiber,
			DaemonUserAgent: c.config.Node.userAgent,
			BlockPublisher:  c.config.Node.RunBlockPublisher,
		},
		Username: c.config.Node.WebInterfaceUsername,
		Password: c.config.Node.WebInterfacePassword,
	}

	var s *api.Server
	if c.config.Node.WebInterfaceHTTPS {
		// Verify cert/key parameters, and if neither exist, create them
		exists, err := checkCertFiles(c.config.Node.WebInterfaceCert, c.config.Node.WebInterfaceKey)
		if err != nil {
			c.logger.WithError(err).Error("checkCertFiles failed")
			return nil, err
		}

		if !exists {
			c.logger.Infof("Autogenerating HTTP certificate and key files %s, %s", c.config.Node.WebInterfaceCert, c.config.Node.WebInterfaceKey)
			if err := createCertFiles(c.config.Node.WebInterfaceCert, c.config.Node.WebInterfaceKey); err != nil {
				c.logger.WithError(err).Error("createCertFiles failed")
				return nil, err
			}

			c.logger.Infof("Created cert file %s", c.config.Node.WebInterfaceCert)
			c.logger.Infof("Created key file %s", c.config.Node.WebInterfaceKey)
		}

		s, err = api.CreateHTTPS(host, config, gw, c.config.Node.WebInterfaceCert, c.config.Node.WebInterfaceKey)
		if err != nil {
			c.logger.WithError(err).Error("Failed to start web failed")
			return nil, err
		}
	} else {
		var err error
		s, err = api.Create(host, config, gw)
		if err != nil {
			c.logger.WithError(err).Error("Failed to start web failed")
			return nil, err
		}
	}

	return s, nil
}

// checkCertFiles returns true if both cert and key files exist, false if neither exist,
// or returns an error if only one does not exist
func checkCertFiles(cert, key string) (bool, error) {
	doesFileExist := func(f string) (bool, error) {
		if _, err := os.Stat(f); err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}

	certExists, err := doesFileExist(cert)
	if err != nil {
		return false, err
	}

	keyExists, err := doesFileExist(key)
	if err != nil {
		return false, err
	}

	switch {
	case certExists && keyExists:
		return true, nil
	case !certExists && !keyExists:
		return false, nil
	case certExists && !keyExists:
		return false, fmt.Errorf("certfile %s exists but keyfile %s does not", cert, key)
	case !certExists && keyExists:
		return false, fmt.Errorf("keyfile %s exists but certfile %s does not", key, cert)
	default:
		log.Panic("unreachable code")
		return false, errors.New("unreachable code")
	}
}

func createCertFiles(certFile, keyFile string) error {
	org := "laqpay daemon autogenerated cert"
	validUntil := time.Now().Add(10 * 365 * 24 * time.Hour)
	cert, key, err := certutil.NewTLSCertPair(org, validUntil, nil)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(certFile, cert, 0600); err != nil {
		return err
	}
	if err := ioutil.WriteFile(keyFile, key, 0600); err != nil {
		os.Remove(certFile)
		return err
	}

	return nil
}

// ParseConfig prepare the config
func (c *Coin) ParseConfig() error {
	return c.config.postProcess()
}

// InitTransaction creates the genesis transaction
func InitTransaction(uxID string, genesisSecKey cipher.SecKey, dist params.Distribution) coin.Transaction {
	dist.MustValidate()

	var txn coin.Transaction

	output := cipher.MustSHA256FromHex(uxID)
	if err := txn.PushInput(output); err != nil {
		log.Panic(err)
	}

	for _, addr := range dist.AddressesDecoded() {
		if err := txn.PushOutput(addr, dist.AddressInitialBalance()*droplet.Multiplier, 1); err != nil {
			log.Panic(err)
		}
	}

	seckeys := make([]cipher.SecKey, 1)
	seckey := genesisSecKey.Hex()
	seckeys[0] = cipher.MustSecKeyFromHex(seckey)
	txn.SignInputs(seckeys)

	if err := txn.UpdateHeader(); err != nil {
		log.Panic(err)
	}

	if err := txn.Verify(); err != nil {
		log.Panic(err)
	}

	log.Printf("signature= %s", txn.Sigs[0].Hex())
	return txn
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(dir, 0750)
}

func shouldVerifyDB(appVersion, dbVersion *semver.Version) bool {
	// If the dbVersion is not set, verify
	if dbVersion == nil {
		return true
	}

	// If the dbVersion is less than the verification checkpoint version
	// and the appVersion is greater than or equal to the checkpoint version,
	// verify
	if dbVersion.LT(dbVerifyCheckpointVersionParsed) && appVersion.GTE(dbVerifyCheckpointVersionParsed) {
		return true
	}

	return false
}

func init() {
	dbVerifyCheckpointVersionParsed = semver.MustParse(DBVerifyCheckpointVersion)
}
