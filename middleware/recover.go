package middleware

import (
	"fmt"

	"runtime"

	"github.com/santoshanand/spark"
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() spark.MiddlewareFunc {
	return func(h spark.HandlerFunc) spark.HandlerFunc {
		return func(c *spark.Context) error {
			defer func() {
				if err := recover(); err != nil {
					trace := make([]byte, 1<<16)
					n := runtime.Stack(trace, true)
					c.Error(fmt.Errorf("echo => panic recover\n %v\n stack trace %d bytes\n %s",
						err, n, trace[:n]))
				}
			}()
			return h(c)
		}
	}
}
