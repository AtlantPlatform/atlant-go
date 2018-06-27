package version

// Version components
const (
	Maj = "1"
	Min = "0"
	Fix = "0-rc2"
)

var (
	// Version is the current version of Tendermint
	// Must be a string because scripts like dist.sh read this file.
	Version = "1.0.0-rc2"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}
