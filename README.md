# Laqpay

[![Build Status](https://travis-ci.com/laqpay/laqpay.svg)](https://travis-ci.com/laqpay/laqpay)
[![GoDoc](https://godoc.org/github.com/laqpay/laqpay?status.svg)](https://godoc.org/github.com/laqpay/laqpay)
[![Go Report Card](https://goreportcard.com/badge/github.com/laqpay/laqpay)](https://goreportcard.com/report/github.com/laqpay/laqpay)

Laqpay is a next-generation cryptocurrency.

Laqpay was written from scratch and designed over four years to realize the
ideal of Bitcoin and represents the apex of cryptocurrency design.
Laqpay is not designed to add features to Bitcoin,
but rather improves Bitcoin by increasing simplicity,
security and stripping out everything non-essential.

Some people have hyped the Laqpay Project as leading into "Bitcoin 3.0".
The coin itself is not "Bitcoin 3.0",
but is rather "Bitcoin 1.0". Bitcoin is a prototype crypto-coin.
Laqpay was designed to be what Bitcoin would look like if it were built from
scratch, to remedy the rough edges in the Bitcoin design.

- no duplicate coin-base outputs
- enforced checks for hash collisions
- simple deterministic wallets
- no transaction malleability
- no signature malleability
- removal of the scripting language
- CoinJoin and normal transactions are indistinguishable
- elimination of edge-cases that prevent independent node implementations
- <=10 second transaction times
- elimination of the need for mining to achieve blockchain consensus

## Links

* [laqpay.net](https://www.laqpay.net)
* [Laqpay Blog](https://www.laqpay.net/blog)
* [Laqpay Docs](https://www.laqpay.net/docs)
* [Laqpay Blockchain Explorer](https://explorer.laqpay.net)
* [Laqpay Development Telegram Channel](https://t.me/laqpaydev)
* [Laqpay Github Wiki](https://github.com/laqpay/laqpay/wiki)

## Table of Contents

<!-- MarkdownTOC levels="1,2,3,4,5" autolink="true" bracket="round" -->

- [Changelog](#changelog)
- [Installation](#installation)
	- [Go 1.10+ Installation and Setup](#go-110-installation-and-setup)
	- [Go get laqpay](#go-get-laqpay)
	- [Run Laqpay from the command line](#run-laqpay-from-the-command-line)
	- [Show Laqpay node options](#show-laqpay-node-options)
	- [Run Laqpay with options](#run-laqpay-with-options)
	- [Docker image](#docker-image)
	- [Building your own images](#building-your-own-images)
	- [Development image](#development-image)
- [API Documentation](#api-documentation)
	- [REST API](#rest-api)
	- [Laqpay command line interface](#laqpay-command-line-interface)
- [Integrating Laqpay with your application](#integrating-laqpay-with-your-application)
- [Contributing a node to the network](#contributing-a-node-to-the-network)
- [Creating a new coin](#creating-a-new-coin)
- [Daemon CLI Options](#daemon-cli-options)
- [URI Specification](#uri-specification)
- [Wire protocol user agent](#wire-protocol-user-agent)
- [Development](#development)
	- [Modules](#modules)
	- [Client libraries](#client-libraries)
	- [Running Tests](#running-tests)
	- [Running Integration Tests](#running-integration-tests)
		- [Stable Integration Tests](#stable-integration-tests)
		- [Live Integration Tests](#live-integration-tests)
		- [Debugging Integration Tests](#debugging-integration-tests)
		- [Update golden files in integration testdata](#update-golden-files-in-integration-testdata)
	- [Test coverage](#test-coverage)
		- [Test coverage for the live node](#test-coverage-for-the-live-node)
	- [Formatting](#formatting)
	- [Code Linting](#code-linting)
	- [Profiling](#profiling)
	- [Fuzzing](#fuzzing)
		- [base58](#base58)
		- [encoder](#encoder)
	- [Dependencies](#dependencies)
		- [Rules](#rules)
		- [Management](#management)
	- [Configuration Modes](#configuration-modes)
		- [Development Desktop Client Mode](#development-desktop-client-mode)
		- [Server Daemon Mode](#server-daemon-mode)
		- [Electron Desktop Client Mode](#electron-desktop-client-mode)
		- [Standalone Desktop Client Mode](#standalone-desktop-client-mode)
	- [Wallet GUI Development](#wallet-gui-development)
		- [Translations](#translations)
	- [Releases](#releases)
		- [Update the version](#update-the-version)
		- [Pre-release testing](#pre-release-testing)
		- [Creating release builds](#creating-release-builds)
		- [Release signing](#release-signing)
- [Responsible Disclosure](#responsible-disclosure)

<!-- /MarkdownTOC -->

## Changelog

[CHANGELOG.md](CHANGELOG.md)

## Installation

Laqpay supports go1.10+.

### Go 1.10+ Installation and Setup

[Golang 1.10+ Installation/Setup](./INSTALLATION.md)

### Go get laqpay

```sh
$ go get github.com/laqpay/laqpay/cmd/...
```

This will download `github.com/laqpay/laqpay` to `$GOPATH/src/github.com/laqpay/laqpay`.

You can also clone the repo directly with `git clone https://github.com/laqpay/laqpay`,
but it must be cloned to this path: `$GOPATH/src/github.com/laqpay/laqpay`.

### Run Laqpay from the command line

```sh
$ cd $GOPATH/src/github.com/laqpay/laqpay
$ make run-client
```

### Show Laqpay node options

```sh
$ cd $GOPATH/src/github.com/laqpay/laqpay
$ make run-help
```

### Run Laqpay with options

Example:

```sh
$ cd $GOPATH/src/github.com/laqpay/laqpay
$ make ARGS="--launch-browser=false -data-dir=/custom/path" run
```

### Docker image

This is the quickest way to start using Laqpay using Docker.

```sh
$ docker volume create laqpay-data
$ docker volume create laqpay-wallet
$ docker run -ti --rm \
    -v laqpay-data:/data/.laqpay \
    -v laqpay-wallet:/wallet \
    -p 6000:6000 \
    -p 6420:6420 \
    laqpay/laqpay
```

This image has a `laqpay` user for the laqpay daemon to run, with UID and GID 10000.
When you mount the volumes, the container will change their owner, so you
must be aware that if you are mounting an existing host folder any content you
have there will be own by 10000.

The container will run with some default options, but you can change them
by just appending flags at the end of the `docker run` command. The following
example will show you the available options.

```sh
$ docker run --rm laqpay/laqpay -help
```

Access the dashboard: [http://localhost:6420](http://localhost:6420).

Access the API: [http://localhost:6420/version](http://localhost:6420/version).

### Building your own images

[Building your own images](docker/images/mainnet/README.md).

### Development image

The [laqpay/laqpaydev-cli docker image](docker/images/dev-cli/README.md) is provided in order to make
easy to start developing Laqpay. It comes with the compiler, linters, debugger
and the vim editor among other tools.

The [laqpay/laqpaydev-dind docker image](docker/images/dev-docker/README.md) comes with docker installed
and all tools available on `laqpay/laqpaydev-cli:develop` docker image.

Also, the [laqpay/laqpaydev-vscode docker image](docker/images/dev-vscode/README.md) is provided
to facilitate the setup of the development process with [Visual Studio Code](https://code.visualstudio.com)
and useful tools included in `laqpay/laqpaydev-cli`.

## API Documentation

### REST API

[REST API](src/api/README.md).

### Laqpay command line interface

[CLI command API](cmd/cli/README.md).

## Integrating Laqpay with your application

[Laqpay Integration Documentation](INTEGRATION.md)

## Contributing a node to the network

Add your node's `ip:port` to the [peers.txt](peers.txt) file.
This file will be periodically uploaded to https://downloads.laqpay.net/blockchain/peers.txt
and used to seed client with peers.

*Note*: Do not add Laqwire nodes to `peers.txt`.
Only add Laqpay nodes with high uptime and a static IP address (such as a Laqpay node hosted on a VPS).

## Creating a new coin

See the [newcoin tool README](./cmd/newcoin/README.md)

## Daemon CLI Options

See the [Laqpay Daemon CLI options](./cmd/laqpay/README.md)

## URI Specification

Laqpay URIs obey the same rules as specified in Bitcoin's [BIP21](https://github.com/bitcoin/bips/blob/master/bip-0021.mediawiki).
They use the same fields, except with the addition of an optional `hours` parameter, specifying the coin hours.

Example Laqpay URIs:

* `laqpay:2hYbwYudg34AjkJJCRVRcMeqSWHUixjkfwY`
* `laqpay:2hYbwYudg34AjkJJCRVRcMeqSWHUixjkfwY?amount=123.456&hours=70`
* `laqpay:2hYbwYudg34AjkJJCRVRcMeqSWHUixjkfwY?amount=123.456&hours=70&label=friend&message=Birthday%20Gift`

Additonally, if no `laqpay:` prefix is present when parsing, the string may be treated as an address:

* `2hYbwYudg34AjkJJCRVRcMeqSWHUixjkfwY`

However, do not use this URI in QR codes displayed to the user, because the address can't be disambiguated from other Laqfiber coins.

## Wire protocol user agent

[Wire protocol user agent description](https://github.com/laqpay/laqpay/wiki/Wire-protocol-user-agent)

## Development

We have two branches: `master` and `develop`.

`develop` is the default branch and will have the latest code.

`master` will always be equal to the current stable release on the website, and should correspond with the latest release tag.

### Modules

* `api` - REST API interface
* `cipher` - cryptographic library (key generation, addresses, hashes)
* `cipher/base58` - Base58 encoding
* `cipher/encoder` - reflect-based deterministic runtime binary encoder
* `cipher/encrypt` - at-rest data encryption (chacha20poly1305+scrypt)
* `cipher/go-bip39` - BIP-39 seed generation
* `cli` - CLI library
* `coin` - blockchain data structures (blocks, transactions, unspent outputs)
* `daemon` - top-level application manager, combining all components (networking, database, wallets)
* `daemon/gnet` - networking library
* `daemon/pex` - peer management
* `params` - configurable transaction verification parameters
* `readable` - JSON-encodable representations of internal structures
* `laqpay` - core application initialization and configuration
* `testutil` - testing utility methods
* `transaction` - methods for creating transactions
* `util` - miscellaneous utilities
* `visor` - top-level blockchain database layer
* `visor/blockdb` - low-level blockchain database layer
* `visor/historydb` - low-level blockchain database layer for historical blockchain metadata
* `wallet` - wallet file management

### Client libraries

Laqpay implements client libraries which export core functionality for usage from
other programming languages.

* [liblaqpay C client library and SWIG interface](https://github.com/laqpay/liblaqpay)
* [laqpay-lite: Javascript and mobile bindings](https://github.com/laqpay/laqpay-lite)

### Running Tests

```sh
$ make test
```

### Running Integration Tests

There are integration tests for the CLI and HTTP API interfaces. They have two
run modes, "stable" and "live".

The stable integration tests will use a laqpay daemon
whose blockchain is synced to a specific point and has networking disabled so that the internal
state does not change.

The live integration tests should be run against a synced or syncing node with networking enabled.

#### Stable Integration Tests

```sh
$ make integration-test-stable
```

or

```sh
$ ./ci-scripts/integration-test-stable.sh -v -w
```

The `-w` option, run wallet integrations tests.

The `-v` option, show verbose logs.

#### Live Integration Tests

The live integration tests run against a live runnning laqpay node, so before running the test, we
need to start a laqpay node:

```sh
$ ./run-daemon.sh
```

After the laqpay node is up, run the following command to start the live tests:

```sh
$ make integration-test-live
```

The above command will run all tests except the wallet-related tests. To run wallet tests, we
need to manually specify a wallet file, and it must have at least `2 coins` and `256 coinhours`,
it also must have been loaded by the node.

We can specify the wallet by setting two environment variables:

* `API_WALLET_ID`, which is the filename (without path), that is loaded by the daemon to test against.
  This is the `"id"` field in API requests. It is used by the API integration tests.
  The wallet directory that the daemon uses can be controlled with the `-wallet-dir` option.
* `CLI_WALLET_FILE`, which is the filename (with path), to be used by the CLI integration tests

If the wallet is encrypted, also set `WALLET_PASSWORD`.

Example of running the daemon with settings for integration tests:

```sh
$ export API_WALLET_ID="$valid_wallet_filename"
$ export CLI_WALLET_FILE="$HOME/.laqpay/wallets/$valid_wallet_filename"
$ export WALLET_PASSWORD="$wallet_password"
$ make run-integration-test-live
```

Then run the tests with the following command:

```sh
$ make integration-test-live-wallet
```

There are two other live integration test modes for CSRF disabled and networking disabled.

To run the CSRF disabled tests:

```sh
$ export API_WALLET_ID="$valid_wallet_filename"
$ export CLI_WALLET_FILE="$HOME/.laqpay/wallets/$valid_wallet_filename"
$ export WALLET_PASSWORD="$wallet_password"
$ make run-integration-test-live-disable-csrf
```

```sh
$ make integration-test-live-disable-csrf
```

To run the networking disabled tests, which require a live wallet:

```sh
$ export API_WALLET_ID="$valid_wallet_filename"
$ export CLI_WALLET_FILE="$HOME/.laqpay/wallets/$valid_wallet_filename"
$ export WALLET_PASSWORD="$wallet_password"
$ make run-integration-test-live-disable-networking
```

Then run the tests with the following command:

```sh
$ make integration-test-live-wallet
```

#### Debugging Integration Tests

Run specific test case:

It's annoying and a waste of time to run all tests to see if the test we real care
is working correctly. There's an option: `-r`, which can be used to run specific test case.
For example: if we only want to test `TestStableAddressBalance` and see the result, we can run:

```sh
$ ./ci-scripts/integration-test-stable.sh -v -r TestStableAddressBalance
```

#### Update golden files in integration testdata

Golden files are expected data responses from the CLI or HTTP API saved to disk.
When the tests are run, their output is compared to the golden files.

To update golden files, use the provided `make` command:

```sh
$ make update-golden-files
```

We can also update a specific test case's golden file with the `-r` option.
For example:
```sh
$ ./ci-scripts/integration-test-stable.sh -v -u -r TestStableAddressBalance
```

### Test coverage

Coverage is automatically generated for `make test` and integration tests run against a stable node.
This includes integration test coverage. The coverage output files are placed in `coverage/`.

To merge coverage from all tests into a single HTML file for viewing:

```sh
$ make check
$ make merge-coverage
```

Then open `coverage/all-coverage.html` in the browser.

#### Test coverage for the live node

Some tests can only be run with a live node, for example wallet spending tests.
To generate coverage for this, build and run the laqpay node in test mode before running the live integration tests.

In one shell:

```sh
$ make run-integration-test-live-cover
```

In another shell:

```sh
$ make integration-test-live
```

After the tests have run, CTRL-C to exit the process from the first shell.
A coverage file will be generated at `coverage/laqpay-live.coverage.out`.

Merge the coverage with `make merge-coverage` then open the `coverage/all-coverage.html` file to view it,
or generate the HTML coverage in isolation with `go tool cover -html`

### Formatting

All `.go` source files should be formatted `goimports`.  You can do this with:

```sh
$ make format
```

### Code Linting

Install prerequisites:

```sh
$ make install-linters
```

Run linters:

```sh
$ make lint
```

### Profiling

A full CPU profile of the program from start to finish can be obtained by running the node with the `-profile-cpu` flag.
Once the node terminates, a profile file is written to `-profile-cpu-file` (defaults to `cpu.prof`).
This profile can be analyzed with

```sh
$ go tool pprof cpu.prof
```

The HTTP interface for obtaining more profiling data or obtaining data while running can be enabled with `-http-prof`.
The HTTP profiling interface can be controlled with `-http-prof-host` and listens on `localhost:6060` by default.

See https://golang.org/pkg/net/http/pprof/ for guidance on using the HTTP profiler.

Some useful examples include:

```sh
$ go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10
$ go tool pprof http://localhost:6060/debug/pprof/heap
```

A web page interface is provided by http/pprof at http://localhost:6060/debug/pprof/.

### Fuzzing

Fuzz tests are run with [go-fuzz](https://github.com/dvyukov/go-fuzz).
[Follow the instructions on the go-fuzz page](https://github.com/dvyukov/go-fuzz) to install it.

Fuzz tests are written for the following packages:

#### base58

To fuzz the `cipher/base58` package,

```sh
$ make fuzz-base58
```

#### encoder

To fuzz the `cipher/encoder` package,

```sh
$ make fuzz-encoder
```

### Dependencies

#### Rules

Dependencies must not require `cgo`.  This means dependencies cannot be wrappers around C libraries.
Requiring `cgo` breaks cross compilation and interferes with repeatable (deterministic) builds.

Critical cryptographic dependencies used by code in package `cipher` are archived inside the `cipher` folder,
rather than in the `vendor` folder.  This prevents a user of the `cipher` package from accidentally using a
different version of the `cipher` dependencies than were developed, which could have catastrophic but hidden problems.

#### Management

Dependencies are managed with [dep](https://github.com/golang/dep).

To [install `dep` for development](https://github.com/golang/dep/blob/master/docs/installation.md#development):

```sh
$ go get -u github.com/golang/dep/cmd/dep
```

`dep` vendors all dependencies into the repo.

If you change the dependencies, you should update them as needed with `dep ensure`.

Use `dep help` for instructions on vendoring a specific version of a dependency, or updating them.

When updating or initializing, `dep` will find the latest version of a dependency that will compile.

Examples:

Initialize all dependencies:

```sh
$ dep init
```

Update all dependencies:

```sh
$ dep ensure -update -v
```

Add a single dependency (latest version):

```sh
$ dep ensure github.com/foo/bar
```

Add a single dependency (more specific version), or downgrade an existing dependency:

```sh
$ dep ensure github.com/foo/bar@tag
```

### Configuration Modes
There are 4 configuration modes in which you can run a laqpay node:
- Development Desktop Daemon
- Server Daemon
- Electron Desktop Client
- Standalone Desktop Client

#### Development Desktop Client Mode
This mode is configured via `run-client.sh`
```bash
$ ./run-client.sh
```

#### Server Daemon Mode
The default settings for a laqpay node are chosen for `Server Daemon`, which is typically run from source.
This mode is usually preferred to be run with security options, though `-disable-csrf` is normal for server daemon mode, it is left enabled by default.

```bash
$ ./run-daemon.sh
```

To disable CSRF:

```bash
$ ./run-daemon.sh -disable-csrf
```

#### Electron Desktop Client Mode
This mode configures itself via electron-main.js

#### Standalone Desktop Client Mode
This mode is configured by compiling with `STANDALONE_CLIENT` build tag.
The configuration is handled in `cmd/laqpay/laqpay.go`

### Wallet GUI Development

The compiled wallet source should be checked in to the repo, so that others do not need to install node to run the software.

Instructions for doing this:

[Wallet GUI Development README](src/gui/static/README.md)

#### Translations

You can find information about how to work with translation files in the [Translations README](./src/gui/static/src/assets/i18n/README.md).

### Releases

#### Update the version

0. If the `master` branch has commits that are not in `develop` (e.g. due to a hotfix applied to `master`), merge `master` into `develop`
0. Make sure the translations are up to date. See the [i18n README](./src/gui/static/src/assets/i18n/README.md) for instructions on how to update translations and how to check if they are up to date.
0. Compile the `src/gui/static/dist/` to make sure that it is up to date (see [Wallet GUI Development README](src/gui/static/README.md))
0. Update version strings to the new version in the following files: `electron/package-lock.json`, `electron/package.json`, `electron/laqpay/current-laqpay.json`, `src/cli/cli.go`, `src/gui/static/src/current-laqpay.json`, `src/cli/integration/testdata/status*.golden`, `template/coin.template`, `README.md` files .
0. If changes require a new database verification on the next upgrade, update `src/laqpay/laqpay.go`'s `DBVerifyCheckpointVersion` value
0. Update `CHANGELOG.md`: move the "unreleased" changes to the version and add the date
0. Update the files in https://github.com/laqpay/repo-info by following the [metadata update procedure](https://github.com/laqpay/repo-info/#updating-laqpay-repository-metadate),
0. Merge these changes to `develop`
0. Follow the steps in [pre-release testing](#pre-release-testing)
0. Make a PR merging `develop` into `master`
0. Review the PR and merge it
0. Tag the `master` branch with the version number. Version tags start with `v`, e.g. `v0.20.0`.
    Sign the tag. If you have your GPG key in github, creating a release on the Github website will automatically tag the release.
    It can be tagged from the command line with `git tag -as v0.20.0 $COMMIT_ID`, but Github will not recognize it as a "release".
0. Make sure that the client runs properly from the `master` branch
0. Release builds are created and uploaded by travis. To do it manually, checkout the `master` branch and follow the [create release builds](electron/README.md) instructions.

If there are problems discovered after merging to `master`, start over, and increment the 3rd version number.
For example, `v0.20.0` becomes `v0.20.1`, for minor fixes.

#### Pre-release testing

Performs these actions before releasing:

* `make check`
* `make integration-test-live`
* `make integration-test-live-disable-networking` (requires node run with `-disable-networking`)
* `make integration-test-live-disable-csrf` (requires node run with `-disable-csrf`)
* `make intergration-test-live-wallet` (see [live integration tests](#live-integration-tests)) 6 times: with an unencrypted and encrypted wallet for each wallet type: `deterministic`, `bip44` and `collection`
* `go run cmd/cli/cli.go checkdb` against a fully synced database
* `go run cmd/cli/cli.go checkDBDecoding` against a fully synced database
* On all OSes, make sure that the client runs properly from the command line (`./run-client.sh` and `./run-daemon.sh`)
* Build the releases and make sure that the Electron client runs properly on Windows, Linux and macOS.
    * Use a clean data directory with no wallets or database to sync from scratch and verify the wallet setup wizard.
    * Load a test wallet with nonzero balance from seed to confirm wallet loading works
    * Send coins to another wallet to confirm spending works
    * Restart the client, confirm that it reloads properly
* For both the Android and iOS mobile wallets, configure the node url to be https://staging.node.laqpay.net
  and test all operations to ensure it will work with the new node version.

#### Creating release builds

[Create Release builds](electron/README.md).

#### Release signing

Releases are signed with this PGP key:

`0x5801631BD27C7874`

The fingerprint for this key is:

```
pub   ed25519 2017-09-01 [SC] [expires: 2023-03-18]
      10A7 22B7 6F2F FE7B D238  0222 5801 631B D27C 7874
uid                      GZ-C LAQPAY <token@protonmail.com>
sub   cv25519 2017-09-01 [E] [expires: 2023-03-18]
```

Keybase.io account: https://keybase.io/gzc

Follow the [Tor Project's instructions for verifying signatures](https://www.torproject.org/docs/verifying-signatures.html.en).

If you can't or don't want to import the keys from a keyserver, the signing key is available in the repo: [gz-c.asc](gz-c.asc).

Releases and their signatures can be found on the [releases page](https://github.com/laqpay/laqpay/releases).

Instructions for generating a PGP key, publishing it, signing the tags and binaries:
https://gist.github.com/gz-c/de3f9c43343b2f1a27c640fe529b067c

## Responsible Disclosure

Security flaws in laqpay source or infrastructure can be sent to security@laqpay.net.
Bounties are available for accepted critical bug reports.

PGP Key for signing:

```
-----BEGIN PGP PUBLIC KEY BLOCK-----

mDMEWaj46RYJKwYBBAHaRw8BAQdApB44Kgde4Kiax3M9Ta+QbzKQQPoUHYP51fhN
1XTSbRi0I0daLUMgU0tZQ09JTiA8dG9rZW5AcHJvdG9ubWFpbC5jb20+iJYEExYK
AD4CGwMFCwkIBwIGFQgJCgsCBBYCAwECHgECF4AWIQQQpyK3by/+e9I4AiJYAWMb
0nx4dAUCWq/TNwUJCmzbzgAKCRBYAWMb0nx4dKzqAP4tKJIk1vV2bO60nYdEuFB8
FAgb5ITlkj9PyoXcunETVAEAhigo4miyE/nmE9JT3Q/ZAB40YXS6w3hWSl3YOF1P
VQq4OARZqPjpEgorBgEEAZdVAQUBAQdAa8NkEMxo0dr2x9PlNjTZ6/gGwhaf5OEG
t2sLnPtYxlcDAQgHiH4EGBYKACYCGwwWIQQQpyK3by/+e9I4AiJYAWMb0nx4dAUC
Wq/TTQUJCmzb5AAKCRBYAWMb0nx4dFPAAQD7otGsKbV70UopH+Xdq0CDTzWRbaGw
FAoZLIZRcFv8zwD/Z3i9NjKJ8+LS5oc8rn8yNx8xRS+8iXKQq55bDmz7Igw=
=5fwW
-----END PGP PUBLIC KEY BLOCK-----
```

Key ID: [0x5801631BD27C7874](https://pgp.mit.edu/pks/lookup?search=0x5801631BD27C7874&op=index)

The fingerprint for this key is:

```
pub   ed25519 2017-09-01 [SC] [expires: 2023-03-18]
      10A7 22B7 6F2F FE7B D238  0222 5801 631B D27C 7874
uid                      GZ-C LAQPAY <token@protonmail.com>
sub   cv25519 2017-09-01 [E] [expires: 2023-03-18]
```

Keybase.io account: https://keybase.io/gzc
