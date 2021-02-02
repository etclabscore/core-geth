package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
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
