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

func TestTraverse(t *testing.T) {
	cases := []struct {
		doc              string
		rawSchema        string
		nodeTestMutation func(s *spec.Schema) *spec.Schema
		onNode           func(s *spec.Schema) error
		onNodeCallWantN  int
	}{
		{
			doc:             `empty schema`,
			rawSchema:       `{}`,
			onNodeCallWantN: 1,
		},

		{
			doc: `schema with prop`,
			rawSchema: `{
				"type": "object",
				"properties": {
					"foo": {}
				}
			}`,
			onNodeCallWantN: 2,
		},

		{
			doc: `simplest cyclical schema, literal`,
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
			nodeTestMutation: nil,
			onNode:           nil,
			onNodeCallWantN:  3,
		},

		{
			doc: `simplest cyclical schema, programmatic`,
			rawSchema: `{
			"type": "object",
			"properties": {
				"foo": {}
			}
			}`,
			nodeTestMutation: func(s *spec.Schema) *spec.Schema {
				// Test a one-deep cycle
				s.Properties["foo"] = copySchema(s)
				return s
			},
			onNodeCallWantN: 3,
		},

		{
			doc: "chained cycles",
			rawSchema: `{
			   "title": "1",
			   "type": "object",
			   "properties": {
				 "foo": {
				   "title": "2",
				   "items": [
					 {
					   "title": "3",
					   "type": "array",
					   "items": { "title": "4" }
					 }
				   ]
				 }
			   }
			  }`,
			nodeTestMutation: func(s *spec.Schema) *spec.Schema {
				*s.Properties["foo"].Items.Schemas[0].Items.Schema = copySchema(s)
				return s
			},
			onNode:          nil,
			onNodeCallWantN: 3,
		},
	}

	for i, c := range cases {
		a := NewAnalysisT()
		sch := mustReadSchema(c.rawSchema)

		fmt.Printf("%d: %s %s\n", i, c.doc, mustWriteJSON(sch))

		// Run programmatic test-schema mutation, if any.
		if c.nodeTestMutation != nil {
			sch.AsWritable()
			c.nodeTestMutation(sch)
		}

		// n is the mutator fn (ie onNodeCallbackWrapper) call counter.
		n := 0

		testController := NewController(t)
		mockMutator := NewMockMutator(testController)
		mockMutator.EXPECT().OnSchema(Any()).Return(nil).Times(c.onNodeCallWantN)

		// Wrap the node call fn for call count, and to handle nil check.
		a.Traverse(sch, func(s *spec.Schema) error {
			n++
			fmt.Printf("%s%straverse_n=%d cb_n=%d schema=%s\n", strings.Repeat(".", a.recurseIter), strings.Repeat(" ", n), a.recurseIter, n, mustWriteJSON(s))
			if c.onNode != nil {
				c.onNode(s)
			}
			return mockMutator.OnSchema(s)
		})

		testController.Finish()

		//if n != c.onNodeCallWantN {
		//	t.Errorf("fail, testcase=%d \"%s\" got=%d want=%d ,schema=%s", i, c.doc, n, c.onNodeCallWantN, mustWriteJSONIndent(sch))
		//}

		fmt.Println()
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
