package consistentcommit

import (
	"github.com/theothertomelliott/consistentcommit/executor"
	"github.com/theothertomelliott/consistentcommit/files"
)

type Builder interface {
	Build(config BuildConfig) (string, error)
}

type builder struct {
	EnvProvider
	Executor executor.Executor
	FileRepo files.Repo
}

func (b *builder) Build(config BuildConfig) (string, error) {
	_, err := b.Executor.Run(
		config.BuildCommand.Executable,
		config.BuildCommand.Args,
		"",
		b.env(config.BuildCommand),
	)
	if err != nil {
		return "", err
	}

	return b.FileRepo.CopyToTempDir(config.BuildOutputDir)
}
