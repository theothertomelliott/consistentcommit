package files

type Repo interface {
	CopyToTempDir(string) (string, error)
	MakeTempDir() (string, error)
	RmDir(string) error
	DirContent(string) ([]File, error)
}

type File interface {
	Path() string
	Content() []byte
}
