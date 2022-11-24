package spark

import (
	"fmt"
	"net/http"
	"testing"
)

func TestSparkWithHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	inst := New(Options{Handler: mux})
	if err := inst.Shutdown(); err != nil {
		fmt.Println("error shutdown")
	}
}
