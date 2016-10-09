package discovery

import "testing"

func TestAll(t *testing.T) {

	Nodes()
	err := Register()
	if err != nil {
		t.Error(err)
	}

}
