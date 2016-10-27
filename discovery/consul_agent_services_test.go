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

func TestServicerInterface(t *testing.T) {
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
			expectedResp: Services{"redis": {ID: "redis", Service: "redis", Tags: []string{}, Address: "", Port: 8000}},
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

func (f *FakeServicer) Register(addr, path string, s *Service) error {
	if f.Err != nil {
		return f.Err
	}
	return nil
}

func (f *FakeServicer) Unregister(addr, path, serviceID string) error {
	if f.Err != nil {
		return f.Err
	}
	return nil
}

func TestDoServiceRegister(t *testing.T) {

	var tests = []struct {
		f           *FakeServicer
		addr        string
		service     *Service
		expectedErr error
	}{
		{
			f: &FakeServicer{
				Response: Services{"redis": {ID: "redis", Service: "redis", Tags: []string{"db"}, Address: "", Port: 8080}},
				Err:      nil,
			},
			addr:        "localhost",
			service:     &Service{ID: "redis", Service: "redis", Tags: []string{"db"}, Address: "", Port: 8080},
			expectedErr: nil,
		},
		{
			f: &FakeServicer{
				Response: nil,
				Err:      errors.New("TCP timeout"),
			},
			addr:        "localhost",
			service:     nil,
			expectedErr: errors.New("Error registering service Consul API: TCP timeout"),
		},
	}
	for _, test := range tests {
		err := DoServiceRegister(test.f, test.addr, test.service)
		assert.Equal(t, test.expectedErr, err)
	}
}

func TestServiceRegisterInterface(t *testing.T) {
	var tests = []struct {
		hf          http.HandlerFunc
		service     *Service
		expectedErr error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ID": "redis","Service": "redis","Tags": [],"Address": "","Port": 8000}`))
			},
			service:     &Service{ID: "redis", Service: "redis", Tags: []string{"db"}, Address: "", Port: 8080},
			expectedErr: nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{error: description}`))
			},
			service:     nil,
			expectedErr: errors.New("Error registering service Consul API: {error: description}"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		err := DoServiceRegister(&DefaultServiceRegister{}, ts.Listener.Addr().String(), test.service)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}

func TestDoServiceUnregister(t *testing.T) {

	var tests = []struct {
		f           *FakeServicer
		addr        string
		serviceID   string
		expectedErr error
	}{
		{
			f: &FakeServicer{
				Response: Services{"redis": {ID: "redis", Service: "redis", Tags: []string{"db"}, Address: "", Port: 8080}},
				Err:      nil,
			},
			addr:        "localhost",
			serviceID:   "redis",
			expectedErr: nil,
		},
		{
			f: &FakeServicer{
				Response: nil,
				Err:      errors.New("TCP timeout"),
			},
			addr:        "localhost",
			serviceID:   "",
			expectedErr: errors.New("Error unregistering service Consul API: TCP timeout"),
		},
	}
	for _, test := range tests {
		err := DoServiceUnregister(test.f, test.addr, test.serviceID)
		assert.Equal(t, test.expectedErr, err)
	}
}

func TestServiceUnregisterInterface(t *testing.T) {
	var tests = []struct {
		hf          http.HandlerFunc
		serviceID   string
		expectedErr error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ID": "redis","Service": "redis","Tags": [],"Address": "","Port": 8000}`))
			},
			serviceID:   "redis",
			expectedErr: nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{error: description}`))
			},
			serviceID:   "",
			expectedErr: errors.New("Error unregistering service Consul API: {error: description}"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		err := DoServiceUnregister(&DefaultServiceRegister{}, ts.Listener.Addr().String(), test.serviceID)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}
