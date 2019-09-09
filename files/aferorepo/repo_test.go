package aferorepo

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
	"github.com/theothertomelliott/consistentcommit/files"
)

func TestCopyToTempDir(t *testing.T) {
	fs := afero.NewMemMapFs()
	file1Content := []byte("this is file 1")
	afero.WriteFile(fs, "/source/file1.txt", file1Content, os.ModePerm)
	file2Content := []byte("this is file 2")
	afero.WriteFile(fs, "/source/file2.txt", file2Content, os.ModePerm)

	repo := New(fs)
	dst, err := repo.CopyToTempDir("/source")
	if err != nil {
		t.Fatalf("error copying files: %v", err)
	}

	// Check content of dst directory
	fmt.Println(dst)
	file1Got, err := afero.ReadFile(fs, filepath.Join(dst, "file1.txt"))
	if err != nil {
		t.Errorf("error reading files: %v", err)
	}
	file2Got, err := afero.ReadFile(fs, filepath.Join(dst, "file2.txt"))
	if err != nil {
		t.Errorf("error reading files: %v", err)
	}

	if string(file1Content) != string(file1Got) {
		t.Errorf("file1.txt mismatch:\n%v", cmp.Diff(string(file1Content), string(file1Got)))
	}
	if string(file2Content) != string(file2Got) {
		t.Errorf("file2.txt mismatch:\n%v", cmp.Diff(string(file2Content), string(file2Got)))
	}
}

func TestDirContent(t *testing.T) {
	fs := afero.NewMemMapFs()
	file1Content := []byte("this is file 1")
	afero.WriteFile(fs, "/source/file1.txt", file1Content, os.ModePerm)
	file2Content := []byte("this is file 2")
	afero.WriteFile(fs, "/source/file2.txt", file2Content, os.ModePerm)

	repo := New(fs)
	got, err := repo.DirContent("/source")
	if err != nil {
		t.Fatalf("error copying files: %v", err)
	}

	expected := []files.File{
		&file{
			VPath:    "/source/file1.txt",
			VContent: file1Content,
		},
		&file{
			VPath:    "/source/file2.txt",
			VContent: file2Content,
		},
	}

	if !cmp.Equal(expected, got) {
		t.Errorf("file content didn't match (-expected, +got):\n%v", cmp.Diff(expected, got))
	}
}
