package middleware

import (
	"net/http"

	"github.com/santoshanand/spark"
)

type (
	// RedirectToSlashOptions -
	RedirectToSlashOptions struct {
		Code int
	}
)

// StripTrailingSlash returns a middleware which removes trailing slash from request
// path.
func StripTrailingSlash() spark.HandlerFunc {
	return func(c *spark.Context) error {
		p := c.Request().URL.Path
		l := len(p)
		if p[l-1] == '/' {
			c.Request().URL.Path = p[:l-1]
		}
		return nil
	}
}

// RedirectToSlash returns a middleware which redirects requests without trailing
// slash path to trailing slash path.
func RedirectToSlash(opts ...RedirectToSlashOptions) spark.HandlerFunc {
	code := http.StatusMovedPermanently

	if len(opts) > 0 {
		o := opts[0]
		code = o.Code
	}

	return func(c *spark.Context) error {
		p := c.Request().URL.Path
		l := len(p)
		if p[l-1] != '/' {
			c.Redirect(code, p+"/")
		}
		return nil
	}
}
