package virtualenv

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/theothertomelliott/consistentcommit/executor"
)

var _ executor.Executor = &Environment{}

type program func(fs afero.Fs, args []string, workingDir string, env func(string) string) (executor.Output, error)

type execution struct {
	programs map[string]program
}

func (e *execution) RegisterProgram(name string, r program) {
	e.programs[name] = r
}

func (e *Environment) Run(command executor.Command, env func(string) string) (executor.Output, error) {
	content, err := afero.ReadFile(e.fs, filepath.Join(command.WorkingDir, command.Executable))
	if err != nil {
		return nil, err
	}
	program, exists := e.programs[string(content)]
	if !exists {
		return nil, fmt.Errorf("program not registered for name %q", string(content))
	}
	return program(e.fs, command.Args, command.WorkingDir, env)
}
