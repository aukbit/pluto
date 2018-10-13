package router

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	context "golang.org/x/net/context"

	"github.com/aukbit/pluto/reply"
	"github.com/rs/zerolog"
)

//
// HANDLER
//

// Handler is a function type like "net/http" Handler
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// HandlerFunc is a function type like "net/http" HandlerFunc
type HandlerFunc func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r)
func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn(w, r)
}

// Err struct containing an error and helper fields
type Err struct {
	Err      error             `json:"-"`
	Status   int               `json:"-"`
	Type     string            `json:"type"`
	Message  string            `json:"message"`
	Code     string            `json:"code,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// WrapErr reduces the repetition of dealing with errors in Handlers
// returning an error
// https://blog.golang.org/error-handling-and-go
type WrapErr func(http.ResponseWriter, *http.Request) *Err

func (fn WrapErr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *Err, not os.Error.
		zerolog.Ctx(r.Context()).Error().Msgf("%v", e.Err)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		reply.Json(w, r, e.Status, e)
	}
}

//
// ROUTER
//

// Router ..
type Router struct {
	trie *trie
}

// NewRouter creates a new router instance
func NewRouter() *Router {
	return &Router{trie: newTrie()}
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
	if data == nil {
		data = newData()
	}
	data.value = value
	data.prefix = prefix
	data.vars = vars
	data.methods[method] = handler.ServeHTTP
	r.trie.Put(key, data)
}

// HandleFunc registers the handler function for the given pattern.
func (r *Router) HandleFunc(method, path string, handlerFn HandlerFunc) {
	r.Handle(method, path, HandlerFunc(handlerFn))
}

// GET is a shortcut for Handle with method "GET"
func (r *Router) GET(path string, handlerFn HandlerFunc) {
	r.HandleFunc("GET", path, handlerFn)
}

// POST is a shortcut for Handle with method "POST"
func (r *Router) POST(path string, handlerFn HandlerFunc) {
	r.HandleFunc("POST", path, handlerFn)
}

// PUT is a shortcut for Handle with method "PUT"
func (r *Router) PUT(path string, handler HandlerFunc) {
	r.HandleFunc("PUT", path, handler)
}

// DELETE is a shortcut for Handle with method "DELETE"
func (r *Router) DELETE(path string, handlerFn HandlerFunc) {
	r.HandleFunc("DELETE", path, handlerFn)
}

// ServeHTTP
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m := r.findMatch(req)
	if m == nil {
		NotFoundHandler(w, req)
		return
	}
	m.ServeHTTP(w, req)
}

// WrapperMiddleware ..
func (r *Router) WrapperMiddleware(mids ...Middleware) {
	for _, k := range r.trie.Keys() {
		data := r.trie.Get(k)
		for m, h := range data.methods {
			data.methods[m] = Wrap(h, mids...)
			r.trie.Put(k, data)
		}
	}
}

// func (r *Router) findMatch(req *http.Request) *Match {
// 	path := req.URL.Path
// 	method := req.Method
// 	data, values := findData(r, method, path, "", "", "", []string{})
// 	if data != nil {
// 		ctx := setContext(req.Context(), data.vars, values)
// 		handler := data.methods[req.Method]
// 		return &Match{handler: handler, ctx: ctx}
// 	}
// 	return nil
// }

func (r *Router) findMatch(req *http.Request) *Match {
	path := req.URL.Path
	paths := validPaths(path, "", "", "", nil, nil)
	for key, values := range paths {
		data := r.trie.Get(key)
		if data != nil {
			ctx := setContext(req.Context(), data.vars, values)
			handler, ok := data.methods[req.Method]
			if !ok {
				return nil
			}
			return &Match{handler: handler, ctx: ctx}
		}
	}
	return nil
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
		if s[0] == ':' {
			params = append(params, s[1:])
			path = strings.Replace(path, s, ":", 1)
		}
	}
	key = path
	value = path[strings.LastIndex(key, "/"):]
	if i := strings.LastIndex(key, "/"); i != 0 {
		prefix = path[:i]
	}
	return key, value, prefix, params
}

// validPaths returns candidate paths to be searched on trie
// eg. /home/123/room -> {"/home/123/room":[], "/home/123/:":["room"],
// "/home/:/room":["123"], "/:/123/room":["home"], "/home/:/:":["123","room"],
// "/:/123/:":["home","room"], "/:/:/room":["home","123"], "/:/:/:":["home","123","room"]}
func validPaths(path, key, segment, cutset string, track []string, out map[string][]string) map[string][]string {
	// fmt.Printf("validPaths path %v, key %v, segment %v, cutset %v, track %v, out %v\n", path, key, segment, cutset, track, out)
	if out == nil {
		out = make(map[string][]string)
	}
	if path == "/" {
		out[path] = []string{}
		return out
	}
	// first iteration
	if path != "" && len(key) == 0 {
		// remove trailing slash `/`
		if len(path) > 1 && path[len(path)-1:] == `/` {
			path = path[:len(path)-1]
		}
		key = path
		cutset = path
		out[key] = []string{}
	}
	// all done end iteration
	if strings.Count(path, "/") == strings.Count(key, ":") {
		return out
	}

	// TODO maybe try to use Regex
	i := strings.Index(cutset[1:], "/")
	if i == -1 {
		segment = cutset[1:]
		key = path[:strings.LastIndex(path, segment)] + strings.Replace(path[strings.LastIndex(path, segment):], segment, ":", 1)
		tmp := []string{segment}
		track = append(tmp, track...)
		out[key] = append(out[key], track...)
		// prepare inputs
		cutset = key[:strings.Index(key, "/:")]
		return validPaths(key, key, segment, cutset, track, out)
	}
	segment = cutset[1 : i+1]
	key = strings.Replace(path, segment, ":", 1)
	out[key] = append(out[key], segment)
	out[key] = append(out[key], track...)
	cutset = cutset[i+1:]
	return validPaths(path, key, segment, cutset, track, out)

}

// findData.. deprecated
func findData(r *Router, method, path, suffix, key, segment string, values []string) (*data, []string) {
	// log.Printf("findData method:%v path:%v suffix:%v key:%v segment:%v values:%v\n", method, path, suffix, key, segment, values)
	// initialize
	if path != "" && suffix == "" && key == "" {
		// remove trailing slash `/`
		if len(path) > 1 && strings.LastIndex(path, "/") == len(path)-1 {
			path = path[:len(path)-1]
		}
		suffix = path
		key = path
	}

	// test key
	if key != "" {
		// if key is valid returns

		if d := r.trie.Get(key); d != nil {
			return d, values
		}
		// if d := r.trie.Get(key); d.methods[method] != nil {
		// 	return d, values
		// }

		// Nothing found, returns nil
		if c := strings.Count(key, ":"); c != 0 && c == strings.Count(key, "/") {
			return nil, []string{}
		}

		// restore values
		if len(values) > 1 {
			// remove last segment in the slice
			values = append(values[:len(values)-1], values[:len(values)-2]...)
		} else {
			values = []string{}
		}
	}

	// initialize for the inner loop
	// : stays fixed in the last segment
	if suffix == "" {
		x := strings.Index(key, "/:")
		suffix = key[:x]
		values = append(values, segment)
		path = key
	}

	// TODO maybe try to use Regex
	i := strings.Index(suffix[1:], "/")
	if i == -1 {
		segment = suffix[1:]
		key = path[:strings.LastIndex(path, segment)] + strings.Replace(path[strings.LastIndex(path, segment):], segment, ":", 1)
		// key = strings.Replace(path, segment, ":", 1)
		suffix = ""
	} else {
		segment = suffix[1 : i+1]
		key = strings.Replace(path, segment, ":", 1)
		suffix = suffix[i+1:]
	}
	values = append(values, segment)

	d, v := findData(r, method, path, suffix, key, segment, values)
	return d, v
}

// func setContext(ctx context.Context, vars, values []string) context.Context {
// 	if len(vars) != len(values) {
// 		return ctx
// 	}
// 	for i, value := range values {
// 		// pick opposite var
// 		ctx = context.WithValue(ctx, vars[len(vars)-1-i], value)
// 	}
// 	return ctx
// }

func setContext(ctx context.Context, vars, values []string) context.Context {
	if len(vars) != len(values) {
		return ctx
	}
	for i, value := range values {
		ctx = context.WithValue(ctx, contextKey{vars[i]}, value)
	}
	return ctx
}

// Middleware wraps an http.HandlerFunc with additional
// functionality.
type Middleware func(HandlerFunc) HandlerFunc

// Wrap h with all specified middlewares.
func Wrap(h HandlerFunc, middlewares ...Middleware) HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

// Match wraps an Handler and context
type Match struct {
	handler Handler
	ctx     context.Context
}

func (m *Match) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.handler.ServeHTTP(w, req.WithContext(m.ctx))
}

// DefaultRootHandler hello world handler
func DefaultRootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!\n")
}

// NotFoundHandler default not found resource json handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	reply.Json(w, r, http.StatusNotFound, &Err{
		Type:    "invalid_request_error",
		Message: fmt.Sprintf("Invalid request errors arise when your request has invalid parameters. path: %v query: %v", r.URL.EscapedPath(), r.URL.RawQuery),
	})
}
