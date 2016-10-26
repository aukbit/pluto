package discovery

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paulormart/assert"
)

type FakeServicer struct {
	Response Services
	Err      error
}

func (f *FakeServicer) GetServices(addr, path string) (Services, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Response, nil
}

func TestGetServices(t *testing.T) {

	var tests = []struct {
		f            *FakeServicer
		addr         string
		expectedResp Services
		expectedErr  error
	}{
		{
			f: &FakeServicer{
				Response: Services{"redis": {ID: "redis", Service: "redis", Tags: []string{"db"}, Address: "", Port: 8080}},
				Err:      nil,
			},
			addr:         "localhost",
			expectedResp: Services{"redis": {ID: "redis", Service: "redis", Tags: []string{"db"}, Address: "", Port: 8080}},
			expectedErr:  nil,
		},
		{
			f: &FakeServicer{
				Response: nil,
				Err:      errors.New("TCP timeout"),
			},
			addr:         "localhost",
			expectedResp: nil,
			expectedErr:  errors.New("Error querying Consul API: TCP timeout"),
		},
	}
	for _, test := range tests {
		resp, err := GetServices(test.f, test.addr)
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
	}
}

func TestServicer(t *testing.T) {
	var tests = []struct {
		hf           http.HandlerFunc
		expectedResp Services
		expectedErr  error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"redis": {"ID": "redis","Service": "redis","Tags": [],"Address": "","Port": 8000}}`))
			},
			expectedResp: Services{"redis": {"redis", "redis", []string{}, "", 8000}},
			expectedErr:  nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			},
			expectedResp: nil,
			expectedErr:  errors.New("Error querying Consul API: json: cannot unmarshal array into Go value of type discovery.Services"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		resp, err := GetServices(&DefaultServicer{}, ts.Listener.Addr().String())
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}
