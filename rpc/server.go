// Copyright 2015 The go-ethereum Authors
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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync/atomic"

	"github.com/aws/aws-sdk-go/private/util"
	"github.com/davecgh/go-spew/spew"
	mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-openapi/spec"
	goopenrpcT "github.com/gregdhill/go-openrpc/types"
)

const MetadataApi = "rpc"

var (
	// defaultOpenRPCSchemaRaw can be used to establish a default (package-wide) OpenRPC schema from raw JSON.
	// Methods will be cross referenced with actual registered method names in order to serve
	// only server-enabled methods, enabling user and on-the-fly server endpoint availability configuration.
	defaultOpenRPCSchemaRaw string

	errOpenRPCDiscoverUnavailable   = errors.New("openrpc discover data unavailable")
	errOpenRPCDiscoverSchemaInvalid = errors.New("openrpc discover data invalid")
)

// CodecOption specifies which type of messages a codec supports.
//
// Deprecated: this option is no longer honored by Server.
type CodecOption int

const (
	// OptionMethodInvocation is an indication that the codec supports RPC method calls
	OptionMethodInvocation CodecOption = 1 << iota

	// OptionSubscriptions is an indication that the codec suports RPC notifications
	OptionSubscriptions = 1 << iota // support pub sub
)

// Server is an RPC server.
type Server struct {
	services         serviceRegistry
	idgen            func() ID
	run              int32
	codecs           mapset.Set
	OpenRPCSchemaRaw string
}

// NewServer creates a new server instance with no registered handlers.
func NewServer() *Server {
	server := &Server{
		idgen:            randomIDGenerator(),
		codecs:           mapset.NewSet(),
		run:              1,
		OpenRPCSchemaRaw: defaultOpenRPCSchemaRaw,
	}

	// Register the default service providing meta information about the RPC service such
	// as the services and methods it offers.
	rpcService := &RPCService{server: server, doc: NewOpenRPCDescription(server)}
	server.RegisterName(MetadataApi, rpcService)
	return server
}

func validateOpenRPCSchemaRaw(schemaJSON string) error {
	if schemaJSON == "" {
		return errOpenRPCDiscoverSchemaInvalid
	}
	var schema goopenrpcT.OpenRPCSpec1
	if err := json.Unmarshal([]byte(schemaJSON), &schema); err != nil {
		return fmt.Errorf("%v: %v", errOpenRPCDiscoverSchemaInvalid, err)
	}
	return nil
}

// SetDefaultOpenRPCSchemaRaw validates and sets the package-wide OpenRPC schema data.
func SetDefaultOpenRPCSchemaRaw(schemaJSON string) error {
	if err := validateOpenRPCSchemaRaw(schemaJSON); err != nil {
		return err
	}
	defaultOpenRPCSchemaRaw = schemaJSON
	return nil
}

// SetOpenRPCSchemaRaw validates and sets the raw OpenRPC schema data for a server.
func (s *Server) SetOpenRPCSchemaRaw(schemaJSON string) error {
	if err := validateOpenRPCSchemaRaw(schemaJSON); err != nil {
		return err
	}
	s.OpenRPCSchemaRaw = schemaJSON
	return nil
}

// RegisterName creates a service for the given receiver type under the given name. When no
// methods on the given receiver match the criteria to be either a RPC method or a
// subscription an error is returned. Otherwise a new service is created and added to the
// service collection this server provides to clients.
func (s *Server) RegisterName(name string, receiver interface{}) error {
	return s.services.registerName(name, receiver)
}

// ServeCodec reads incoming requests from codec, calls the appropriate callback and writes
// the response back using the given codec. It will block until the codec is closed or the
// server is stopped. In either case the codec is closed.
//
// Note that codec options are no longer supported.
func (s *Server) ServeCodec(codec ServerCodec, options CodecOption) {
	defer codec.close()

	// Don't serve if server is stopped.
	if atomic.LoadInt32(&s.run) == 0 {
		return
	}

	// Add the codec to the set so it can be closed by Stop.
	s.codecs.Add(codec)
	defer s.codecs.Remove(codec)

	c := initClient(codec, s.idgen, &s.services)
	<-codec.closed()
	c.Close()
}

