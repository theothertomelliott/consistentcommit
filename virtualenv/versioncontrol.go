package virtualenv

import "github.com/theothertomelliott/consistentcommit/versioncontrol"

var _ versioncontrol.VersionControl = &Environment{}

func (e *Environment) Checkout(commit string) error {
	panic("not implemented")
}
