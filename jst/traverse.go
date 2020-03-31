package jst

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
)

type AnalysisT struct {
	OpenMetaDescription string
	schemaTitles        map[string]string

	recurseIter   int
	recursorStack []spec.Schema
	mutatedStack  []spec.Schema

	/*
		@BelfordZ could modify 'prePostMap' to just postArray,
		and have isCycle just be "findSchema", returning the mutated schema if any.
		Look up orig--mutated by index/uniquetitle.
	*/
}

func NewAnalysisT() *AnalysisT {
	return &AnalysisT{
		OpenMetaDescription: "Analysisiser",
		schemaTitles:        make(map[string]string),
		recurseIter:         0,
		recursorStack:       []spec.Schema{},
		mutatedStack:        []spec.Schema{},
	}
}

func mustReadSchema(jsonStr string) spec.Schema {
	s := spec.Schema{}
	json.Unmarshal([]byte(jsonStr), &s)
	return s
}

func mustWriteJSON(v interface{}) string {
	b, _ := json.Marshal(v)
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

//type AnalysedSchema struct {
//	j string
//	spec.Schema
//}
//
//func AsAnalysedSchema(s spec.Schema) AnalysedSchema {
//	b, _ := json.Marshal(s)
//	return AnalysedSchema{j: string(b), Schema: s}
//}
//
//func (a AnalysedSchema) SchemaDeterministicID() string {
//	// We need to guarantee orders.
//	return a.j
//}

//func SchemasAreEquivalent(s1, s2 *spec.Schema) bool {
//	s1.
//}

func (a *AnalysisT) rec(s *spec.Schema, onNode func(node *spec.Schema) *spec.Schema) *spec.Schema {
	if s == nil {
		panic("nil schema")
	}
	loop := -1
	for i, st := range a.recursorStack {
		//if mustWriteJSON(st) == mustWriteJSON(s) {
		if reflect.DeepEqual(a.recursorStack[i], s) {

			fmt.Println("same", mustWriteJSON(st))

			loop = i
			break

		}
		//loop++
		fmt.Println("notsame", mustWriteJSON(st), mustWriteJSON(s))
		if a.recurseIter > 100 {
			panic("gotcha")
		}
	}

	if loop >= 0 {
		return nil
	}

	//// If the stack of mutated schemas is not yet long enough
	//// as this index, then append to it.
	//// There is no way of getting the eventual length ahead of time.
	//for (len(a.mutatedStack))-1 < loop {
	//	a.mutatedStack = append(a.mutatedStack, nil)
	//}
	//a.mutatedStack[loop] = s

	return a.Traverse(s, onNode)
}

func (a *AnalysisT) seen(sch spec.Schema) bool {
	for i := range a.recursorStack {
		if reflect.DeepEqual(a.recursorStack[i], sch) {
			return true
		}
	}
	return false
}

// analysisOnNode runs a callback function on each leaf of a the JSON schema tree.
// It will return the first error it encounters.
func (a *AnalysisT) Traverse(sch *spec.Schema, onNode func(node *spec.Schema) *spec.Schema) *spec.Schema {

	a.recurseIter++

	if sch == nil {
		return nil
	}

	cop := &spec.Schema{}
	*cop = *sch
	if a.seen(*cop) {
		return sch
	}
	a.recursorStack = append(a.recursorStack, *cop)

	// Slices.
	for i := 0; i < len(sch.OneOf); i++ {
		it := sch.OneOf[i]

		sch.OneOf[i] = *a.Traverse(&it, onNode)
	}
	for i := 0; i < len(sch.AnyOf); i++ {
		it := sch.AnyOf[i]

		sch.AnyOf[i] = *a.Traverse(&it, onNode)
	}
	for i := 0; i < len(sch.AllOf); i++ {
		it := sch.AllOf[i]

		sch.AllOf[i] = *a.Traverse(&it, onNode)
	}
	// Maps.

	// FIXME: Handle as "$ref" instead.
	for k := range sch.Definitions {
		v := sch.Definitions[k]
		v.Title = k

		sch.Definitions[k] = *a.Traverse(&v, onNode)
	}

	for k := range sch.Properties {
		v := sch.Properties[k]
		//v.Title = k // PTAL: Is this right?
		fmt.Println("in props", a.recursorStack)

		sch.Properties[k] = *a.Traverse(&v, onNode)
	}
	for k := range sch.PatternProperties {
		v := sch.PatternProperties[k]
		//v.Title = k // PTAL: Ditto?

		sch.PatternProperties[k] = *a.Traverse(&v, onNode)
	}
	if sch.Items == nil {
		//onNode(sch)
		return onNode(sch)
		//return nil
	}
	if sch.Items.Len() > 1 {
		for i := range sch.Items.Schemas {
			sch.Items.Schemas[i] = *a.Traverse(&sch.Items.Schemas[i], onNode) // PTAL: Is this right, onNode)?
		}
	} else {
		*sch.Items.Schema = *a.Traverse(sch.Items.Schema, onNode)
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
