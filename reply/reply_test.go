package reply_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aukbit/pluto/v6/reply"
	"github.com/paulormart/assert"
)

func TestReply(t *testing.T) {

	// Table tests
	var tests = []struct {
		Fn func(w http.ResponseWriter, r *http.Request)
		S  int    // status code
		B  string //body
	}{{
		Fn: func(w http.ResponseWriter, r *http.Request) {
			reply.Json(w, r, http.StatusOK, "Hello World")
		},
		S: 200,
		B: "\"Hello World\"\n",
	}, {
		Fn: func(w http.ResponseWriter, r *http.Request) {
			reply.Json(w, r, http.StatusNotFound, "not found")
		},
		S: http.StatusNotFound,
		B: "\"not found\"\n",
	}, {
		Fn: func(w http.ResponseWriter, r *http.Request) {
			data := map[string]interface{}{"url": "golang.org"}
			reply.Json(w, r, http.StatusCreated, data)
		},
		S: http.StatusCreated,
		B: "{\"url\":\"golang.org\"}\n",
	}}

	for _, test := range tests {
		w := httptest.NewRecorder()
		test.Fn(w, nil)
		assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
		assert.Equal(t, test.S, w.Code)
		assert.Equal(t, test.B, w.Body.String())
	}

}
