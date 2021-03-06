package spark

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	gpath "path"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"encoding/xml"

	"golang.org/x/net/websocket"
)

type (
	// Spark - struct to create object
	Spark struct {
		prefix                  string
		middleware              []MiddlewareFunc
		maxParam                *int
		notFoundHandler         HandlerFunc
		defaultHTTPErrorHandler HTTPErrorHandler
		httpErrorHandler        HTTPErrorHandler
		binder                  BindFunc
		renderer                Renderer
		pool                    sync.Pool
		debug                   bool
		router                  *Router
	}

	// Route - route struct
	Route struct {
		Method  string
		Path    string
		Handler Handler
	}

	// HTTPError - http error struct
	HTTPError struct {
		code    int
		message string
	}

	// Middleware - middleware interface
	Middleware interface{}
	//MiddlewareFunc - middleware function
	MiddlewareFunc func(HandlerFunc) HandlerFunc
	// Handler - handler interface
	Handler interface{}
	// HandlerFunc - handler function
	HandlerFunc func(*Context) error

	// HTTPErrorHandler is a centralized HTTP error handler.
	HTTPErrorHandler func(error, *Context)

	// BindFunc - bind function
	BindFunc func(*http.Request, interface{}) error

	// Renderer is the interface that wraps the Render method.
	//
	// Render renders the HTML template with given name and specified data.
	// It writes the output to w.
	Renderer interface {
		Render(w io.Writer, name string, data interface{}) error
	}
)

const (
	// CONNECT HTTP method
	CONNECT = "CONNECT"
	// DELETE HTTP method
	DELETE = "DELETE"
	// GET HTTP method
	GET = "GET"
	// HEAD HTTP method
	HEAD = "HEAD"
	// OPTIONS HTTP method
	OPTIONS = "OPTIONS"
	// PATCH HTTP method
	PATCH = "PATCH"
	// POST HTTP method
	POST = "POST"
	// PUT HTTP method
	PUT = "PUT"
	// TRACE HTTP method
	TRACE = "TRACE"

	//-------------
	// Media types
	//-------------

	// ApplicationJSON - content type application/json
	ApplicationJSON = "application/json; charset=utf-8"
	// ApplicationXML - application/xml utf8
	ApplicationXML = "application/xml; charset=utf-8"
	// ApplicationForm - application/x-www-form-urlencoded
	ApplicationForm = "application/x-www-form-urlencoded"
	// ApplicationProtobuf - application/protobuf
	ApplicationProtobuf = "application/protobuf"
	// ApplicationMsgpack - application/msgpack
	ApplicationMsgpack = "application/msgpack"
	// TextHTML - text/html; charset=utf-8
	TextHTML = "text/html; charset=utf-8"
	//TextPlain - text/plain utf8
	TextPlain = "text/plain; charset=utf-8"
	//MultipartForm - multipart/form-data
	MultipartForm = "multipart/form-data"

	//---------
	// Headers
	//---------

	// Accept - header accept
	Accept = "Accept"
	//AcceptEncoding -
	AcceptEncoding = "Accept-Encoding"
	//Authorization -
	Authorization = "Authorization"
	//ContentDisposition -
	ContentDisposition = "Content-Disposition"
	//ContentEncoding -
	ContentEncoding = "Content-Encoding"
	//ContentLength -
	ContentLength = "Content-Length"
	//ContentType -
	ContentType = "Content-Type"
	//Location -
	Location = "Location"
	//Upgrade -
	Upgrade = "Upgrade"
	// Vary -
	Vary = "Vary"

	//-----------
	// Protocols
	//-----------

	// WebSocket - websocket
	WebSocket = "websocket"

	indexFile = "index.html"
)

