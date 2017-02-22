package router

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"context"

	"bitbucket.org/aukbit/pluto/reply"
)

//
// HANDLER
//

// Handler is a function type like "net/http" Handler
type Handler func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
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
	data.methods[method] = handler
	r.trie.Put(key, data)
}

// GET is a shortcut for Handle with method "GET"
func (r *Router) GET(path string, handler Handler) {
	r.Handle("GET", path, handler)
}

// POST is a shortcut for Handle with method "GET"
func (r *Router) POST(path string, handler Handler) {
	r.Handle("POST", path, handler)
}

// PUT is a shortcut for Handle with method "GET"
func (r *Router) PUT(path string, handler Handler) {
	r.Handle("PUT", path, handler)
}

// DELETE is a shortcut for Handle with method "GET"
func (r *Router) DELETE(path string, handler Handler) {
	r.Handle("DELETE", path, handler)
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

func (r *Router) findMatch(req *http.Request) *Match {
	path := req.URL.Path
	method := req.Method
	data, values := findData(r, method, path, "", "", "", []string{})
	// if data != nil {
	// 	ctx := setContext(req.Context(), data.vars, values)
	// 	handler := data.methods[req.Method]
	// 	return &Match{handler: handler, ctx: ctx}
	// }
	fmt.Printf("findMatch  %v %v\n", data, values)
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
func validPaths(path, key, segment, cutset string, c int, track []string, out map[string][]string) map[string][]string {
	// fmt.Printf("validPaths %v %v %v %v\n", path, segment, track, out)
	if out == nil {
		out = make(map[string][]string)
	}
	if path == "/" {
		out[path] = []string{}
		return out
	}
	if path != "" && len(key) == 0 {
		// remove trailing slash `/`
		if len(path) > 1 && path[len(path)-1:] == `/` {
			path = path[:len(path)-1]
		}
		key = path
		cutset = path

	}

	// initialize for the inner loop
	// : stays fixed in the last segment
	if cutset == "" {
		x := strings.Index(key, "/:")
		cutset = key[:x]
		path = key
	}

	// TODO maybe try to use Regex
	// i := strings.Index(cutset[1:], "/")
	// if i == -1 {
	// 	segment = cutset[1:]
	// 	key = path[:strings.LastIndex(path, segment)] + strings.Replace(path[strings.LastIndex(path, segment):], segment, ":", 1)
	// 	track = append(track, segment)
	// 	fmt.Printf("format 1 %v %v\n", track, cutset)
	// 	out[key] = append(out[key], track...)
	// 	cutset = ""
	// } else {
	// 	segment = cutset[1 : i+1]
	// 	key = strings.Replace(path, segment, ":", 1)
	// 	cutset = cutset[i+1:]
	// 	fmt.Printf("format 2 %v %v\n", segment, cutset)
	// 	out[key] = append(out[key], segment)
	// 	// fmt.Printf("*** %v %v %v\n", segment, key, cutset)
	// }

	i := strings.Index(cutset[1:], "/")
	fmt.Printf("i %v\n", i)
	if i == -1 {
		segment = cutset[1:]
		key = path[:strings.LastIndex(path, segment)] + strings.Replace(path[strings.LastIndex(path, segment):], segment, ":", 1)
		fmt.Printf("path %v key %v\n", path, key)
		track = append(track, segment)
		out[key] = append(out[key], track...)
		cutset = ""
		track = []string{segment}
		c++
	} else {
		segment = cutset[1 : i+1]
		key = strings.Replace(path, segment, ":", 1)
		cutset = cutset[i+1:]
		if c > 0 {
			track = append(track, segment)
			out[key] = append(out[key], track...)
			// track = []string{}
		} else {
			out[key] = append(out[key], segment)
		}
	}

	//

	fmt.Printf("out %v, path %v\n", out, path)
	out = validPaths(path, key, segment, cutset, c, track, out)
	return out
	// if segment != "" {
	// 	// replace path by segment
	// 	path = strings.Replace(path, `/`+segment, "/-", 1)
	// }
	// out[test] = append(out[test], track...)
	// //
	// segments := strings.Split(test, "/")[1:]
	// for _, s := range segments {
	// 	if s != `-` {
	// 		track = append(track, s)
	// 		return validPaths(path, s, test, track, out)
	// 	}
	//
	// 	// if s[0] == ':' {
	// 	// 	params = append(params, s[1:])
	// 	// 	path = strings.Replace(path, s, ":", 1)
	// }
	// fmt.Printf("out %v %v %v %v\n", out, path, cutset, segment)
	// return out
}

func findData(r *Router, method, path, suffix, key, segment string, values []string) (*data, []string) {
	log.Printf("findData method:%v path:%v suffix:%v key:%v segment:%v values:%v\n", method, path, suffix, key, segment, values)
	// initialize
	if path != "" && suffix == "" && key == "" {
		// remove trailing slash `/`
		if len(path) > 1 && strings.LastIndex(path, "/") == len(path)-1 {
			path = path[:len(path)-1]
		}
		suffix = path
		key = path
	}

	fmt.Printf("key 1 %v\n", key)

	// test key
	if key != "" {
		// if key is valid returns

		if d := r.trie.Get(key); d != nil {
			return d, values
		}
		// if d := r.trie.Get(key); d.methods[method] != nil {
		// 	return d, values
		// }
		fmt.Printf("key 2 %v\n", key)

		// Nothing found, returns nil
		if c := strings.Count(key, ":"); c != 0 && c == strings.Count(key, "/") {
			fmt.Printf("key 2.1 %v %v\n", key, values)
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

	fmt.Printf("key 3 %v\n", key)

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
	fmt.Printf("key 4 %v %v\n", d, v)
	return d, v
}

func setContext(ctx context.Context, vars, values []string) context.Context {
	if len(vars) != len(values) {
		return ctx
	}
	for i, value := range values {
		// pick opposite var
		ctx = context.WithValue(ctx, vars[len(vars)-1-i], value)
	}
	return ctx
}

// Middleware wraps an http.Handler with additional
// functionality.
type Middleware func(Handler) Handler

// Wrap h with all specified middlewares.
func Wrap(h Handler, middlewares ...Middleware) Handler {
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
	reply.Json(w, r, http.StatusNotFound, "404 page not found")
}
