package consistentcommit

import (
	"github.com/theothertomelliott/consistentcommit/executor"
	"github.com/theothertomelliott/consistentcommit/files"
)

var (
	_ files.Repo        = &mockFileRepo{}
	_ executor.Executor = &mockExecutor{}
)

type mockExecutor struct {
	run func(executable string, args []string, workingDir string, env func(string) string) (executor.Output, error)
}

func (m *mockExecutor) Run(executable string, args []string, workingDir string, env func(string) string) (executor.Output, error) {
	return m.run(executable, args, workingDir, env)
}

type mockFileRepo struct {
	copyToTempDir func(string) (string, error)
	rmDir         func(string) error
	dirContent    func(string) ([]files.File, error)
}

func (m *mockFileRepo) CopyToTempDir(src string) (string, error) {
	return m.copyToTempDir(src)
}

func (m *mockFileRepo) RmDir(dir string) error {
	return m.rmDir(dir)
}

func (m *mockFileRepo) DirContent(dir string) ([]files.File, error) {
	return m.dirContent(dir)
}
