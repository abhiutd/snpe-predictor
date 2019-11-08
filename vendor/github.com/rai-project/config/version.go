package config

// VersionInfo ...
type VersionInfo struct {
	// Version is populated at compile time by govvv from ./VERSION
	Version string
	// GitCommit is populated at compile time by govvv.
	BuildDate string
	// GitState is populated at compile time by govvv.
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
}
