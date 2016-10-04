package frontend_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/auth"
	pba "bitbucket.org/aukbit/pluto/auth/proto"
	"bitbucket.org/aukbit/pluto/examples/auth/backend/service"
	"bitbucket.org/aukbit/pluto/examples/auth/frontend/service"
	pbu "bitbucket.org/aukbit/pluto/examples/user/proto"
	"bitbucket.org/aukbit/pluto/reply"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
	"github.com/paulormart/assert"
)

type Error struct {
	string
}

const (
	USER_URL = "http://localhost:8080"
	AUTH_URL = "http://localhost:8081"
)

var wg sync.WaitGroup

func TestMain(m *testing.M) {
	if !testing.Short() {
		wg.Add(4)
		go MockUserBackend()
		time.Sleep(time.Millisecond * 50)
		go MockUserFrontend()
		time.Sleep(time.Millisecond * 50)
		go RunAuthBackend()
		time.Sleep(time.Millisecond * 50)
		go RunAuthFrontend()
		time.Sleep(time.Millisecond * 1000)
	}
	result := m.Run()
	if !testing.Short() {
		wg.Wait()
	}
	os.Exit(result)
}

func TestAll(t *testing.T) {
	defer syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	r, err := http.NewRequest("POST", AUTH_URL+"/authenticate", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	r.SetBasicAuth("firstgopher@email.com", "123456")
	// call handler
	response, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	actualBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	token := &pba.Token{}
	err = json.Unmarshal(actualBody, token)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, response.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, true, len(token.Jwt) > 0)

	// Test access to private resources
	r, err = http.NewRequest("POST", USER_URL+"/user",
		strings.NewReader(`{"name":"Gopher", "email": "secondgopher@email.com", "password":"123456"}`))
	if err != nil {
		t.Fatal(err)
	}
	// set Bearer authorization header
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.Jwt))
	// call handler
	response, err = http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	actualBody, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, response.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, `"ok"`, string(actualBody))
}

// Helper functions

func RunAuthBackend() {
	defer wg.Done()
	if err := backend.Run(); err != nil {
		log.Fatal(err)
	}
}

func RunAuthFrontend() {
	defer wg.Done()
	if err := frontend.Run(); err != nil {
		log.Fatal(err)
	}
}

func MockUserBackend() {
	defer wg.Done()
	// GRPC server
	// Define gRPC server and register
	grpcServer := grpc.NewServer()
	// Register grpc Server
	pbu.RegisterUserServiceServer(grpcServer, &MockUser{})
	// Define Pluto Server
	grpcSrv := server.NewServer(server.Addr(":65080"), server.GRPCServer(grpcServer))
	// Define Pluto Service
	s := pluto.NewService(pluto.Name("MockUserBackend"), pluto.Servers(grpcSrv))
	// Run service
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

func MockUserFrontend() {
	defer wg.Done()
	// Define handlers
	mux := router.NewMux()
	mux.POST("/user", PostHandler)
	mux.AddMiddleware(auth.MiddlewareBearerAuth())
	// define http server
	srv := server.NewServer(
		server.Name("api"),
		server.Addr(":8080"),
		server.Mux(mux))
	// define authentication client
	clt := auth.NewClientAuth("127.0.0.1:65081")
	// Define Pluto service
	s := pluto.NewService(
		pluto.Name("MockUserFrontend"),
		pluto.Servers(srv),
		pluto.Clients(clt))
	// Run service
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

// User frontend views
func PostHandler(w http.ResponseWriter, r *http.Request) {
	// ...
	// create user with data sent on user backend
	// check examples/user/frontend/views
	// ...
	reply.Json(w, r, http.StatusCreated, "ok")
}

// User backend views
type MockUser struct{}

func (s *MockUser) ReadUser(ctx context.Context, nu *pbu.User) (*pbu.User, error) {
	// ...
	return &pbu.User{}, nil
}

func (s *MockUser) CreateUser(ctx context.Context, nu *pbu.NewUser) (*pbu.User, error) {
	// ...
	return &pbu.User{}, nil
}

func (s *MockUser) UpdateUser(ctx context.Context, nu *pbu.User) (*pbu.User, error) {
	// ...
	return &pbu.User{}, nil
}

func (s *MockUser) DeleteUser(ctx context.Context, nu *pbu.User) (*pbu.User, error) {
	// ...
	return &pbu.User{}, nil
}

func (s *MockUser) FilterUsers(ctx context.Context, nu *pbu.Filter) (*pbu.Users, error) {
	// ...
	return &pbu.Users{}, nil
}

func (s *MockUser) VerifyUser(ctx context.Context, nu *pbu.Credentials) (*pbu.Verification, error) {
	// ...
	// verify user with data persisted
	// check examples/user/backend/views
	// ...
	return &pbu.Verification{IsValid: true}, nil
}
