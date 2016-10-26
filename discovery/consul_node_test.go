package discovery

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"assert"
)

type FakeNoder struct {
	Response Nodes
	Err      error
}

func (f *FakeNoder) GetNodes(addr, path string) (Nodes, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Response, nil
}

func TestGetNodes(t *testing.T) {

	var tests = []struct {
		f            *FakeNoder
		addr         string
		expectedResp Nodes
		expectedErr  error
	}{
		{
			f: &FakeNoder{
				Response: Nodes{{"baz", "10.1.10.11", map[string]string{"lan": "10.1.10.11", "wan": "10.1.10.11"}}},
				Err:      nil,
			},
			addr:         "localhost",
			expectedResp: Nodes{{"baz", "10.1.10.11", map[string]string{"lan": "10.1.10.11", "wan": "10.1.10.11"}}},
			expectedErr:  nil,
		},
		{
			f: &FakeNoder{
				Response: nil,
				Err:      errors.New("TCP timeout"),
			},
			addr:         "localhost",
			expectedResp: nil,
			expectedErr:  errors.New("Error querying Consul API: TCP timeout"),
		},
	}
	for _, test := range tests {
		resp, err := GetNodes(test.f, test.addr)
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
	}
}

func TestNoder(t *testing.T) {
	var tests = []struct {
		hf           http.HandlerFunc
		expectedResp Nodes
		expectedErr  error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[{"Node": "baz","Address": "10.1.10.11","TaggedAddresses": {"lan": "10.1.10.11","wan": "10.1.10.11"}},{"Node": "foobar","Address": "10.1.10.12","TaggedAddresses": {"lan": "10.1.10.11","wan": "10.1.10.12"}}]`))
			},
			expectedResp: Nodes{{"baz", "10.1.10.11", map[string]string{"lan": "10.1.10.11", "wan": "10.1.10.11"}}, {"foobar", "10.1.10.12", map[string]string{"lan": "10.1.10.11", "wan": "10.1.10.12"}}},
			expectedErr:  nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{}`))
			},
			expectedResp: nil,
			expectedErr:  errors.New("Error querying Consul API: json: cannot unmarshal object into Go value of type []discovery.Node"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		resp, err := GetNodes(&DefaultNoder{}, ts.Listener.Addr().String())
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}
