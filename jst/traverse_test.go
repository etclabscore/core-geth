package jst

import (
	"fmt"
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

type schemaMatcher struct {
	s *spec.Schema
}

func newSchemaMatcherFromJSON(s *spec.Schema) schemaMatcher {
	return schemaMatcher{s: s}
}

func (s schemaMatcher) Matches(s2i interface{}) bool {
	s2, cast := s2i.(*spec.Schema)
	if !cast {
		return false
	}
	return schemasAreEquivalent(s.s, s2)
}

func (s schemaMatcher) String() string {
	return fmt.Sprintf("%v", s.s)
}

func TestAnalysisT_Traverse(t *testing.T) {
	test := func(prop string, sch *spec.Schema, expectCallTotal int) {
		a := NewAnalysisT()
		testController := NewController(t)
		mockMutator := NewMockMutator(testController)

		mutatorCalledN := 0
		fmt.Printf("'%s' mutatorCalledN=%d/want_N=%d, a.recurseIter=%d %s@ schema=\n%s\n",  prop, mutatorCalledN,  expectCallTotal, a.recurseIter, strings.Repeat(" .", a.recurseIter-mutatorCalledN), mustWriteJSONIndent(sch))
		a.Traverse(sch, func(s *spec.Schema) error {
			mutatorCalledN++
			fmt.Printf("'%s' mutatorCalledN=%d, a.recurseIter=%d %s %s\n", prop, mutatorCalledN, a.recurseIter, strings.Repeat(" .", a.recurseIter-mutatorCalledN), mustWriteJSON(s))
			matcher := newSchemaMatcherFromJSON(s)
			mockMutator.EXPECT().OnSchema(matcher).Times(1)
			return mockMutator.OnSchema(s)
		})
		if mutatorCalledN != expectCallTotal {
			t.Errorf("bad mutator call total: '%s' want=%d got=%d", prop, expectCallTotal, mutatorCalledN)
		} else {
			fmt.Printf("'%s' OK: mutatorCalledN %d/%d\n", prop, mutatorCalledN, expectCallTotal)
		}
		fmt.Println()
		testController.Finish()
	}

	t.Run("basic functionality", func(t *testing.T) {
		newBasicSchema := func(prop string, any interface{}) *spec.Schema { // Ternary default set
			var s *spec.Schema
			if any == nil {
				s = mustReadSchema(fmt.Sprintf(`{"%s": [{}, {}]}`, prop))
			} else {
				s = mustReadSchema(fmt.Sprintf(`{"%s": %s}`, prop, mustWriteJSON(any)))
			}
			return s
		}

		for _, s := range []string{"anyOf", "allOf", "oneOf"} {
			test(s, newBasicSchema(s, nil), 2)
		}
		test("traverses items when items is ordered list", newBasicSchema("items", nil), 2)
		test("traverses items when items constrained to single schema", mustReadSchema(`{"items": {"items": {"a": {}, "b": {}}}}`), 3)

		// This test is a good example of behavior.
		// The schema at properties.a is equivalent to the one at properties.b
		//
		test("traverses properties", mustReadSchema(`{"properties": {"a": {}, "b": {}}}`), 2)
	})

	t.Run("cycle detection", func(t *testing.T) {
		test("basic", func() *spec.Schema {

			raw := `{"type": "object", "properties": {"foo": {}}}`
			s := mustReadSchema(raw)
			s.Properties["foo"] = *mustReadSchema(raw)
			return s

		}(), 3)

		test("chained", func() *spec.Schema {

			raw := `{
			  "title": "1-top",
			  "type": "object",
			  "properties": {
				"foo": {
				  "title": "2",
				  "items": [
					{
					  "title": "3",
					  "type": "array",
					  "items": {
						"title": "4-maxdepth"
					  }
					}
				  ]
				}
			  }
			}`

			s := mustReadSchema(raw)
			s.Properties["foo"].Items.Schemas[0].Items.Schema = mustReadSchema(raw)
			return s

		}(), 2)

		test("chained in media res", func() *spec.Schema {
			raw := `{
			  "title": "1",
			  "type": "object",
			  "properties": {
				"foo": {
				  "title": "2",
				  "anyOf": [
					{
					  "title": "3",
					  "type": "array",
					  "items": {
						"title": "4",
						"properties": {
						  "baz": {
							"title": "5"
						  }
						}
					  }
					}
				  ]
				}
			  }
			}`

			s := mustReadSchema(raw)
			s.Properties["foo"].AnyOf[0].Items.Schema.Properties["baz"] = mustReadSchema(raw).Properties["foo"]
			return s
		}(), 8)

		test("chained in media res different branch", func() *spec.Schema {
			raw := `{
  "title": "1",
  "type": "object",
  "properties": {
    "foo": {
      "title": "2",
      "anyOf": [
        {
          "title": "3",
          "type": "array",
          "items": {
            "title": "4",
            "properties": {
              "baz": {
                "title": "5"
              }
            }
          }
        }
      ]
    },
    "bar": {
      "title": "6",
      "type": "object",
      "allOf": [
        {
          "title": "7",
          "type": "object",
          "properties": {
            "baz": {
              "title": "8"
            }
          }
        }
      ]
    }
  }
}`
			s := mustReadSchema(raw)
			s.Properties["foo"].AnyOf[0].Items.Schema.Properties["baz"] = *mustReadSchema(raw)
			s.Properties["bar"].AllOf[0].Properties["baz"] = mustReadSchema(raw).Properties["foo"].AnyOf[0]
			return s
		}(), 16)

		test("multiple cycles", func() *spec.Schema {

			raw := `{
  "title": "1",
  "type": "object",
  "properties": {
    "foo": {
      "title": "2",
      "anyOf": [
        {
          "title": "3",
          "type": "array",
          "items": {
            "title": "4",
            "properties": {
              "baz": {
                "title": "5"
              }
            }
          }
        }
      ]
    },
    "bar": {
      "title": "6",
      "type": "object",
      "allOf": [
        {
          "title": "7",
          "type": "object",
          "properties": {
            "baz": {
              "title": "8"
            }
          }
        }
      ]
    }
  }
}`

			s := mustReadSchema(raw)
			s.Properties["bar"].AllOf[0].Properties["baz"] = mustReadSchema(raw).Properties["foo"].AnyOf[0].Items.Schema.Properties["baz"]
			bar := s.Properties["bar"]
			bar.WithAllOf(*mustReadSchema(raw))
			bar.WithAllOf(mustReadSchema(raw).Properties["foo"].AnyOf[0].Items.Schemas...)
			return s
		}(), 7)
	})

	t.Run("mutated schema refs", func(t *testing.T) {

	})
}

func TestSchemaExpand(t *testing.T) {
	raw := `{
	  "title": "1",
	  "type": "object",
	  "properties": {
		"foo": {
		  "$ref": "#/definitions/foo"
		}
	  },
	  "definitions": {
		"foo": {
		  "title": "myfoo"
		}
	  }
	}`

	s := mustReadSchema(raw)

	fmt.Println("Before", mustWriteJSONIndent(s))

	s.AsWritable()
	err := spec.ExpandSchema(s, s, nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("After", mustWriteJSONIndent(s))
}

/*











































*/