package auth

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

	"go.uber.org/zap"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/auth"
	pba "github.com/aukbit/pluto/auth/proto"
	"github.com/aukbit/pluto/examples/auth/backend/service"
	"github.com/aukbit/pluto/examples/auth/frontend/service"
	pbu "github.com/aukbit/pluto/examples/user/proto"
	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server"
	"github.com/aukbit/pluto/server/router"
	"github.com/paulormart/assert"
)

type Error struct {
	string
}

const (
	USER_URL = "http://localhost:8088"
	AUTH_URL = "http://localhost:8089"
)

var wg sync.WaitGroup

func TestMain(m *testing.M) {
	if !testing.Short() {
		wg.Add(4)
		go MockUserBackend()
		time.Sleep(time.Millisecond * 500)
		go MockUserFrontend()
		time.Sleep(time.Millisecond * 500)
		go RunAuthBackend()
		time.Sleep(time.Millisecond * 500)
		go RunAuthFrontend()
		time.Sleep(time.Millisecond * 1000)
	}
	result := m.Run()
	if !testing.Short() {
		wg.Wait()
	}
	os.Exit(result)
}

func TestExampleAuth(t *testing.T) {
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

// type A func(s *grpc.Server, srv pbu.UserServiceServer)

func MockUserBackend() {
	defer wg.Done()
	// Define Pluto Server
	grpcSrv := server.New(
		server.Addr(":65080"),
		server.GRPCRegister(func(g *grpc.Server) {
			pbu.RegisterUserServiceServer(g, &MockUser{})
		}),
	)
	// logger
	logger, _ := zap.NewDevelopment()
	// Define Pluto Service
	s := pluto.New(
		pluto.Name("MockUserBackend"),
		pluto.Servers(grpcSrv),
		pluto.Logger(logger),
		pluto.HealthAddr(":9094"),
	)
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
	// define http server
	srv := server.New(
		server.Name("user_api"),
		server.Addr(":8088"),
		server.Mux(mux),
		server.Middlewares(auth.MiddlewareBearerAuth()),
	)
	// define authentication client
	clt := auth.NewClientAuth("127.0.0.1:65081")
	// Logger
	logger, _ := zap.NewDevelopment()
	// Define Pluto service
	s := pluto.New(
		pluto.Name("MockUserFrontend"),
		pluto.Servers(srv),
		pluto.Clients(clt),
		pluto.Logger(logger),
		pluto.HealthAddr(":9095"),
	)
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
