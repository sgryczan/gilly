package lib

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
)

// HandleRoot handles requests to the root path
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello %q", html.EscapeString(r.URL.Path))
}

//HandleMutate handles requests to the /mutate path
func HandleMutate(w http.ResponseWriter, r *http.Request) {

	// read the body / request
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
	}

	// mutate the request
	mutated, err := Mutate(body, false)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
	}

	// and write it back
	w.WriteHeader(http.StatusOK)
	w.Write(mutated)
}
