package views

import (
	"errors"
	"net/http"
	"time"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/client"
	pb "github.com/aukbit/pluto/examples/dist/user_bff/proto"
	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server/router"
	"github.com/golang/protobuf/jsonpb"
)

var (
	errClientUserNotAvailable = errors.New("Client user not available")
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// get context
	ctx := r.Context()
	// new user
	newUser := &pb.NewUser{}
	if err := jsonpb.Unmarshal(r.Body, newUser); err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("client")
	if !ok {
		reply.Json(w, r, http.StatusInternalServerError, errClientUserNotAvailable)
		return
	}
	// dial
	conn, err := c.Dial(client.Timeout(5 * time.Second))
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	defer conn.Close()
	// make a call the backend service
	user, err := c.Stub(conn).(pb.UserServiceClient).CreateUser(ctx, newUser)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusCreated, user)
}

func GetHandlerDetail(w http.ResponseWriter, r *http.Request) {
	// get context
	ctx := r.Context()
	// get id context
	id := router.FromContext(ctx, "id")
	// set proto user
	user := &pb.User{Id: id}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("client")
	if !ok {
		reply.Json(w, r, http.StatusInternalServerError, errClientUserNotAvailable)
		return
	}
	// dial
	conn, err := c.Dial(client.Timeout(5 * time.Second))
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	defer conn.Close()
	// make a call the backend service
	user, err = c.Stub(conn).(pb.UserServiceClient).ReadUser(ctx, user)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, user)
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	// get context
	ctx := r.Context()
	// get id context
	id := router.FromContext(ctx, "id")
	// set proto user
	user := &pb.User{Id: id}
	// unmarshal body
	if err := jsonpb.Unmarshal(r.Body, user); err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("client")
	if !ok {
		reply.Json(w, r, http.StatusInternalServerError, errClientUserNotAvailable.Error())
		return
	}
	// dial
	conn, err := c.Dial(client.Timeout(5 * time.Second))
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	defer conn.Close()
	// make a call the backend service
	user, err = c.Stub(conn).(pb.UserServiceClient).UpdateUser(ctx, user)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, user)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// get context
	ctx := r.Context()
	// get id context
	id := router.FromContext(ctx, "id")
	// set proto user
	user := &pb.User{Id: id}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("client")
	if !ok {
		reply.Json(w, r, http.StatusInternalServerError, errClientUserNotAvailable)
		return
	}
	// dial
	conn, err := c.Dial(client.Timeout(5 * time.Second))
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	defer conn.Close()
	// make a call the backend service
	user, err = c.Stub(conn).(pb.UserServiceClient).DeleteUser(ctx, user)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, user)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	// get context
	ctx := r.Context()
	// get parameters
	n := r.URL.Query().Get("name")
	// set proto filter
	filter := &pb.Filter{Name: n}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("client")
	if !ok {
		reply.Json(w, r, http.StatusInternalServerError, errClientUserNotAvailable.Error())
		return
	}
	// dial
	conn, err := c.Dial(client.Timeout(5 * time.Second))
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	defer conn.Close()
	// make a call the backend service
	users, err := c.Stub(conn).(pb.UserServiceClient).FilterUsers(ctx, filter)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, users)
}
