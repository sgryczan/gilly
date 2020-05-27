package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"bitbucket.ngage.netapp.com/scm/hcit/gilly/lib"
)

var version string
var DefaultRegistry string

func main() {

	DefaultRegistry := os.Getenv("GILLY_REGISTRY")
	if DefaultRegistry == "" {
		DefaultRegistry = "docker.repo.eng.netapp.com"
	}
	if version == "" {
		version = "0.0.0"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", lib.HandleRoot)
	mux.HandleFunc("/mutate", lib.HandleMutate)

	s := &http.Server{
		Addr:           ":8443",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}

	log.Printf("[Gilly]  started v%s", version)
	log.Fatal(s.ListenAndServeTLS("./ssl/gilly.pem", "./ssl/gilly.key"))
}
