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

// analysisOnNode runs a callback function on each leaf of a the JSON schema tree.
// It will return the first error it encounters.
func (a *AnalysisT) Traverse(sch *spec.Schema, onNode func(node *spec.Schema) error) error {

	a.recurseIter++

	if sch == nil {
		return errors.New("traverse called on nil schema")
	}

	sch.AsWritable()

	if a.TraverseOptions.ExpandAtNode {
		spec.ExpandSchema(sch, sch, nil)
	}

	a.recursorStack = append(a.recursorStack, sch)
	defer func() {
		a.mutatedStack = append(a.mutatedStack, sch)
	}()

	rec := func(s *spec.Schema, fn func(n *spec.Schema) error) error {
		if a.TraverseOptions.UniqueOnly {
			if a.seen(s) {
				return nil
			}
		}
		return a.Traverse(s, fn)
	}

	// Slices.
	for i := 0; i < len(sch.AnyOf); i++ {
		it := sch.AnyOf[i]
		rec(&it, onNode)
		sch.AnyOf[i] = it
	}
	for i := 0; i < len(sch.AllOf); i++ {
		it := sch.AllOf[i]
		rec(&it, onNode)
		sch.AllOf[i] = it
	}
	for i := 0; i < len(sch.OneOf); i++ {
		it := sch.OneOf[i]
		rec(&it, onNode)
		sch.OneOf[i] = it
	}

	for k := range sch.Properties {
		v := sch.Properties[k]
		//v.ID = "prop:"+k // PTAL: Is this right?
		rec(&v, onNode)
		sch.Properties[k] = v
	}
	for k := range sch.PatternProperties {
		v := sch.PatternProperties[k]
		//v.Title = k // PTAL: Ditto?
		rec(&v, onNode)
		sch.PatternProperties[k] = v
	}
	if sch.Items == nil {
		return onNode(sch)
	}
	if sch.Items.Len() > 1 {
		for i := range sch.Items.Schemas {
			// PTAL: Is this right, onNode)?
			rec(&sch.Items.Schemas[i], onNode)
		}
	} else {
		rec(sch.Items.Schema, onNode)
	}
	return onNode(sch)
}

//as := AsAnalysedSchema(*sch)
//a.schemaTitles[as.j] = as.j

////a.schemaTitles[mustWriteJSON(sch)] = mustWriteJSON(sch)
//for i, st := range a.recursorStack {
//	if reflect.DeepEqual(st , sch) { // } mustWriteJSON(st) == mustWriteJSON(sch) {
//		//if reflect.DeepEqual(st, s) {
//		// If the stack of mutated schemas is not yet long enough
//		// as this index, then append to it.
//		// There is no way of getting the eventual length ahead of time.
//		for (len(a.mutatedStack))-1 < i {
//			a.mutatedStack = append(a.mutatedStack, nil)
//		}
//		a.mutatedStack[i] = sch
//		fmt.Println("same", mustWriteJSON(st))
//		return nil
//	} else {
//		fmt.Println("notsame", mustWriteJSON(st), mustWriteJSON(sch))
//		//panic("samsambutdiff")
//	}
//	if a.recurseIter > 100 {
//		panic("gotcha")
//	}
//}
////a.recursorStack = append(a.recursorStack, sch)
//if len(a.recursorStack) > 100 {
//	panic("lngstck")
//}
////a.mutatedStack[a.recurseIter] = *sch
