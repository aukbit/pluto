// A string symbol table for extended ASCII strings, implemented
// using a 256-way trie.
// http://algs4.cs.princeton.edu/code/edu/princeton/cs/algs4/TrieST.java.html
//
package router

import (
	"bytes"
)

// R extended ASCII
const R = 256

type Trie struct {
	// root of trie
	root 			*Node
	// number of keys in trie
	n				int
}

// R-way trie node
type Node struct {
	data			*Data
	next			[R]*Node
}

// NewNode
func NewNode() *Node {
	//return &Node{data: &Data{}}
	return &Node{data: NewData()}
}

// newTrie
func NewTrie() *Trie {
	return &Trie{root: NewNode()}
}

// Get returns the value associated with the given key
func (t *Trie) Get (key string) *Data {
	n := get(t.root, key, 0)
	if n == nil {
		//return &Data{}
		return NewData()
		//return nil
	}
	return n.data
}

// get
func get (n *Node, key string, d int) *Node {
	if n == nil {
		return nil
	}
	if d == len(key) {
		return n
	}
	return get(n.next[key[d]], key, d+1)
}

// contains verify if a key exists in the Trie
func (t *Trie) Contains (key string) bool {
	return t.Get(key).value != ""
}

// Put Inserts the key-value pair into the symbol table, overwriting the old value
// with the new value if the key is already in the symbol table.
// If the value is nil, this effectively deletes the key from the symbol table.
func (t *Trie) Put(key string, data *Data) {
	if data.value == "" {
		t.Remove(key)
	} else {
		t.root = t.put(t.root, key, data, 0)
	}
}

// put
func (t *Trie) put(n *Node, key string, data *Data, d int) *Node{
	//log.Printf("put n=%v key=%v data=%v d=%v", &n, key, data, d)
	if n == nil {
		n = NewNode()
	}
	//log.Printf("put n.data=%v len(n.next)=%v", n.data, len(n.next))
	if d == len(key) {
		//log.Printf("put d==len(key) %v, n.data.value=%v", d, n.data.value)
		if n.data.value == "" {
			t.n++
		}
		//log.Printf("put t.n %v", t.n)
		n.data = data
		return n
	}
	c := key[d]
	//log.Printf("put n.next[c]=%v", n.next[c])
	n.next[c] = t.put(n.next[c], key, data, d+1);
	return n
}

// Remove the key from the trie if present
func (t *Trie) Remove(key string) {
	t.root = t.remove(t.root, key, 0)
}

func (t *Trie) remove(n *Node, key string, d int) *Node {
	if n == nil {
		return nil
	}
	if (d == len(key)) {
		if n.data.value != "" {
			t.n--
		}
		n.data.value = ""
		//n.data = &Data{}
	} else {
		c := key[d]
		n.next[c] = t.remove(n.next[c], key, d+1)
	}

	// remove sub-trie rooted at n if its completely empty
	//if !reflect.DeepEqual(n.data, &Data{}) {
	if n.data.value != "" {
		return n
	}
	for i:=0; i < R; i++ {
		if n.next[i] != nil {
			return n
		}
	}
	return nil
}

// Size returns the number of key-value pairs in this trie
func (t *Trie) Size() int {
	return t.n
}

// Size returns the number of key-value pairs in this trie
func (t *Trie) IsEmpty() bool {
	return t.n == 0
}

// Keys returns all keys in the trie
func (t *Trie) Keys() []string {
	return (t.KeysWithPrefix(""))
}

// KeysWithPrefix returns all keys in the trie
// that start with a prefix
func (t *Trie) KeysWithPrefix(prefix string) []string {
	n := get(t.root, prefix, 0)
	results := &[]string{}
	t.collectKeysWithPrefix(n, prefix, results)
	return *results
}

func (t *Trie) collectKeysWithPrefix(n *Node, prefix string, results *[]string) {
	if n == nil {
		return
	}
	if n.data.value != "" {
		*results = append(*results, prefix)
	}
	for c:=0; c < R; c++ {
		buffer := bytes.NewBufferString(prefix)
		buffer.WriteByte(byte(c))
		t.collectKeysWithPrefix(n.next[c], buffer.String(), results)
	}
}

// KeysThatMatch returns all keys in the trie that match a pattern
// where . symbol is treated as a wildcard character
func (t *Trie) KeysThatMatch(pattern string) []string {
	results := &[]string{}
	t.collectKeysThatMatch(t.root, "", pattern, results)
	return *results
}

func (t *Trie) collectKeysThatMatch(n *Node, prefix, pattern string, results *[]string) {
	if n == nil {
		return
	}
	if len(prefix) == len(pattern) && n.data.value != ""{
		*results = append(*results, prefix)
	}
	if len(prefix) == len(pattern) {
		return
	}
	p := pattern[len(prefix)]
	if  p == '.' {
		for c := 0; c < R; c++ {
			buffer := bytes.NewBufferString(prefix)
			buffer.WriteByte(byte(c))
			t.collectKeysThatMatch(n.next[c], buffer.String(), pattern, results)
		}
	} else {
		buffer := bytes.NewBufferString(prefix)
		buffer.WriteByte(byte(p))
		t.collectKeysThatMatch(n.next[p], buffer.String(), pattern, results)
	}
}

// LongestPrefixOf returns the string that is the longest prefix
// or null, if no such string
func (t *Trie) LongestPrefixOf(query string) string {
	length := t.longestPrefixOf(t.root, query, 0, -1)
	if length == -1 {
		return ""
	} else {
		return query[0:length]
	}
}

// longestPrefixOf returns the length of the longest string key in the subtrie
// rooted at x that is a prefix of the query string,
// assuming the first d character match and we have already
// found a prefix match of given length (-1 if no such match)
func (t *Trie) longestPrefixOf(n *Node, query string, d, length int) int {
	if n == nil {
		return length
	}
	if n.data.value != "" {
		length = d
	}
	if d == len(query) {
		return length
	}
	c := query[d]
	return t.longestPrefixOf(n.next[c], query, d+1, length)
}