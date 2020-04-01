package jst

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-openapi/spec"
	. "github.com/golang/mock/gomock"
	_ "github.com/test-go/testify"
)

func BenchmarkJSONMarshalDeterminism(b *testing.B) {
	v := `{
		"type": "object",
		"properties": {
		"foo": {},
		"bar": {},
		"baz": {}
	}`
	sch := mustReadSchema(v)
	v = mustWriteJSON(sch)
	for i := 0; i < b.N; i++ {
		got := mustWriteJSON(sch)
		if v != got {
			b.Fatal("bad")
		}
		sch = mustReadSchema(got)
	}
}

func noMutate(s *spec.Schema) *spec.Schema {
	return s
}

func assertMockMutatorCalledTimes(n int) func(m *MockMutator) {
	return func(m *MockMutator) {
		m.EXPECT().OnSchema(Any()).Times(n)
	}
}

func TestTraverse(t *testing.T) {
	type mutateExpect struct {
		mutate func(s *spec.Schema) *spec.Schema
		expect func(m *MockMutator)
	}

	cases := []struct {
		doc              string
		rawSchema        string
		nodeTestMutation func(s *spec.Schema) *spec.Schema
		onNode           func(s *spec.Schema) error
		tests            []mutateExpect
	}{
		{
			doc:       `empty schema`,
			rawSchema: `{}`,
			tests: []mutateExpect{
				{
					mutate: noMutate,
					expect: assertMockMutatorCalledTimes(1),
				},
			},
		},

		{
			doc: `schema with prop`,
			rawSchema: `{
				"type": "object",
				"properties": {
					"foo": {}
				}
			}`,
			tests: []mutateExpect{
				{
					mutate: noMutate,
					expect: assertMockMutatorCalledTimes(2),
				},
			},
		},

		{
			doc: `basic cyclical schema, literal`,
			rawSchema: `{
			"type": "object",
			"properties": {
				"foo": {
					"type": "object",
					"properties": {
						"foo": {}
					}
				}
			}
			}`,
			tests: []mutateExpect{
				{
					mutate: noMutate,
					expect: assertMockMutatorCalledTimes(3),
				},
			},
		},

		{
			doc: `basic cyclical schema, programmatic`,
			rawSchema: `{
			"type": "object",
			"properties": {
				"foo": {}
			}
			}`,
			tests: []mutateExpect{
				{
					mutate: func(s *spec.Schema) *spec.Schema {
						// Test a one-deep cycle
						s.Properties["foo"] = copySchema(s)
						return s
					},
					expect: assertMockMutatorCalledTimes(3),
				},
			},
		},

		{
			doc: "chained cycles",
			rawSchema: `{
			   "title": "1-top",
			   "type": "object",
			   "properties": {
				 "foo": {
				   "title": "2",
				   "items": [
					 {
					   "title": "3",
					   "type": "array",
					   "items": { "title": "4-maxdepth" }
					 }
				   ]
				 }
			   }
			  }`,
			tests: []mutateExpect{
				{
					mutate: func(s *spec.Schema) *spec.Schema {
						// Test a one-deep cycle
						s.Properties["foo"] = copySchema(s)
						return s
					},
					expect: assertMockMutatorCalledTimes(3),
				},
			},
		},
	}

	for i, c := range cases {
		for j, k := range c.tests {
			a := NewAnalysisT()
			sch := mustReadSchema(c.rawSchema)

			// Run programmatic test-schema mutation, if any.
			if k.mutate != nil {
				sch.AsWritable()
				k.mutate(sch)
			}
			fmt.Printf("* %d/%d: %s\nschema=%s\n", i, j, c.doc, mustWriteJSONIndent(sch))
			fmt.Println()

			// n is the mutator fn (ie onNodeCallbackWrapper) call counter.
			n := 0

			testController := NewController(t)
			mockMutator := NewMockMutator(testController)
			k.expect(mockMutator)

			// Wrap the node call fn for call count, and to handle nil check.
			a.Traverse(sch, func(s *spec.Schema) error {
				n++
				fmt.Printf("a.recurseIter=%d/n=%d]%s> schema=\n%s\n", a.recurseIter, n, strings.Repeat("=", a.recurseIter-n), mustWriteJSONIndent(s))
				if c.onNode != nil {
					c.onNode(s)
				}
				return mockMutator.OnSchema(s)
			})

			testController.Finish()

			fmt.Println()
			fmt.Println()
		}
	}
}

func TestSchemaEq(t *testing.T) {
	orig := mustReadSchema(`
		{
			"type": "object",
			"properties": {
				"foo": {}
			}
		}`)

	list := []*spec.Schema{}
	list = append(list, orig)

	seen1 := false
	for _, l := range list {
		if reflect.DeepEqual(l, &orig) {
			seen1 = true
			break
			//return
		}
	}
	if !seen1 {
		t.Fatal("not seen 1")
	}

	cop := &spec.Schema{}
	*cop = *orig

	seen2 := false
	for _, l := range list {
		if reflect.DeepEqual(l, cop) {
			seen2 = true
			break
		}
	}
	if !seen2 {
		t.Fatal("not seen 2")
	}
}

// GOTCHA: This panics at JSON marshaling.
//func TestMarshalOverflow(t *testing.T) {
//	str := `{
//			"type": "object",
//			"properties": {
//				"foo": {}
//			}
//			}`
//	sch := mustReadSchema(str)
//	sch.Properties["foo"] = *sch
//
//	output := mustWriteJSONIndent(sch)
//	t.Log(output)
//}
