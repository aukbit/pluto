package discovery

import (
	"errors"
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
				Response: Services{},
				Err:      nil,
			},
			addr:         "localhost",
			expectedResp: Services{},
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
