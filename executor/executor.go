package executor

type Executor interface {
	Run(executable string, args []string, workingDir string, env func(string) string) (Output, error)
}

type Output interface {
	Stdout() []byte
	Stderr() []byte
}
