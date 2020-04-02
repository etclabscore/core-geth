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

func traverseTest(t *testing.T, prop string, options *TraverseOptions, sch *spec.Schema, expectCallTotal int) {

	runTest := func(a *AnalysisT, prop string, sch *spec.Schema, expectCallTotal int) {
		testController := NewController(t)
		mockMutator := NewMockMutator(testController)

		mutatorCalledN := 0
		fmt.Printf("'%s' unique=%v schema=\n%s\n", prop, a.TraverseOptions.UniqueOnly, mustWriteJSONIndent(sch))
		a.Traverse(sch, func(s *spec.Schema) error {
			mutatorCalledN++
			fmt.Printf("'%s' unique=%v mutatorCalledN=%d, a.recurseIter=%d %s %s\n", prop, a.TraverseOptions.UniqueOnly, mutatorCalledN, a.recurseIter, strings.Repeat(" .", a.recurseIter-mutatorCalledN), mustWriteJSON(s))
			matcher := newSchemaMatcherFromJSON(s)
			mockMutator.EXPECT().OnSchema(matcher).Times(1)
			return mockMutator.OnSchema(s)
		})
		if mutatorCalledN != expectCallTotal {
			t.Errorf("bad mutator call total: '%s' unique=%v want=%d got=%d", prop, a.TraverseOptions.UniqueOnly, expectCallTotal, mutatorCalledN)
		} else {
			fmt.Printf("'%s' unique=%v OK: mutatorCalledN %d/%d\n", prop, a.TraverseOptions.UniqueOnly, mutatorCalledN, expectCallTotal)
		}
		fmt.Println()
		testController.Finish()
	}

	a := NewAnalysisT()
	if options != nil {
		a.TraverseOptions = *options
	}
	prop = t.Name() + ":" + prop
	runTest(a, prop, sch, expectCallTotal)
}

func TestAnalysisT_Traverse(t *testing.T) {

	type traverseTestsOptsExpect map[*int]*TraverseOptions
	traverseUniqueOpts := &TraverseOptions{UniqueOnly: true}
	pint := func(i int) *int {
		return &i
	}

	runTraverseTestWithOptsExpect := func(t *testing.T, prop string, oe traverseTestsOptsExpect, sch *spec.Schema) {
		for k := range oe {
			v := oe[k]
			traverseTest(t, prop, v, sch, *k)
		}
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
			runTraverseTestWithOptsExpect(t, s, traverseTestsOptsExpect{
				pint(3): nil,
				pint(2): traverseUniqueOpts,
			}, newBasicSchema(s, nil))
		}

		runTraverseTestWithOptsExpect(t, "traverses items when items is ordered list", traverseTestsOptsExpect{
			pint(3): nil,
			pint(2): traverseUniqueOpts,
		}, newBasicSchema("items", nil))

		runTraverseTestWithOptsExpect(t, "traverses items when items constrained to single schema", traverseTestsOptsExpect{
			pint(3): nil,
			pint(3): traverseUniqueOpts,
		}, mustReadSchema(`{"items": {"items": {"a": {}, "b": {}}}}`))

		// This test is a good example of behavior.
		// The schema at properties.a is equivalent to the one at properties.b
		// If the Traverser uses schema equivalence to abort traversal, then
		// equivalent schemas at multiple paths will not be called.
		runTraverseTestWithOptsExpect(t, "traverses properties", traverseTestsOptsExpect{
			pint(3): nil,
			pint(2): traverseUniqueOpts,
		}, mustReadSchema(`{"properties": {"a": {}, "b": {}}}`))
	})

	t.Run("cycle detection", func(t *testing.T) {
		runTraverseTestWithOptsExpect(t, "basic", traverseTestsOptsExpect{
			pint(3): nil,
			pint(3): traverseUniqueOpts,
		}, func() *spec.Schema {

			raw := `{"type": "object", "properties": {"foo": {}}}`
			s := mustReadSchema(raw)
			s.Properties["foo"] = *mustReadSchema(raw)
			return s

		}())

		runTraverseTestWithOptsExpect(t, "basic", traverseTestsOptsExpect{
			pint(3): nil,
			pint(3): traverseUniqueOpts,
		}, func() *spec.Schema {

			raw := `{"type": "object", "properties": {"foo": {}}}`
			s := mustReadSchema(raw)
			s.Properties["foo"] = *mustReadSchema(raw)
			return s

		}())

		runTraverseTestWithOptsExpect(t, "chained", traverseTestsOptsExpect{
			pint(2): nil,
			pint(2): traverseUniqueOpts,
		}, func() *spec.Schema {

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

		}())

		runTraverseTestWithOptsExpect(t, "chained in media res", traverseTestsOptsExpect{
			pint(8): nil,
			pint(8): traverseUniqueOpts,
		}, func() *spec.Schema {
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
		}())

		runTraverseTestWithOptsExpect(t, "chained in media res different branch", traverseTestsOptsExpect{
			pint(17): nil,
			pint(14): traverseUniqueOpts,
		}, func() *spec.Schema {
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
		}())

		runTraverseTestWithOptsExpect(t, "multiple cycles", traverseTestsOptsExpect{
			pint(8): nil,
			pint(7): traverseUniqueOpts,
		}, func() *spec.Schema {

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
		}())
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
