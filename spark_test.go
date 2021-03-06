package spark

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"reflect"
	"strings"

	"errors"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
)

type (
	user struct {
		ID   string `json:"id" xml:"id"`
		Name string `json:"name" xml:"name"`
	}
)

func TestSpark(t *testing.T) {
	e := New()
	req, _ := http.NewRequest(GET, "/", nil)
	rec := httptest.NewRecorder()
	c := NewContext(req, NewResponse(rec), e)

	// Router
	assert.NotNil(t, e.Router())

	// Debug
	e.SetDebug(true)
	assert.True(t, e.Debug())

	// DefaultHTTPErrorHandler
	e.DefaultHTTPErrorHandler(errors.New("error"), c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestSparkIndex(t *testing.T) {
	e := New()
	e.Index("examples/website/public/index.html")
	c, b := request(GET, "/", e)
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)
}

func TestSparkFavicon(t *testing.T) {
	e := New()
	e.Favicon("examples/website/public/favicon.ico")
	c, b := request(GET, "/favicon.ico", e)
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)
}

func TestSparkStatic(t *testing.T) {
	e := New()

	// OK
	e.Static("/scripts", "examples/website/public/scripts")
	c, b := request(GET, "/scripts/main.js", e)
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)

	// No file
	e.Static("/scripts", "examples/website/public/scripts")
	c, _ = request(GET, "/scripts/index.js", e)
	assert.Equal(t, http.StatusNotFound, c)

	// Directory
	e.Static("/scripts", "examples/website/public/scripts")
	c, _ = request(GET, "/scripts", e)
	assert.Equal(t, http.StatusForbidden, c)

	// Directory with index.html
	e.Static("/", "examples/website/public")
	c, r := request(GET, "/", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, true, strings.HasPrefix(r, "<!doctype html>"))

	// Sub-directory with index.html
	c, r = request(GET, "/folder", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "sub directory", r)
}

func TestSparkMiddleware(t *testing.T) {
	e := New()
	buf := new(bytes.Buffer)

	// spark.MiddlewareFunc
	e.Use(MiddlewareFunc(func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			buf.WriteString("a")
			return h(c)
		}
	}))

	// func(spark.HandlerFunc) spark.HandlerFunc
	e.Use(func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			buf.WriteString("b")
			return h(c)
		}
	})

	// spark.HandlerFunc
	e.Use(HandlerFunc(func(c *Context) error {
		buf.WriteString("c")
		return nil
	}))

	// func(*spark.Context) error
	e.Use(func(c *Context) error {
		buf.WriteString("d")
		return nil
	})

	// func(http.Handler) http.Handler
	e.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf.WriteString("e")
			h.ServeHTTP(w, r)
		})
	})

	// http.Handler
	e.Use(http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf.WriteString("f")
	})))

	// http.HandlerFunc
	e.Use(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf.WriteString("g")
	}))

	// func(http.ResponseWriter, *http.Request)
	e.Use(func(w http.ResponseWriter, r *http.Request) {
		buf.WriteString("h")
	})

	// Unknown
	assert.Panics(t, func() {
		e.Use(nil)
	})

	// Route
	e.Get("/", func(c *Context) error {
		return c.String(http.StatusOK, "Hello!")
	})

	c, b := request(GET, "/", e)
	assert.Equal(t, "abcdefgh", buf.String())
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "Hello!", b)

	// Error
	e.Use(func(*Context) error {
		return errors.New("error")
	})
	c, b = request(GET, "/", e)
	assert.Equal(t, http.StatusInternalServerError, c)
}

func TestSparkHandler(t *testing.T) {
	e := New()

	// HandlerFunc
	e.Get("/1", HandlerFunc(func(c *Context) error {
		return c.String(http.StatusOK, "1")
	}))

	// func(*spark.Context) error
	e.Get("/2", func(c *Context) error {
		return c.String(http.StatusOK, "2")
	})

	// http.Handler/http.HandlerFunc
	e.Get("/3", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("3"))
	}))

	// func(http.ResponseWriter, *http.Request)
	e.Get("/4", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("4"))
	})

	for _, p := range []string{"1", "2", "3", "4"} {
		c, b := request(GET, "/"+p, e)
		assert.Equal(t, http.StatusOK, c)
		assert.Equal(t, p, b)
	}

	// Unknown
	assert.Panics(t, func() {
		e.Get("/5", nil)
	})
}

func TestSparkConnect(t *testing.T) {
	e := New()
	testMethod(t, CONNECT, "/", e)
}

func TestSparkDelete(t *testing.T) {
	e := New()
	testMethod(t, DELETE, "/", e)
}

func TestSparkGet(t *testing.T) {
	e := New()
	testMethod(t, GET, "/", e)
}

func TestSparkHead(t *testing.T) {
	e := New()
	testMethod(t, HEAD, "/", e)
}

func TestSparkOptions(t *testing.T) {
	e := New()
	testMethod(t, OPTIONS, "/", e)
}

