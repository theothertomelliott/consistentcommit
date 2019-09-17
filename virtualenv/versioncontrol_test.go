package virtualenv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TODO: Need to distinguish between the repo dir and a tmp dir

func TestCheckout(t *testing.T) {
	var tests = []struct {
		name          string
		env           *Environment
		workingDir    string
		commit        string
		expectedFiles map[string][]byte
		expectErr     bool
	}{
		{
			name:      "commit not found",
			env:       New(),
			commit:    "missing",
			expectErr: true,
		},
		{
			name: "commit with one file",
			env: func() *Environment {
				e := New()
				e.AddCommit("c1", tree{
					"file1": []byte("content1"),
					"file2": []byte("content2"),
				})
				return e
			}(),
			commit: "c1",
			expectedFiles: tree{
				"file1": []byte("content1"),
				"file2": []byte("content2"),
			},
		},
		{
			name: "selects correct commit with multiple",
			env: func() *Environment {
				e := New()
				e.AddCommit("c1", tree{
					"file1": []byte("content1"),
					"file2": []byte("content2"),
				})
				e.AddCommit("c2", tree{
					"file3": []byte("content3"),
					"file4": []byte("content4"),
				})
				return e
			}(),
			commit: "c1",
			expectedFiles: tree{
				"file1": []byte("content1"),
				"file2": []byte("content2"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.env.Checkout(test.workingDir, test.commit)
			if (err != nil) != test.expectErr {
				t.Errorf("error not as expected:%v", err)
			}
			for fn, expected := range test.expectedFiles {
				got, err := test.env.ReadFile(fn)
				if err != nil {
					t.Errorf("could not read file %q: %v", fn, err)
				} else {
					if !cmp.Equal(expected, got) {
						t.Errorf("content for %q not as expected", fn)
					}
				}
			}
		})
	}
}
