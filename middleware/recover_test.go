package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/santoshanand/spark"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	e := spark.New()
	e.SetDebug(true)
	req, _ := http.NewRequest(spark.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := spark.NewContext(req, spark.NewResponse(rec), e)
	h := func(c *spark.Context) error {
		panic("test")
	}
	Recover()(h)(c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "panic recover")
}
