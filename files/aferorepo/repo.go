package aferorepo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/afero"
	"github.com/theothertomelliott/consistentcommit/files"
)

var _ files.Repo = &Repo{}

func New(fs afero.Fs) files.Repo {
	return &Repo{
		fs: fs,
	}
}

type Repo struct {
	fs afero.Fs
}

func (r *Repo) CopyToTempDir(sourceDir string) (string, error) {
	destination, err := r.MakeTempDir()
	if err != nil {
		return "", err
	}

	err = afero.Walk(r.fs, sourceDir, func(
		path string,
		info os.FileInfo,
		err error,
	) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		path = strings.Replace(path, sourceDir, "", 1)
		reader, err := r.fs.Open(filepath.Join(sourceDir, path))
		if err != nil {
			return err
		}
		err = afero.WriteReader(r.fs, filepath.Join(destination, path), reader)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return destination, nil
}

func (r *Repo) MakeTempDir() (string, error) {
	randomName := uuid.New()
	return afero.GetTempDir(r.fs, fmt.Sprintf("consistentcommit-%v", randomName.String())), nil
}

func (r *Repo) RmDir(dir string) error {
	return r.fs.RemoveAll(dir)
}

func (r *Repo) DirContent(dir string) ([]files.File, error) {
	var files []files.File
	err := afero.Walk(r.fs, dir, func(
		path string,
		info os.FileInfo,
		err error,
	) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		content, err := afero.ReadFile(r.fs, path)
		if err != nil {
			return nil
		}
		files = append(files, &file{
			VPath:    path,
			VContent: content,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
