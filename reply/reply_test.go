package reply_test

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/paulormart/assert"
	"pluto/reply"
)

func TestAll(t *testing.T){

	// Table tests
	var tests = []struct{
		Fn func(w http.ResponseWriter, r *http.Request)
		S int // status code
		B string //body
	}{{
		Fn: func(w http.ResponseWriter, r *http.Request){
			reply.Json(w, r, http.StatusOK, "Hello World")
		},
		S: 200,
		B: `"Hello World"`,
	},{
		Fn: func(w http.ResponseWriter, r *http.Request){
			reply.Json(w, r, http.StatusNotFound, "not found")
		},
		S: http.StatusNotFound,
		B: `"not found"`,
	},{
		Fn: func(w http.ResponseWriter, r *http.Request){
			data := map[string]interface{}{"url": "golang.org"}
			reply.Json(w, r, http.StatusCreated, data)
		},
		S: http.StatusCreated,
		B: `{"url":"golang.org"}`,
	}}

	for _, test := range tests {
		w := httptest.NewRecorder()
		test.Fn(w, nil)

		assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
		assert.Equal(t, test.S, w.Code)
		assert.Equal(t, test.B, w.Body.String())
	}



}
