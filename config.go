package main

import (
	"strconv"
	"strings"
	"time"

	cli "github.com/jawher/mow.cli"
)

// defaultLogLevel might be overridden by testing.go
var defaultLogLevel = "4"

var (
	goMaxProcs = app.String(cli.StringOpt{
		Name:   "p go-procs",
		Desc:   "The maximum number of CPUs that can be used simultaneously by Go runtime.",
		EnvVar: "AN_GOMAXPROCS",
		Value:  "128",
	})
	// logLevel is set in main func
	logLevel *string
)

var (
	stateDir = app.String(cli.StringOpt{
		Name:   "S state-dir",
		Desc:   "Directory prefix for state indexed storage.",
		EnvVar: "AN_STATE_DIR",
		Value:  "var/state",
	})
	stateGcInterval = app.String(cli.StringOpt{
		Name:   "state-gcinterval",
		Desc:   "Set a default GC interval for the state DB. Setting it lower will result in increased CPU load.",
		EnvVar: "AN_STATE_GCINTERVAL",
		Value:  "5m",
	})
	fsDir = app.String(cli.StringOpt{
		Name:   "F fs-dir",
		Desc:   "Directory prefix for IPFS filesystem storage.",
		EnvVar: "AN_FS_DIR",
		Value:  "var/fs",
	})
	logDir = app.String(cli.StringOpt{
		Name:   "log-dir",
		Desc:   "Directory prefix for logs",
		EnvVar: "AN_LOG_DIR",
		Value:  "var/log",
	})
	fsBootstrapPeers = app.Strings(cli.StringsOpt{
		Name:      "B bootstrap-peers",
		Desc:      "Append to the list of IPFS bootstrap peers.",
		EnvVar:    "AN_FS_BOOTSTRAP_PEERS",
		Value:     []string{},
		HideValue: true,
	})
	fsRelayEnabled = app.String(cli.StringOpt{
		Name:   "R relay-enabled",
		Desc:   "Enables IPFS relay support, may implicitly use extra network bandwidth.",
		EnvVar: "AN_FS_RELAY_ENABLED",
		Value:  "true",
	})
	fsWarmupDur = app.String(cli.StringOpt{
		Name:   "warmup",
		Desc:   "Allocate some time for IPFS to warmup and find peers.",
		EnvVar: "AN_FS_WARMUP_DUR",
		Value:  "5s",
	})
	fsListenAddr = app.String(cli.StringOpt{
		Name:   "L fs-listen-addr",
		Desc:   "Sets IPFS listen address to communicate with peers.",
		EnvVar: "AN_FS_LISTEN_ADDR",
		Value:  "0.0.0.0:33770",
	})
	webListenAddr = app.String(cli.StringOpt{
		Name:   "W web-listen-addr",
		Desc:   "Sets webserver listen address for public API.",
		EnvVar: "AN_WEB_LISTEN_ADDR",
		Value:  "0.0.0.0:33780",
	})
	clusterEnabled = app.String(cli.StringOpt{
		Name:   "cluster-enabled",
		Desc:   "Enable cluster discovery (experimental).",
		EnvVar: "AN_CLUSTER_ENABLED",
		Value:  "false",
	})
	clusterName = app.String(cli.StringOpt{
		Name:   "C cluster-name",
		Desc:   "Specifies cluster name.",
		EnvVar: "AN_CLUSTER_NAME",
		Value:  "",
	})
	fsNetworkProfile = app.String(cli.StringOpt{
		Name:   "N fs-network-profile",
		Desc:   "Sets IPFS network profile. Available: default, server, no-modify.",
		EnvVar: "AN_FS_NETWORK_PROFILE",
		Value:  "default",
	})
	envTestnet = app.Bool(cli.BoolOpt{
		Name:   "T testnet",
		Desc:   "Switch node into testing mode, it runs in a seprate testnet environment.",
		EnvVar: "AN_TESTNET_ENABLED",
		Value:  false,
	})
	envTestnetKey = app.String(cli.StringOpt{
		Name:      "testnet-key",
		Desc:      "Override the default testnet key with yours (generate it using atlant-keygen).",
		EnvVar:    "AN_TESTNET_KEY",
		Value:     testKey,
		HideValue: true,
	})
	envTestnetDomains = app.Strings(cli.StringsOpt{
		Name:      "testnet-auth-domains",
		Desc:      "Specify additional DNS authority domains for a testnet environment.",
		EnvVar:    "AN_TESTNET_DOMAINS",
		Value:     nil,
		HideValue: true,
	})
)

var (
	ethAddress = app.String(cli.StringOpt{
		Name:      "E ethereum-wallet",
		Desc:      "Specify Ethereum wallet to associate with work done in the session.",
		EnvVar:    "AN_ETHEREUM_WALLET",
		Value:     "",
		HideValue: true,
	})
	ethSignPath = app.String(cli.StringOpt{
		Name:      "sign",
		Desc:      "Specify private key location.",
		EnvVar:    "AN_TX_SIGN",
		Value:     "",
		HideValue: true,
	})
	ethPass = app.String(cli.StringOpt{
		Name:      "password",
		Desc:      "Specify passphrase for eth account.",
		EnvVar:    "AN_ETHEREUM_PASS",
		Value:     "",
		HideValue: true,
	})
)

// use atlant-keygen to generate a custom key
var (
	mainKey = "cc54285ef67145dd2d0f46f1ad1150cc1f1e22e0806e065c02746621a9ccda3d"
	testKey = "eb4d8d6b43697aff9ee41bcd05e7125838e428ca44c7bbc65cc6cf1c9b881d0b"
)

var (
	mainBootstrapPeers = []string{}
	testBootstrapPeers = []string{
		"/dns4/node-dev1.atlant.io/tcp/33770/ipfs/14V8BdHqHhExw4645xB3Xa2iheBrjYCMr7StXWUA9hBTqp8cM",
		"/dns4/node-dev2.atlant.io/tcp/33770/ipfs/14V8Bds64aUZJx6ag2TUXozS78Sko6fJ8kbkHF4bgvv9zgR6j",
		"/dns4/node-dev3.atlant.io/tcp/33770/ipfs/14V8BVs2FyU5qREKd68SgPqccrChiWX2uKdeeMtUhGfqJZjyK",
		"/dns4/node-dev4.atlant.io/tcp/33770/ipfs/14V8BTKR9MjKhfqgT4ybBjSb7kZHXmwvgBba7ujhN6ecTXJji",
	}
)

func duration(s string, defaults time.Duration) time.Duration {
	dur, err := time.ParseDuration(s)
	if err != nil {
		dur = defaults
	}
	return dur
}

func toList(s string) []string {
	return strings.Split(s, ",")
}

func toBool(s string) bool {
	switch strings.ToLower(s) {
	case "true", "1", "t", "yes":
		return true
	default:
		return false
	}
}

func toNatural(s string, defaults uint64) int {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		// defaults in case of incorrect or empty "" value
		return int(defaults)
	} else if i < 0 {
		// not defaults, because nobody expects +100 while specifying -100
		return 0
	}
	return int(i)
}
