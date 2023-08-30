package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// meowsbits:
//
// The VMTests were removed entirely by
// - https://github.com/ethereum/go-ethereum/commit/fb4007bb2208a5b76f76287c03001ef906261691#diff-59830ebc3a4184110566bf1a290d08473dfdcbd492ce498b14cd1a5e2fa2e441
// - https://github.com/ethereum/go-ethereum/pull/23350
//
// // BenchmarkVM runs benchmarks against the JSON VM test suite cases.
// // If the go test -short flag is passed, only the FIRST file in each subdirectory
// // (which describes related groups of tests) will be run.
// func BenchmarkVM(b *testing.B) {
// 	vmt := new(testMatcher)
// 	vmt.skipLoad("^vmSystemOperationsTest.json")
// 	vmt.walkB(b, vmTestDir, func(b *testing.B, name string, test *VMTest) {
// 		b.ReportAllocs()
// 		vmconfig := vm.Config{EVMInterpreter: *testEVM, EWASMInterpreter: *testEWASM}
// 		var statedb = &state.StateDB{}
// 		_, sdb := MakePreState(rawdb.NewMemoryDatabase(), test.json.Pre, false)
// 		*statedb = *sdb
// 		start := time.Now()
// 		b.ResetTimer()
// 		for i := 0; i < b.N; i++ {
// 			test.exec(statedb, vmconfig)
// 			b.StopTimer()
// 			*statedb = *sdb
// 			b.StartTimer()
// 		}
// 		b.StopTimer()
// 		gasRemaining := uint64(0)
// 		if test.json.GasRemaining != nil {
// 			gasRemaining = uint64(*test.json.GasRemaining)
// 		}
// 		gasUsed := test.json.Exec.GasLimit - gasRemaining
// 		elapsed := uint64(time.Since(start))
// 		if elapsed < 1 {
// 			elapsed = 1
// 		}
// 		mgasps := (100 * 1000 * gasUsed * uint64(b.N)) / elapsed
// 		b.ReportMetric(float64(mgasps)/100, "mgas/s")
// 	})
// }

// walkB invokes its runTest argument for all subtests in the given directory.
//
// runTest should be a function of type func(t *testing.T, name string, x <TestType>),
// where TestType is the type of the test contained in test files.
// nolint:unused
// nolint:goimports
func (tm *testMatcher) walkB(b *testing.B, dir string, runTest interface{}) {
	// Walk the directory.
	dirinfo, err := os.Stat(dir)
	if os.IsNotExist(err) || !dirinfo.IsDir() {
		fmt.Fprintf(os.Stderr, "can'b find test files in %s, did you clone the tests submodule?\n", dir)
		b.Fatal("missing test files")
	}
	// shortSelectFiles is used as a lookup under the -short testing option
	// to benchmark only the first of each species of test (instead of all iterations of each JSON VM test).
	shortSelectFiles := make(map[string]bool)
	shouldSkip := func(path string) bool {
		if !testing.Short() {
			return false
		}
		fname := strings.TrimSuffix(filepath.Base(path), ".json")
		onlyWords := regexp.MustCompile(`^[a-zA-Z_]+`)
		matches := onlyWords.FindStringSubmatch(fname)
		match1 := fname
		if len(matches) > 0 {
			match1 = matches[0]
		}
		if _, ok := shortSelectFiles[match1]; ok {
			return true
		}
		b.Logf("select test: %s", match1)
		shortSelectFiles[match1] = true
		return false
	}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		name := filepath.ToSlash(strings.TrimPrefix(path, dir+string(filepath.Separator)))
		if info.IsDir() {
			if _, skipload := tm.findSkip(name + "/"); skipload {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".json" {
			if shouldSkip(path) {
				return nil
			}
			b.Run(name, func(b *testing.B) { tm.runTestFileB(b, path, name, runTest) })
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}
}

//nolint:unused
func (tm *testMatcher) runTestFileB(b *testing.B, path, name string, runTest interface{}) {
	if r, _ := tm.findSkip(name); r != "" {
		b.Skip(r)
	}
	if tm.runonlylistpat != nil {
		if !tm.runonlylistpat.MatchString(name) {
			b.Skip("Skipped by whitelist")
		}
	}

	// Load the file as map[string]<testType>.
	m := makeMapFromTestFuncB(runTest)
	if err := readJSONFile(path, m.Addr().Interface()); err != nil {
		b.Fatal(err)
	}

	// Run all tests from the map. Don't wrap in a subtest if there is only one test in the file.
	keys := sortedMapKeys(m)
	if len(keys) == 1 {
		runTestFuncB(runTest, b, name, m, keys[0])
	} else {
		for _, key := range keys {
			name := name + "/" + key
			b.Run(key, func(b *testing.B) {
				if r, _ := tm.findSkip(name); r != "" {
					b.Skip(r)
				}
				runTestFuncB(runTest, b, name, m, key)
			})
		}
	}
}

//nolint:unused
func makeMapFromTestFuncB(f interface{}) reflect.Value {
	stringT := reflect.TypeOf("")
	testingT := reflect.TypeOf((*testing.B)(nil))
	ftyp := reflect.TypeOf(f)
	if ftyp.Kind() != reflect.Func || ftyp.NumIn() != 3 || ftyp.NumOut() != 0 || ftyp.In(0) != testingT || ftyp.In(1) != stringT {
		panic(fmt.Sprintf("bad test function type: want func(*testing.B, string, <TestType>), have %s", ftyp))
	}
	testType := ftyp.In(2)
	mp := reflect.New(reflect.MapOf(stringT, testType))
	return mp.Elem()
}

//nolint:unused
func runTestFuncB(runTest interface{}, b *testing.B, name string, m reflect.Value, key string) {
	reflect.ValueOf(runTest).Call([]reflect.Value{
		reflect.ValueOf(b),
		reflect.ValueOf(name),
		m.MapIndex(reflect.ValueOf(key)),
	})
}