// serveSingleRequest reads and processes a single RPC request from the given codec. This
// is used to serve HTTP connections. Subscriptions and reverse calls are not allowed in
// this mode.
func (s *Server) serveSingleRequest(ctx context.Context, codec ServerCodec) {
	// Don't serve if server is stopped.
	if atomic.LoadInt32(&s.run) == 0 {
		return
	}

	h := newHandler(ctx, codec, s.idgen, &s.services)
	h.allowSubscribe = false
	defer h.close(io.EOF, nil)

	reqs, batch, err := codec.readBatch()
	if err != nil {
		if err != io.EOF {
			codec.writeJSON(ctx, errorMessage(&invalidMessageError{"parse error"}))
		}
		return
	}
	if batch {
		h.handleBatch(reqs)
	} else {
		h.handleMsg(reqs[0])
	}
}

// Stop stops reading new requests, waits for stopPendingRequestTimeout to allow pending
// requests to finish, then closes all codecs which will cancel pending requests and
// subscriptions.
func (s *Server) Stop() {
	if atomic.CompareAndSwapInt32(&s.run, 1, 0) {
		log.Debug("RPC server shutting down")
		s.codecs.Each(func(c interface{}) bool {
			c.(ServerCodec).close()
			return true
		})
	}
}

// RPCService gives meta information about the server.
// e.g. gives information about the loaded modules.
type RPCService struct {
	server *Server
	doc    *OpenRPCDescription
}

// Modules returns the list of RPC services with their version number
func (s *RPCService) Modules() map[string]string {
	s.server.services.mu.Lock()
	defer s.server.services.mu.Unlock()

	modules := make(map[string]string)
	for name := range s.server.services.services {
		modules[name] = "1.0"
	}
	return modules
}

func (s *RPCService) ModuleMethods(mod string) []string {
	s.server.services.mu.Lock()
	defer s.server.services.mu.Unlock()

	list := []string{}

	for name, ser := range s.server.services.services {
		if name == mod {
			for cname := range ser.callbacks {
				list = append(list, cname)
			}
		}
	}
	return list
}

func (s *RPCService) methods() map[string][]string {
	s.server.services.mu.Lock()
	defer s.server.services.mu.Unlock()

	methods := make(map[string][]string)
	for name, ser := range s.server.services.services {
		for s := range ser.callbacks {
			_, ok := methods[name]
			if !ok {
				methods[name] = []string{s}
			} else {
				methods[name] = append(methods[name], s)
			}
		}
	}
	return methods
}

func (s *RPCService) SetOpenRPCDiscoverDocument(documentPath string) error {
	bs, err := ioutil.ReadFile(documentPath)
	if err != nil {
		return err
	}
	doc := string(bs)
	return s.server.SetOpenRPCSchemaRaw(doc)
}

type OpenRPCCheck struct {
	// Over is the methods in the document which are not actually available at the server.
	Over []string

	// Under is the methods on the server not available in the document.
	Under []goopenrpcT.Method // OpenRPCCheckUnderSet
}

type GoOpenRPCMethodSet []goopenrpcT.Method

type OpenRPCCheckUnderSet []OpenRPCCheckUnder

func (o OpenRPCCheckUnderSet) Len() int {
	return len(o)
}

func (o OpenRPCCheckUnderSet) Less(i, j int) bool {
	var si string = o[i].Name
	var sj string = o[j].Name
	var si_lower = strings.ToLower(si)
	var sj_lower = strings.ToLower(sj)
	if si_lower == sj_lower {
		return si < sj
	}
	return si_lower < sj_lower
	return false
}

func (o OpenRPCCheckUnderSet) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

type OpenRPCCheckUnder struct {
	Name                    string `json:"name"`
	Description             string `json:"description"`
	Summary                 string `json:"summary"`
	Fn                      OpenRPCCheckUnderFn
	IsSubscribe, HasContext bool
	ErrPos                  int
	Args                    []OpenRPCCheckUnderArg
}

type OpenRPCCheckUnderFn struct {
	Str         string
	Name        string
	File        string
	Line        int
	Doc         string
	Body        []string
	ParamsList  []string
	ResultsList []string
}

