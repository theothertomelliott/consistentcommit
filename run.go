package consistentcommit

import (
	"fmt"

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

	outputDir, err := t.FileRepo.MakeTempDir()
	if err != nil {
		return nil, err
	}

	defaultEnv := t.env(command)
	env := func(key string) string {
		if key == "OUTPUT_DIR" {
			return outputDir
		}
		return defaultEnv(key)
	}

	_, err = t.Executor.Run(command.Executable, command.Args, testDir, env)
	if err != nil {
		return nil, err
	}

	files, err := t.FileRepo.DirContent(outputDir)
	if err != nil {
		return nil, err
	}
	return fileResult(files), nil
}

var _ BuildResult = fileResult{}

type fileResult []files.File

func (f fileResult) Compare(res BuildResult, cfg ComparisonConfig) ([]Difference, error) {
	b, isFr := res.(fileResult)
	if !isFr {
		return nil, fmt.Errorf("result not of same type")
	}

	var (
		myFilesByPath    = make(map[string]files.File)
		theirFilesByPath = make(map[string]files.File)
	)
	for _, file := range f {
		myFilesByPath[file.Path()] = file
	}
	for _, file := range b {
		theirFilesByPath[file.Path()] = file
	}

	var diffs []Difference
	for path, myFile := range myFilesByPath {
		theirFile := theirFilesByPath[path]
		if string(theirFile.Content()) != string(myFile.Content()) {
			diffs = append(diffs, fileDiff{
				description: path,
			})
		}
	}
	for path, _ := range theirFilesByPath {
		if _, exists := myFilesByPath[path]; !exists {
			diffs = append(diffs, fileDiff{
				description: path,
			})
		}
	}

	return diffs, nil
}

var _ Difference = fileDiff{}

type fileDiff struct {
	description string
}

func (f fileDiff) Describe() string {
	return f.description
}
