package frontend

import (
	"context"
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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	errClientUserNotAvailable = errors.New("Client user not available")
)

// PostHandler ...
func PostHandler(w http.ResponseWriter, r *http.Request) *router.HandlerErr {
	// get context
	ctx := r.Context()
	// new user
	nu := &pb.NewUser{}
	if err := json.NewDecoder(r.Body).Decode(nu); err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	defer r.Body.Close()
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return &router.HandlerErr{
			Error:   errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	// dial
	conn, err := c.Dial(client.Timeout(2 * time.Second))
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	defer conn.Close()
	// make call
	user, err := c.Stub(conn).(pb.UserServiceClient).CreateUser(ctx, nu)
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	zerolog.Ctx(ctx).Info().Msg(fmt.Sprintf("POST user %s created", user.Id))
	// set header location
	w.Header().Set("Location", r.URL.Path+"/"+user.Id)
	reply.Json(w, r, http.StatusCreated, user)
	return nil
}

// GetHandlerDetail ...
func GetHandlerDetail(w http.ResponseWriter, r *http.Request) *router.HandlerErr {
	// get context
	ctx := r.Context()
	// get id context
	id := router.FromContext(ctx, "id")
	validID, err := uuid.Parse(id)
	if err != nil {
		return &router.HandlerErr{
			Error:   fmt.Errorf("Id %v not found", id),
			Message: fmt.Errorf("Id %v not found", id).Error(),
			Code:    http.StatusNotFound,
		}
	}
	// set proto user
	user := &pb.User{Id: validID.String()}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return &router.HandlerErr{
			Error:   errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	// dial
	conn, err := c.Dial()
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	defer conn.Close()
	// make a call the backend service
	user, err = c.Stub(conn).(pb.UserServiceClient).ReadUser(ctx, user)
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	log.Ctx(r.Context()).Info().Msg(fmt.Sprintf("GET user %s", user.Id))
	// set header location
	w.Header().Add("Location", r.URL.Path)
	reply.Json(w, r, http.StatusOK, user)
	return nil
}

// PutHandler ...
func PutHandler(w http.ResponseWriter, r *http.Request) *router.HandlerErr {
	// get context
	ctx := r.Context()
	// get id context
	id := router.FromContext(ctx, "id")
	validID, err := uuid.Parse(id)
	if err != nil {
		return router.NewHandlerErr(
			fmt.Errorf("Id %v not found", id),
			http.StatusNotFound,
		)
	}
	// set proto user
	user := &pb.User{Id: validID.String()}
	// unmarshal body
	if err = jsonpb.Unmarshal(r.Body, user); err != nil {
		return router.NewHandlerErr(
			err,
			http.StatusInternalServerError,
		)
	}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return router.NewHandlerErr(
			errClientUserNotAvailable,
			http.StatusInternalServerError,
		)
	}
	// dial
	conn, err := c.Dial()
	if err != nil {
		return router.NewHandlerErr(
			err,
			http.StatusInternalServerError,
		)
	}
	defer conn.Close()
	// make a call the backend service
	user, err = c.Stub(conn).(pb.UserServiceClient).UpdateUser(ctx, user)
	if err != nil {
		return router.NewHandlerErr(
			err,
			http.StatusInternalServerError,
		)
	}
	log.Ctx(r.Context()).Info().Msg(fmt.Sprintf("PUT user %s updated", user.Id))
	// set header location
	w.Header().Set("Location", r.URL.Path)
	reply.Json(w, r, http.StatusOK, user)
	return nil
}

// DeleteHandler ...
func DeleteHandler(w http.ResponseWriter, r *http.Request) *router.HandlerErr {
	// get context
	ctx := r.Context()
	// get id context
	id := router.FromContext(ctx, "id")
	validID, err := uuid.Parse(id)
	if err != nil {
		return &router.HandlerErr{
			Error:   fmt.Errorf("Id %v not found", id),
			Message: fmt.Errorf("Id %v not found", id).Error(),
			Code:    http.StatusNotFound,
		}
	}
	// set proto user
	user := &pb.User{Id: validID.String()}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return &router.HandlerErr{
			Error:   errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	// dial
	conn, err := c.Dial()
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	defer conn.Close()
	// make a call the backend service
	user, err = c.Stub(conn).(pb.UserServiceClient).DeleteUser(ctx, user)
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	log.Ctx(r.Context()).Info().Msg(fmt.Sprintf("DELETE user %s deleted", user.Id))
	reply.Json(w, r, http.StatusOK, user)
	return nil
}

// GetHandler ...
func GetHandler(w http.ResponseWriter, r *http.Request) *router.HandlerErr {
	// get context
	ctx := r.Context()
	// get parameters
	n := r.URL.Query().Get("name")
	// set proto filter
	filter := &pb.Filter{Name: n}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return &router.HandlerErr{
			Error:   errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	// dial
	conn, err := c.Dial()
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	defer conn.Close()
	// make a call the backend service
	users, err := c.Stub(conn).(pb.UserServiceClient).FilterUsers(ctx, filter)
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	log.Ctx(r.Context()).Info().Msg(fmt.Sprintf("GET users %v", users))
	reply.Json(w, r, http.StatusOK, users)
	return nil
}

// GetStreamHandler ...
func GetStreamHandler(w http.ResponseWriter, r *http.Request) *router.HandlerErr {
	// get context
	ctx := r.Context()
	// get parameters
	n := r.URL.Query().Get("name")
	// set proto filter
	filter := &pb.Filter{Name: n}
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return &router.HandlerErr{
			Error:   errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	// dial
	conn, err := c.Dial()
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	defer conn.Close()
	// make call
	stream, err := c.Stub(conn).(pb.UserServiceClient).StreamUsers(context.Background(), filter)
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	users := &pb.Users{}
	for {
		u, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &router.HandlerErr{
				Error:   fmt.Errorf("%v.StreamUsers(_) = _, %v", c.Stub(conn), err),
				Message: fmt.Errorf("%v.StreamUsers(_) = _, %v", c.Stub(conn), err).Error(),
				Code:    http.StatusInternalServerError,
			}
		}
		users.Data = append(users.Data, u)
	}
	log.Ctx(r.Context()).Info().Msg(fmt.Sprintf("GET Stream users %v", users))
	reply.Json(w, r, http.StatusOK, users)
	return nil
}
