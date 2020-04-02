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

func assertMockMutatorCalledTimes(n int) func(m *MockMutator, s *spec.Schema) {
	return func(m *MockMutator, s *spec.Schema) {
		m.EXPECT().OnSchema(Any()).Times(n)
	}
}

func TestTraverse(t *testing.T) {
	type mutateExpect struct {
		mutate func(s *spec.Schema) *spec.Schema
		expect func(m *MockMutator, s *spec.Schema)
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
		//
		//{
		//	doc: `properties`,
		//	tests: func () (mes []mutateExpect) {
		//		for _, a := range []string{"anyOf", "allOf", "oneOf"} {
		//			a := a
		//			sch := mustReadSchema(fmt.Sprintf(`{"%s": [{}, {}]}`, a))
		//
		//			mes = append(mes, mutateExpect{
		//				mutate: func(s *spec.Schema) *spec.Schema {
		//					*s = *sch
		//					//*saddr = *s
		//					return s
		//				},
		//				expect:	func(m *MockMutator, s *spec.Schema) {
		//					//m.EXPECT().OnSchema(Any()).Times(2) // PASS
		//					//m.EXPECT().OnSchema(sch).Times(2) // FAIL
		//					//m.EXPECT().OnSchema(*sch).Times(2) // FAIL
		//					//m.EXPECT().OnSchema(saddr).Times(2) // FAIL
		//					//m.EXPECT().OnSchema(saddr).Times(2) // FAIL
		//					//m.EXPECT().OnSchema(Any()).Times(2) // FAIL
		//					m.EXPECT().OnSchema(mock.MatchedBy(func(in interface{}) bool {
		//						//fmt.Println(spew.Sdump(in))
		//						jsoon := mustWriteJSON(in)
		//						fmt.Println("jsoooonn===>>>>", jsoon)
		//						return jsoon  == "{}"
		//					})).MinTimes(1)
		//
		//					//switch a {
		//					//case "anyOf":
		//					//
		//					//
		//					//
		//					//	//m.EXPECT().OnSchema(sch.AnyOf[0]).MinTimes(1)
		//					//	//m.EXPECT().OnSchema(&sch.AnyOf[0]).Times(0) // PASS
		//					//	//m.EXPECT().OnSchema(mock.MatchedBy(func(in interface{}) bool {
		//					//	//	v := in.(spec.Schema)
		//					//	//	fmt.Println("in", spew.Sdump(v))
		//					//	//	fmt.Println("want", spew.Sdump(sch.AnyOf[0]))
		//					//	//	return reflect.DeepEqual(v, spec.Schema{})
		//					//	//	return reflect.DeepEqual(in.(spec.Schema), spec.Schema{})
		//					//	//}).Matches(sch.AnyOf[0])).MinTimes()
		//					//
		//					//	//m.EXPECT().OnSchema(mock.MatchedBy(func(in interface{}) bool {
		//					//	//	v := in.(*spec.Schema)
		//					//	//	return reflect.DeepEqual(v, s)
		//					//	//}).Matches(&sch.AnyOf[0])).MinTimes(1)
		//					//
		//					//	// FUCK THIS
		//					//
		//					//case "oneOf":
		//					//	fmt.Println("!!!!!", mustWriteJSONIndent(sch))
		//					//	//m.EXPECT().OnSchema(sch.Items.Schema).MinTimes(1)
		//					//	//m.EXPECT().OnSchema(sch.OneOf[0]).MinTimes(1)
		//					//	//m.EXPECT().OnSchema(sch.OneOf[0]).Times(2)
		//					//case "allOf":
		//					//	//m.EXPECT().OnSchema(sch.AllOf[0]).Times(2)
		//					//}
		//				},
		//			})
		//		}
		//		return
		//	}(),
		//},

		//{
		//	doc: `anyOf`,
		//	tests: func() (mes []mutateExpect) {
		//		mu := func(s *spec.Schema) *spec.Schema {
		//			*s = *mustReadSchema(fmt.Sprintf(`{"%s": [{}, {}]}`, "anyOf"))
		//			return s
		//		}
		//		mes = append(mes, mutateExpect{
		//			mutate: mu,
		//			expect: func(m *MockMutator, s *spec.Schema) {
		//				m.EXPECT().OnSchema(mock.MatchedBy(func(in interface{}) bool {
		//					jsoon := mustWriteJSON(in)
		//					fmt.Println("jsoooonn===>>>>", jsoon)
		//					return jsoon == "{}"
		//				})).MinTimes(0)
		//			},
		//		})
		//		return
		//	}(),
		//},

		{
			doc: `cyclical schema: basic, literal`,
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
			doc: `cyclical schema: basic, programmatic`,
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
			doc: "cyclical schema: chained, programmatic",
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
			var sch *spec.Schema
			if c.rawSchema != "" {
				sch = mustReadSchema(c.rawSchema)
			} else {
				sch = &spec.Schema{}
			}

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
			k.expect(mockMutator, sch)

			// Wrap the node call fn for call count, and to handle nil check.
			a.Traverse(sch, func(s *spec.Schema) error {
				n++
				fmt.Printf("n=%d/a.recurseIter=%d]%s> schema=\n%s\n", a.recurseIter, n, strings.Repeat("=", a.recurseIter-n), mustWriteJSONIndent(s))
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

		n := 0
		fmt.Printf("'%s' n=%d/a.recurseIter=%d %s@ schema=%s\n", prop, n, a.recurseIter, strings.Repeat(" .", a.recurseIter-n), mustWriteJSONIndent(sch))
		a.Traverse(sch, func(s *spec.Schema) error {
			n++
			fmt.Printf("'%s' n=%d/a.recurseIter=%d %s| schema=%s\n", prop, n, a.recurseIter, strings.Repeat(" .", a.recurseIter-n), mustWriteJSON(s))
			matcher := newSchemaMatcherFromJSON(s)
			mockMutator.EXPECT().OnSchema(matcher).Times(1)
			return mockMutator.OnSchema(s)
		})
		if n != expectCallTotal {
			t.Errorf("bad mutator call total: '%s' want=%d got=%d", prop, expectCallTotal, n)
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
		test("traverses properties", mustReadSchema(`{"properties": {"a": {}, "b": {}}}`), 3)
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

/*
space

































































*/
