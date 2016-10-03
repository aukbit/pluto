package frontend_test

import (
	"encoding/json"
	"io"
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
	"bitbucket.org/aukbit/pluto/examples/auth/backend/service"
	"bitbucket.org/aukbit/pluto/examples/auth/frontend/service"
	pb "bitbucket.org/aukbit/pluto/examples/auth/proto"
	pbu "bitbucket.org/aukbit/pluto/examples/user/proto"
	"bitbucket.org/aukbit/pluto/server"
	"github.com/paulormart/assert"
)

type Error struct {
	string
}

const URL = "http://localhost:8081"

var wg sync.WaitGroup

func RunBackend() {
	defer wg.Done()
	if err := backend.Run(); err != nil {
		log.Fatal(err)
	}
}

func RunFrontend() {
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
	pbu.RegisterUserServiceServer(grpcServer, &User{})

	// Define Pluto Server
	grpcSrv := server.NewServer(server.Addr(":65080"), server.GRPCServer(grpcServer))

	// Define Pluto Service
	s := pluto.NewService(pluto.Servers(grpcSrv))
	// Run service
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	if !testing.Short() {
		wg.Add(3)
		// mock user backend
		go MockUserBackend()
		time.Sleep(time.Millisecond * 100)
		go RunBackend()
		time.Sleep(time.Millisecond * 100)
		go RunFrontend()
	}
	result := m.Run()
	if !testing.Short() {
		wg.Wait()
	}
	os.Exit(result)
}

func TestAll(t *testing.T) {
	defer syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	var tests = []struct {
		Method string
		Path   string
		Body   io.Reader
		// BodyContains func(string) string
		Status int
	}{
		{
			Method: "POST",
			Path:   URL + "/authenticate",
			Body:   strings.NewReader(`{"email": "gopher@email.com", "password":"123456"}`),
			// BodyContains: func(id string) string { return `{"id":"` + id + `","name":"Gopher","email":"gopher@email.com"}` },
			Status: http.StatusOK,
		},
	}

	token := &pb.Token{}
	for _, test := range tests {

		r, err := http.NewRequest(test.Method, test.Path, test.Body)
		if err != nil {
			t.Fatal(err)
		}
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
		err = json.Unmarshal(actualBody, token)
		if err != nil {
			assert.Equal(t, response.Header.Get("Content-Type"), "application/json")
			assert.Equal(t, test.Status, response.StatusCode)
		} else {
			assert.Equal(t, response.Header.Get("Content-Type"), "application/json")
			assert.Equal(t, test.Status, response.StatusCode)
			// assert.Equal(t, test.BodyContains(user.Id), string(actualBody))
		}
		assert.Equal(t, true, len(token.Jwt) > 0)
	}

}

// backend user views
type User struct{}

func (s *User) ReadUser(ctx context.Context, nu *pbu.User) (*pbu.User, error) {
	// user object
	u := &pbu.User{Id: "123456", Name: "Gopher", Email: "gopher@email.com"}
	return u, nil
}

func (s *User) CreateUser(ctx context.Context, nu *pbu.NewUser) (*pbu.User, error) {
	return &pbu.User{}, nil
}

func (s *User) UpdateUser(ctx context.Context, nu *pbu.User) (*pbu.User, error) {
	return &pbu.User{}, nil
}

func (s *User) DeleteUser(ctx context.Context, nu *pbu.User) (*pbu.User, error) {
	return &pbu.User{}, nil
}

func (s *User) FilterUsers(ctx context.Context, nu *pbu.Filter) (*pbu.Users, error) {
	return &pbu.Users{}, nil
}
