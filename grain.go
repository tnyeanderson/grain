package grain

import (
	"net/http"
)

// Middleware is usually part of a middleware chain.
type Middleware func(*Context)

// Context is passed between the middlewares to each request. The original
// [http.Request] and [http.ResponseWriter] are embedded.
type Context struct {
	*http.Request
	http.ResponseWriter

	// Data can be used to store arbitrary data to be passed between the
	// middleware functions.
	Data map[string]any

	done bool
}

// Done should be called from a middleware when it is ready to close the
// connection, usually meaning once it has written the response in full. After
// this is called, subsequent middlewares will not be called.
func (c *Context) Done() {
	c.done = true
}

// Handler returns an [http.Handler] which will call the middleware chain
// for each request it handles.
func Handler(middlewares ...Middleware) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{
			Request:        r,
			ResponseWriter: w,
			Data:           map[string]any{},
		}

		for _, middleware := range middlewares {
			if ctx.done == true {
				break
			}

			middleware(ctx)
		}
	}

	return http.HandlerFunc(f)
}
