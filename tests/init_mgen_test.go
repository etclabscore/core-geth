package tests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// walkFullName invokes its runTest argument for all subtests in the given directory.
// THE ONLY DIFFERENCE BETWEEN THIS METHOD AND ITS BROTHER tm.gen IS THAT THIS
// ONE DOESNT TRIM THE PATH PREFIX, THUS LEAVING THE FULL FILE AS NAME INTACT.
//
// runTest should be a function of type func(t *testing.T, name string, x <TestType>),
// where TestType is the type of the test contained in test files.
func (tm *testMatcher) walkFullName(t *testing.T, dir string, runTest interface{}) {
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

// walk invokes its runTest argument for all subtests in the given directory.
//
// runTest should be a function of type func(t *testing.T, name string, x <TestType>),
// where TestType is the type of the test contained in test files.
func (tm *testMatcher) walkScanNDJSON(t *testing.T, dir string, runTest interface{}) {
	// Walk the directory.
	dirinfo, err := os.Stat(dir)
	if os.IsNotExist(err) || !dirinfo.IsDir() {
		fmt.Fprintf(os.Stderr, "can't find test files in %s, did you clone the tests submodule?\n", dir)
		t.Skip("missing test files")
	}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		name := filepath.ToSlash(strings.TrimPrefix(path, dir+string(filepath.Separator)))
		if info.IsDir() {
			if _, skipload := tm.findSkip(name + "/"); skipload {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".ndjson" {
			t.Run(name, func(t *testing.T) { tm.runTestFileNDJSON(t, path, name, runTest) })
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func (tm *testMatcher) runTestFileNDJSON(t *testing.T, path, name string, runTest interface{}) {
	if r, _ := tm.findSkip(name); r != "" {
		t.Skip(r)
	}
	if tm.whitelistpat != nil {
		if !tm.whitelistpat.MatchString(name) {
			t.Skip("Skipped by whitelist")
		}
	}
	t.Parallel()

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		m := makeVOfTestFunc(runTest)
		ma := m.Addr().Interface()
		err = json.Unmarshal(scanner.Bytes(), &ma)
		if err != nil {
			t.Fatal(err)
		}
		t.Run(name, func(t *testing.T) {
			if r, _ := tm.findSkip(name); r != "" {
				t.Skip(r)
			}
			runTestFuncNotMap(runTest, t, name, m)
		})
	}
}

func makeVOfTestFunc(f interface{}) reflect.Value {
	stringT := reflect.TypeOf("")
	testingT := reflect.TypeOf((*testing.T)(nil))
	ftyp := reflect.TypeOf(f)
	if ftyp.Kind() != reflect.Func || ftyp.NumIn() != 3 || ftyp.NumOut() != 0 || ftyp.In(0) != testingT || ftyp.In(1) != stringT {
		panic(fmt.Sprintf("bad test function type: want func(*testing.T, string, <TestType>), have %s", ftyp))
	}
	testType := ftyp.In(2)
	mp := reflect.New(testType)
	return mp.Elem()
}

func runTestFuncNotMap(runTest interface{}, t *testing.T, name string, m reflect.Value) {
	reflect.ValueOf(runTest).Call([]reflect.Value{
		reflect.ValueOf(t),
		reflect.ValueOf(name),
		m,
	})
}
