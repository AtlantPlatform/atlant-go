## ATLANT Node

Build status: [![CircleCI](https://circleci.com/gh/AtlantPlatform/atlant-go/tree/master.svg?style=svg&circle-token=c3fb3de524334e4566ccdae222192e585607f164)]()

### Intro

`atlant-go` is a proof-of-concept ATLANT Node implementation, which contains the distributed store of all information pertaining to PTOs (Property Token Offerings), as described in the [ATLANT white paper](https://atlant.io/assets/documents/en/Atlant_WP_publish.pdf) (beginning on page 33).

The node is designed to poll two main smart contract sets in order to fetch secure data: PTO contracts, handling property token offerings and property tokens behaviour, and KYC contracts, enabling fully compliant and transparent property token trading. The smart contract also queries the list of ATL Platform token holders and their balances.
All of this together helps to create smooth, compliant ATLANT Platform operations and performance, thus facilitating optimal experience for all platform participants.
As per the ATLANT White Paper, both ATL token holders and PTO token holders can derive value from operations conducted on the Platform, provided that two main requirements are met: being a verified individual (having KYC details completed and registered with ATLANT) and running an ATLANT node, which helps to secure the network by distributing PTO data globally, thus increasing redundancy.

### Setup

Get a pre-built binary for your platform in [releases](https://github.com/AtlantPlatform/atlant-go/releases) section or compile manually after installing [Go](https://golang.org/dl/):

```
$ go get github.com/AtlantPlatform/atlant-go
```

This operation will consume ~1 GB of hard drive space, the most heavy download will be go-ethereum dependency.

### Initialisation

Prior to node startup, you should initialise it:

```
$ atlant-go -T init
INFO[0000] atlant-go init
INFO[0000] initilizing within ATLANT Node TestNet

$ atlant-go -h
Usage: atlant-go [OPTIONS] COMMAND [arg...]

ATLANT Node

Options:
  -p, --go-procs               The maximum number of CPUs that can be used simultaneously by Go runtime. (env $AN_GOMAXPROCS) (default "128")
  -S, --state-dir              Directory prefix for state indexed storage. (env $AN_STATE_DIR) (default "var/state")
  -F, --fs-dir                 Directory prefix for IPFS filesystem storage. (env $AN_FS_DIR) (default "var/fs")
      --log-dir                Directory prefix for logs (env $AN_LOG_DIR) (default "var/log")
  -B, --bootstrap-peers        The list of IPFS bootstrap peers. (env $AN_FS_BOOTSTRAP_PEERS)
  -R, --relay-enabled          Enables IPFS relay support, may implicitly use extra network bandwidth. (env $AN_FS_RELAY_ENABLED) (default "true")
      --warmup                 Allocate some time for IPFS to warmup and find peers. (env $AN_FS_WARMUP_DUR) (default "5s")
  -L, --fs-listen-addr         Sets IPFS listen address to communicate with peers. (env $AN_FS_LISTEN_ADDR) (default "0.0.0.0:33770")
  -W, --web-listen-addr        Sets webserver listen address for public API. (env $AN_WEB_LISTEN_ADDR) (default "0.0.0.0:33780")
  -N, --fs-network-profile     Sets IPFS network profile. Available: default, server, no-modify. (env $AN_FS_NETWORK_PROFILE) (default "default")
  -T, --testnet                Switch node into testing mode, it runs in a seprate testnet environment. (env $AN_TESTNET_ENABLED)
      --testnet-key            Override the default testnet key with yours (generate it using atlant-keygen). (env $AN_TESTNET_KEY)
      --testnet-auth-domains   Specify additional DNS authority domains for a testnet environment. (env $AN_TESTNET_DOMAINS)
  -E, --ethereum-wallet        Specify Ethereum wallet to associate with work done in the session. (env $AN_ETHEREUM_WALLET)
  -l, --log-level              Logging verbosity (0 = minimum, 1...4, 5 = debug). (env $AN_LOG_LEVEL) (default "4")

Commands:
  init                         Initialize node and its IPFS repo.
  version                      Show version info.

Run 'atlant-go COMMAND --help' for more information on a command.
```

### Running in a testnet

The node must be initialised with `-T` flag beforehand. When running a node, specify your Ethereum address to participate in receiving a bonus from each successful PTO. The `-T` flag is not required, the testnet state will be detected from configs.

```
$ atlant-go -E 0xa936055b4c9b4a1213e64b7fc8c7ff295939ce71
INFO[0000] ATLANT TestNet welcomes you!
INFO[0000] atlant-go node is starting
```

If you want to start your own net, generate a new key with `atlant-keygen`:

```
$ go get github.com/AtlantPlatform/cmd/atlant-keygen
$ atlant-keygen net
0467345b7004890e5b7d325316d3ab15a6a20e1e34ec78b0d3abbdcd86793859
```

Then supply it when initialising a node:

```
export AN_TESTNET_KEY=0467345b7004890e5b7d325316d3ab15a6a20e1e34ec78b0d3abbdcd86793859
$ atlant-go -T init
$ atlant-go -E 0xa936055b4c9b4a1213e64b7fc8c7ff295939ce71
```

### API

The web server by default runs at http://localhost:33780
To browse all content within your browser, go to http://localhost:33780/index for an Apache2-styled autoindex.

* `POST /api/v1/put/:path` — writes a document to a path, overwriting if exists, you can specify HTTP Headers:
    - `X-Meta-UserMeta` — JSON encoded user-meta data blob;
* `POST /api/v1/delete/:id` — deletes a specific record by its ID;
* `GET /api/v1/content/:path` — access content located at path, returns meta info in HTTP Headers:
    - `X-Meta-ID` — record ID;
    - `X-Meta-Version` — current record version;
    - `X-Meta-Previous` — previous record version, if exists;
    - `X-Meta-Path` — record path;
    - `X-Meta-UserMeta` — user meta data;
    - `X-Meta-Deleted` — specifies whether record has been deleted.
* `GET /api/v1/listVersions/:path` — list all available versions of a record.
* `GET /api/v1/listAll/:prefix` — list all records with matching prefix (might be a lot of record).
* `GET /api/v1/meta/:path` — access record meta only, example JSON response:
```json
{
    "id": "01CBKY9WEHMS2XFY7KMED1XAPH",
    "path": "/files/file2",
    "createdAt": 1524308969465054914,
    "version": "QmYhNy5gWjBEGr6kZcgyHhrnjTzuVS525yR4K3gRRZmBXu",
    "versionPrevious": "QmXs854VAXyanT8QiHbx8NkvgjrCC56nnyQhqf2g1Dpv4z",
    "isDeleted": false,
    "size": 5,
    "userMeta": "eyJmb28iOiJiYXIifQ=="
}
```

Both `meta` and `content` accessors allow to pass a specfic version in query params, e.g. `?ver=QmXs854VAXyanT8QiHbx8NkvgjrCC56nnyQhqf2g1Dpv4z`.

* `GET /api/v1/ethBalance` — returns ETH balance of default account (specified during node startup with `-E` flag);
* `GET /api/v1/atlBalance` — returns ATL balance in ATLANT Tokens;
* `GET /api/v1/ptoBalance/:name` — returns PTO coin balance, each PTO token has different name; Example: `/ptoBalance/atl123`.
* `GET /api/v1/kycStatus` — returns Know Your Customer status info;

For all Ethereum info methods above, you can specify any specific account address in query params, e.g. `?account=0xa936055b4c9b4a1213e64b7fc8c7ff295939ce71`.

* `GET /api/v1/stats` — returns various internal stats.
* `GET /api/v1/ping`
* `GET /api/v1/env`
* `GET /api/v1/session`
* `GET /api/v1/version`
* `GET /api/v1/logs` — lists all available log files, each log file is rotated daily;
* `GET /api/v1/log/:year/:month/:day` — access a specific log file by day, e.g. `/2018/04/23`.

### Lite CLI

`atlant-lite` implements a light client for the network, allowing to browse and upload files by calling the API of a full node.

```
$ go get github.com/AtlantPlatform/atlant-go/cmd/atlant-lite

$ atlant-lite

Usage: atlant-lite [OPTIONS] COMMAND [arg...]

Options:
  -A, --addr   Full node address (default "localhost:33780")

Commands:
  ping         Ping node and get its ID
  version      Get node version
  put          Put an object into the store
  get          Get object contents from the store
  meta         Get object meta data from the store
  delete       Delete object from a store
  versions     List all object versions
  ls           List all objects and sub-directories in a prefix

Run 'atlant-lite COMMAND --help' for more information on a command.
```

To make calls into the official testet, use either node address or an alias ("testnet" for any node, "testnetN" for node N), e.g.

```
$ atlant-lite -A testnet version
1.0.0-rc4
```


### License

Copyright 2017-2019 Tensigma Ltd. All rights reserved.
Use of this source code is governed by Microsoft Reference Source
License (MS-RSL) that can be found in the [LICENSE](/LICENSE) file.
