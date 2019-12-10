package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	bindAddr        string
	defaultResponse string
)

func init() {
	flag.StringVar(&bindAddr, "bind-address", ":8080", "ip:port where http requests are served")
	flag.StringVar(&defaultResponse, "default-response", "yes\n", "what to respond when recpinged")
	flag.Parse()
}

func main() {
	started := time.Now()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		counter := r.URL.Query().Get("counter")
		if counter != "" {
			fmt.Printf("{\"message\": \"got request with counter: %s\", \"counter\": \"%s\"}\n", counter, counter)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, defaultResponse)
	})
	http.HandleFunc("/for", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fmt.Sprintf("%.0f seconds\n", math.RoundToEven(time.Now().Sub(started).Seconds())))
	})

	fmt.Println("running @", bindAddr)
	go func() {
		log.Fatal((&http.Server{Addr: bindAddr}).ListenAndServe())
	}()

	<-interrupt

	fmt.Println("shutting down")

	(&http.Server{Addr: bindAddr}).Shutdown(context.Background())
}
