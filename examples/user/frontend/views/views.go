package frontend

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/client"
	pb "github.com/aukbit/pluto/examples/user/proto"
	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server/router"
	"github.com/golang/protobuf/jsonpb"
	"github.com/google/uuid"
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
	c, ok := ctx.Value(pluto.Key("pluto")).(*pluto.Service).Client("user")
	if !ok {
		return &router.HandlerErr{
			Error:   errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	// dial
	i, err := c.Dial(client.Timeout(2 * time.Second))
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	defer c.Close()
	// make call
	user, err := i.(pb.UserServiceClient).CreateUser(ctx, nu)
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

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
	id := ctx.Value("id").(string)
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
	c, ok := ctx.Value(pluto.Key("pluto")).(*pluto.Service).Client("user")
	if !ok {
		return &router.HandlerErr{
			Error:   errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	// dial
	i, err := c.Dial()
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	defer c.Close()
	// make a call the backend service
	user, err = i.(pb.UserServiceClient).ReadUser(ctx, user)
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
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
	id := ctx.Value("id").(string)
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
	// unmarshal body
	if err = jsonpb.Unmarshal(r.Body, user); err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	// get gRPC client from service
	c, ok := ctx.Value(pluto.Key("pluto")).(*pluto.Service).Client("user")
	if !ok {
		return &router.HandlerErr{
			Error:   errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	// dial
	i, err := c.Dial()
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	defer c.Close()
	// make a call the backend service
	user, err = i.(pb.UserServiceClient).UpdateUser(ctx, user)
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	fmt.Printf("$$$$$$$$$ %v\n", user)
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
	id := ctx.Value("id").(string)
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
	c, ok := ctx.Value(pluto.Key("pluto")).(*pluto.Service).Client("user")
	if !ok {
		return &router.HandlerErr{
			Error:   errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	// dial
	i, err := c.Dial()
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	defer c.Close()
	// make a call the backend service
	user, err = i.(pb.UserServiceClient).DeleteUser(ctx, user)
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
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
	c, ok := ctx.Value(pluto.Key("pluto")).(*pluto.Service).Client("user")
	if !ok {
		return &router.HandlerErr{
			Error:   errClientUserNotAvailable,
			Message: errClientUserNotAvailable.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	// dial
	i, err := c.Dial()
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	defer c.Close()
	// make a call the backend service
	users, err := i.(pb.UserServiceClient).FilterUsers(ctx, filter)
	if err != nil {
		return &router.HandlerErr{
			Error:   err,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	reply.Json(w, r, http.StatusOK, users)
	return nil
}
