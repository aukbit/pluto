package discovery

import (
	"assert"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCatalogServicesPath(t *testing.T) {
	assert.Equal(t, "/v1/catalog/service", catalogServicePath)
}

type FakeNodeServicer struct {
	Response NodeServices
	Err      error
}

func (f *FakeNodeServicer) GetNodeServices(addr, path string) (NodeServices, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Response, nil
}

func TestGetNodeServices(t *testing.T) {

	var tests = []struct {
		f            *FakeNodeServicer
		addr         string
		expectedResp NodeServices
		expectedErr  error
	}{
		{
			f: &FakeNodeServicer{
				Response: NodeServices{{Node: "foobar", Address: "10.1.10.12", ServiceID: "redis", ServiceName: "redis", ServiceAddress: "", ServicePort: 8000}},
				Err:      nil,
			},
			addr:         "localhost",
			expectedResp: NodeServices{{Node: "foobar", Address: "10.1.10.12", ServiceID: "redis", ServiceName: "redis", ServiceAddress: "", ServicePort: 8000}},
			expectedErr:  nil,
		},
		{
			f: &FakeNodeServicer{
				Response: nil,
				Err:      errors.New("TCP timeout"),
			},
			addr:         "localhost",
			expectedResp: nil,
			expectedErr:  errors.New("Error querying Consul API: TCP timeout"),
		},
	}
	for _, test := range tests {
		resp, err := GetNodeServices(test.f, test.addr)
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
	}
}

func TestNodeServicerInterface(t *testing.T) {
	var tests = []struct {
		hf           http.HandlerFunc
		expectedResp NodeServices
		expectedErr  error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[{"Node": "foobar","Address": "10.1.10.12","ServiceID": "redis","ServiceName": "redis","ServiceTags": null,"ServiceAddress": "","ServicePort": 8000}]`))
			},
			expectedResp: NodeServices{{Node: "foobar", Address: "10.1.10.12", ServiceID: "redis", ServiceName: "redis", ServiceAddress: "", ServicePort: 8000}},
			expectedErr:  nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{}`))
			},
			expectedResp: nil,
			expectedErr:  errors.New("Error querying Consul API: json: cannot unmarshal object into Go value of type discovery.NodeServices"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		resp, err := GetNodeServices(&DefaultNodeServicer{}, ts.Listener.Addr().String())
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}