var (
	methods = [...]string{
		CONNECT,
		DELETE,
		GET,
		HEAD,
		OPTIONS,
		PATCH,
		POST,
		PUT,
		TRACE,
	}

	//--------
	// Errors
	//--------

	// ErrUnsupportedMediaType -
	ErrUnsupportedMediaType = errors.New("spark ??? unsupported media type")
	// ErrRendererNotRegistered -
	ErrRendererNotRegistered = errors.New("spark ??? renderer not registered")
)

// New creates an instance of Spark.
func New() (e *Spark) {
	e = &Spark{maxParam: new(int)}
	e.pool.New = func() interface{} {
		return NewContext(nil, new(Response), e)
	}
	e.router = NewRouter(e)

	//----------
	// Defaults
	//----------

	e.notFoundHandler = func(c *Context) error {
		return NewHTTPError(http.StatusNotFound)
	}
	e.defaultHTTPErrorHandler = func(err error, c *Context) {
		code := http.StatusInternalServerError
		msg := http.StatusText(code)
		if he, ok := err.(*HTTPError); ok {
			code = he.code
			msg = he.message
		}
		if e.debug {
			msg = err.Error()
		}
		http.Error(c.response, msg, code)
	}
	e.SetHTTPErrorHandler(e.defaultHTTPErrorHandler)
	e.SetBinder(func(r *http.Request, v interface{}) error {
		ct := r.Header.Get(ContentType)
		err := ErrUnsupportedMediaType
		if strings.HasPrefix(ApplicationJSON, ct) {
			err = json.NewDecoder(r.Body).Decode(v)
		} else if strings.HasPrefix(ApplicationXML, ct) {
			err = xml.NewDecoder(r.Body).Decode(v)
		}
		return err
	})
	return
}

// Router returns router.
func (e *Spark) Router() *Router {
	return e.router
}

// DefaultHTTPErrorHandler invokes the default HTTP error handler.
func (e *Spark) DefaultHTTPErrorHandler(err error, c *Context) {
	e.defaultHTTPErrorHandler(err, c)
}

// SetHTTPErrorHandler registers a custom Spark.HTTPErrorHandler.
func (e *Spark) SetHTTPErrorHandler(h HTTPErrorHandler) {
	e.httpErrorHandler = h
}

// SetBinder registers a custom binder. It's invoked by Context.Bind().
func (e *Spark) SetBinder(b BindFunc) {
	e.binder = b
}

// SetRenderer registers an HTML template renderer. It's invoked by Context.Render().
func (e *Spark) SetRenderer(r Renderer) {
	e.renderer = r
}

// SetDebug sets debug mode.
func (e *Spark) SetDebug(on bool) {
	e.debug = on
}

// Debug returns debug mode.
func (e *Spark) Debug() bool {
	return e.debug
}

// Use adds handler to the middleware chain.
func (e *Spark) Use(m ...Middleware) {
	for _, h := range m {
		e.middleware = append(e.middleware, wrapMiddleware(h))
	}
}

// Connect adds a CONNECT route > handler to the router.
func (e *Spark) Connect(path string, h Handler) {
	e.add(CONNECT, path, h)
}

// Delete adds a DELETE route > handler to the router.
func (e *Spark) Delete(path string, h Handler) {
	e.add(DELETE, path, h)
}

// Get adds a GET route > handler to the router.
func (e *Spark) Get(path string, h Handler) {
	e.add(GET, path, h)
}

// Head adds a HEAD route > handler to the router.
func (e *Spark) Head(path string, h Handler) {
	e.add(HEAD, path, h)
}

// Options adds an OPTIONS route > handler to the router.
func (e *Spark) Options(path string, h Handler) {
	e.add(OPTIONS, path, h)
}

// Patch adds a PATCH route > handler to the router.
func (e *Spark) Patch(path string, h Handler) {
	e.add(PATCH, path, h)
}

// Post adds a POST route > handler to the router.
func (e *Spark) Post(path string, h Handler) {
	e.add(POST, path, h)
}

// Put adds a PUT route > handler to the router.
func (e *Spark) Put(path string, h Handler) {
	e.add(PUT, path, h)
}

