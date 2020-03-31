package jst

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
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
			nodeTestMutation: func(s *spec.Schema) *spec.Schema {
				return s
			},
			onNode: func(s *spec.Schema) error {
				return nil
			},
			onNodeCallWantN: 3,
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

				//PASSING:
				//Programmatically modify test value.
				ps := make(map[string]spec.Schema)
				for k, v := range s.Properties {
					ps[k] = v
				}
				ps["foo"] = *s
				s.WithProperties(ps)
				s.Properties = ps
				return s

				//// FAILING:
				//s.Properties["foo"] = *s
				//return s
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

				// a: "works"
				news := mustReadSchema(mustWriteJSON(s))
				*s.Properties["foo"].Items.Schemas[0].Items.Schema = *news

				// b: overflow
				//*s.Properties["foo"].Items.Schemas[0].Items.Schema = *s

				fmt.Println("@debug/chained-cycles", mustWriteJSON(s))

				return s
			},
			onNode:          nil,
			onNodeCallWantN: 3,
		},
	}

	for i, c := range cases {
		a := NewAnalysisT()
		sch := mustReadSchema(c.rawSchema)

		// Develop:skip
		failing := []int{}
		isAFailure := func() bool {
			for _, f := range failing {
				if i == f {
					return true
				}
			}
			return false
		}
		if isAFailure() {
			continue
		}

		if c.nodeTestMutation != nil {
			sch.AsWritable()
			c.nodeTestMutation(sch)
		}

		if i == 4 {
			fmt.Println("postmutation", mustWriteJSON(sch))
		}

		// n is the mutator fn (ie onNodeCallback) call counter.
		n := 0

		// Wrap the node call fn for call count, and to handle nil check.
		onNodeCallback := func(s *spec.Schema) error {
			n++

			// Set a default callback for test cases not specifying a mutator function.
			if c.onNode == nil {
				c.onNode = func(s *spec.Schema) error {
					fmt.Println("~node mutation:", mustWriteJSON(s))
					return nil
				}
			}
			return c.onNode(s)
		}
		a.Traverse(sch, onNodeCallback)

		if n != c.onNodeCallWantN {
			t.Errorf("fail, testcase=%d \"%s\" got=%d want=%d ,schema=%s", i, c.doc, n, c.onNodeCallWantN, mustWriteJSONIndent(sch))
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
