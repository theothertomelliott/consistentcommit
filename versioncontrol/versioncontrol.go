package versioncontrol

type VersionControl interface {
	Checkout(workingDir string, commit string) error
}
