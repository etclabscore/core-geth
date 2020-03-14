// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"

	goopenrpcT "github.com/gregdhill/go-openrpc/types"
)

func (s *RPCService) Describe() (*goopenrpcT.OpenRPCSpec1, error) {

	for module, list := range s.methods() {
		if module == "rpc" {
			continue
		}

		for _, methodName := range list {
			fullName := strings.Join([]string{module, methodName}, serviceMethodSeparators[0])
			method := s.server.services.services[module].callbacks[methodName]

			// FIXME: Development only.
			if method.isSubscribe {
				continue
			}
			if err := s.doc.RegisterMethod(fullName, method); err != nil {
				return nil, err
			}
		}
	}
	return s.doc.Doc, nil
}

// ---

type OpenRPCDescription struct {
	Doc *goopenrpcT.OpenRPCSpec1
}

func NewOpenRPCDescription(server *Server) *OpenRPCDescription {

	doc := &goopenrpcT.OpenRPCSpec1{
		OpenRPC:      "v1",
		Info:         goopenrpcT.Info{},
		Servers:      []goopenrpcT.Server{},
		Methods:      []goopenrpcT.Method{},
		Components:   goopenrpcT.Components{},
		ExternalDocs: goopenrpcT.ExternalDocs{},
		Objects:      nil,
	}

	return &OpenRPCDescription{Doc: doc}
}

