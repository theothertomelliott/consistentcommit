package consistentcommit

import (
	"fmt"
	"testing"

	"github.com/theothertomelliott/consistentcommit/executor"
)

func TestBuild(t *testing.T) {
	var tests = []struct {
		name      string
		builder   *builder
		config    BuildConfig
		expected  string
		expectErr bool
	}{
		{
			name: "simple",
			builder: &builder{
				Executor: &mockExecutor{
					run: func(executable string, args []string, workingDir string, env func(string) string) (executor.Output, error) {
						if executable != "executable" {
							t.Errorf("incorrect executable: %v", executable)
						}
						return nil, nil
					},
				},
				FileRepo: &mockFileRepo{
					copyToTempDir: func(src string) (string, error) {
						if src != "output" {
							t.Errorf("incorrect source dir: %v", src)
						}
						return "tmp", nil
					},
				},
			},
			config: BuildConfig{
				BuildCommand: Command{
					Executable: "executable",
				},
				BuildOutputDir: "output",
			},
			expected: "tmp",
		},
		{
			name: "executor errors out",
			builder: &builder{
				Executor: &mockExecutor{
					run: func(executable string, args []string, workingDir string, env func(string) string) (executor.Output, error) {
						return nil, fmt.Errorf("executable failure")
					},
				},
				FileRepo: &mockFileRepo{
					copyToTempDir: func(src string) (string, error) {
						t.Error("files should not have been copied")
						return "", nil
					},
				},
			},
			config: BuildConfig{
				BuildCommand: Command{
					Executable: "executable",
				},
				BuildOutputDir: "output",
			},
			expectErr: true,
		},
		{
			name: "copy errors out",
			builder: &builder{
				Executor: &mockExecutor{
					run: func(executable string, args []string, workingDir string, env func(string) string) (executor.Output, error) {
						return nil, nil
					},
				},
				FileRepo: &mockFileRepo{
					copyToTempDir: func(src string) (string, error) {
						return "", fmt.Errorf("copy failure")
					},
				},
			},
			config: BuildConfig{
				BuildCommand: Command{
					Executable: "executable",
				},
				BuildOutputDir: "output",
			},
			expectErr: true,
		},
	}
	for _, test := range tests {
		got, err := test.builder.Build(test.config)
		if got != test.expected {
			t.Errorf("Expected %q, got %q", test.expected, got)
		}
		if test.expectErr != (err != nil) {
			t.Errorf("Error value was: %v", err)
		}
	}
}
