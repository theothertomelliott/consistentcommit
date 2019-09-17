package virtualenv

import (
	"github.com/spf13/afero"
	"github.com/theothertomelliott/consistentcommit/files"
	"github.com/theothertomelliott/consistentcommit/files/aferorepo"
)

var (
	_ files.Repo = &Environment{}
)

func New() *Environment {
	fs := afero.NewMemMapFs()
	return &Environment{
		execution: execution{
			programs: make(map[string]program),
		},
		Repo: aferorepo.New(fs),
		fs:   fs,
	}
}

type Environment struct {
	execution

	files.Repo
	fs afero.Fs
}
