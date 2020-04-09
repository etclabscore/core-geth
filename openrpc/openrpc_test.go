package openrpc

//
//func mustMarshalJSON(v interface{}) string {
//	b, _ := json.MarshalIndent(v, "", "    ")
//	return string(b)
//}
//
//
//func TestServer(t *testing.T) {
//	files, err := ioutil.ReadDir("testdata")
//	if err != nil {
//		t.Fatal("where'd my testdata go?")
//	}
//	for _, f := range files {
//		if f.IsDir() || strings.HasPrefix(f.Name(), ".") {
//			continue
//		}
//		path := filepath.Join("testdata", f.Name())
//		name := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
//		t.Run(name, func(t *testing.T) {
//			runTestScript(t, path)
//		})
//	}
//}
//
//func runTestScript(t *testing.T, file string) {
//	server := rpc.NewServer()
//	content, err := ioutil.ReadFile(file)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	clientConn, serverConn := net.Pipe()
//	defer clientConn.Close()
//	go server.ServeCodec(rpc.NewCodec(serverConn), 0)
//	readbuf := bufio.NewReader(clientConn)
//	for _, line := range strings.Split(string(content), "\n") {
//		line = strings.TrimSpace(line)
//		switch {
//		case len(line) == 0 || strings.HasPrefix(line, "//"):
//			// skip comments, blank lines
//			continue
//		case strings.HasPrefix(line, "--> "):
//			t.Log(line)
//			// write to connection
//			clientConn.SetWriteDeadline(time.Now().Add(5 * time.Second))
//			if _, err := io.WriteString(clientConn, line[4:]+"\n"); err != nil {
//				t.Fatalf("write error: %v", err)
//			}
//		case strings.HasPrefix(line, "<-- "):
//			t.Log(line)
//			want := line[4:]
//			// read line from connection and compare text
//			clientConn.SetReadDeadline(time.Now().Add(5 * time.Second))
//			sent, err := readbuf.ReadString('\n')
//			if err != nil {
//				t.Fatalf("read error: %v", err)
//			}
//			sent = strings.TrimRight(sent, "\r\n")
//			if sent != want {
//				t.Errorf("wrong line from server\ngot:  %s\nwant: %s", sent, want)
//			}
//		default:
//			panic("invalid line in test script: " + line)
//		}
//	}
//}
//
//type Pet struct {
//	Name         string
//	Age          int
//	Fluffy, Fast bool
//}
//
//type PetStoreService struct {
//	pets []*Pet
//}
//
//// GetPets returns all the pets the store has.
//func (s *PetStoreService) GetPets() ([]*Pet, error) {
//	// Returns all pets.
//	return s.pets, nil
//}
//
//// AddPet adds a pet to the store.
//func (s *PetStoreService) AddPet(p Pet) error {
//	if s.pets == nil {
//		s.pets = []*Pet{}
//	}
//	s.pets = append(s.pets, &p)
//	return nil
//}
//
//func TestOpenRPCDiscover(t *testing.T) {
//	server := rpc.NewServer()
//	defer server.Stop()
//
//	//rpcService := &rpc.RPCService{server: server, doc: rpc.NewOpenRPCDescription(server)}
//	err := server.RegisterName("eth", &ethapi.PublicBlockChainAPI{})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	store := &PetStoreService{
//		pets: []*Pet{
//			{
//				Name:   "Lindy",
//				Age:    7,
//				Fluffy: true,
//			},
//		},
//	}
//	err = server.RegisterName("store", store)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	opts := &openrpcReflect.DocumentDiscoverOpts{
//		Inline:          false,
//		SchemaMutations: []MutateType{SchemaMutateType_Expand, SchemaMutateType_RemoveDefinitions},
//	}
//	doc := Wrap(server, opts)
//	err = server.RegisterName("rpc", doc)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	listener, err := net.Listen("tcp", "127.0.0.1:0")
//	if err != nil {
//		t.Fatal("can't listen:", err)
//	}
//	defer listener.Close()
//	go server.ServeListener(listener)
//
//	requests := []string{
//		`{"jsonrpc":"2.0","id":1,"method":"rpc_modules"}` + "\n",
//		`{"jsonrpc":"2.0","id":1,"method":"rpc_discover"}` + "\n",
//	}
//
//
//	for _, request := range requests {
//		makeRequest(t, request, listener)
//	}
//
//
//}
//
//const maxReadSize = 1024*1024
//func makeRequest(t *testing.T, request string, listener net.Listener) {
//	deadline := time.Now().Add(10 * time.Second)
//	conn, err := net.Dial("tcp", listener.Addr().String())
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer conn.Close()
//	conn.SetDeadline(deadline)
//	conn.Write([]byte(request))
//	conn.(*net.TCPConn).CloseWrite()
//
//	buf := make([]byte, maxReadSize)
//	n, err := conn.Read(buf)
//	if err != nil {
//		t.Fatal(err)
//	}
//	pretty := make(map[string]interface{})
//	err = json.Unmarshal(buf[:n], &pretty)
//	if err != nil {
//		t.Fatal(err)
//	}
//	bufPretty, err := json.MarshalIndent(pretty, "", "  ")
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(string(bufPretty))
//	fmt.Println()
//}
//
//func testOnNode(node *spec.Schema) error {
//	b, err := json.MarshalIndent(node, "", "    ")
//	if err != nil {
//		return err
//	}
//	fmt.Println(string(b))
//	return nil
//}
//
//func TestAnalysisOnNode(t *testing.T) {
//	schemaJSON := `
//{
//	"type": "object",
//	"properties": {
//		"foo": {}
//	}
//}`
//
//	schema := spec.Schema{}
//	err := json.Unmarshal([]byte(schemaJSON), &schema)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	aa := go_jsonschema_traverse.NewAnalysisT()
//	err = aa.WalkDepthFirst(&schema, testOnNode)
//	if err != nil {
//		t.Error(err)
//	}
//
//	schema.Properties["foo"] = schema
//	err = aa.WalkDepthFirst(&schema, testOnNode)
//	if err != nil {
//		t.Error(err)
//	}
//}
