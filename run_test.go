package consistentcommit

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/theothertomelliott/consistentcommit/executor"
	"github.com/theothertomelliott/consistentcommit/files"
)

func TestEmptyRunSucceeds(t *testing.T) {
	testDir := "workingDir"
	testCommand := executor.Command{}

	runner := &testRunner{
		Executor: &mockExecutor{
			run: func(command executor.Command, env func(string) string) (executor.Output, error) {
				if !cmp.Equal(testCommand, command) {
					t.Errorf("command not as expected:\n%v", cmp.Diff(testCommand, command))
				}
				outputEnv := env("OUTPUT_DIR")
				if outputEnv != "tmp" {
					t.Errorf("unexpected output dir env: %v", outputEnv)
				}
				return nil, nil
			},
		},
		FileRepo: &mockFileRepo{
			makeTempDir: func() (string, error) {
				return "tmp", nil
			},
			dirContent: func(string) ([]files.File, error) {
				return nil, nil
			},
		},
	}

	result, err := runner.Run(testDir, testCommand)
	if err != nil {
		t.Errorf("test run: %v", err)
	}
	if result == nil {
		t.Errorf("expected a result")
	}
}

func TestResultComparison(t *testing.T) {
	testDir := "workingDir"
	testCommand := executor.Command{}

	runner := &testRunner{
		Executor: &mockExecutor{
			run: func(command executor.Command, env func(string) string) (executor.Output, error) {
				return nil, nil
			},
		},
		FileRepo: &mockFileRepo{
			makeTempDir: func() (string, error) {
				return "tmp", nil
			},
			dirContent: func(string) ([]files.File, error) {
				return []files.File{
					mockFile{
						path:    "file1",
						content: []byte("content1"),
					},
				}, nil
			},
		},
	}

	result1, err := runner.Run(testDir, testCommand)
	if err != nil {
		t.Errorf("test run 1: %v", err)
	}

	runner.FileRepo = &mockFileRepo{
		makeTempDir: func() (string, error) {
			return "tmp", nil
		},
		dirContent: func(string) ([]files.File, error) {
			return []files.File{
				mockFile{
					path:    "file1",
					content: []byte("content2"),
				},
			}, nil
		},
	}

	result2, err := runner.Run(testDir, testCommand)
	if err != nil {
		t.Errorf("test run 1: %v", err)
	}

	diff, err := result1.Compare(result2, ComparisonConfig{})
	if err != nil {
		t.Errorf("comparison: %v", err)
	}
	if len(diff) != 1 {
		t.Errorf("unexpected differences: %v", len(diff))
	}
}
