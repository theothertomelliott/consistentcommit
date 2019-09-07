package consistentcommit

import (
	"fmt"
	"testing"
)

func TestCheckRun(t *testing.T) {
	currentCommit := "master"

	type br struct {
		BuildResult
		id string
	}

	simpleBuildConfig := BuildConfig{
		BuildCommand: Command{
			Executable: "build",
		},
		BuildOutputDir: "/build/output/dir",
	}

	builder := &mockBuilder{
		build: func(config BuildConfig) (string, error) {
			// TODO: More complete comparison
			if config.BuildOutputDir != simpleBuildConfig.BuildOutputDir {
				t.Errorf("config did not match: %v", config)
			}
			return fmt.Sprintf("buildOutput/%v", currentCommit), nil
		},
	}

	runner := &mockTestRunner{
		run: func(testDir string, command Command) (BuildResult, error) {
			if testDir != "buildOutput/candidate" && testDir != "buildOutput/golden" {
				t.Errorf("unexpected test dir: %v", testDir)
			}
			return &br{
				id: testDir,
				BuildResult: &mockBuildResult{
					compare: func(res BuildResult, cfg ComparisonConfig) ([]Difference, error) {
						if asBr, isBr := res.(*br); isBr {
							if asBr.id == testDir {
								t.Errorf("comparing build result with self")
							}
							// Confirm if the build result ids are among the expected ones
							if asBr.id != "buildOutput/candidate" && asBr.id != "buildOutput/golden" {
								t.Errorf("unexpected build result for comparison: %v", asBr.id)
							}
						} else {
							t.Errorf("build result wasn't the expected type")
						}
						return nil, nil
					},
				},
			}, nil
		},
	}

	vcs := &mockVersionControl{
		checkout: func(commit string) error {
			if commit != "golden" && commit != "candidate" {
				t.Errorf("unexpected commit: %v", commit)
			}
			currentCommit = commit
			return nil
		},
	}

	check := &Check{
		builder: builder,
		runner:  runner,
		vcs:     vcs,
	}

	cfg := Config{
		GoldenCommit:    "golden",
		CandidateCommit: "candidate",
		Build:           simpleBuildConfig,
		Tests: []Command{
			Command{
				Executable: "run",
			},
		},
	}

	err := check.Run(cfg)
	if err != nil {
		t.Errorf("Error value was: %v", err)
	}
}
