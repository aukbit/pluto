// Package router trie based on a 256-way trie expressed on the textbook
// Algorithms, 4th Edition by Robert Sedgewick and Kevin Wayne
// A string symbol table for extended ASCII strings, implemented
// using a 256-way trie.
// http://algs4.cs.princeton.edu/code/edu/princeton/cs/algs4/TrieST.java.html
//
package router

import "bytes"

// R extended ASCII
const R = 256

//
// TRIE
//

// trie struct
type trie struct {
	// root of trie
	root *node
	// number of keys in trie
	n int
}

// newTrie creates new instace trie
func newTrie() *trie {
	return &trie{root: newNode()}
}

// Get returns the value associated with the given key
func (t *trie) Get(key string) *data {
	n := get(t.root, key, 0)
	if n == nil {
		return nil
	}
	return n.data
}

// get
func get(n *node, key string, index int) *node {
	if n == nil {
		// n = newNode()
		return nil
	}
	if index == len(key) {
		return n
	}
	ascii := key[index]
	return get(n.next[ascii], key, index+1)
}

// Contains verify if a key exists in the Trie
func (t *trie) Contains(key string) bool {
	return t.Get(key).value != ""
}

// Put Inserts the key-value pair into the symbol table, overwriting the old value
// with the new value if the key is already in the symbol table.
// If the value is nil, this effectively deletes the key from the symbol table.
func (t *trie) Put(key string, data *data) {
	if data.value == "" {
		t.Remove(key)
	}
	var ok bool

	_, ok = put(t.root, key, data, 0)
	if ok {
		t.n++
	}
}

// put
func put(n *node, key string, data *data, index int) (*node, bool) {
	if n == nil {
		n = newNode()
	}
	var ok bool
	if index == len(key) {
		if n.data == nil {
			ok = true
		}
		n.data = data
		return n, ok
	}
	ascii := key[index]
	n.next[ascii], ok = put(n.next[ascii], key, data, index+1)
	return n, ok
}

// Remove the key from the trie if present
func (t *trie) Remove(key string) {
	t.root = t.remove(t.root, key, 0)
}

func (t *trie) remove(n *node, key string, d int) *node {
	if n == nil {
		return nil
	}
	if d == len(key) {
		if n.data.value != "" {
			t.n--
		}
		n.data.value = ""
		//n.data = &data{}
	} else {
		c := key[d]
		n.next[c] = t.remove(n.next[c], key, d+1)
	}

	// remove sub-trie rooted at n if its completely empty
	//if !reflect.DeepEqual(n.data, &data{}) {
	if n.data.value != "" {
		return n
	}
	for i := 0; i < R; i++ {
		if n.next[i] != nil {
			return n
		}
	}
	return nil
}

// Size returns the number of key-value pairs in this trie
func (t *trie) Size() int {
	return t.n
}

// IsEmpty returns the number of key-value pairs in this trie
func (t *trie) IsEmpty() bool {
	return t.n == 0
}

// Keys returns all keys in the trie
func (t *trie) Keys() []string {
	return (t.KeysWithPrefix(""))
}

// KeysWithPrefix returns all keys in the trie
// that start with a prefix
func (t *trie) KeysWithPrefix(prefix string) []string {
	n := get(t.root, prefix, 0)
	results := &[]string{}
	collectKeysWithPrefix(n, prefix, results)
	return *results
}

func collectKeysWithPrefix(n *node, prefix string, results *[]string) {
	if n == nil {
		return
	}
	if n.data != nil {
		*results = append(*results, prefix)
	}

	for c := 0; c < R; c++ {
		buffer := bytes.NewBufferString(prefix)
		buffer.WriteByte(byte(c))
		collectKeysWithPrefix(n.next[c], buffer.String(), results)
	}
}

// KeysThatMatch returns all keys in the trie that match a pattern
// where . symbol is treated as a wildcard character
func (t *trie) KeysThatMatch(pattern string) []string {
	results := &[]string{}
	collectKeysThatMatch(t.root, "", pattern, results)
	return *results
}

func collectKeysThatMatch(n *node, prefix, pattern string, results *[]string) {
	if n == nil {
		return
	}
	if len(prefix) == len(pattern) && n.data != nil {
		*results = append(*results, prefix)
	}
	if len(prefix) == len(pattern) {
		return
	}
	p := pattern[len(prefix)]
	if p == '.' {
		for c := 0; c < R; c++ {
			buffer := bytes.NewBufferString(prefix)
			buffer.WriteByte(byte(c))
			collectKeysThatMatch(n.next[c], buffer.String(), pattern, results)
		}
	} else {
		buffer := bytes.NewBufferString(prefix)
		buffer.WriteByte(byte(p))
		collectKeysThatMatch(n.next[p], buffer.String(), pattern, results)
	}
}

// LongestPrefixOf returns the string that is the longest prefix
// or null, if no such string
func (t *trie) LongestPrefixOf(query string) string {
	length := longestPrefixOf(t.root, query, 0, -1)
	if length == -1 {
		return ""
	}
	return query[0:length]
}

// longestPrefixOf returns the length of the longest string key in the subtrie
// rooted at x that is a prefix of the query string,
// assuming the first d character match and we have already
// found a prefix match of given length (-1 if no such match)
func longestPrefixOf(n *node, query string, d, length int) int {
	if n == nil {
		return length
	}
	if n.data != nil {
		length = d
	}
	if d == len(query) {
		return length
	}
	c := query[d]
	return longestPrefixOf(n.next[c], query, d+1, length)
}

//
// NODE
//

// node a representation of each trie node
type node struct {
	data *data
	next [R]*node
}

// NewNode creates new instace node
func newNode() *node {
	return &node{}
}

//
// DATA
//

// Data a data struct that each node can handle
type data struct {
	value   string
	prefix  string
	vars    []string
	methods map[string]Handler
}

// newData returns a new data instance
func newData() *data {
	return &data{
		vars:    []string{},
		methods: make(map[string]Handler),
	}
}

// GetValue returns data value
func (d *data) Value() string {
	return d.value
}

// SetValue sets data value
func (d *data) SetValue(val string) {
	d.value = val
}
