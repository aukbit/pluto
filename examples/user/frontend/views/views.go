package frontend

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/client"
	pb "github.com/aukbit/pluto/examples/user/proto"
	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server/router"
	"github.com/golang/protobuf/jsonpb"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	errClientUserNotAvailable = errors.New("Client user not available")
)

// PostHandler ...
func PostHandler(w http.ResponseWriter, r *http.Request) *router.Err {
	// get context
	ctx := r.Context()
	// new user
	nu := &pb.NewUser{}
	if err := json.NewDecoder(r.Body).Decode(nu); err != nil {
		return &router.Err{
			Err:     err,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer r.Body.Close()
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return &router.Err{
			Err:     errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	// dial
	conn, err := c.Dial(client.Timeout(2 * time.Second))
	if err != nil {
		return &router.Err{
			Err:     err,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer conn.Close()
	// make call
	user, err := c.Stub(conn).(pb.UserServiceClient).CreateUser(ctx, nu)
	if err != nil {
		return &router.Err{
			Err:     err,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	log.Ctx(ctx).Info().Msg(fmt.Sprintf("POST user %s created", user.Id))
	// set header location
	w.Header().Set("Location", r.URL.Path+"/"+user.Id)
	reply.Jsonpb(w, r, http.StatusCreated, &jsonpb.Marshaler{}, user)
	return nil
}

// GetHandlerDetail ...
func GetHandlerDetail(w http.ResponseWriter, r *http.Request) *router.Err {
	// get context
	ctx := r.Context()
	// get id context
	id := router.FromContext(ctx, "id")
	validID, err := uuid.Parse(id)
	if err != nil {
		return &router.Err{
			Err:     fmt.Errorf("Id %v not found", id),
			Message: fmt.Errorf("Id %v not found", id).Error(),
			Status:  http.StatusNotFound,
		}
	}
	// set proto user
	user := &pb.User{Id: validID.String()}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return &router.Err{
			Err:     errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	// dial
	conn, err := c.Dial()
	if err != nil {
		return &router.Err{
			Err:     err,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer conn.Close()
	// make a call the backend service
	user, err = c.Stub(conn).(pb.UserServiceClient).ReadUser(ctx, user)
	if err != nil {
		return &router.Err{
			Err:     err,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	log.Ctx(r.Context()).Info().Msg(fmt.Sprintf("GET user %s", user.Id))
	// set header location
	w.Header().Add("Location", r.URL.Path)
	reply.Jsonpb(w, r, http.StatusOK, &jsonpb.Marshaler{}, user)
	return nil
}

// PutHandler ...
func PutHandler(w http.ResponseWriter, r *http.Request) *router.Err {
	// get context
	ctx := r.Context()
	// get id context
	id := router.FromContext(ctx, "id")
	validID, err := uuid.Parse(id)
	if err != nil {
		return &router.Err{
			Err:    fmt.Errorf("Id %v not found", id),
			Status: http.StatusNotFound,
		}
	}
	// set proto user
	user := &pb.User{Id: validID.String()}
	// unmarshal body
	if err = jsonpb.Unmarshal(r.Body, user); err != nil {
		return &router.Err{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return &router.Err{
			Err:    errClientUserNotAvailable,
			Status: http.StatusInternalServerError,
		}
	}
	// dial
	conn, err := c.Dial()
	if err != nil {
		return &router.Err{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}
	defer conn.Close()
	// make a call the backend service
	user, err = c.Stub(conn).(pb.UserServiceClient).UpdateUser(ctx, user)
	if err != nil {
		return &router.Err{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}
	log.Ctx(r.Context()).Info().Msg(fmt.Sprintf("PUT user %s updated", user.Id))
	// set header location
	w.Header().Set("Location", r.URL.Path)
	reply.Jsonpb(w, r, http.StatusOK, &jsonpb.Marshaler{}, user)
	return nil
}

// DeleteHandler ...
func DeleteHandler(w http.ResponseWriter, r *http.Request) *router.Err {
	// get context
	ctx := r.Context()
	// get id context
	id := router.FromContext(ctx, "id")
	validID, err := uuid.Parse(id)
	if err != nil {
		return &router.Err{
			Err:     fmt.Errorf("Id %v not found", id),
			Message: fmt.Errorf("Id %v not found", id).Error(),
			Status:  http.StatusNotFound,
		}
	}
	// set proto user
	user := &pb.User{Id: validID.String()}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return &router.Err{
			Err:     errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	// dial
	conn, err := c.Dial()
	if err != nil {
		return &router.Err{
			Err:     err,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer conn.Close()
	// make a call the backend service
	user, err = c.Stub(conn).(pb.UserServiceClient).DeleteUser(ctx, user)
	if err != nil {
		return &router.Err{
			Err:     err,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	log.Ctx(r.Context()).Info().Msg(fmt.Sprintf("DELETE user %s deleted", user.Id))
	reply.Jsonpb(w, r, http.StatusOK, &jsonpb.Marshaler{}, user)
	return nil
}

// GetHandler ...
func GetHandler(w http.ResponseWriter, r *http.Request) *router.Err {
	// get context
	ctx := r.Context()
	// get parameters
	n := r.URL.Query().Get("name")
	// set proto filter
	filter := &pb.Filter{Name: n}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return &router.Err{
			Err:     errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	// dial
	conn, err := c.Dial()
	if err != nil {
		return &router.Err{
			Err:     err,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer conn.Close()
	// make a call the backend service
	users, err := c.Stub(conn).(pb.UserServiceClient).FilterUsers(ctx, filter)
	if err != nil {
		return &router.Err{
			Err:     err,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	log.Ctx(r.Context()).Info().Msg(fmt.Sprintf("GET users %v", users))
	reply.Jsonpb(w, r, http.StatusOK, &jsonpb.Marshaler{}, users)
	return nil
}

// GetStreamHandler ...
func GetStreamHandler(w http.ResponseWriter, r *http.Request) *router.Err {
	// get context
	ctx := r.Context()
	// get parameters
	n := r.URL.Query().Get("name")
	// set proto filter
	filter := &pb.Filter{Name: n}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return &router.Err{
			Err:     errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	// dial
	conn, err := c.Dial()
	if err != nil {
		return &router.Err{
			Err:     err,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer conn.Close()
	// make call
	stream, err := c.Stub(conn).(pb.UserServiceClient).StreamUsers(ctx, filter)
	if err != nil {
		return &router.Err{
			Err:     err,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	users := &pb.Users{}
	for {
		u, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &router.Err{
				Err:     fmt.Errorf("%v.StreamUsers(_) = _, %v", c.Stub(conn), err),
				Message: fmt.Errorf("%v.StreamUsers(_) = _, %v", c.Stub(conn), err).Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		users.Data = append(users.Data, u)
	}
	log.Ctx(r.Context()).Info().Msg(fmt.Sprintf("GET Stream users %v", users))
	reply.Jsonpb(w, r, http.StatusOK, &jsonpb.Marshaler{}, users)
	return nil
}
