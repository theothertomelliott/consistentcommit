package virtualenv

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/theothertomelliott/consistentcommit/versioncontrol"
)

var _ versioncontrol.VersionControl = &Environment{}

type tree map[string][]byte

type vcs struct {
	commits map[string]tree
}

func (v *vcs) AddCommit(name string, files tree) {
	if v.commits == nil {
		v.commits = make(map[string]tree)
	}
	v.commits[name] = files
}

func (e *Environment) Checkout(workingDir string, commit string) error {
	if e == nil {
		return nil
	}
	files, exists := e.commits[commit]
	if !exists {
		return fmt.Errorf("commit not found: %v", commit)
	}
	e.fs = afero.NewMemMapFs()
	for path, content := range files {
		err := e.AddFile(path, content)
		if err != nil {
			return err
		}
	}
	return nil
}
