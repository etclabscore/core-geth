package tracers

import (
	"testing"
)

// BenchmarkTraceResultsAppend1 compares performance against BenchmarkTraceResultsAppend2,
// comparing the performance of different ways of appending items to slices.
// This is used in PrivateTraceAPI#Block appending results to the traceResults value.
func BenchmarkTraceResultsAppend1(b *testing.B) {
	traceResults := []interface{}{}
	for i := 0; i < 10000; i++ {
		traceResults = append(traceResults, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// results := make([]interface{}, len(traceResults))
		results := []interface{}{}
		for _, it := range traceResults { // nolint:gosimple
			results = append(results, it) // nolint:staticcheck
		}
	}
}

// BenchmarkTraceResultsAppend2 (see comment for BenchmarkTraceResultsAppend1).
func BenchmarkTraceResultsAppend2(b *testing.B) {
	traceResults := []interface{}{}
	for i := 0; i < 10000; i++ {
		traceResults = append(traceResults, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// results := make([]interface{}, len(traceResults))
		results := []interface{}{}
		results = append(results, traceResults...) // nolint:ineffassign,staticcheck
	}
}
