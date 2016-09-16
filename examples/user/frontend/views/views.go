package frontend

import (
	"pluto/reply"
	"net/http"
	pb "pluto/examples/user/proto"
	"github.com/golang/protobuf/proto"
	//"io/ioutil"
	//"encoding/json"
	"log"
	"io/ioutil"
)

func GetHandler (w http.ResponseWriter, r *http.Request){
	reply.Json(w, r, http.StatusOK, "Hello World")
}
func PostHandler (w http.ResponseWriter, r *http.Request){

	// new user
	user := &pb.NewUser{}

	// read http body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
	}
	if err := proto.Unmarshal(body, user); err != nil {
		reply.Json(w, r, http.StatusNotAcceptable, err.Error())
	}
	log.Printf("TESTE %v", user)
	//
	//// TODO send to backend
	//
	//data, err := proto.Marshal(user)
	//if err != nil {
	//	reply.Json(w, r, http.StatusInternalServerError, err.Error())
	//}
	//
	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	reply.Json(w, r, http.StatusInternalServerError, err.Error())
	//}
	//data := make(map[string]string)
	//if err := json.Unmarshal(body, &data); err != nil {
	//	log.Println(err.Error())
	//}
	reply.Json(w, r, http.StatusCreated, user)
}
func GetHandlerDetail (w http.ResponseWriter, r *http.Request){
	reply.Json(w, r, http.StatusOK, "Hello World")
}
func PutHandler (w http.ResponseWriter, r *http.Request){
	reply.Json(w, r, http.StatusOK, "Hello World")
}
func DeleteHandler (w http.ResponseWriter, r *http.Request){
	reply.Json(w, r, http.StatusOK, "Hello World")
}