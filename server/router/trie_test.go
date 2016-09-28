package router_test

import (
	"reflect"
	"testing"

	"bitbucket.org/aukbit/pluto/server/router"
	"github.com/paulormart/assert"
)

func TestStructs(t *testing.T) {
	// comparing structs
	type A struct {
		val int
	}
	assert.Equal(t, false, &A{} == &A{})
	assert.Equal(t, false, &A{} == new(A))
	assert.Equal(t, false, new(A) == new(A))
	assert.Equal(t, true, reflect.DeepEqual(&A{}, &A{}))
	assert.Equal(t, true, reflect.DeepEqual(&A{}, new(A)))
	assert.Equal(t, true, reflect.DeepEqual(new(A), new(A)))
	cc := &A{}
	assert.Equal(t, true, reflect.DeepEqual(cc, new(A)))

	//fmt.Printf("new(A) = %v", new(A))
	assert.Equal(t, true, reflect.DeepEqual(*(new(A)), A{}))
}

func TestTrie(t *testing.T) {

	data := router.NewData()
	assert.Equal(t, reflect.TypeOf(&router.Data{}), reflect.TypeOf(data))
	node := router.NewNode()
	assert.Equal(t, reflect.TypeOf(&router.Node{}), reflect.TypeOf(node))
	trie := router.NewTrie()
	assert.Equal(t, reflect.TypeOf(&router.Trie{}), reflect.TypeOf(trie))

	g := trie.Get("home")
	assert.Equal(t, data, g)
	assert.Equal(t, false, trie.Contains("home"))
	//
	d := &router.Data{}
	d.SetValue("hello")
	trie.Put("home", d)
	assert.Equal(t, true, trie.Contains("home"))
	assert.Equal(t, 1, trie.Size())
	assert.Equal(t, false, trie.IsEmpty())
	//
	trie.Remove("home")
	assert.Equal(t, false, trie.Contains("home"))
	assert.Equal(t, 0, trie.Size())
	assert.Equal(t, true, trie.IsEmpty())
	//
	a := &router.Data{}
	a.SetValue("hello home")
	trie.Put("home", a)
	b := &router.Data{}
	b.SetValue("hello room")
	trie.Put("room", b)
	c := &router.Data{}
	c.SetValue("hello bedroom")
	trie.Put("bedroom", c)
	k := trie.Keys()
	assert.Equal(t, 3, len(k))
	//
	kwp := trie.KeysWithPrefix("hom")
	assert.Equal(t, 1, len(kwp))
	//
	ktm := trie.KeysThatMatch("")
	assert.Equal(t, 0, len(ktm))
	ktm1 := trie.KeysThatMatch("home")
	assert.Equal(t, 1, len(ktm1))
	ktm2 := trie.KeysThatMatch("room")
	assert.Equal(t, 1, len(ktm2))
	ktm3 := trie.KeysThatMatch(".o..")
	assert.Equal(t, 2, len(ktm3))
	ktm4 := trie.KeysThatMatch("....")
	assert.Equal(t, 2, len(ktm4))
	ktm5 := trie.KeysThatMatch(".")
	assert.Equal(t, 0, len(ktm5))

	lpo := trie.LongestPrefixOf("")
	assert.Equal(t, "", lpo)
	lpo1 := trie.LongestPrefixOf("home")
	assert.Equal(t, "home", lpo1)

}
