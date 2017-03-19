package router

import (
	"reflect"
	"testing"

	"github.com/paulormart/assert"
)

func TestNewTrie(t *testing.T) {
	tr := newTrie()
	assert.Equal(t, reflect.TypeOf(&trie{}), reflect.TypeOf(tr))
	assert.Equal(t, 0, tr.n)
	assert.Equal(t, 0, tr.Size())
	assert.Equal(t, true, tr.IsEmpty())
	assert.Equal(t, reflect.TypeOf(&node{}), reflect.TypeOf(tr.root))
}

func TestNewData(t *testing.T) {
	d := newData()
	d.value = "home"
	d.prefix = "/"
	assert.Equal(t, reflect.TypeOf(&data{}), reflect.TypeOf(d))
	assert.Equal(t, "home", d.value)
	assert.Equal(t, "/", d.prefix)
	assert.Equal(t, []string{}, d.vars)
	assert.Equal(t, make(map[string]HandlerFunc), d.methods)
}

func TestPut(t *testing.T) {
	d := newData()
	d.value = "home"
	d.prefix = "/"
	tr := newTrie()
	tr.Put("/home", d)
	assert.Equal(t, 1, tr.Size())
}

func TestGet(t *testing.T) {
	tr := newTrie()
	// empty
	assert.Equal(t, 0, tr.Size())
	assert.Equal(t, nil, tr.Get("/home"))

	d := newData()
	d.value = "/home"
	d.prefix = "/"
	tr.Put("/home", d)
	assert.Equal(t, 1, tr.Size())
	assert.Equal(t, true, tr.Contains("/home"))
	// nil
	assert.Equal(t, nil, tr.Get("/"))
}

func TestKeys(t *testing.T) {
	// key /home
	d := newData()
	d.value = "/home"
	d.prefix = "/"
	tr := newTrie()
	tr.Put("/home", d)
	// key /room
	d = newData()
	d.value = "/room"
	d.prefix = "/"
	tr.Put("/room", d)
	assert.Equal(t, 2, tr.Size())
	// KeysWithPrefix
	kwp := tr.KeysWithPrefix("/")
	assert.Equal(t, 2, len(kwp))
	// KeysThatMatch
	ktm := tr.KeysThatMatch("")
	assert.Equal(t, 0, len(ktm))
	ktm1 := tr.KeysThatMatch("/home")
	assert.Equal(t, 1, len(ktm1))
	ktm2 := tr.KeysThatMatch("/room")
	assert.Equal(t, 1, len(ktm2))
	ktm3 := tr.KeysThatMatch("..o..")
	assert.Equal(t, 2, len(ktm3))
	ktm4 := tr.KeysThatMatch(".....")
	assert.Equal(t, 2, len(ktm4))
	ktm5 := tr.KeysThatMatch(".")
	assert.Equal(t, 0, len(ktm5))
	// Keys
	k := tr.Keys()
	assert.Equal(t, 2, len(k))
	// LongestPrefixOf
	lpo := tr.LongestPrefixOf("")
	assert.Equal(t, "", lpo)
	lpo1 := tr.LongestPrefixOf("/home")
	assert.Equal(t, "/home", lpo1)
}

func BenchmarkPutTrie(b *testing.B) {
	d := newData()
	d.value = "home"
	d.prefix = "/"

	// run the Put function b.N times
	for n := 0; n < b.N; n++ {
		tr := newTrie()
		tr.Put("/test", d)
	}
}