func TestSparkPatch(t *testing.T) {
	e := New()
	testMethod(t, PATCH, "/", e)
}

func TestSparkPost(t *testing.T) {
	e := New()
	testMethod(t, POST, "/", e)
}

func TestSparkPut(t *testing.T) {
	e := New()
	testMethod(t, PUT, "/", e)
}

func TestSparkTrace(t *testing.T) {
	e := New()
	testMethod(t, TRACE, "/", e)
}

func TestSparkWebSocket(t *testing.T) {
	e := New()
	e.WebSocket("/ws", func(c *Context) error {
		c.socket.Write([]byte("test"))
		return nil
	})
	srv := httptest.NewServer(e)
	defer srv.Close()
	addr := srv.Listener.Addr().String()
	origin := "http://localhost"
	url := fmt.Sprintf("ws://%s/ws", addr)
	ws, err := websocket.Dial(url, "", origin)
	if assert.NoError(t, err) {
		ws.Write([]byte("test"))
		defer ws.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(ws)
		assert.Equal(t, "test", buf.String())
	}
}

func TestSparkURL(t *testing.T) {
	e := New()

	static := func(*Context) error { return nil }
	getUser := func(*Context) error { return nil }
	getFile := func(*Context) error { return nil }

	e.Get("/static/file", static)
	e.Get("/users/:id", getUser)
	g := e.Group("/group")
	g.Get("/users/:uid/files/:fid", getFile)

	assert.Equal(t, "/static/file", e.URL(static))
	assert.Equal(t, "/users/:id", e.URL(getUser))
	assert.Equal(t, "/users/1", e.URL(getUser, "1"))
	assert.Equal(t, "/group/users/1/files/:fid", e.URL(getFile, "1"))
	assert.Equal(t, "/group/users/1/files/1", e.URL(getFile, "1", "1"))
}

func TestSparkRoutes(t *testing.T) {
	e := New()
	h := func(*Context) error { return nil }
	routes := []Route{
		{GET, "/users/:user/events", h},
		{GET, "/users/:user/events/public", h},
		{POST, "/repos/:owner/:repo/git/refs", h},
		{POST, "/repos/:owner/:repo/git/tags", h},
	}
	for _, r := range routes {
		e.add(r.Method, r.Path, h)
	}

	for i, r := range e.Routes() {
		assert.Equal(t, routes[i].Method, r.Method)
		assert.Equal(t, routes[i].Path, r.Path)
	}
}

func TestSparkGroup(t *testing.T) {
	e := New()
	buf := new(bytes.Buffer)
	e.Use(func(*Context) error {
		buf.WriteString("0")
		return nil
	})
	h := func(*Context) error { return nil }

	//--------
	// Routes
	//--------

	e.Get("/users", h)

	// Group
	g1 := e.Group("/group1")
	g1.Use(func(*Context) error {
		buf.WriteString("1")
		return nil
	})
	g1.Get("/", h)

	// Group with no parent middleware
	g2 := e.Group("/group2", func(*Context) error {
		buf.WriteString("2")
		return nil
	})
	g2.Get("/", h)

	// Nested groups
	g3 := e.Group("/group3")
	g4 := g3.Group("/group4")
	g4.Get("/", func(c *Context) error {
		return c.NoContent(http.StatusOK)
	})

	request(GET, "/users", e)
	// println(len(e.middleware))
	assert.Equal(t, "0", buf.String())

	buf.Reset()
	request(GET, "/group1/", e)
	// println(len(g1.spark.middleware))
	assert.Equal(t, "01", buf.String())

	buf.Reset()
	request(GET, "/group2/", e)
	assert.Equal(t, "2", buf.String())

	buf.Reset()
	c, _ := request(GET, "/group3/group4/", e)
	assert.Equal(t, http.StatusOK, c)
}

func TestSparkNotFound(t *testing.T) {
	e := New()
	r, _ := http.NewRequest(GET, "/files", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSparkHTTPError(t *testing.T) {
	m := http.StatusText(http.StatusBadRequest)
	he := NewHTTPError(http.StatusBadRequest, m)
	assert.Equal(t, http.StatusBadRequest, he.Code())
	assert.Equal(t, m, he.Error())
}

func TestSparkServer(t *testing.T) {
	e := New()
	s := e.Server(":1323")
	assert.IsType(t, &http.Server{}, s)
}

func testMethod(t *testing.T, method, path string, e *Spark) {
	m := fmt.Sprintf("%c%s", method[0], strings.ToLower(method[1:]))
	p := reflect.ValueOf(path)
	h := reflect.ValueOf(func(c *Context) error {
		c.String(http.StatusOK, method)
		return nil
	})
	i := interface{}(e)
	reflect.ValueOf(i).MethodByName(m).Call([]reflect.Value{p, h})
	_, body := request(method, path, e)
	if body != method {
		t.Errorf("expected body `%s`, got %s.", method, body)
	}
}

func request(method, path string, e *Spark) (int, string) {
	r, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)
	return w.Code, w.Body.String()
}
