package reply

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
)

func Json(w http.ResponseWriter, r *http.Request, status int, data interface{}){

	d, err := json.Marshal(data)
	if err != nil {
		//With(w, r, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(d); err != nil {
		log.Fatal(fmt.Sprintf("ERROR w.Write(d) %v", err))
	}


}