type OpenRPCCheckUnderArg struct {
	Name, Kind string
}

func packageNameFromRuntimePCFuncName(runtimeFuncForPCName string) string {
	re := regexp.MustCompile(`(?im)^(?P<pkgdir>.*/)(?P<pkgbase>[a-zA-Z0-9\-_]*)`)
	match := re.FindStringSubmatch(runtimeFuncForPCName)
	pmap := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(match) {
			pmap[name] = match[i]
		}
	}
	return pmap["pkgdir"] + pmap["pkgbase"]
}

func (s *RPCService) DescribeOpenRPC() (*OpenRPCCheck, error) {
	var err error
	if s.server.OpenRPCSchemaRaw == "" {
		return nil, errOpenRPCDiscoverUnavailable
	}
	doc := &goopenrpcT.OpenRPCSpec1{
		Servers: []goopenrpcT.Server{},
	}
	err = json.Unmarshal([]byte(s.server.OpenRPCSchemaRaw), doc)
	if err != nil {
		log.Crit("openrpc json umarshal", "error", err)
	}

	check := &OpenRPCCheck{
		Over:  []string{},
		Under: []goopenrpcT.Method{}, //[]OpenRPCCheckUnder{},
	}

	// Audit documented doc methods vs. actual server availability
	// This removes methods described in the OpenRPC JSON document
	// which are not currently exposed on the server's API.
	// This is done on the fly (as opposed to at servre init or doc setting)
	// because it's possible that exposed APIs could be modified in proc.
	docMethodsAvailable := []goopenrpcT.Method{}
	serverMethodsAvailable := s.methods()

	// Find Over methods.
	// These are methods described in the document, but which are not available
	// at the server.
outer:
	for _, m := range doc.Methods {
		// Get the module/method name from the document.
		methodName := m.Name
		module, path, err := elementizeMethodName(methodName)
		if err != nil {
			return nil, err
		}

		// Check if the server has this module available.
		paths, ok := serverMethodsAvailable[module]
		if !ok {
			check.Over = append(check.Over, methodName)
			continue
		}

		// Check if the server has this module+path(=full method name).
		for _, pa := range paths {
			if pa == path {
				docMethodsAvailable = append(docMethodsAvailable, m)
				continue outer
			}
		}

		// Were not continued over; path was not found.
		check.Over = append(check.Over, methodName)
	}

	// Find under methods.
	// These are methods which are available on the server, but not described in the document.
	for mod, list := range serverMethodsAvailable {
		if mod == "rpc" {
			continue
		}
	modmethodlistloop:
		for _, item := range list {
			name := strings.Join([]string{mod, item}, serviceMethodSeparators[0])
			for _, m := range doc.Methods {
				if m.Name == name {
					continue modmethodlistloop
				}
			}

			it := s.server.services.services[mod].callbacks[item]
			it.makeArgTypes()
			it.makeRetTypes()

			fnp := runtime.FuncForPC(it.fn.Pointer())

			orpcM := goopenrpcT.Method{
				Name:           name,
				Tags:           nil,
				Summary:        "",
				Description:    "",
				ExternalDocs:   goopenrpcT.ExternalDocs{},
				Params:         nil,
				Result:         nil,
				Deprecated:     false,
				Servers:        nil,
				Errors:         nil,
				Links:          nil,
				ParamStructure: "",
				Examples:       nil,
			}

			fnFile, fnLine := fnp.FileLine(fnp.Entry())

			/*
			   {
			     "Name": "admin_importChain",
			     "Fn": {
			       "Str": "\u003cfunc(*eth.PrivateAdminAPI, string) (bool, error) Value\u003e",
			       "Name": "github.com/ethereum/go-ethereum/eth.(*PrivateAdminAPI).ImportChain",
			       "File": "/home/ia/go/src/github.com/ethereum/go-ethereum/eth/api.go",
			       "Line": 219,
			       "Doc": "ImportChain imports a blockchain from a local file.\n",
			       "Body": [],
			       "ParamsList": ["file"]
			     },
			     "IsSubscribe": false,
			     "ErrPos": 1,
			     "Args": [
			       {
			         "Name": "string",
			         "Kind": "string"
			       }
			     ]
			   },
			*/
			fns := OpenRPCCheckUnderFn{
				Str:         it.fn.String(),
				Name:        fnp.Name(),
				File:        fnFile,
				Line:        fnLine,
				Body:        []string{},
				ParamsList:  []string{},
				ResultsList: []string{},
			}

			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, fnFile, nil, parser.ParseComments)
			if err != nil {
				panic(err)
			}

			fp, err := parser.ParseFile(fset, fnFile, nil, parser.PackageClauseOnly)
			if err != nil {
				panic(err)
			}

			pkgName := "N/A"
			if fp.Name != nil {
				pkgName = fp.Name.Name
			}

			for _, decl := range f.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}

				// If this is the function we're looking for.
				spl := strings.Split(fnp.Name(), ".")

				// FIXME: We can't assume that the first matching function Name
				// in the file is the correct one. Best to use file/+pos measurements.
				if fn.Name.Name == spl[len(spl)-1] {
					fns.Doc = fn.Doc.Text()

					orpcM.Summary = fn.Doc.Text()
					orpcM.Description = fmt.Sprintf(`%s
%s:%d
%s`, fnp.Name(), fnFile, fnLine, it.fn.String())

					//for _, l := range fn.Body.List {
					//	l.Pos()
					//}

					// for _, l := range fn.Body.List {

					// }
					// spew.Config.DisablePointerAddresses = true
					// spew.Config.Dis
					if fn.Type.Params != nil {
						j := 0
						for i, p := range fn.Type.Params.List {
							if p == nil {
								continue
							}
							if it.hasCtx && strings.Contains(fmt.Sprintf("%s", p.Type), "context") {
								continue
							}
							//fns.ParamsList = append(fns.ParamsList, spew.Sdump(p))

							if orpcM.Params == nil {
								orpcM.Params = []*goopenrpcT.ContentDescriptor{}
							}
							for _, n := range p.Names {
								//names = append(names, n.Name)

								spew.Config.Indent = "    "
								spew.Config.DisablePointerAddresses = true
								spew.Config.DisableCapacities = true

								fmt.Println("********")
								fmt.Println(mod, item, "i=", i, "j=", j)
								fmt.Println("__it__=> ", spew.Sdump(it))
								fmt.Println()
								fmt.Println("__fn__=> ", spew.Sdump(fn))
								//fmt.Println()
								//fmt.Println("__pkgPath__=>", it.fn.Type().Elem().PkgPath())
								//fmt.Println()
								//fmt.Println("__f.Package__=>", spew.Sdump(f.Package))
								fmt.Println()
								fmt.Println("__packageName__=>", pkgName)
								fmt.Println()
								fmt.Println("__packageNameR__=>", fnp.Name(), "->", packageNameFromRuntimePCFuncName(fnp.Name()))
								fmt.Println()
								fmt.Println("__p.type__=>", spew.Sdump(p.Type))

								//
								switch tt := p.Type.(type) {
								case *ast.SelectorExpr:

									fmt.Println("p.Type(selector_expr).name=", tt.X, tt.Sel)
								case *ast.StarExpr:
									fmt.Println("p.Type(star_expr).name=", tt.X, tt.Star)
								default:
									fmt.Println("p.Type(default).name=", tt)
								}

								pname := n.Name
								//pname := strings.Join(names, "+")
								//resType := fmt.Sprintf("type:%v", p.Type)
								//if p.Tag != nil {
								//	resType = "tag:" + p.Tag.Value
								//}

								// use other types for type
								ts := []string{}
								for _, a := range it.argTypes {
									ts = append(ts, a.String())
								}
								resType := strings.Join(ts, ",") + "@" + fmt.Sprintf("%d", j)
								if len(it.argTypes) > 0 && j <= len(it.argTypes)-1 {
									rt := it.argTypes[j]

									if rtname := rt.Name(); rtname != "" {
										//resType = "name:"+rtname
										resType = rtname
									} else if rtname = rt.String(); rtname != "" {
										//resType = "string:"+rtname
										resType = rtname
									} else if rtname = rt.Kind().String(); rtname != "" {
										//resType = "kind:"+rtname
										resType = rtname
									}
									if strings.HasPrefix(resType, "*") {
										//resType = strings.TrimPrefix(resType, "*")
										//resType = util.Capitalize(resType)
										//resType += "OrNull"
									}
									j++
								}

								tit := fmt.Sprintf("%s_%s:Arg%d", mod, item, i)
								if pname != "" {
									tit = fmt.Sprintf("%s_%s:%s", mod, item, util.Capitalize(pname))
								}

								orpcM.Params = append(orpcM.Params, &goopenrpcT.ContentDescriptor{
									Content: goopenrpcT.Content{
										Name:        pname,
										Summary:     p.Doc.Text(),
										Description: p.Comment.Text(), // p.Tag.Value,
										Required:    false,            // FIXME
										Deprecated:  false,            // FIXME
										Schema: spec.Schema{
											SchemaProps: spec.SchemaProps{
												Title: tit,
												Type:  spec.StringOrArray{resType}, // FIXME
											},
										},
									},
								})
							}

							// for _, n := range p.Names {
							// 	fns.ParamsList = append(fns.ParamsList, n.Name) // n.String()
							// }
						}
					}
					if fn.Type.Results != nil {
						j := 0
						for _, p := range fn.Type.Results.List {
							if p == nil {
								continue
							}
							if strings.Contains(fmt.Sprintf("%v", p.Type), "error") {
								continue
							}
							//fns.ResultsList = append(fns.ResultsList, spew.Sdump(p))

							if len(p.Names) > 0 {
								//for _, n := range p.Names {
								//names = append(names, n.String())
								pname := p.Names[0].Name
								//fmt.Sprintf("%v", p.Type)
								//resType := fmt.Sprintf("type:%v", p.Type)
								//if p.Tag != nil {
								//	resType = "tag:" + p.Tag.Value
								//}
								// use other types for type
								ts := []string{}
								for _, a := range it.retTypes {
									ts = append(ts, a.String())
								}
								resType := strings.Join(ts, "/")
								if len(it.retTypes) > 0 && j <= len(it.retTypes)-1 {
									rt := it.retTypes[j]
									if rtname := rt.Name(); rtname != "" {
										//resType = "name:"+rtname
										resType = rtname
									} else if rtname = rt.String(); rtname != "" {
										//resType = "string:"+rtname
										resType = rtname
									} else if rtname = rt.Kind().String(); rtname != "" {
										//resType = "kind:"+rtname
										resType = rtname
									}
									if strings.HasPrefix(resType, "*") {
										//resType = strings.TrimPrefix(resType, "*")
										////resType = util.Capitalize(resType)
										//resType += "OrNull"
									}
									j++
								}

								tit := fmt.Sprintf("%s_%s:Result", mod, item)
								if pname != "" {
									tit = fmt.Sprintf("%s_%s:%s", mod, item, util.Capitalize(pname))
								}
								orpcM.Result = &goopenrpcT.ContentDescriptor{
									Content: goopenrpcT.Content{
										Name:        pname,
										Summary:     p.Doc.Text(),
										Description: p.Comment.Text(), // p.Tag.Value,
										Required:    false,            // FIXME
										Deprecated:  false,            // FIXME
										Schema: spec.Schema{
											SchemaProps: spec.SchemaProps{
												Title: tit,
												Type:  spec.StringOrArray{resType}, // FIXME
											},
										},
									},
								}

								//}
							} else {
								//names = append(names, n.String())
								pname := fmt.Sprintf("%s", it.retTypes[0].String())
								//resType := fmt.Sprintf("type:%v", p.Type)
								//if p.Tag != nil {
								//	resType = "tag:" + p.Tag.Value
								//}

								// use other types for type
								ts := []string{}
								for _, a := range it.retTypes {
									ts = append(ts, a.String())
								}
								resType := strings.Join(ts, "/")
								if len(it.retTypes) > 0 && j <= len(it.retTypes)-1 {
									rt := it.retTypes[j]
									if rtname := rt.Name(); rtname != "" {
										//resType = "name:"+rtname
										resType = rtname
									} else if rtname = rt.String(); rtname != "" {
										//resType = "string:"+rtname
										resType = rtname
									} else if rtname = rt.Kind().String(); rtname != "" {
										//resType = "kind:"+rtname
										resType = rtname
									}
									if strings.HasPrefix(resType, "*") {
										//resType = strings.TrimPrefix(resType, "*")
										////resType = util.Capitalize(resType)
										//resType += "OrNull"
									}
									j++
								}

								tit := fmt.Sprintf("%s_%s:Result", mod, item)
								//if pname != "" {
								//	tit = fmt.Sprintf("%s%_s", item, util.Capitalize(pname))
								//}
								orpcM.Result = &goopenrpcT.ContentDescriptor{
									Content: goopenrpcT.Content{
										Name:        pname,
										Summary:     p.Doc.Text(),
										Description: p.Comment.Text(), // p.Tag.Value,
										Required:    false,            // FIXME
										Deprecated:  false,            // FIXME
										Schema: spec.Schema{
											SchemaProps: spec.SchemaProps{
												Title: tit,
												Type:  spec.StringOrArray{resType}, // FIXME
											},
										},
									},
								}
							}

							//pname := strings.Join(names, "+")

							// for _, n := range p.Names {
							// 	fns.ParamsList = append(fns.ParamsList, n.Name) // n.String()
							// }
						}
					}
					//if fn.Type.Results != nil {
					//	for _, p := range fn.Type.Results.List {
					//		if p == nil {
					//			continue
					//		}
					//
					//		fns.ResultsList = append(fns.ResultsList, spew.Sdump(p))
					//		// for _, n := range p.Names {
					//		// 	fns.ParamsList = append(fns.ResultsList, n.Name) // n.String()
					//		// }
					//	}
					//}

				} else {
					// fmt.Println("NONMATCH", fn.Name.Name)
				}
			}

			cu := OpenRPCCheckUnder{
				Name:        name,
				Fn:          fns,
				IsSubscribe: it.isSubscribe,
				HasContext:  it.hasCtx,
				ErrPos:      it.errPos,
				Args:        []OpenRPCCheckUnderArg{},
			}

			argTypes := it.argTypes

			for _, a := range argTypes {
				cu.Args = append(cu.Args, OpenRPCCheckUnderArg{
					Name: a.Name(),
					Kind: a.Kind().String(),
				})
			}

			check.Under = append(check.Under, orpcM)
		}
	}
	//sort.Sort(check.Under)
	sort.Slice(check.Under, func(i, j int) bool {
		return check.Under[i].Name < check.Under[j].Name
	})
	return check, nil
}

