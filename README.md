# Spark

A micro web framework for Golang.


## Features

- Fast HTTP router which smartly prioritize routes.
- Extensible middleware, supports:
	- `spark.MiddlewareFunc`
	- `func(spark.HandlerFunc) spark.HandlerFunc`
	- `spark.HandlerFunc`
	- `func(*spark.Context) error`
	- `func(http.Handler) http.Handler`
	- `http.Handler`
	- `http.HandlerFunc`
	- `func(http.ResponseWriter, *http.Request)`
- Extensible handler, supports:
    - `spark.HandlerFunc`
    - `func(*spark.Context) error`
    - `http.Handler`
    - `http.HandlerFunc`
    - `func(http.ResponseWriter, *http.Request)`
- Sub-router/Groups
- Handy functions to send variety of HTTP response:
    - HTML
    - HTML via templates
    - String 
    - JSON
    - XML
    - NoContent
    - Redirect
    - Error
- Build-in support for:
	- Favicon
	- Index file
	- Static files
	- WebSocket
- Centralized HTTP error handling.
- Customizable HTTP request binding function.
- Customizable HTTP response rendering function, allowing you to use any HTML template engine.

## Installation

```sh
$ go get github.com/santoshanand/spark
```


Inspired with Echo Framework.