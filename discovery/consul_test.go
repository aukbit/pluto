package discovery

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paulormart/assert"
)

func TestNewConsulDefault(t *testing.T) {
	cd := newConsulDefault()
	assert.Equal(t, "localhost:8500", cd.cfg.Addr)
	assert.Equal(t, "http://localhost:8500", cd.cfg.URL())
	assert.Equal(t, Services{}, cd.cfg.Services)
	assert.Equal(t, Checks{}, cd.cfg.Checks)
}

func TestIsAvailable(t *testing.T) {
	var tests = []struct {
		hf           http.HandlerFunc
		expectedResp bool
		expectedErr  error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`"10.1.10.12:8300"`))
			},
			expectedResp: true,
			expectedErr:  nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`""`))
			},
			expectedResp: false,
			expectedErr:  nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			},
			expectedResp: false,
			expectedErr:  errors.New("Error querying Consul API: json: cannot unmarshal array into Go value of type string"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		cd := newConsulDefault(Addr(ts.Listener.Addr().String()))
		resp, err := cd.IsAvailable()
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}

func TestService(t *testing.T) {
	var tests = []struct {
		hf           http.HandlerFunc
		expectedResp []string
		expectedErr  error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[{"Node": "foobar","Address": "10.1.10.12","ServiceID": "redis","ServiceName": "redis","ServiceTags": null,"ServiceAddress": "","ServicePort": 8000}]`))
			},
			expectedResp: []string{"10.1.10.12:8000"},
			expectedErr:  nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			},
			expectedResp: nil,
			expectedErr:  errors.New("Error service: redis is not available in any of the nodes"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		cd := newConsulDefault(Addr(ts.Listener.Addr().String()))
		resp, err := cd.Service("redis")
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}

func TestRegister(t *testing.T) {
	var tests = []struct {
		hf           http.HandlerFunc
		expectedResp Service
		expectedErr  error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ID": "redis","Name": "redis","Address": "","Port": 8000}`))
			},
			expectedErr: nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ID": "service:redis","Name": "Service 'redis' check","Status": "passing","Notes": "","Output": "","ServiceID": "redis","ServiceName": "redis"}`))
			},
			expectedErr: nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{}`))
			},
			expectedErr: nil,
			// expectedErr:  errors.New("Error service: redis is not available in any of the nodes"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		s := Service{
			ID:      "redis",
			Service: "redis",
			Port:    8000,
		}
		c := Check{
			ID:    "test_check",
			Name:  "TCP check",
			Notes: "Ensure the server is listening on the specific port",
			DeregisterCriticalServiceAfter: "1m",
			TCP:       ":60500",
			Interval:  "10s",
			Timeout:   "1s",
			ServiceID: "test",
		}
		cd := NewDiscovery(Addr(ts.Listener.Addr().String()))
		err := cd.Register(ServicesCfg(s), ChecksCfg(c))
		assert.Equal(t, test.expectedErr, err)
		err = cd.Unregister()
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}
