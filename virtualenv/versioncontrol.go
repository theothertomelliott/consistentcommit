package virtualenv

import (
	"fmt"
	"path"
	"strings"

	"github.com/theothertomelliott/consistentcommit/versioncontrol"
)

var _ versioncontrol.VersionControl = &Environment{}

type (
	tree      map[string][]byte
	commitRef struct {
		repo   string
		commit string
	}
	vcs struct {
		repoDirs map[string]struct{}
		commits  map[commitRef]tree
	}
)

func (v *vcs) AddCommit(workingDir string, name string, files tree) {
	if v.commits == nil {
		v.commits = make(map[commitRef]tree)
	}
	if v.repoDirs == nil {
		v.repoDirs = make(map[string]struct{})
	}
	v.commits[commitRef{
		repo:   workingDir,
		commit: name,
	}] = files
	v.repoDirs[workingDir] = struct{}{}
}

func (e *Environment) Checkout(workingDir string, commit string) error {
	if e == nil {
		return nil
	}
	var repoName *string
	for repo := range e.repoDirs {
		if strings.HasPrefix(workingDir, repo) {
			r := repo
			repoName = &r
		}
	}
	if repoName == nil {
		return fmt.Errorf("no repo found at: %v", workingDir)
	}
	files, exists := e.commits[commitRef{
		repo:   *repoName,
		commit: commit,
	}]
	if !exists {
		return fmt.Errorf("commit not found: %v", commit)
	}
	e.fs.RemoveAll(*repoName)
	for filePath, content := range files {
		err := e.AddFile(path.Join(workingDir, filePath), content)
		if err != nil {
			return err
		}
	}
	return nil
}
