package executor

type Executor interface {
	Run(command Command, defaultEnv func(string) string) (Output, error)
}

type Output interface {
	Stdout() []byte
	Stderr() []byte
}