// Trace adds a TRACE route > handler to the router.
func (e *Spark) Trace(path string, h Handler) {
	e.add(TRACE, path, h)
}

// WebSocket adds a WebSocket route > handler to the router.
func (e *Spark) WebSocket(path string, h HandlerFunc) {
	e.Get(path, func(c *Context) (err error) {
		wss := websocket.Server{
			Handler: func(ws *websocket.Conn) {
				c.socket = ws
				c.response.status = http.StatusSwitchingProtocols
				err = h(c)
			},
		}
		wss.ServeHTTP(c.response, c.request)
		return err
	})
}

func (e *Spark) add(method, path string, h Handler) {
	path = e.prefix + path
	e.router.Add(method, path, wrapHandler(h), e)
	r := Route{
		Method:  method,
		Path:    path,
		Handler: runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name(),
	}
	e.router.routes = append(e.router.routes, r)
}

// Index serves index file.
func (e *Spark) Index(file string) {
	e.ServeFile("/", file)
}

// Favicon serves the default favicon - GET /favicon.ico.
func (e *Spark) Favicon(file string) {
	e.ServeFile("/favicon.ico", file)
}

// Static serves static files from a directory. It's an alias for `Sark.ServeDir`
func (e *Spark) Static(path, dir string) {
	e.ServeDir(path, dir)
}

// ServeDir serves files from a directory.
func (e *Spark) ServeDir(path, dir string) {
	e.Get(path+"*", func(c *Context) error {
		return serveFile(dir, c.P(0), c) // Param `_name`
	})
}

// ServeFile serves a file.
func (e *Spark) ServeFile(path, file string) {
	e.Get(path, func(c *Context) error {
		dir, file := gpath.Split(file)
		return serveFile(dir, file, c)
	})
}

func serveFile(dir, file string, c *Context) error {
	fs := http.Dir(dir)
	f, err := fs.Open(file)
	if err != nil {
		return NewHTTPError(http.StatusNotFound)
	}

	fi, _ := f.Stat()
	if fi.IsDir() {
		file = gpath.Join(file, indexFile)
		f, err = fs.Open(file)
		if err != nil {
			return NewHTTPError(http.StatusForbidden)
		}
		fi, _ = f.Stat()
	}

	http.ServeContent(c.response, c.request, fi.Name(), fi.ModTime(), f)
	return nil
}

// Group creates a new sub router with prefix. It inherits all properties from
// the parent. Passing middleware overrides parent middleware.
func (e *Spark) Group(prefix string, m ...Middleware) *Group {
	g := &Group{*e}
	g.spark.prefix += prefix
	if len(m) > 0 {
		g.spark.middleware = nil
		g.Use(m...)
	}
	return g
}

// URI generates a URI from handler.
func (e *Spark) URI(h Handler, params ...interface{}) string {
	uri := new(bytes.Buffer)
	pl := len(params)
	n := 0
	hn := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	for _, r := range e.router.routes {
		if r.Handler == hn {
			for i, l := 0, len(r.Path); i < l; i++ {
				if r.Path[i] == ':' && n < pl {
					for ; i < l && r.Path[i] != '/'; i++ {
					}
					uri.WriteString(fmt.Sprintf("%v", params[n]))
					n++
				}
				if i < l {
					uri.WriteByte(r.Path[i])
				}
			}
			break
		}
	}
	return uri.String()
}

// URL is an alias for `URI` function.
func (e *Spark) URL(h Handler, params ...interface{}) string {
	return e.URI(h, params...)
}

// Routes returns the registered routes.
func (e *Spark) Routes() []Route {
	return e.router.routes
}

// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
func (e *Spark) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := e.pool.Get().(*Context)
	h, spark := e.router.Find(r.Method, r.URL.Path, c)
	if spark != nil {
		e = spark
	}
	c.reset(r, w, e)
	if h == nil {
		h = e.notFoundHandler
	}

	// Chain middleware with handler in the end
	for i := len(e.middleware) - 1; i >= 0; i-- {
		h = e.middleware[i](h)
	}

	// Execute chain
	if err := h(c); err != nil {
		e.httpErrorHandler(err, c)
	}

	e.pool.Put(c)
}

