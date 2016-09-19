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
	// get service from context by service name
	ctx := r.Context()
	s := ctx.Value("pluto_frontend")
	// get gRPC client from service
	c := s.(pluto.Service).Clients()["client_user"]
	// make a call the backend service
	user, err := c.Call().(pb.UserServiceClient).CreateUser(ctx, newUser)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
	}
	reply.Json(w, r, http.StatusCreated, user)
}
func GetHandlerDetail (w http.ResponseWriter, r *http.Request){
	// get id from context
	ctx := r.Context()
	id := ctx.Value("id").(string)
	// set proto user
	user := &pb.User{Id: id}
	// get service from context by service name
	s := ctx.Value("pluto_frontend")
	// get gRPC client from service
	c := s.(pluto.Service).Clients()["client_user"]
	// make a call the backend service
	user, err := c.Call().(pb.UserServiceClient).ReadUser(ctx, user)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
	}
	reply.Json(w, r, http.StatusOK, user)
}
func PutHandler (w http.ResponseWriter, r *http.Request){
	reply.Json(w, r, http.StatusOK, "Hello World")
}
func DeleteHandler (w http.ResponseWriter, r *http.Request){
	reply.Json(w, r, http.StatusOK, "Hello World")
}