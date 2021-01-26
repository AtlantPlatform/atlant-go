## Atlant Node [![CircleCI](https://circleci.com/gh/AtlantPlatform/atlant-go/tree/master.svg?style=shield&circle-token=c3fb3de524334e4566ccdae222192e585607f164)](https://circleci.com/gh/AtlantPlatform/atlant-go/tree/master) ![](https://img.shields.io/badge/version-1.0.1--rc1-blue.svg)

<p align="center">
<img src="https://avatars3.githubusercontent.com/u/30299272?s=400&u=b11d6a41091e04d7e133a758e6efb917371b981d&v=4" width="175">
</p>

### Intro

`atlant-go` is the reference Atlant Node implementation, which contains the distributed store of all information pertaining to Properties and PTOs (Property Token Offerings).

The node is designed to poll two main smart contract sets in order to fetch secure data: PTO contracts, handling property token offerings and property tokens behaviour, and KYC contracts, enabling fully compliant and transparent property token trading. The smart contract also queries the list of ATL Platform token holders and their balances.
All of this together helps to create smooth and compliant Platform operations and performance, thus facilitating optimal experience for all Platform members.

Both ATL and PTO token holders can derive value from operations conducted on the Platform, provided that two main requirements are met: 
- being a verified individual (having KYC details completed and registered with ATLANT)
- running Atlant Node, which helps to secure the network by distributing PTO data globally

### Setup

Get a pre-built binary for your platform [in the Releases](https://github.com/AtlantPlatform/atlant-node/releases) section

### Initialization

Prior to node startup, you should initialize it:

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
      --cluster-enabled        Enable cluster discovery (experimental). (env $AN_CLUSTER_ENABLED) (default "false")
  -C, --cluster-name           Specifies cluster name. (env $AN_CLUSTER_NAME)
  -N, --fs-network-profile     Sets IPFS network profile. Available: default, server, no-modify. (env $AN_FS_NETWORK_PROFILE) (default "default")
  -T, --testnet                Switch node into testing mode, it runs in a separate testnet environment. (env $AN_TESTNET_ENABLED)
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

The node must be initialized with `-T` flag beforehand. When running a node, specify your Ethereum address to participate in receiving a bonus from each successful PTO. The `-T` flag is not required, the testnet state will be detected from configs.

```
$ atlant-go -E 0xa936055b4c9b4a1213e64b7fc8c7ff295939ce71
INFO[0000] ATLANT TestNet welcomes you!
INFO[0000] atlant-go node is starting
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

### License

Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
Use of this source code is governed by BSD-3-Clause "New" or "Revised"
License (BSD-3-Clause) that can be found in the [LICENSE](/LICENSE) file.
