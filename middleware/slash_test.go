package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/santoshanand/spark"
	"github.com/stretchr/testify/assert"
)

func TestStripTrailingSlash(t *testing.T) {
	req, _ := http.NewRequest(spark.GET, "/users/", nil)
	rec := httptest.NewRecorder()
	c := spark.NewContext(req, spark.NewResponse(rec), spark.New())
	StripTrailingSlash()(c)
	assert.Equal(t, "/users", c.Request().URL.Path)
}

func TestRedirectToSlash(t *testing.T) {
	req, _ := http.NewRequest(spark.GET, "/users", nil)
	rec := httptest.NewRecorder()
	c := spark.NewContext(req, spark.NewResponse(rec), spark.New())
	RedirectToSlash(RedirectToSlashOptions{Code: http.StatusTemporaryRedirect})(c)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/users/", c.Response().Header().Get("Location"))
}
