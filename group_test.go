package spark

import "testing"

func TestGroup(t *testing.T) {
	g := New().Group("/group")
	h := func(*Context) error { return nil }
	g.Connect("/", h)
	g.Delete("/", h)
	g.Get("/", h)
	g.Head("/", h)
	g.Options("/", h)
	g.Patch("/", h)
	g.Post("/", h)
	g.Put("/", h)
	g.Trace("/", h)
	g.WebSocket("/ws", h)
	g.Static("/scripts", "scripts")
	g.ServeDir("/scripts", "scripts")
	g.ServeFile("/scripts/main.js", "scripts/main.js")
}
