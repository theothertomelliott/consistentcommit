package consistentcommit

import (
	"github.com/theothertomelliott/consistentcommit/executor"
	"github.com/theothertomelliott/consistentcommit/files"
)

type TestRunner interface {
	Run(testDir string, command Command) (BuildResult, error)
}

type testRunner struct {
	EnvProvider
	Executor executor.Executor
	FileRepo files.Repo
}

func (t *testRunner) Run(testDir string, command Command) (BuildResult, error) {
	_, err := t.Executor.Run(command.Executable, command.Args, testDir, t.env(command))
	if err != nil {
		return nil, err
	}

	panic("output directory not established")
	outputDir := ""

	files, err := t.FileRepo.DirContent(outputDir)
	if err != nil {
		return nil, err
	}

	return fileResult(files), nil
}

var _ BuildResult = fileResult{}

type fileResult []files.File

func (f fileResult) Compare(BuildResult, ComparisonConfig) ([]Difference, error) {
	panic("not implemented")
	return nil, nil
}
