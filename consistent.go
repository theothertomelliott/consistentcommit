package consistentcommit

import (
	"fmt"

	"github.com/theothertomelliott/consistentcommit/versioncontrol"
)

type Check struct {
	builder Builder
	runner  TestRunner
	vcs     versioncontrol.VersionControl
}

func (c *Check) Run(cfg Config) error {
	// Get the golden build
	goldenDir, err := c.getDirForCommit(cfg.GoldenCommit, cfg.Build)
	if err != nil {
		return err
	}

	// Get the candidate build
	candidateDir, err := c.getDirForCommit(cfg.CandidateCommit, cfg.Build)
	if err != nil {
		return err
	}

	for _, test := range cfg.Tests {
		// 	Execute the test for the golden commit
		goldenResult, err := c.runner.Run(goldenDir, test)
		if err != nil {
			return err
		}

		// 	Execute the test for the candidate commit
		candidateResult, err := c.runner.Run(candidateDir, test)
		if err != nil {
			return err
		}

		// 	Compare the results
		differences, err := goldenResult.Compare(candidateResult, cfg.Comparison)
		if err != nil {
			return err
		}
		if len(differences) > 0 {
			return fmt.Errorf("found %d differences", len(differences))
		}
	}

	return nil
}

func (c *Check) getDirForCommit(commit string, bCfg BuildConfig) (string, error) {
	// 	Get the commit
	err := c.vcs.Checkout(commit)
	if err != nil {
		return "", nil
	}

	// 	Run the build and get the artifact
	return c.builder.Build(bCfg)
}
