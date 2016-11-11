package discovery

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paulormart/assert"
)

func TestAgentSelfPath(t *testing.T) {
	assert.Equal(t, "/v1/status/leader", statusLeaderPath)
}

type FakeLeader struct {
	Response string
	Err      error
}

func (f *FakeLeader) Status(addr, path string) (string, error) {
	if f.Err != nil {
		return "", f.Err
	}
	return f.Response, nil
}

func TestStatus(t *testing.T) {

	var tests = []struct {
		f            *FakeLeader
		addr         string
		expectedResp string
		expectedErr  error
	}{
		{
			f: &FakeLeader{
				Response: "10.1.10.12:8300",
				Err:      nil,
			},
			addr:         "localhost",
			expectedResp: "10.1.10.12:8300",
			expectedErr:  nil,
		},
		{
			f: &FakeLeader{
				Response: "",
				Err:      errors.New("TCP timeout"),
			},
			addr:         "localhost",
			expectedResp: "",
			expectedErr:  errors.New("Error querying Consul API: TCP timeout"),
		},
	}
	for _, test := range tests {
		resp, err := GetStatus(test.f, test.addr)
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
	}
}

func TestLeaderInterface(t *testing.T) {
	var tests = []struct {
		hf           http.HandlerFunc
		expectedResp string
		expectedErr  error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`"10.1.10.12:8300"`))
			},
			expectedResp: "10.1.10.12:8300",
			expectedErr:  nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`""`))
			},
			expectedResp: "",
			expectedErr:  nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			},
			expectedResp: "",
			expectedErr:  errors.New("Error querying Consul API: json: cannot unmarshal array into Go value of type string"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		resp, err := GetStatus(&DefaultLeader{}, ts.Listener.Addr().String())
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}

func TestAvailability(t *testing.T) {
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
		resp, err := isAvailable(ts.Listener.Addr().String())
		assert.Equal(t, test.expectedResp, resp)
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}
