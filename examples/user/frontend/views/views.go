package frontend

import (
	"net/http"

	"bitbucket.org/aukbit/pluto"
	pb "bitbucket.org/aukbit/pluto/examples/user/proto"
	"bitbucket.org/aukbit/pluto/reply"
	"github.com/golang/protobuf/jsonpb"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// new user
	newUser := &pb.NewUser{}
	if err := jsonpb.Unmarshal(r.Body, newUser); err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	// get service from context by service name
	ctx := r.Context()
	s := ctx.Value("pluto_frontend")
	// get gRPC client from service
	c := s.(pluto.Service).Client("client_user")
	// make a call the backend service
	user, err := c.Call().(pb.UserServiceClient).CreateUser(ctx, newUser)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusCreated, user)
}

func GetHandlerDetail(w http.ResponseWriter, r *http.Request) {
	// get id from context
	ctx := r.Context()
	id := ctx.Value("id").(string)
	// set proto user
	user := &pb.User{Id: id}
	// get service from context by service name
	s := ctx.Value("pluto_frontend")
	// get gRPC client from service
	c := s.(pluto.Service).Client("client_user")
	// make a call the backend service
	user, err := c.Call().(pb.UserServiceClient).ReadUser(ctx, user)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, user)
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	// get id from context
	ctx := r.Context()
	id := ctx.Value("id").(string)
	// set proto user
	user := &pb.User{Id: id}
	// unmarshal body
	if err := jsonpb.Unmarshal(r.Body, user); err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	// get service from context by service name
	s := ctx.Value("pluto_frontend")
	// get gRPC client from service
	c := s.(pluto.Service).Client("client_user")
	// make a call the backend service
	user, err := c.Call().(pb.UserServiceClient).UpdateUser(ctx, user)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, user)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// get id from context
	ctx := r.Context()
	id := ctx.Value("id").(string)
	// set proto user
	user := &pb.User{Id: id}
	// get service from context by service name
	s := ctx.Value("pluto_frontend")
	// get gRPC client from service
	c := s.(pluto.Service).Client("client_user")
	// make a call the backend service
	user, err := c.Call().(pb.UserServiceClient).DeleteUser(ctx, user)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, user)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	// get parameters
	n := r.URL.Query().Get("name")
	// set proto filter
	filter := &pb.Filter{Name: n}
	// get context
	ctx := r.Context()
	// get service from context by service name
	s := ctx.Value("pluto_frontend")
	// get gRPC client from service
	c := s.(pluto.Service).Client("client_user")
	// make a call the backend service
	users, err := c.Call().(pb.UserServiceClient).FilterUsers(ctx, filter)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, users)
}
