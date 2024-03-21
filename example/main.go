/*
This example shows what grain can do.

Define middleware chains, wrap the grain handler with additional logging
capabilites, easily make a function to "trap" errors from any middleware in
your chain.

Using grain makes these middlewares easy to understand, because there just
isn't much to it! Again, this is for educational purposes. It is not meant to
be particularly useful.
*/
package main

import (
	"fmt"
	"grain"
	"log"
	"net/http"
)

func main() {
	// Register handlers
	http.Handle("/", handlerAll())

	http.Handle("/saywhat", grain.Handler(
		func(c *grain.Context) { c.ResponseWriter.Write([]byte("what")) },
	))

	http.Handle("/log", logger(
		func(c *grain.Context) { log.Println("called /log") },
	))

	http.Handle("/fail", handlerFail())

	// Start the server
	http.ListenAndServe(":8989", nil)
}

func handlerAll() http.Handler {
	return grain.Handler(
		func(c *grain.Context) { log.Println("called middleware1") },
		func(c *grain.Context) { log.Println("called middleware2") },
		func(c *grain.Context) {
			log.Println("setting val=hello in context")
			c.Data["val"] = "hello"
		},
		func(c *grain.Context) {
			log.Printf("val is now %s", c.Data["val"])
		},
		func(c *grain.Context) {
			log.Println("responding and finishing...")
			c.ResponseWriter.Write([]byte("woot"))
			c.Done()
		},
		func(c *grain.Context) {
			log.Println("should not be called")
		},
	)
}

func handlerFail() http.Handler {
	return grain.Handler(
		func(c *grain.Context) { log.Println("called /fail") },
		func(c *grain.Context) {
			log.Println("failing")
			if true {
				throw(c, fmt.Errorf("whoops"))
				return
			}
			log.Println("should not print")
		},
		func(c *grain.Context) {
			log.Println("should not be called")
		},
	)
}

// logger returns an [http.Handler] that logs at the beginning and end of a middleware
// chain
func logger(middlewares ...grain.Middleware) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		log.Println("got request")
		grain.Handler(middlewares...).ServeHTTP(w, r)
		log.Println("finished request")
	}
	return http.HandlerFunc(f)
}

// throw can be called by any other middleware to respond with an error and
// call [c.Done]. Return after calling this!
func throw(c *grain.Context, err error) {
	c.ResponseWriter.WriteHeader(http.StatusInternalServerError)
	c.ResponseWriter.Write([]byte(err.Error()))
	c.Done()
}
