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
	"io"
	"io/ioutil"
	"net"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"

	mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
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

func NewServerWithListener(listener net.Listener) *Server {
	server := &Server{
		idgen:            randomIDGenerator(),
		codecs:           mapset.NewSet(),
		run:              1,
		OpenRPCSchemaRaw: defaultOpenRPCSchemaRaw,
	}

	// Register the default service providing meta information about the RPC service such
	// as the services and methods it offers.
	rpcService := &RPCService{server: server, doc: NewOpenRPCDescription(server)}
	rpcService.doc.Doc.Servers = append(rpcService.doc.Doc.Servers, goopenrpcT.Server{
		Name:        listener.Addr().Network(),
		URL:         listener.Addr().String(),
		Summary:     "",
		Description: params.VersionName+"/v"+params.VersionWithMeta,
		Variables:   nil,
	})
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
	Over []goopenrpcT.Method

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

//func (s *RPCService) Describe() (*goopenrpcT.OpenRPCSpec1, error) {
//	describedDoc, err := s.Describe()
//	if err != nil {
//		return nil, err
//	}
//	return describedDoc, nil
//}

func (s *RPCService) DescribeOpenRPC() (*OpenRPCCheck, error) {
	var err error
	if s.server.OpenRPCSchemaRaw == "" {
		return nil, errOpenRPCDiscoverUnavailable
	}
	referenceDoc := &goopenrpcT.OpenRPCSpec1{
		Servers: []goopenrpcT.Server{},
	}
	err = json.Unmarshal([]byte(s.server.OpenRPCSchemaRaw), referenceDoc)
	if err != nil {
		log.Crit("openrpc json umarshal", "error", err)
	}

	describedDoc, err := s.Describe()
	if err != nil {
		return nil, err
	}

	check := &OpenRPCCheck{
		Over:  []goopenrpcT.Method{},
		Under: []goopenrpcT.Method{}, //[]OpenRPCCheckUnder{},
	}

	referenceSuper := func(reference, target *goopenrpcT.OpenRPCSpec1) (methods []goopenrpcT.Method) {
		referenceLoop:
		for _, r := range reference.Methods {
			for _, t := range target.Methods {
				if r.Name == t.Name {
					continue referenceLoop
				}
			}
			methods = append(methods, r)
		}
		return
	}

	check.Over = referenceSuper(referenceDoc, describedDoc)
	check.Under = referenceSuper(describedDoc, referenceDoc)



//	// Audit documented doc methods vs. actual server availability
//	// This removes methods described in the OpenRPC JSON document
//	// which are not currently exposed on the server's API.
//	// This is done on the fly (as opposed to at servre init or doc setting)
//	// because it's possible that exposed APIs could be modified in proc.
//	docMethodsAvailable := []goopenrpcT.Method{}
//	serverMethodsAvailable := s.methods()
//
//	// Find Over methods.
//	// These are methods described in the document, but which are not available
//	// at the server.
//outer:
//	for _, m := range doc.Methods {
//		// Get the module/method name from the document.
//		methodName := m.Name
//		module, path, err := elementizeMethodName(methodName)
//		if err != nil {
//			return nil, err
//		}
//
//		// Check if the server has this module available.
//		paths, ok := serverMethodsAvailable[module]
//		if !ok {
//			check.Over = append(check.Over, methodName)
//			continue
//		}
//
//		// Check if the server has this module+path(=full method name).
//		for _, pa := range paths {
//			if pa == path {
//				docMethodsAvailable = append(docMethodsAvailable, m)
//				continue outer
//			}
//		}
//
//		// Were not continued over; path was not found.
//		check.Over = append(check.Over, methodName)
//	}


	//
	//copy(check.Under, described.Methods)
	//
	//for i, u := range check.Under {
	//	for _, d := range docMethodsAvailable {
	//		if u.Name == d.Name {
	//			check.Under = append(check.Under[:i], check.Under[i+1:]...)
	//		}
	//	}
	//}


	// Find under methods.
	// These are methods which are available on the server, but not described in the document.
	//for mod, list := range check.Under {
	//	if mod == "rpc" {
	//		continue
	//	}
	//modmethodlistloop:
	//	for _, item := range list {
	//		name := strings.Join([]string{mod, item}, serviceMethodSeparators[0])
	//		for _, m := range doc.Methods {
	//			if m.Name == name {
	//
	//				continue modmethodlistloop
	//			}
	//		}
	//		check.Under = append(check.Under, orpcM)
	//	}
	//}
	//sort.Sort(check.Under)
	sort.Slice(check.Over, func(i, j int) bool {
		return check.Over[i].Name < check.Over[j].Name
	})
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
