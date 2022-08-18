package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	bindAddr            string
	expectedSecretValue string
)

const (
	metricTemplate = `
# HELP vault_secret_ok Vault secret is OK 
# TYPE vault_secret_ok gauge
vault_secret_ok %s`
	vaultFilePath = "/var/run/secrets/nais.io/vault-up/secret"
)

func init() {
	flag.StringVar(&bindAddr, "bind-address", ":8080", "ip:port where http requests are served")
	flag.StringVar(&expectedSecretValue, "expected-secret-value", "expected", "expected vault secret value")
	flag.Parse()
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	secret, err := os.ReadFile(vaultFilePath)
	if err != nil {
		fmt.Println(err)
	}

	ok := "0"
	secretString := string(secret)
	if secretString == expectedSecretValue {
		ok = "1"
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metric := fmt.Sprintf(metricTemplate, ok)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, metric)
	})

	fmt.Println("running @", bindAddr)
	go func() {
		log.Fatal((&http.Server{Addr: bindAddr}).ListenAndServe())
	}()

	<-interrupt

	fmt.Println("shutting down")

	(&http.Server{Addr: bindAddr}).Shutdown(context.Background())
}
