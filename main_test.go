package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckSecret(t *testing.T) {
	const expected = "expected-private-value"
	for _, tt := range []struct {
		name    string
		content string
		wantErr bool
	}{
		{name: "matching", content: expected},
		{name: "wrong value", content: "secret-not-to-be-logged", wantErr: true},
		{name: "empty file", wantErr: true},
		{name: "trailing newline", content: expected + "\n", wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "secret")
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatal(err)
			}
			err := checkSecret(path, expected)
			if (err != nil) != tt.wantErr {
				t.Fatalf("checkSecret() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil && (strings.Contains(err.Error(), "secret-not-to-be-logged") || strings.Contains(err.Error(), expected)) {
				t.Fatalf("error contains a secret value: %v", err)
			}
		})
	}
}

func TestCheckSecretMissingFile(t *testing.T) {
	err := checkSecret(filepath.Join(t.TempDir(), "missing"), "expected")
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("checkSecret() error = %v, want os.ErrNotExist", err)
	}
}

func TestCheckSecretUnreadableFile(t *testing.T) {
	if err := checkSecret(t.TempDir(), "expected"); err == nil {
		t.Fatal("checkSecret() succeeded reading a directory")
	}
}
