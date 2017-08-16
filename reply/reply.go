package reply

import (
	"encoding/json"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func Json(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Jsonpb(w http.ResponseWriter, r *http.Request, status int, m *jsonpb.Marshaler, pb proto.Message) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := m.Marshal(w, pb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
