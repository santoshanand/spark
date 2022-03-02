package middleware

import (
	"log"
	"time"

	"github.com/labstack/gommon/color"
	"github.com/santoshanand/spark"
)

// Logger - common logger
func Logger() spark.MiddlewareFunc {
	return func(h spark.HandlerFunc) spark.HandlerFunc {
		return func(c *spark.Context) error {
			start := time.Now()
			if err := h(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			method := c.Request().Method
			path := c.Request().URL.Path
			if path == "" {
				path = "/"
			}
			size := c.Response().Size()

			n := c.Response().Status()
			code := color.Green(n)
			switch {
			case n >= 500:
				code = color.Red(n)
			case n >= 400:
				code = color.Yellow(n)
			case n >= 300:
				code = color.Cyan(n)
			}

			log.Printf("%s %s %s %s %d", method, path, code, stop.Sub(start), size)
			return nil
		}
	}
}
