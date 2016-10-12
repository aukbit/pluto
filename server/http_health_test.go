package server

import (
	"log"
	"testing"
	"time"

	"github.com/paulormart/assert"
)

func TestHealthHTTP(t *testing.T) {
	s := NewServer(Addr(":8082"))
	go func(s Server) {
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}(s)
	defer s.Stop()
	time.Sleep(time.Second)
	h := s.Health()
	assert.Equal(t, "SERVING", h.Status.String())
}
