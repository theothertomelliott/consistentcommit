package consistentcommit

import (
	"github.com/theothertomelliott/consistentcommit/executor"
	"github.com/theothertomelliott/consistentcommit/files"
	"github.com/theothertomelliott/consistentcommit/versioncontrol"
)

var (
	_ files.Repo                    = &mockFileRepo{}
	_ executor.Executor             = &mockExecutor{}
	_ Builder                       = &mockBuilder{}
	_ TestRunner                    = &mockTestRunner{}
	_ versioncontrol.VersionControl = &mockVersionControl{}
	_ BuildResult                   = &mockBuildResult{}
	_ files.File                    = &mockFile{}
)

type mockExecutor struct {
	run func(command executor.Command, env func(string) string) (executor.Output, error)
}

func (m *mockExecutor) Run(command executor.Command, env func(string) string) (executor.Output, error) {
	return m.run(command, env)
}

type mockFileRepo struct {
	copyToTempDir func(string) (string, error)
	makeTempDir   func() (string, error)
	rmDir         func(string) error
	dirContent    func(string) ([]files.File, error)
}

func (m *mockFileRepo) CopyToTempDir(src string) (string, error) {
	return m.copyToTempDir(src)
}

func (m *mockFileRepo) MakeTempDir() (string, error) {
	return m.makeTempDir()
}

func (m *mockFileRepo) RmDir(dir string) error {
	return m.rmDir(dir)
}

func (m *mockFileRepo) DirContent(dir string) ([]files.File, error) {
	return m.dirContent(dir)
}

type mockBuilder struct {
	build func(config BuildConfig) (string, error)
}

func (m *mockBuilder) Build(config BuildConfig) (string, error) {
	return m.build(config)
}

type mockVersionControl struct {
	checkout func(workingDir string, commit string) error
}

func (m *mockVersionControl) Checkout(workingDir string, commit string) error {
	return m.checkout(workingDir, commit)
}

type mockTestRunner struct {
	run func(testDir string, command executor.Command) (BuildResult, error)
}

func (m *mockTestRunner) Run(testDir string, command executor.Command) (BuildResult, error) {
	return m.run(testDir, command)
}

type mockBuildResult struct {
	compare func(BuildResult, ComparisonConfig) ([]Difference, error)
}

func (m *mockBuildResult) Compare(res BuildResult, cfg ComparisonConfig) ([]Difference, error) {
	return m.compare(res, cfg)
}

type mockFile struct {
	path    string
	content []byte
}

func (m mockFile) Path() string {
	return m.path
}

func (m mockFile) Content() []byte {
	return m.content
}
