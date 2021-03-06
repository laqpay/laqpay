# Laqpay Exchange Integration

A Laqpay node offers a REST API on port 6420 (when running from source; if you are using the releases downloaded from the website, the port is randomized)

A CLI tool is provided in `cmd/laqpay-wallet-cli/laqpay-wallet-cli.go`. This tool communicates over the REST API.

The API interfaces do not support authentication or encryption so they should only be used over localhost.

If your application is written in Go, you can use these client libraries to interface with the node:

* [Laqpay REST API Client Godoc](https://godoc.org/github.com/laqpay/laqpay/src/api#Client)
* [Laqpay CLI Godoc](https://godoc.org/github.com/laqpay/laqpay/src/cli)

*Note*: The CLI interface will be deprecated and replaced with a better one in the future.

The wallet APIs in the REST API operate on wallets loaded from and saved to `~/.laqpay/wallets`.
Use the CLI tool to perform seed generation and transaction signing outside of the Laqpay node.

The Laqpay node's wallet APIs can be enabled from the command line.
`-enable-all-api-sets` will enable all of the APIs which includes the wallet APIs,
or for more control it can specified in a list of API sets, e.g. `-enable-api-sets=READ,STATUS,WALLET`.
See the [REST API](src/api/README.md) for information on API sets.

For a node used to support another application,
it is recommended to use the REST API for blockchain queries and disable the wallet APIs,
and to use the CLI tool for wallet operations (seed and address generation, transaction signing).

<!-- MarkdownTOC autolink="true" bracket="round" levels="1,2,3,4,5,6" -->

- [Running the laqpay node](#running-the-laqpay-node)
- [API Documentation](#api-documentation)
	- [Wallet REST API](#wallet-rest-api)
	- [Laqpay command line interface](#laqpay-command-line-interface)
	- [Laqpay REST API Client Documentation](#laqpay-rest-api-client-documentation)
	- [Laqpay Go Library Documentation](#laqpay-go-library-documentation)
	- [liblaqpay Documentation](#liblaqpay-documentation)
- [Implementation guidelines](#implementation-guidelines)
	- [Scanning deposits](#scanning-deposits)
		- [Using the CLI](#using-the-cli)
		- [Using the REST API](#using-the-rest-api)
		- [Using laqpay as a library in a Go application](#using-laqpay-as-a-library-in-a-go-application)
	- [Sending coins](#sending-coins)
		- [General principles](#general-principles)
		- [Using the CLI](#using-the-cli-1)
		- [Using the REST API](#using-the-rest-api-1)
		- [Using laqpay as a library in a Go application](#using-laqpay-as-a-library-in-a-go-application-1)
		- [Coinhours](#coinhours)
			- [REST API](#rest-api)
			- [CLI](#cli)
	- [Verifying addresses](#verifying-addresses)
		- [Using the CLI](#using-the-cli-2)
		- [Using the REST API](#using-the-rest-api-2)
		- [Using laqpay as a library in a Go application](#using-laqpay-as-a-library-in-a-go-application-2)
		- [Using laqpay as a library in other applications](#using-laqpay-as-a-library-in-other-applications)
	- [Checking Laqpay node connections](#checking-laqpay-node-connections)
		- [Using the CLI](#using-the-cli-3)
		- [Using the REST API](#using-the-rest-api-3)
		- [Using laqpay as a library in a Go application](#using-laqpay-as-a-library-in-a-go-application-3)
	- [Checking Laqpay node status](#checking-laqpay-node-status)
		- [Using the CLI](#using-the-cli-4)
		- [Using the REST API](#using-the-rest-api-4)
		- [Using laqpay as a library in a Go application](#using-laqpay-as-a-library-in-a-go-application-4)
- [xpub wallets](#xpub-wallets)
	- [Create a bip44 HD wallet](#create-a-bip44-hd-wallet)
		- [Using the CLI](#using-the-cli-5)
		- [Using the REST API](#using-the-rest-api-5)
	- [Export an xpub key from a bip44 wallet](#export-an-xpub-key-from-a-bip44-wallet)
		- [Using the CLI](#using-the-cli-6)
		- [Using the REST API](#using-the-rest-api-6)
	- [Create an xpub wallet](#create-an-xpub-wallet)
		- [Using the CLI](#using-the-cli-7)
		- [Using the REST API](#using-the-rest-api-7)

<!-- /MarkdownTOC -->

## API Documentation

### Wallet REST API

[Wallet REST API](src/api/README.md).

### Laqpay command line interface

[CLI command API](cmd/cli/README.md).

### Laqpay REST API Client Documentation

[Laqpay REST API Client](https://godoc.org/github.com/laqpay/laqpay/src/api#Client)

### Laqpay Go Library Documentation

[Laqpay Godoc](https://godoc.org/github.com/laqpay/laqpay)

### liblaqpay Documentation

`liblaqpay` provides a C library for Laqpay's cryptographic operations.
This allows python, ruby, C#, Java, etc applications to perform operations
such as transaction signing and address verification, without using the Laqpay daemon API
or calling the CLI tool from the shell.

[liblaqpay documentation](https://github.com/laqpay/liblaqpay)

## Implementation guidelines

### Scanning deposits

There are multiple approaches to scanning for deposits, depending on your implementation.

One option is to watch for incoming blocks and check them for deposits made to a list of known deposit addresses.
Another option is to check the unspent outputs for a list of known deposit addresses.

#### Using the CLI

To scan the blockchain, use `laqpay-wallet-cli lastBlocks` or `laqpay-wallet-cli blocks`. These will return block data as JSON
and new unspent outputs sent to an address can be detected.

To check address outputs, use `laqpay-wallet-cli addressOutputs`. If you only want the balance, you can use `laqpay-wallet-cli addressBalance`.

#### Using the REST API

To scan the blockchain, call `GET /api/v1/last_blocks?num=` or `GET /api/v1/blocks?start=&end=`. There will return block data as JSON
and new unspent outputs sent to an address can be detected.

To check address outputs, call `GET /api/v1/outputs?addrs=`. If you only want the balance, you can call `GET /api/v1/balance?addrs=`.

* [`GET /api/v1/last_blocks` docs](src/api/README.md#get-last-n-blocks)
* [`GET /api/v1/blocks` docs](src/api/README.md#get-blocks-in-specific-range)
* [`GET /api/v1/outputs` docs](src/api/README.md#get-unspent-output-set-of-address-or-hash)
* [`GET /api/v1/balance` docs](src/api/README.md#get-balance-of-addresses)

#### Using laqpay as a library in a Go application

We recommend using the [Laqpay REST API Client](https://godoc.org/github.com/laqpay/laqpay/src/api#Client).

### Sending coins

#### General principles

After each spend, wait for the transaction to confirm before trying to spend again.

For higher throughput, combine multiple spends into one transaction.

Laqpay uses "coin hours" to ratelimit transactions.
The total number of coinhours in a transaction's outputs must be 50% or less than the number of coinhours in a transaction's inputs,
or else the transaction is invalid and will not be accepted. A transaction must have at least 1 input with at least 1 coin hour.
Sending too many transactions in quick succession will use up all available coinhours.
Coinhours are earned at a rate of 1 coinhour per coin per hour, calculated per second.
This means that 3600 coins will earn 1 coinhour per second.
However, coinhours are only updated when a new block is published to the blockchain.
New blocks are published every 10 seconds, but only if there are pending transactions in the network.

To avoid running out of coinhours in situations where the application may frequently send,
the sender should batch sends into a single transaction and send them on a
30 second to 1 minute interval.

There are other strategies to minimize the likelihood of running out of coinhours, such
as splitting up balances into many unspent outputs and having a large balance which generates
coinhours quickly.

#### Using the CLI

When sending coins from the CLI tool, a wallet file local to the caller is used.
The CLI tool allows you to specify the wallet file on disk to use for operations.

See [CLI command API](cmd/cli/README.md) for documentation of the CLI interface.

To perform a send, the preferred method follows these steps in a loop:

* `laqpay-wallet-cli createRawTransaction $WALLET_FILE -m '[{"addr:"$addr1,"coins:"$coins1"}, ...]` - `-m` flag is send-to-many
* `laqpay-wallet-cli broadcastTransaction` - returns `txid`
* `laqpay-wallet-cli transaction $txid` - repeat this command until `"status"` is `"confirmed"`

That is, create a raw transaction, broadcast it, and wait for it to confirm.

#### Using the REST API

The wallet APIs must be enabled with `-enable-api-sets=WALLET,READ`.

Create a transaction with [POST /wallet/transaction](https://github.com/laqpay/laqpay/blob/develop/src/api/README.md#create-transaction),
then inject it to the network with [POST /injectTransaction](https://github.com/laqpay/laqpay/blob/develop/src/api/README.md#inject-raw-transaction).

When using `POST /wallet/transaction`, a wallet file local to the laqpay node is used.
The wallet file is specified by wallet ID, and all wallet files are in the
configured data directory (which is `$HOME/.laqpay/wallets` by default).

#### Using laqpay as a library in a Go application

If your application is written in Go, you can interface with the CLI library
directly, see [Laqpay CLI Godoc](https://godoc.org/github.com/laqpay/laqpay/src/cli).

A REST API client is also available: [Laqpay REST API Client Godoc](https://godoc.org/github.com/laqpay/laqpay/src/api#Client).

#### Coinhours

Transaction fees in laqpay is paid in coinhours and is currently set to `50%`,
every transaction created burns `50%` of the total coinhours in all the input
unspents.

You need a minimum of `1` of coinhour to create a transaction.

Coinhours are generated at a rate of `1 coinsecond` per `second`
which are then converted to `coinhours`, `1` coinhour = `3600` coinseconds.

> Note: Coinhours don't have decimals and only show up in whole numbers.

##### REST API

When using the REST API, the coin hours sent to the destination and change can be controlled.
The 50% burn fee is still required.

See the [POST /wallet/transaction](https://github.com/laqpay/laqpay/blob/develop/src/api/README.md#create-transaction)
documentation for more information on how to control the coin hours.

We recommend sending at least 1 coin hour to each destination, otherwise the receiver will have to
wait for another coin hour to accumulate before they can make another transaction.

##### CLI

When using the CLI the amount of coinhours sent to the receiver is capped to
the number of coins they receive with a minimum of `1` coinhour for transactions
with `<1` laqpay being sent.

The coinhours left after burning `50%` and sending to receivers are sent to the change address.

For eg. If an address has `10` laqpays and `50` coinhours and only `1` unspent.
If we send `5` laqpays to another address then that address will receive
`5` laqpays and `5` coinhours, `26` coinhours will be burned.
The sending address will be left with `5` laqpays and `19` coinhours which
will then be sent to the change address.


### Verifying addresses

#### Using the CLI

```sh
laqpay-wallet-cli verifyAddress $addr
```

#### Using the REST API

Not directly supported, but API calls that have an address argument will return `400 Bad Request` if they receive an invalid address.

#### Using laqpay as a library in a Go application

https://godoc.org/github.com/laqpay/laqpay/src/cipher#DecodeBase58Address

```go
if _, err := cipher.DecodeBase58Address(address); err != nil {
    fmt.Println("Invalid address:", err)
    return
}
```

#### Using laqpay as a library in other applications

Address validation is available through a C wrapper, `liblaqpay`.

See the [liblaqpay documentation](/lib/cgo/README.md) for usage instructions.

### Checking Laqpay node connections

#### Using the CLI

Not implemented

#### Using the REST API

* `GET /api/v1/network/connections`

#### Using laqpay as a library in a Go application

Use the [Laqpay REST API Client](https://godoc.org/github.com/laqpay/laqpay/src/api#Client)

### Checking Laqpay node status

#### Using the CLI

```sh
laqpay-wallet-cli status
```

#### Using the REST API

A method similar to `laqpay-wallet-cli status` is not implemented, but these endpoints can be used:

* `GET /api/v1/health`
* `GET /api/v1/version`
* `GET /api/v1/blockchain/metadata`
* `GET /api/v1/blockchain/progress`

#### Using laqpay as a library in a Go application

Use the [Laqpay CLI package](https://godoc.org/github.com/laqpay/laqpay/src/cli)

## xpub wallets

You can create a wallet from an `xpub` key, which allows you to generate new
addresses without exposing secret keys. You can use any `xpub` key from an
existing Bitcoin or Ethereum HD wallet you may already have, or you can create
a new HD wallet in Laqpay.

To create an xpub wallet from scratch:

1. Create a bip44 wallet
2. Export an xpub key from the bip44 wallet
3. Create the xpub wallet

If you already have an xpub key, you can skip to step 3.

### Create a bip44 HD wallet

#### Using the CLI

```sh
laqpay-wallet-cli bip44-wallet.wlt -t bip44
```

#### Using the REST API

* `POST /api/v1/wallet/create`

### Export an xpub key from a bip44 wallet

#### Using the CLI

```sh
laqpay-wallet-cli walletKeyExport bip44-wallet.wlt -k xpub --path "0/0"
```

#### Using the REST API

Not possible

### Create an xpub wallet

#### Using the CLI

```sh
laqpay-wallet-cli walletCreate xpub-wallet.wlt -t xpub --xpub $MY_XPUB_KEY
```

#### Using the REST API

* `POST /api/v1/wallet/create`
