package consistentcommit

import "github.com/theothertomelliott/consistentcommit/executor"

type Config struct {
	GoldenCommit    string
	CandidateCommit string
	Build           BuildConfig
	Tests           []executor.Command
	Comparison      ComparisonConfig
}
type ComparisonConfig struct {
}

type BuildConfig struct {
	BuildCommand   executor.Command
	BuildOutputDir string
}
