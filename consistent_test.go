package consistentcommit

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/theothertomelliott/consistentcommit/executor"
)

func TestCheckRun(t *testing.T) {
	currentCommit := "master"

	type br struct {
		BuildResult
		id string
	}

	simpleBuildConfig := BuildConfig{
		BuildCommand: executor.Command{
			Executable: "build",
		},
		BuildOutputDir: "/build/output/dir",
	}

	builder := &mockBuilder{
		build: func(config BuildConfig) (string, error) {
			if !cmp.Equal(config, simpleBuildConfig) {
				t.Errorf("config mismatch (-want +got):\n%v", cmp.Diff(config, simpleBuildConfig))
			}
			return fmt.Sprintf("buildOutput/%v", currentCommit), nil
		},
	}

	runner := &mockTestRunner{
		run: func(testDir string, command executor.Command) (BuildResult, error) {
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
		checkout: func(workingDir string, commit string) error {
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
		Tests: []executor.Command{
			executor.Command{
				Executable: "run",
			},
		},
	}

	err := check.Run(cfg)
	if err != nil {
		t.Errorf("Error value was: %v", err)
	}
}
