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

	// Keep a pristine copy of the value on the recurse stack.
	// The incoming pointer value will be mutated.
	cp := &spec.Schema{}
	*cp = *sch

	a.setHistory(a.recursorStack, cp, iter)

	final := func(s *spec.Schema) error {
		err := onNode(s)
		a.setHistory(a.mutatedStack, s, iter)
		return err
	}

	// jsonschema slices.
	for i := 0; i < len(sch.AnyOf); i++ {
		a.Traverse(&sch.AnyOf[i], onNode)
	}
	for i := 0; i < len(sch.AllOf); i++ {
		a.Traverse(&sch.AllOf[i], onNode)
	}
	for i := 0; i < len(sch.OneOf); i++ {
		a.Traverse(&sch.OneOf[i], onNode)
	}

	// jsonschemama maps
	for k := range sch.Properties {
		v := sch.Properties[k]
		a.Traverse(&v, onNode)
		sch.Properties[k] = v
	}
	for k := range sch.PatternProperties {
		v := sch.PatternProperties[k]
		a.Traverse(&v, onNode)
		sch.PatternProperties[k] = v
	}

	// jsonschema special type
	if sch.Items == nil {
		return final(sch)
	}
	if sch.Items.Len() > 1 {
		for i := range sch.Items.Schemas {
			a.Traverse(&sch.Items.Schemas[i], onNode)
		}
	} else {
		a.Traverse(sch.Items.Schema, onNode)
	}

	return final(sch)
}