package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/santoshanand/spark"
)

func TestLogger(t *testing.T) {
	e := spark.New()
	req, _ := http.NewRequest(spark.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := spark.NewContext(req, spark.NewResponse(rec), e)

	// Status 2xx
	h := func(c *spark.Context) error {
		return c.String(http.StatusOK, "test")
	}
	Logger()(h)(c)

	// Status 3xx
	rec = httptest.NewRecorder()
	c = spark.NewContext(req, spark.NewResponse(rec), e)
	h = func(c *spark.Context) error {
		return c.String(http.StatusTemporaryRedirect, "test")
	}
	Logger()(h)(c)

	// Status 4xx
	rec = httptest.NewRecorder()
	c = spark.NewContext(req, spark.NewResponse(rec), e)
	h = func(c *spark.Context) error {
		return c.String(http.StatusNotFound, "test")
	}
	Logger()(h)(c)

	// Status 5xx with empty path
	req, _ = http.NewRequest(spark.GET, "", nil)
	rec = httptest.NewRecorder()
	c = spark.NewContext(req, spark.NewResponse(rec), e)
	h = func(c *spark.Context) error {
		return errors.New("error")
	}
	Logger()(h)(c)
}
