## Spark


```
  go get github.com/santoshanand/spark/spark

```
``` go 
package main

import (
	"fmt"
	"net/http"

	"github.com/santoshanand/spark/spark"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	inst := spark.New(spark.Options{Handler: mux})
	if err := inst.Shutdown(); err != nil {
		fmt.Println("error shutdown")
	}
}

```