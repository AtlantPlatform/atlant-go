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
	// 1) auth mode:  trust all, trust local (future: trust RSA key)
	authMode = app.String(cli.StringOpt{
		Name:   "a auth",
		Desc:   "Whom to trust while adding new nodes (trust-all|trust-local)",
		EnvVar: "AUTH_MODE",
		Value:  "trust-local",
	})
	// pubKeys = app.String(cli.StringOpt{
	// 	Name:   "pk",
	// 	Desc:   "In case of RSA-auth mode - public parts of keys that should sign the message",
	// 	EnvVar: "KEYS",
	// 	Value:  "",
	// })

	// 2) persistence (storage) mode: memory or disk. keeping it simple so far
	storagePath = app.String(cli.StringOpt{
		Name:   "s storage",
		Desc:   "Storage folder. Leave empty to keep nodes permissions in memory until restart",
		EnvVar: "STORAGE_PATH",
		Value:  "",
	})

	webListenAddr = app.String(cli.StringOpt{
		Name:   "L listen",
		Desc:   "Sets listen address for web server",
		EnvVar: "LISTEN_ADDR",
		Value:  "0.0.0.0:33700",
	})
	enableMetrics = app.Bool(cli.BoolOpt{
		Name:   "m metrics",
		Desc:   "Enable prometheus metrics",
		EnvVar: "ENABLE_PROMETHEUS",
		Value:  true,
	})
	goMaxProcs = app.String(cli.StringOpt{
		Name:   "p go-procs",
		Desc:   "The maximum number of CPUs that can be used simultaneously by Go runtime.",
		EnvVar: "AN_GOMAXPROCS",
		Value:  "128",
	})
	// logLevel is set in main func
	logLevel *string
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
