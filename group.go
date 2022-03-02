package spark

type (
	// Group - Group stuct
	Group struct {
		spark Spark
	}
)

// Use - it is used to add middleware
func (g *Group) Use(m ...Middleware) {
	for _, h := range m {
		g.spark.middleware = append(g.spark.middleware, wrapMiddleware(h))
	}
}

// Connect - connect method
func (g *Group) Connect(path string, h Handler) {
	g.spark.Connect(path, h)
}

// Delete - delele method
func (g *Group) Delete(path string, h Handler) {
	g.spark.Delete(path, h)
}

// Get - get method
func (g *Group) Get(path string, h Handler) {
	g.spark.Get(path, h)
}

// Head - head method
func (g *Group) Head(path string, h Handler) {
	g.spark.Head(path, h)
}

// Options - options method
func (g *Group) Options(path string, h Handler) {
	g.spark.Options(path, h)
}

// Patch - patch method
func (g *Group) Patch(path string, h Handler) {
	g.spark.Patch(path, h)
}

// Post - post method
func (g *Group) Post(path string, h Handler) {
	g.spark.Post(path, h)
}

// Put - put method
func (g *Group) Put(path string, h Handler) {
	g.spark.Put(path, h)
}

// Trace - trace method
func (g *Group) Trace(path string, h Handler) {
	g.spark.Trace(path, h)
}

// WebSocket - websocket method
func (g *Group) WebSocket(path string, h HandlerFunc) {
	g.spark.WebSocket(path, h)
}

// Static - static method
func (g *Group) Static(path, root string) {
	g.spark.Static(path, root)
}

// ServeDir - servedir method
func (g *Group) ServeDir(path, root string) {
	g.spark.ServeDir(path, root)
}

// ServeFile - serve file method
func (g *Group) ServeFile(path, file string) {
	g.spark.ServeFile(path, file)
}

// Group - group with some prefix
func (g *Group) Group(prefix string, m ...Middleware) *Group {
	return g.spark.Group(prefix, m...)
}
