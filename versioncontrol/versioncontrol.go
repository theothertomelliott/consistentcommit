package versioncontrol

type VersionControl interface {
	Checkout(commit string) error
}
