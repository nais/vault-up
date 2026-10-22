package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	secretFile := flag.String("secret-file", "/var/run/secrets/nais.io/vault-up/secret", "path to the injected Vault secret")
	expected := flag.String("expected-secret-value", "expected", "expected Vault secret value")
	flag.Parse()

	if err := checkSecret(*secretFile, *expected); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Print("Vault secret verified")
}

func checkSecret(path, expected string) error {
	secret, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read Vault secret: %w", err)
	}
	if string(secret) != expected {
		return errors.New("Vault secret does not match the expected value")
	}
	return nil
}