func (d *OpenRPCDescription) RegisterMethod(name string, cb *callback) error {

	cb.makeArgTypes()
	cb.makeRetTypes()

	rtFunc := runtime.FuncForPC(cb.fn.Pointer())
	cbFile, _ := rtFunc.FileLine(rtFunc.Entry())

	tokenset := token.NewFileSet()
	astFile, err := parser.ParseFile(tokenset, cbFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	astFuncDel := getAstFunc(cb, astFile, rtFunc)

	if astFuncDel == nil {
		return fmt.Errorf("nil ast func: method name: %s", name)
	}

	method := makeMethod(name, cb, rtFunc, astFuncDel)

	d.Doc.Methods = append(d.Doc.Methods, method)
	sort.Slice(d.Doc.Methods, func(i, j int) bool {
		return d.Doc.Methods[i].Name < d.Doc.Methods[j].Name
	})

	return nil
}

type argOrRet struct {
	v reflect.Value
}

func makeMethod(name string, cb *callback, rt *runtime.Func, fn *ast.FuncDecl) goopenrpcT.Method {
	file, line := rt.FileLine(rt.Entry())
	m := goopenrpcT.Method{
		Name:    name,
		Tags:    nil,
		Summary: fn.Doc.Text(),
		Description: fmt.Sprintf(`
%s
%s:%d'`, rt.Name(), file, line),
		ExternalDocs:   goopenrpcT.ExternalDocs{},
		Params:         []*goopenrpcT.ContentDescriptor{},
		//Result:         &goopenrpcT.ContentDescriptor{},
		Deprecated:     false,
		Servers:        nil,
		Errors:         nil,
		Links:          nil,
		ParamStructure: "",
		Examples:       nil,
	}

	if fn.Type.Params != nil {
		j := 0
		for _, field := range fn.Type.Params.List {
			if field == nil {
				continue
			}
			if cb.hasCtx && strings.Contains(fmt.Sprintf("%s", field.Type), "context") {
				continue
			}
			if len(field.Names) > 0 {
				for _, ident := range field.Names {
					if ident == nil {
						continue
					}
					if j > len(cb.argTypes) - 1 {
						log.Println(name, cb.argTypes, field.Names, j)
						continue
					}
					cd := makeContentDescriptor(cb.argTypes[j], field, ident)
					j++
					m.Params = append(m.Params, &cd)
				}
			} else {
				cd := makeContentDescriptor(cb.argTypes[j], field, nil)
				j++
				m.Params = append(m.Params, &cd)
			}

		}
	}
	if fn.Type.Results != nil {
		j := 0
		for _, field := range fn.Type.Results.List {
			if field == nil {
				continue
			}
			if strings.Contains(fmt.Sprintf("%s", field.Type), "error") {
				continue
			}
			if len(field.Names) > 0 {
				for _, ident := range field.Names {
					cd := makeContentDescriptor(cb.retTypes[j], field, ident)
					j++
					m.Result = &cd
				}
			} else {
				cd := makeContentDescriptor(cb.retTypes[j], field, nil)
				j++
				m.Result = &cd
			}

		}
	}

	return m
}

func makeContentDescriptor(v reflect.Type, field *ast.Field, ident *ast.Ident) goopenrpcT.ContentDescriptor {
	cd := goopenrpcT.ContentDescriptor{
		//Content: goopenrpcT.Content{
		//	Name:        "",
		//	Summary:     field.Doc.Text(),
		//	Description: field.Comment.Text(),
		//	Required:    false,
		//	Deprecated:  false,
		//	Schema: spec.Schema{
		//		VendorExtensible: spec.VendorExtensible{
		//			Extensions: nil,
		//		},
		//		SchemaProps: spec.SchemaProps{
		//			ID: "",
		//			Ref: spec.Ref{
		//				Ref: jsonreference.Ref{
		//					HasFullURL:      false,
		//					HasURLPathOnly:  false,
		//					HasFragmentOnly: false,
		//					HasFileScheme:   false,
		//					HasFullFilePath: false,
		//				},
		//			},
		//			Schema:           "",
		//			Description:      "",
		//			Type:             []string{schemaType},
		//			Nullable:         nullable,
		//			Format:           "",
		//			Title:            "",
		//			Default:          nil,
		//			Maximum:          nil,
		//			ExclusiveMaximum: false,
		//			Minimum:          nil,
		//			ExclusiveMinimum: false,
		//			MaxLength:        nil,
		//			MinLength:        nil,
		//			Pattern:          "",
		//			MaxItems:         nil,
		//			MinItems:         nil,
		//			UniqueItems:      false,
		//			MultipleOf:       nil,
		//			Enum:             nil,
		//			MaxProperties:    nil,
		//			MinProperties:    nil,
		//			Required:         nil,
		//			Items: &spec.SchemaOrArray{
		//				Schema:  &spec.Schema{},
		//				Schemas: nil,
		//			},
		//			AllOf:      nil,
		//			OneOf:      nil,
		//			AnyOf:      nil,
		//			Not:        &spec.Schema{},
		//			Properties: nil,
		//			AdditionalProperties: &spec.SchemaOrBool{
		//				Allows: false,
		//				Schema: &spec.Schema{},
		//			},
		//			PatternProperties: nil,
		//			Dependencies:      nil,
		//			AdditionalItems: &spec.SchemaOrBool{
		//				Allows: false,
		//				Schema: &spec.Schema{},
		//			},
		//			Definitions: nil,
		//		},
		//		SwaggerSchemaProps: spec.SwaggerSchemaProps{
		//			Discriminator: "",
		//			ReadOnly:      false,
		//			XML: &spec.XMLObject{
		//				Name:      "",
		//				Namespace: "",
		//				Prefix:    "",
		//				Attribute: false,
		//				Wrapped:   false,
		//			},
		//			ExternalDocs: &spec.ExternalDocumentation{
		//				Description: "",
		//				URL:         "",
		//			},
		//			Example: nil,
		//		},
		//		ExtraProps: nil,
		//	},
		//},
	}

	var schemaType string
	switch tt := field.Type.(type) {
	case *ast.SelectorExpr:
		schemaType = fmt.Sprintf("%v.%v", tt.X, tt.Sel)
	case *ast.StarExpr:
		schemaType = fmt.Sprintf("%v", tt.X)
		v = v.Elem()
		cd.Schema.Nullable = true
	default:
		schemaType = v.Name()
	}
	schemaType = fmt.Sprintf("%s:%s", v.PkgPath(), schemaType)

	cd.Name = schemaType
	if ident != nil {
		cd.Name = ident.Name
	}

	cd.Summary = field.Doc.Text()
	cd.Description = field.Comment.Text()
	cd.Schema.Type = []string{schemaType}


	return cd
}

func getAstFunc(cb *callback, astFile *ast.File, rf *runtime.Func) *ast.FuncDecl {

	rfName := runtimeFuncName(rf)
	for _, decl := range astFile.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Name == nil || fn.Name.Name != rfName {
			continue
		}
		//log.Println("getAstFunc", spew.Sdump(cb), spew.Sdump(fn))
		fnRecName := ""
		for _, l := range fn.Recv.List {
			if fnRecName != "" {
				break
			}
			i, ok := l.Type.(*ast.Ident)
			if ok {
				fnRecName = i.Name
				continue
			}
			s, ok := l.Type.(*ast.StarExpr)
			if ok {
				fnRecName = fmt.Sprintf("%v", s.X)
			}
		}
		log.Println("=>", "recvr=", cb.rcvr.String(), "fn=", fnRecName)
		if !strings.Contains(cb.rcvr.String(), fnRecName) {
			continue
		}
		// FIXME: Ensure that this is the one true function.
		return fn
	}
	return nil
}

func runtimeFuncName(rf *runtime.Func) string {
	spl := strings.Split(rf.Name(), ".")
	return spl[len(spl)-1]
}

func (d *OpenRPCDescription) findMethodByName(name string) (ok bool, method goopenrpcT.Method) {
	for _, m := range d.Doc.Methods {
		if m.Name == name {
			return true, m
		}
	}
	return false, goopenrpcT.Method{}
}

func runtimeFuncPackageName(rf *runtime.Func) string {
	re := regexp.MustCompile(`(?im)^(?P<pkgdir>.*/)(?P<pkgbase>[a-zA-Z0-9\-_]*)`)
	match := re.FindStringSubmatch(rf.Name())
	pmap := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(match) {
			pmap[name] = match[i]
		}
	}
	return pmap["pkgdir"] + pmap["pkgbase"]
}
