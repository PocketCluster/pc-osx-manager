package route

import (
    "strings"

    "github.com/pkg/errors"
)

const (
    get     = "GET"
    head    = "HEAD"
    post    = "POST"
    put     = "PUT"
    patch   = "PATCH"
    deleteh = "DELETE"
)

// Handle is just like "net/http" Handlers, only takes params.
type Handle func(method, path, request string) error

type Router interface {
    GET(path string, handler Handle) error
    HEAD(path string, handler Handle) error
    POST(path string, handler Handle) error
    PUT(path string, handler Handle) error
    PATCH(path string, handler Handle) error
    DELETE(path string, handler Handle) error
    Dispatch(event Request) error
}

// router name says it all.
type router struct {
    tree           *node
    rootHandler    Handle
}

// New creates a new router. Take the root/fall through route
// like how the default mux works. Only difference is in this case,
// you have to specific one.
func NewRouter(rootHandler Handle) Router {
    node := node{component: "/", methods: make(map[string]Handle)}
    return &router{tree: &node, rootHandler: rootHandler}
}

// Handle takes an http handler, method and pattern for a route.
func (r *router) addHandleForPath(method, path string, handler Handle) error {
    if path[0] != '/' {
        return errors.Errorf("Path has to start with a /.")
    }
    r.tree.addNode(method, path, handler)
    return nil
}

// GET same as Handle only the method is already implied.
func (r *router) GET(path string, handler Handle) error {
    return r.addHandleForPath(get, path, handler)
}

// HEAD same as Handle only the method is already implied.
func (r *router) HEAD(path string, handler Handle) error {
    return r.addHandleForPath(head, path, handler)
}

// POST same as Handle only the method is already implied.
func (r *router) POST(path string, handler Handle) error {
    return r.addHandleForPath(post, path, handler)
}

// PUT same as Handle only the method is already implied.
func (r *router) PUT(path string, handler Handle) error {
    return r.addHandleForPath(put, path, handler)
}

// PATCH same as Handle only the method is already implied.
func (r *router) PATCH(path string, handler Handle) error { // might make this and put one.
    return r.addHandleForPath(patch, path, handler)
}

// DELETE same as Handle only the method is already implied.
func (r *router) DELETE(path string, handler Handle) error {
    return r.addHandleForPath(deleteh, path, handler)
}

// Needed by "net/http" to handle http requests and be a mux to http.ListenAndServe.
func (r *router) Dispatch(event Request) error {
    node, _ := r.tree.traverse(strings.Split(event.path, "/")[1:])
    if handler := node.methods[event.method]; handler != nil {
        return handler(event.method, event.path, event.request)
    }

    return r.rootHandler(event.method, event.path, event.request)
}
