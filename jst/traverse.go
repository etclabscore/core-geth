package jst

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
)

type AnalysisT struct {
	OpenMetaDescription string
	TraverseOptions
	schemaTitles map[string]string

	recurseIter   int
	recursorStack []*spec.Schema
	mutatedStack  []*spec.Schema

	/*
		@BelfordZ could modify 'prePostMap' to just postArray,
		and have isCycle just be "findSchema", returning the mutated schema if any.
		Look up orig--mutated by index/uniquetitle.
	*/
}

type TraverseOptions struct {
	ExpandAtNode bool
	UniqueOnly   bool
}

func NewAnalysisT() *AnalysisT {
	return &AnalysisT{
		OpenMetaDescription: "Analysisiser",
		schemaTitles:        make(map[string]string),
		recurseIter:         0,
		recursorStack:       []*spec.Schema{},
		mutatedStack:        []*spec.Schema{},
	}
}

func mustReadSchema(jsonStr string) *spec.Schema {
	s := &spec.Schema{}
	err := json.Unmarshal([]byte(jsonStr), &s)
	if err != nil {
		panic(fmt.Sprintf("read schema error: %v", err))
	}
	return s
}

func mustWriteJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

func mustWriteJSONIndent(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

func (a *AnalysisT) SchemaAsReferenceSchema(sch spec.Schema) (refSchema spec.Schema, err error) {
	b, _ := json.Marshal(sch)
	titleKey, ok := a.schemaTitles[string(b)]
	if !ok {
		bb, _ := json.Marshal(sch)
		return refSchema, fmt.Errorf("schema not available as reference: %s @ %v", string(b), string(bb))
	}
	refSchema.Ref = spec.Ref{
		Ref: jsonreference.MustCreateRef("#/components/schemas/" + titleKey),
	}
	return
}

func (a *AnalysisT) RegisterSchema(sch spec.Schema, titleKeyer func(schema spec.Schema) string) {
	b, _ := json.Marshal(sch)
	a.schemaTitles[string(b)] = titleKeyer(sch)
}

func (a *AnalysisT) SchemaFromRef(psch spec.Schema, ref spec.Ref) (schema spec.Schema, err error) {
	v, _, err := ref.GetPointer().Get(psch)
	if err != nil {
		return
	}
	return v.(spec.Schema), nil
}

func schemasAreEquivalent(s1, s2 *spec.Schema) bool {
	spec.ExpandSchema(s1, nil, nil)
	spec.ExpandSchema(s2, nil, nil)
	return reflect.DeepEqual(s1, s2)
}

func (a *AnalysisT) seen(sch *spec.Schema) bool {
	for i := range a.recursorStack {
		//if mustWriteJSON(a.recursorStack[i]) == mustWriteJSON(sch) {
		if reflect.DeepEqual(a.recursorStack[i], sch) {
			return true
		}
	}
	return false
}

func (a *AnalysisT) setHistory(sl []*spec.Schema, item *spec.Schema, index int) {
	if len(sl) == 0 || sl == nil {
		sl = []*spec.Schema{}
	}
	if len(sl) - 1 >= index {
		*sl[index] = *item
		return
	}
	for len(sl) -1 < index {
		sl = append(sl, nil)
	}
	sl[index] = item
}

// analysisOnNode runs a callback function on each leaf of a the JSON schema tree.
// It will return the first error it encounters.
func (a *AnalysisT) Traverse(sch *spec.Schema, onNode func(node *spec.Schema) error) error {

	a.recurseIter++
	iter := a.recurseIter

	if sch == nil {
		return errors.New("traverse called on nil schema")
	}

	sch.AsWritable()

	if a.ExpandAtNode {
		spec.ExpandSchema(sch, sch, nil)
	}

	// Keep a pristine copy of the value on the recurse stack.
	// The incoming pointer value will be mutated.
	cp := &spec.Schema{}
	*cp = *sch

	// If depth-first, a duplicity check should happen immediately before the onnode function
	// is called.
	//	If the schema (prior to mutation) is a duplicate of any other known
	// 	unmutated schema, then if the matching schema was mutated, it should be reverted.
	//  Otherwise, the schema is unique, and can be mutated, storing it's mutated ref
	//  in the history (in case a later duplicate is found).

	// If breadth-first, a duplicity check should happen immediately.
	//  If a duplicate schema is found (comparing originals), then the duplicate should be
	//  reverted to it's original value. The mutator should not be called.

	// ---

	// If unique only is set, and the recursion iter == 0 (at root schema)
	//  then the Traverse function should call itself using a unique schema registry closure.
	//  Once unique schemas have been collected

	if a.UniqueOnly {
		if a.seen(sch) {
			a.setHistory(a.mutatedStack, cp, iter)
			return nil
		}
	}

	for len(a.recursorStack)-1 < iter {
		a.recursorStack = append(a.recursorStack, nil)
	}
	a.recursorStack[iter] = cp

	rec := func(s *spec.Schema, fn func(n *spec.Schema) error) error {
		return a.Traverse(s, fn)
	}

	// jsonschema slices.
	for i := 0; i < len(sch.AnyOf); i++ {
		rec(&sch.AnyOf[i], onNode)
	}
	for i := 0; i < len(sch.AllOf); i++ {
		rec(&sch.AllOf[i], onNode)
	}
	for i := 0; i < len(sch.OneOf); i++ {
		rec(&sch.OneOf[i], onNode)
	}

	// jsonschemama maps
	for k := range sch.Properties {
		v := sch.Properties[k]
		rec(&v, onNode)
		sch.Properties[k] = v
	}
	for k := range sch.PatternProperties {
		v := sch.PatternProperties[k]
		rec(&v, onNode)
		sch.PatternProperties[k] = v
	}

	// jsonschema special type
	if sch.Items == nil {
		return onNode(sch)
	}
	if sch.Items.Len() > 1 {
		for i := range sch.Items.Schemas {
			rec(&sch.Items.Schemas[i], onNode)
		}
	} else {
		rec(sch.Items.Schema, onNode)
	}

	return onNode(sch)
}

/*
	unique := func() bool {
		if !a.UniqueOnly {
			return true // we don't care, fake uniqueness
		}

		for i := range a.recursorStack {
			for j := range a.recursorStack {
				if i == j {
					continue
				}
				if mustWriteJSON(a.recursorStack[i]) == mustWriteJSON(a.recursorStack[j]) {
					return false
				}
			}
		}
		return true
	}()

	final := func(unique bool, s *spec.Schema) {
		if unique {
			onNode(s)
			for len(a.mutatedStack)-1 < a.recurseIter {
				a.mutatedStack = append(a.mutatedStack, nil)
			}
			a.mutatedStack[a.recurseIter] = s
		} else {
			*s =
		}

		if a.UniqueOnly {
			for i := range a.recursorStack {
				for j := range a.recursorStack {
					if i == j {
						continue
					}
					if mustWriteJSON(a.recursorStack[i]) == mustWriteJSON(a.recursorStack[j]) {
						*a.mutatedStack[i] = *a.recursorStack[i]
						continue
					}
				}
			}
		}
	}()
*/
