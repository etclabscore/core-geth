package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// walkGen invokes its runTest argument for all subtests in the given directory.
// THE ONLY DIFFERENCE BETWEEN THIS METHOD AND ITS BROTHER tm.gen IS THAT THIS
// ONE DOESNT TRIM THE PATH PREFIX, THUS LEAVING THE FULL FILE AS NAME INTACT.
//
// runTest should be a function of type func(t *testing.T, name string, x <TestType>),
// where TestType is the type of the test contained in test files.
func (tm *testMatcher) walkGen(t *testing.T, dir string, runTest interface{}) {
	// Walk the directory.
	dirinfo, err := os.Stat(dir)
	if os.IsNotExist(err) || !dirinfo.IsDir() {
		fmt.Fprintf(os.Stderr, "can't find test files in %s, did you clone the tests submodule?\n", dir)
		t.Skip("missing test files")
	}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		name := filepath.ToSlash(path) // <-- The difference.
		if info.IsDir() {
			if _, skipload := tm.findSkip(name + "/"); skipload {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".json" {
			t.Run(name, func(t *testing.T) { tm.runTestFile(t, path, name, runTest) })
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

