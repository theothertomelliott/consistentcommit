package executor

type Command struct {
	WorkingDir string
	Executable string
	Args       []string
	Env        map[string]string
}