// Server returns the internal *http.Server.
func (e *Spark) Server(addr string) *http.Server {
	s := &http.Server{Addr: addr}
	s.Handler = e
	return s
}

// Run runs a server.
func (e *Spark) Run(addr string) {
	s := e.Server(addr)
	e.run(s)
}

// RunTLS runs a server with TLS configuration.
func (e *Spark) RunTLS(addr, certFile, keyFile string) {
	s := e.Server(addr)
	e.run(s, certFile, keyFile)
}

// RunServer runs a custom server.
func (e *Spark) RunServer(s *http.Server) {
	e.run(s)
}

// RunTLSServer runs a custom server with TLS configuration.
func (e *Spark) RunTLSServer(s *http.Server, certFile, keyFile string) {
	e.run(s, certFile, keyFile)
}

func (e *Spark) run(s *http.Server, files ...string) {
	if len(files) == 0 {
		log.Fatal(s.ListenAndServe())
	} else if len(files) == 2 {
		log.Fatal(s.ListenAndServeTLS(files[0], files[1]))
	} else {
		log.Fatal("spark => invalid TLS configuration")
	}
}

// NewHTTPError - http error
func NewHTTPError(code int, msg ...string) *HTTPError {
	he := &HTTPError{code: code, message: http.StatusText(code)}
	if len(msg) > 0 {
		m := msg[0]
		he.message = m
	}
	return he
}

// SetCode sets code.
func (e *HTTPError) SetCode(code int) {
	e.code = code
}

// Code returns code.
func (e *HTTPError) Code() int {
	return e.code
}

// Error returns message.
func (e *HTTPError) Error() string {
	return e.message
}

// wrapMiddleware wraps middleware.
func wrapMiddleware(m Middleware) MiddlewareFunc {
	switch m := m.(type) {
	case MiddlewareFunc:
		return m
	case func(HandlerFunc) HandlerFunc:
		return m
	case HandlerFunc:
		return wrapHandlerFuncMW(m)
	case func(*Context) error:
		return wrapHandlerFuncMW(m)
	case func(http.Handler) http.Handler:
		return func(h HandlerFunc) HandlerFunc {
			return func(c *Context) (err error) {
				m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.response.writer = w
					c.request = r
					err = h(c)
				})).ServeHTTP(c.response.writer, c.request)
				return
			}
		}
	case http.Handler:
		return wrapHTTPHandlerFuncMW(m.ServeHTTP)
	case func(http.ResponseWriter, *http.Request):
		return wrapHTTPHandlerFuncMW(m)
	default:
		panic("spark => unknown middleware")
	}
}

// wrapHandlerFuncMW wraps HandlerFunc middleware.
func wrapHandlerFuncMW(m HandlerFunc) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			if err := m(c); err != nil {
				return err
			}
			return h(c)
		}
	}
}

// wrapHTTPHandlerFuncMW wraps http.HandlerFunc middleware.
func wrapHTTPHandlerFuncMW(m http.HandlerFunc) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			if !c.response.committed {
				m.ServeHTTP(c.response.writer, c.request)
			}
			return h(c)
		}
	}
}

// wrapHandler wraps handler.
func wrapHandler(h Handler) HandlerFunc {
	switch h := h.(type) {
	case HandlerFunc:
		return h
	case func(*Context) error:
		return h
	case http.Handler, http.HandlerFunc:
		return func(c *Context) error {
			h.(http.Handler).ServeHTTP(c.response, c.request)
			return nil
		}
	case func(http.ResponseWriter, *http.Request):
		return func(c *Context) error {
			h(c.response, c.request)
			return nil
		}
	default:
		panic("spark => unknown handler")
	}
}
