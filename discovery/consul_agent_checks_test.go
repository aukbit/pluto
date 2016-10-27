package discovery

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paulormart/assert"
)

type FakeChecker struct {
	Response Checks
	Err      error
}

func (f *FakeChecker) GetChecks(addr, path string) (Checks, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Response, nil
}

func TestGetChecks(t *testing.T) {

	var tests = []struct {
		f            *FakeChecker
		addr         string
		expectedResp Checks
		expectedErr  error
	}{
		{
			f: &FakeChecker{
				Response: Checks{"service:redis": {Node: "foobar", CheckID: "service:redis", Name: "Service 'redis' check", Status: "passing", Notes: "", ServiceID: "redis"}},
				Err:      nil,
			},
			addr:         "localhost",
			expectedResp: Checks{"service:redis": {Node: "foobar", CheckID: "service:redis", Name: "Service 'redis' check", Status: "passing", Notes: "", ServiceID: "redis"}},
			expectedErr:  nil,
		},
		{
			f: &FakeChecker{
				Response: nil,
				Err:      errors.New("TCP timeout"),
			},
			addr:         "localhost",
			expectedResp: nil,
			expectedErr:  errors.New("Error querying Consul API: TCP timeout"),
		},
	}
	for _, test := range tests {
		resp, err := GetChecks(test.f, test.addr)
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
	}
}

func TestCheckerInterface(t *testing.T) {
	var tests = []struct {
		hf           http.HandlerFunc
		expectedResp Checks
		expectedErr  error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"service:redis": {"Node": "foobar","CheckID": "service:redis","Name": "Service 'redis' check","Status": "passing","Notes": "","Output": "","ServiceID": "redis","ServiceName": "redis"}}`))
			},
			expectedResp: Checks{"service:redis": {Node: "foobar", CheckID: "service:redis", Name: "Service 'redis' check", Status: "passing", Notes: "", ServiceID: "redis"}},
			expectedErr:  nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			},
			expectedResp: nil,
			expectedErr:  errors.New("Error querying Consul API: json: cannot unmarshal array into Go value of type discovery.Checks"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		resp, err := GetChecks(&DefaultChecker{}, ts.Listener.Addr().String())
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}

func (f *FakeChecker) Register(addr, path string, s *Check) error {
	if f.Err != nil {
		return f.Err
	}
	return nil
}

func (f *FakeChecker) Unregister(addr, path, checkID string) error {
	if f.Err != nil {
		return f.Err
	}
	return nil
}

func TestDoCheckRegister(t *testing.T) {

	var tests = []struct {
		f           *FakeChecker
		addr        string
		check       *Check
		expectedErr error
	}{
		{
			f: &FakeChecker{
				Response: Checks{"service:redis": {Node: "foobar", CheckID: "service:redis", Name: "Service 'redis' check", Status: "passing", Notes: "", ServiceID: "redis"}},
				Err:      nil,
			},
			addr:        "localhost",
			check:       &Check{ID: "service:redis", Name: "Service 'redis' check", Status: "passing", Notes: "", ServiceID: "redis"},
			expectedErr: nil,
		},
		{
			f: &FakeChecker{
				Response: nil,
				Err:      errors.New("TCP timeout"),
			},
			addr:        "localhost",
			check:       nil,
			expectedErr: errors.New("Error registering check Consul API: TCP timeout"),
		},
	}
	for _, test := range tests {
		err := DoCheckRegister(test.f, test.addr, test.check)
		assert.Equal(t, test.expectedErr, err)
	}
}

func TestCheckRegisterInterface(t *testing.T) {
	var tests = []struct {
		hf          http.HandlerFunc
		check       *Check
		expectedErr error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"Node": "foobar","CheckID": "service:redis","Name": "Service 'redis' check","Status": "passing","Notes": "","Output": "","ServiceID": "redis","ServiceName": "redis"}`))
			},
			check:       &Check{ID: "service:redis", Name: "Service 'redis' check", Status: "passing", Notes: "", ServiceID: "redis"},
			expectedErr: nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{error: description}`))
			},
			check:       nil,
			expectedErr: errors.New("Error registering check Consul API: {error: description}"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		err := DoCheckRegister(&DefaultCheckRegister{}, ts.Listener.Addr().String(), test.check)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}

func TestDoCheckUnregister(t *testing.T) {

	var tests = []struct {
		f           *FakeChecker
		addr        string
		checkID     string
		expectedErr error
	}{
		{
			f: &FakeChecker{
				Response: Checks{"service:redis": {Node: "foobar", CheckID: "service:redis", Name: "Service 'redis' check", Status: "passing", Notes: "", ServiceID: "redis"}},
				Err:      nil,
			},
			addr:        "localhost",
			checkID:     "redis",
			expectedErr: nil,
		},
		{
			f: &FakeChecker{
				Response: nil,
				Err:      errors.New("TCP timeout"),
			},
			addr:        "localhost",
			checkID:     "",
			expectedErr: errors.New("Error unregistering check Consul API: TCP timeout"),
		},
	}
	for _, test := range tests {
		err := DoCheckUnregister(test.f, test.addr, test.checkID)
		assert.Equal(t, test.expectedErr, err)
	}
}

func TestCheckUnregisterInterface(t *testing.T) {
	var tests = []struct {
		hf          http.HandlerFunc
		checkID     string
		expectedErr error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"Node": "foobar","CheckID": "service:redis","Name": "Service 'redis' check","Status": "passing","Notes": "","Output": "","ServiceID": "redis","ServiceName": "redis"}`))
			},
			checkID:     "redis",
			expectedErr: nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{error: description}`))
			},
			checkID:     "",
			expectedErr: errors.New("Error unregistering check Consul API: {error: description}"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		err := DoCheckUnregister(&DefaultCheckRegister{}, ts.Listener.Addr().String(), test.checkID)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}