// Discover returns a configured schema that is audited for actual server availability.
// Only methods that the server makes available are included in the 'methods' array of
// the discover schema. Components are not audited.
func (s *RPCService) Discover() (schema *goopenrpcT.OpenRPCSpec1, err error) {
	if s.server.OpenRPCSchemaRaw == "" {
		return nil, errOpenRPCDiscoverUnavailable
	}
	schema = &goopenrpcT.OpenRPCSpec1{
		Servers: []goopenrpcT.Server{},
	}
	err = json.Unmarshal([]byte(s.server.OpenRPCSchemaRaw), schema)
	if err != nil {
		log.Crit("openrpc json umarshal", "error", err)
	}

	// Audit documented schema methods vs. actual server availability
	// This removes methods described in the OpenRPC JSON schema document
	// which are not currently exposed on the server's API.
	// This is done on the fly (as opposed to at servre init or schema setting)
	// because it's possible that exposed APIs could be modified in proc.
	schemaMethodsAvailable := []goopenrpcT.Method{}
	serverMethodsAvailable := s.methods()

	for _, m := range schema.Methods {
		module, path, err := elementizeMethodName(m.Name)
		if err != nil {
			return nil, err
		}
		paths, ok := serverMethodsAvailable[module]
		if !ok {
			continue
		}

		// the module exists, does the path exist?
		for _, pa := range paths {
			if pa == path {
				schemaMethodsAvailable = append(schemaMethodsAvailable, m)
				break
			}
		}
	}
	schema.Methods = schemaMethodsAvailable
	return
}
