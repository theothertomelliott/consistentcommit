package afero

import "github.com/theothertomelliott/consistentcommit/files"

var _ files.File = &file{}

type file struct {
	VPath    string
	VContent []byte
}

func (f *file) Path() string {
	return f.VPath
}

func (f *file) Content() []byte {
	return f.VContent
}
