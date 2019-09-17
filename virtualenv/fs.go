package virtualenv

import (
	"os"
	"github.com/spf13/afero"
	"github.com/theothertomelliott/consistentcommit/files"
)

var _ files.Repo = &Environment{}

func (e *Environment) AddFile(path string, content []byte) error {
	return afero.WriteFile(e.fs, path, content, os.ModePerm)
}
