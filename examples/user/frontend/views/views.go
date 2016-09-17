package frontend

import (
	"pluto/reply"
	"net/http"
	pb "pluto/examples/user/proto"
	"github.com/golang/protobuf/jsonpb"
	"pluto"
)

func GetHandler (w http.ResponseWriter, r *http.Request){
	reply.Json(w, r, http.StatusOK, "Hello World")
}

func PostHandler (w http.ResponseWriter, r *http.Request){
	// new user
	newUser := &pb.NewUser{}
	if err := jsonpb.Unmarshal(r.Body, newUser); err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
	}
	// send to backend
	// get service from context by service name
	ctx := r.Context()
	s := ctx.Value("frontend.pluto")
	// get gRPC client
	c := s.(pluto.Service).Clients()["user.client"]
	user, err := c.Call().(pb.UserServiceClient).CreateUser(ctx, newUser)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
	}
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