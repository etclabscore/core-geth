package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func TestMigrateFaucetDirectory(t *testing.T) {
	hardToCollideName := fmt.Sprintf("faucet-migration-test-datadir-%d", time.Now().UnixNano())
	tempDir := filepath.Join(os.TempDir(), hardToCollideName)
	defer func() {
		os.RemoveAll(tempDir)
	}()

	faucetDataDir := filepath.Join(tempDir, "mychain")
	oldFaucetNodeDataDir := filepath.Join(faucetDataDir, multiFaucetNodeName)

	if err := os.MkdirAll(oldFaucetNodeDataDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	filepath.Walk(faucetDataDir, func(path string, info os.FileInfo, err error) error {
		t.Logf("%s", path)
		return nil
	})

	if err := migrateFaucetDirectory(faucetDataDir); err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(faucetDataDir, coreFaucetNodeName)
	d, err := os.Stat(expected)
	if err != nil {
		t.Fatal(err)
	}
	if !d.IsDir() {
		t.Fatal("non-directory")
	}
	filepath.Walk(faucetDataDir, func(path string, info os.FileInfo, err error) error {
		t.Logf("%s", path)
		return nil
	})
}

func TestFacebook(t *testing.T) {
	// TODO: Remove facebook auth or implement facebook api, which seems to require an API key
	t.Skipf("The facebook access is flaky, needs to be reimplemented or removed")
	for _, tt := range []struct {
		url  string
		want common.Address
	}{
		{
			"https://www.facebook.com/fooz.gazonk/posts/2837228539847129",
			common.HexToAddress("0xDeadDeaDDeaDbEefbEeFbEEfBeeFBeefBeeFbEEF"),
		},
	} {
		_, _, gotAddress, err := authFacebook(tt.url)
		if err != nil {
			t.Fatal(err)
		}
		if gotAddress != tt.want {
			t.Fatalf("address wrong, have %v want %v", gotAddress, tt.want)
		}
	}
}
