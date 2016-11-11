package discovery

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paulormart/assert"
)

type Fake struct {
	ID string
}

func TestRegisterHelper(t *testing.T) {
	var tests = []struct {
		hf          http.HandlerFunc
		expectedErr error
	}{
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ID": "something"}`))
			},
			expectedErr: nil,
		},
		{
			hf: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{error: description}`))
			},
			expectedErr: errors.New("{error: description}"),
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(test.hf))
		err := register(ts.Listener.Addr().String(), "", Fake{})
		assert.Equal(t, test.expectedErr, err)
		ts.Close()
	}
}

//
// func TestDoServiceUnregister(t *testing.T) {
//
// 	var tests = []struct {
// 		f           *FakeServicer
// 		addr        string
// 		serviceID   string
// 		expectedErr error
// 	}{
// 		{
// 			f: &FakeServicer{
// 				Response: Services{"redis": {ID: "redis", Service: "redis", Tags: []string{"db"}, Address: "", Port: 8080}},
// 				Err:      nil,
// 			},
// 			addr:        "localhost",
// 			serviceID:   "redis",
// 			expectedErr: nil,
// 		},
// 		{
// 			f: &FakeServicer{
// 				Response: nil,
// 				Err:      errors.New("TCP timeout"),
// 			},
// 			addr:        "localhost",
// 			serviceID:   "",
// 			expectedErr: errors.New("Error unregistering service Consul API: TCP timeout"),
// 		},
// 	}
// 	for _, test := range tests {
// 		err := DoServiceUnregister(test.f, test.addr, test.serviceID)
// 		assert.Equal(t, test.expectedErr, err)
// 	}
// }
//
// func TestServiceUnregisterInterface(t *testing.T) {
// 	var tests = []struct {
// 		hf          http.HandlerFunc
// 		serviceID   string
// 		expectedErr error
// 	}{
// 		{
// 			hf: func(w http.ResponseWriter, r *http.Request) {
// 				w.Header().Set("Content-Type", "application/json")
// 				w.WriteHeader(http.StatusOK)
// 				w.Write([]byte(`{"ID": "redis","Service": "redis","Tags": [],"Address": "","Port": 8000}`))
// 			},
// 			serviceID:   "redis",
// 			expectedErr: nil,
// 		},
// 		{
// 			hf: func(w http.ResponseWriter, r *http.Request) {
// 				w.Header().Set("Content-Type", "application/json")
// 				w.WriteHeader(http.StatusBadRequest)
// 				w.Write([]byte(`{error: description}`))
// 			},
// 			serviceID:   "",
// 			expectedErr: errors.New("Error unregistering service Consul API: {error: description}"),
// 		},
// 	}
// 	for _, test := range tests {
// 		ts := httptest.NewServer(http.HandlerFunc(test.hf))
// 		err := DoServiceUnregister(&DefaultServiceRegister{}, ts.Listener.Addr().String(), test.serviceID)
// 		assert.Equal(t, test.expectedErr, err)
// 		ts.Close()
// 	}
// }
