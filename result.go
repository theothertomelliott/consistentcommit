package consistentcommit

type BuildResult interface {
	Compare(BuildResult, ComparisonConfig) ([]Difference, error)
}

type Difference interface {
	Describe() string
}
