package discovery

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paulormart/assert"
)

func TestCatalogServicesPath(t *testing.T) {
	assert.Equal(t, "/v1/catalog/service", catalogServicePath)
}

type FakeServiceNoder struct {
	Response ServiceNodes
	Err      error
}

func (f *FakeServiceNoder) GetServiceNodes(addr, path, serviceID string) (ServiceNodes, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Response, nil
}

func TestGetServiceNodes(t *testing.T) {

	var tests = []struct {
		f            *FakeServiceNoder
		addr         string
		serviceID    string
		expectedResp ServiceNodes
		expectedErr  error
	}{
		{
			f: &FakeServiceNoder{
				Response: ServiceNodes{{Node: "foobar", Address: "10.1.10.12", ServiceID: "redis", ServiceName: "redis", ServiceAddress: "", ServicePort: 8000}},
				Err:      nil,
			},
			addr:         "localhost",
			serviceID:    "redis",
			expectedResp: ServiceNodes{{Node: "foobar", Address: "10.1.10.12", ServiceID: "redis", ServiceName: "redis", ServiceAddress: "", ServicePort: 8000}},
			expectedErr:  nil,
		},
		{
			f: &FakeServiceNoder{
				Response: nil,
				Err:      errors.New("TCP timeout"),
			},
			addr:         "localhost",
			serviceID:    "",
			expectedResp: nil,
			expectedErr:  errors.New("Error querying Consul API: TCP timeout"),
		},
	}
	for _, test := range tests {
		resp, err := GetServiceNodes(test.f, test.addr, test.serviceID)
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
	}
}

func TestServiceNoderInterface(t *testing.T) {
	var tests = []struct {
		hf           http.HandlerFunc
		serviceID    string
		expectedResp ServiceNodes
		expectedErr  error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[{"Node": "foobar","Address": "10.1.10.12","ServiceID": "redis","ServiceName": "redis","ServiceTags": null,"ServiceAddress": "","ServicePort": 8000}]`))
			},
			serviceID:    "redis",
			expectedResp: ServiceNodes{{Node: "foobar", Address: "10.1.10.12", ServiceID: "redis", ServiceName: "redis", ServiceAddress: "", ServicePort: 8000}},
			expectedErr:  nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{}`))
			},
			serviceID:    "",
			expectedResp: nil,
			expectedErr:  errors.New("Error querying Consul API: json: cannot unmarshal object into Go value of type discovery.ServiceNodes"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		resp, err := GetServiceNodes(&DefaultServiceNoder{}, ts.Listener.Addr().String(), test.serviceID)
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}

func TestGetServiceTargets(t *testing.T) {
	var tests = []struct {
		hf           http.HandlerFunc
		serviceID    string
		expectedResp []string
		expectedErr  error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[{"Node": "foobar","Address": "10.1.10.12","ServiceID": "redis","ServiceName": "redis","ServiceTags": null,"ServiceAddress": "","ServicePort": 8000}]`))
			},
			serviceID:    "redis",
			expectedResp: []string{"10.1.10.12:8000"},
			expectedErr:  nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			},
			serviceID:    "redis",
			expectedResp: nil,
			expectedErr:  errors.New("Error service: redis is not available in any of the nodes"),
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{}`))
			},
			serviceID:    "",
			expectedResp: nil,
			expectedErr:  errors.New("Error querying Consul API: json: cannot unmarshal object into Go value of type discovery.ServiceNodes"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		resp, err := GetServiceTargets(ts.Listener.Addr().String(), test.serviceID)
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}
