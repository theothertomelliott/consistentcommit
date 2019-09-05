package consistentcommit

type Config struct {
	GoldenCommit    string
	CandidateCommit string
	Build           BuildConfig
	Tests           []Command
	Comparison      ComparisonConfig
}

type ComparisonConfig struct {
}

type BuildConfig struct {
	BuildCommand   Command
	BuildOutputDir string
}

type Command struct {
	Executable string
	Args       []string
	Env        map[string]string
}
