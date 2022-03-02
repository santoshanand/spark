package main

import (
	"net/http"

	"github.com/santoshanand/spark"
	"github.com/santoshanand/spark/middleware"
	// "github.com/tylerb/graceful"
)

func main() {
	// Setup
	e := spark.New()
	e.Use(middleware.Logger())
	e.Get("/", func(c *spark.Context) error {
		return c.String(http.StatusOK, "Working")
	})

	// graceful.ListenAndServe(e.Server(":1323"), 5*time.Second)
	e.Run(":3000")
}
