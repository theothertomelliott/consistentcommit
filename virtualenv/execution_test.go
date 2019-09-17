package virtualenv

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
	"github.com/theothertomelliott/consistentcommit/executor"
)

func TestExecution(t *testing.T) {
	var tests = []struct {
		name        string
		programs    map[string]program
		files       map[string][]byte
		command     executor.Command
		env         func(string) string
		expectedOut executor.Output
		expectErr   bool
	}{
		{
			name: "program not found",
			command: executor.Command{
				Executable: "missing",
			},
			expectErr: true,
		},
		{
			name: "valid program",
			command: executor.Command{
				Executable: "app",
			},
			programs: map[string]program{
				"program": func(fs afero.Fs, args []string, workingDir string, env func(string) string) (executor.Output, error) {
					return &mockOutput{
						Out: []byte("out"),
						Err: []byte("err"),
					}, nil
				},
			},
			files: map[string][]byte{
				"app": []byte("program"),
			},
			expectedOut: &mockOutput{
				Out: []byte("out"),
				Err: []byte("err"),
			},
		},
		{
			name: "program returns error",
			command: executor.Command{
				Executable: "app",
			},
			programs: map[string]program{
				"program": func(fs afero.Fs, args []string, workingDir string, env func(string) string) (executor.Output, error) {
					return nil, fmt.Errorf("program failed")
				},
			},
			files: map[string][]byte{
				"app": []byte("program"),
			},
			expectErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			environment := New()
			for name, prog := range test.programs {
				environment.RegisterProgram(name, prog)
			}
			for path, content := range test.files {
				environment.AddFile(path, content)
			}
			output, err := environment.Run(test.command, nil)
			if !cmp.Equal(test.expectedOut, output) {
				t.Errorf("output not as expected:\n%v", cmp.Diff(test.expectedOut, output))
			}
			if (err != nil) != test.expectErr {
				t.Errorf("error not as expected:%v", err)
			}
		})
	}
}

var _ executor.Output = &mockOutput{}

type mockOutput struct {
	Out []byte
	Err []byte
}

func (m *mockOutput) Stdout() []byte {
	return m.Out
}

func (m *mockOutput) Stderr() []byte {
	return m.Err
}
