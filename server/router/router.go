// https://vluxe.io/golang-router.html
package router

import (
	"net/http"
	"fmt"
	"regexp"
	"strings"
	"golang.org/x/net/context"
	"pluto/reply"
	"log"
)

// Handler is a function type like "net/http" Handler
type Handler func (http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

// Match
type Match struct {
	handler 	Handler
	ctx 		context.Context
}

// Mux interface to expose Router struct
type Mux interface {
	GET(string, Handler)
	POST(string, Handler)
	PUT(string, Handler)
	DELETE(string, Handler)
	ServeHTTP(http.ResponseWriter, *http.Request)
	AddContextWith(key interface{}, val interface{})
}

// Router
type Router struct {
	trie 		*Trie
}

// DefaultRootHandler
func DefaultRootHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello World!\n")
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	reply.Json(w, r, http.StatusNotFound, "404 page not found")
}

func NewRouter() Mux {
	return &Router{trie: NewTrie()}
}

// Handle takes a method, pattern, and http handler for a route.
func (r *Router) Handle(method, path string, handler Handler) {
	if matches, err := regexp.MatchString("^[A-Z]+$", method); !matches || err != nil {
		panic("Http method " + method + " is not valid")
	}
	if path[0] != '/' {
		panic("Path must start with /")
	}
	key, value, prefix, vars := transformPath(path)
	data := r.trie.Get(key)
	data.value = value
	data.prefix = prefix
	data.vars = vars
	data.methods[method] = handler
	r.trie.Put(key, data)
}

// Get is a shortcut for Handle with method "GET"
func (r *Router) GET(path string, handler Handler) {
	r.Handle("GET", path, handler)
}
// Post is a shortcut for Handle with method "GET"
func (r *Router) POST(path string, handler Handler) {
	r.Handle("POST", path, handler)
}
// Get is a shortcut for Handle with method "GET"
func (r *Router) PUT(path string, handler Handler) {
	r.Handle("PUT", path, handler)
}
// Get is a shortcut for Handle with method "GET"
func (r *Router) DELETE(path string, handler Handler) {
	r.Handle("DELETE", path, handler)
}

// AddContextWith applies a wrapper to all handlers,
// data is available in the handler context
func (r *Router) AddContextWith(key interface{}, val interface{}) {
	for _, k := range r.trie.Keys() {
		data := r.trie.Get(k)
		for m, h := range data.methods {
			data.methods[m] = wrapper(key, val, h)
			r.trie.Put(k, data)
		}
	}
}

func wrapper(key interface{}, val interface{}, next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, key, val)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// transformPath returns a tuple with key, value, prefix and params for the
// for the presented path
// e.g. /home/:id/room -> /home/:/room, /room, /home/:, [id]
func transformPath(path string) (key, value, prefix string, params []string) {
	params = []string{}
	if path == "/" {
		return "/", "/", "", params
	}
	if path[0] != '/' {
		panic("Path must start with '/'")
	}
	segments := strings.Split(path, "/")[1:]
	for _, s := range segments {
		if s[0] == ':'{
			params = append(params, s[1:])
			path = strings.Replace(path, s, ":", 1)
		}
	}
	key = path
	value = path[strings.LastIndex(key, "/"):]
	if i := strings.LastIndex(key, "/"); i !=0 {
		prefix = path[:i]
	}
	return key, value, prefix, params
}

func (r *Router) findData(method, path, sufix, key, segment string, values []string) (*Data, []string) {
	//log.Printf("findData method=%v, path=%v, sufix=%v, key=%v, segment=%v, values=%v", method, path, sufix, key, segment, values)

	// initialize
	if path != "" && sufix == "" && key == "" {
		// remove trailing slash
		if len(path) > 1 && strings.LastIndex(path, "/") == len(path)-1 {
			path = path[:len(path)-1]
		}
		sufix = path
		key = path
	}

	// test key
	if key != "" {
		// if key is valid returns
		if d := r.trie.Get(key); d.methods[method] != nil {
			return d, values
		}
		// Nothing found, returns nil
		if c := strings.Count(key, ":"); c != 0 && c == strings.Count(key, "/") {
			return nil, []string{}
		}
		// restore values
		if len (values) > 1 {
			// remove last segment in the slice
			values = append(values[:len(values)-1], values[:len(values)-2]...)
		} else {
			values = []string{}
		}
	}

	// initialize for the inner loop
	// : stays fixed in the last segment
	if sufix == "" {
		x := strings.Index(key, "/:")
		sufix = key[:x]
		values = append(values, segment)
		path = key
	}

	// TODO maybe try to use Regex
	i := strings.Index(sufix[1:], "/")
	if i == -1 {
		segment = sufix[1:]
		key = strings.Replace(path, segment, ":", 1)
		sufix = ""
	} else {
		segment = sufix[1:i+1]
		key = strings.Replace(path, segment, ":", 1)
		sufix = sufix[i+1:]
	}
	values = append(values, segment)
	return r.findData(method, path, sufix, key, segment, values)
}

func setContext (ctx context.Context, vars, values []string) context.Context {
	if (len(vars) != len(values)){
		return ctx
	}
	for i, value := range values {
		// pick opposite var
		ctx = context.WithValue(ctx, vars[len(vars)-1-i], value)
	}
	return ctx
}

func (r *Router) FindMatch(req *http.Request) *Match  {
	path := req.URL.Path
	method := req.Method
	data, values := r.findData(method, path, "", "", "", []string{})
	if data != nil {
		ctx := setContext(req.Context(), data.vars, values)
		handler := data.methods[req.Method]
		return &Match{handler: handler, ctx: ctx}
	}
	return nil
}

func (m *Match) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.handler.ServeHTTP(w, req.WithContext(m.ctx))
}

// ServeHTTP
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("----- %v \t%v", req.Method, req.URL)
	m := r.FindMatch(req)
	if m != nil{
		m.ServeHTTP(w, req)
	} else {
		NotFoundHandler(w, req)
	}
}
