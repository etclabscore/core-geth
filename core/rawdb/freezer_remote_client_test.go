package rawdb

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/cmd/ancient-store-mem/lib"
	"github.com/ethereum/go-ethereum/rpc"
)

func newTestServer(t *testing.T) *rpc.Server {
	server := rpc.NewServer()
	mockFreezerServer := lib.NewMemFreezerRemoteServerAPI()
	err := server.RegisterName("freezer", mockFreezerServer)
	if err != nil {
		t.Fatal(err)
	}
	return server
}

func TestClient1(t *testing.T) {
	server := newTestServer(t)
	client := rpc.DialInProc(server)

	frClient := &FreezerRemoteClient{
		client: client,
		quit:   make(chan struct{}),
	}

	ancientTestProgram := func(head *uint64, i int) {

		kinds := []string{FreezerRemoteHashTable, FreezerRemoteHeaderTable,
			FreezerRemoteBodiesTable, FreezerRemoteReceiptTable, FreezerRemoteDifficultyTable}
		testKind := kinds[i%len(kinds)]

		err := frClient.AppendAncient(*head, []byte{uint8(*head)}, []byte{uint8(*head)}, []byte{uint8(*head)}, []byte{uint8(*head)}, []byte{uint8(*head)})
		if err != nil {
			t.Fatalf("append: %v", err)
		}
		*head++

		ok, err := frClient.HasAncient(testKind, *head-1)
		if err != nil || !ok {
			t.Fatalf("has ancient: %v %v", err, ok)
		}

		serverHead, err := frClient.Ancients()
		if err != nil {
			t.Fatalf("ancients: %v", err)
		}
		if serverHead != *head {
			t.Fatalf("mismatch server/local *head: local=%d server=%d", *head, serverHead)
		}

		v, err := frClient.Ancient(testKind, *head-1)
		if err != nil {
			t.Fatalf("ancient: %v", err)
		}
		if !bytes.Equal([]byte{uint8(*head - 1)}, v) {
			t.Fatalf("mismatch store value: want: %x, got: %x", []byte{uint8(*head)}, v)
		}

		got, err := frClient.AncientSize(testKind)
		if err != nil {
			t.Fatalf("ancient size: %v", err)
		}
		if got < 1 {
			t.Fatalf("low ancient size: %d", got)
		}

		if i > 23 && i%23 == 0 {
			err := frClient.Sync()
			if err != nil {
				t.Fatalf("sync: %v", err)
			}
		}
		if i > 42 && i%42 == 0 {
			*head -= 15
			err := frClient.TruncateTail(*head)
			if err != nil {
				t.Fatalf("truncate: %v", err)
			}
		}
	}

	head := uint64(0)
	for i := 0; i < 1000; i++ {
		ancientTestProgram(&head, i)
	}

	n, err := frClient.Ancients()
	if err != nil {
		t.Fatal(err)
	}
	if n != 670 {
		t.Fatalf("got: %d, want: 670", n)
	}
}
