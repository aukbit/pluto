package user

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/aukbit/pluto/examples/user/backend/service"
	"github.com/aukbit/pluto/examples/user/frontend/service"
	pb "github.com/aukbit/pluto/examples/user/proto"
	"github.com/paulormart/assert"
)

type Error struct {
	string
}

const URL = "http://localhost:8087"

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

func TestMain(m *testing.M) {
	if !testing.Short() {
		wg.Add(2)
		go RunBackend()
		// time.Sleep(time.Millisecond * 2000)
		go RunFrontend()
		time.Sleep(2 * time.Second)
	}
	result := m.Run()
	if !testing.Short() {
		// wg.Wait()
	}
	os.Exit(result)
}

func TestExampleUser(t *testing.T) {
	// defer syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	user := &pb.User{}

	var tests = []struct {
		Method         string
		Path           func(string) string
		Body           io.Reader
		Response       func(string) *pb.User
		ResponseHeader func(string) *http.Header
		// ResponseError func(string) *pb.User
		Status int
	}{
		{
			Method: "POST",
			Path:   func(id string) string { return URL + "/user" },
			Body:   strings.NewReader(`{"name":"Gopher", "email": "gopher@email.com", "password":"123456"}`),
			Response: func(id string) *pb.User {
				return &pb.User{
					Id:    id,
					Name:  "Gopher",
					Email: "gopher@email.com",
				}
			},
			ResponseHeader: func(id string) *http.Header {
				h := &http.Header{}
				h.Set("Content-Type", "application/json")
				h.Set("Location", "/user/"+id)
				return h
			},
			Status: http.StatusCreated,
		},
		// {
		// 	Method: "GET",
		// 	Path:   func(id string) string { return URL + "/user/" + id },
		// 	Response: func(id string) *pb.User {
		// 		return &pb.User{
		// 			Id:    id,
		// 			Name:  "Gopher",
		// 			Email: "gopher@email.com",
		// 		}
		// 	},
		// 	ResponseHeader: func(id string) *http.Header {
		// 		h := &http.Header{}
		// 		h.Set("Content-Type", "application/json")
		// 		h.Set("Location", "/user/"+id)
		// 		return h
		// 	},
		// 	Status: http.StatusOK,
		// },
		// {
		// 	Method: "GET",
		// 	Path:   func(id string) string { return URL + "/user/abc" },
		// 	Response: func(id string) *pb.User {
		// 		return &pb.User{}
		// 	},
		// 	Status: http.StatusNotFound,
		// },
		// {
		// 	Method: "PUT",
		// 	Path:   func(id string) string { return URL + "/user/" + id },
		// 	Body:   strings.NewReader(`{"name":"Super Gopher house"}`),
		// 	Response: func(id string) *pb.User {
		// 		return &pb.User{
		// 			Id:    id,
		// 			Name:  "Super Gopher house",
		// 			Email: "gopher@email.com",
		// 		}
		// 	},
		// 	ResponseHeader: func(id string) *http.Header {
		// 		h := &http.Header{}
		// 		h.Set("Content-Type", "application/json")
		// 		h.Set("Location", "/user/"+id)
		// 		return h
		// 	},
		// 	Status: http.StatusOK,
		// },
		// {
		// 	Method: "PUT",
		// 	Path:   func(id string) string { return URL + "/user/abc" },
		// 	Body:   strings.NewReader(`{"name":"Super Gopher house"}`),
		// 	Response: func(id string) *pb.User {
		// 		return &pb.User{}
		// 	},
		// 	Status: http.StatusNotFound,
		// },
		// {
		// 	Method: "DELETE",
		// 	Path:   func(id string) string { return URL + "/user/" + id },
		// 	Response: func(id string) *pb.User {
		// 		return &pb.User{}
		// 	},
		// 	ResponseHeader: func(id string) *http.Header {
		// 		h := &http.Header{}
		// 		h.Set("Content-Type", "application/json")
		// 		return h
		// 	},
		// 	Status: http.StatusOK,
		// },
		// {
		// 	Method: "DELETE",
		// 	Path:   func(id string) string { return URL + "/user/abc" },
		// 	Response: func(id string) *pb.User {
		// 		return &pb.User{}
		// 	},
		// 	Status: http.StatusNotFound,
		// },
	}

	for _, test := range tests {
		req, err := http.NewRequest(test.Method, test.Path(user.Id), test.Body)
		if err != nil {
			t.Fatal(err)
		}
		// call handler
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.Request.Method == "DELETE" {
			user = &pb.User{}
		}
		// decode body into user struct
		err = json.NewDecoder(resp.Body).Decode(user)
		if err != nil {
			assert.Equal(t, test.Status, resp.StatusCode)
		} else {
			assert.Equal(t, test.Status, resp.StatusCode)
			assert.Equal(t, test.ResponseHeader("").Get("Content-Type"), resp.Header.Get("Content-Type"))
			assert.Equal(t, test.ResponseHeader(user.Id).Get("Location"), resp.Header.Get("Location"))
			assert.Equal(t, test.Response(user.Id).String(), user.String())
		}
	}
}

func helperAddUser(name string) error {
	req, err := http.NewRequest("POST", URL+"/user", strings.NewReader(`{"name":"`+name+`", "email": "gopher@email.com", "password":"123456"}`))
	if err != nil {
		return err
	}
	// call handler
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func helperDeleteUser(id string) error {
	req, err := http.NewRequest("DELETE", URL+"/user/"+id, strings.NewReader(`{}`))
	if err != nil {
		return err
	}
	// call handler
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func TestListUsers(t *testing.T) {

	defer syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	err := helperAddUser("abc")
	if err != nil {
		t.Fatal(err)
	}

	err = helperAddUser("def")
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		Method         string
		Path           string
		Response       *pb.Users
		ResponseHeader func() *http.Header
		// ResponseError func(string) *pb.User
		Status int
	}{
		{
			Method: "GET",
			Path:   URL + "/user?name=abc",
			Response: &pb.Users{
				Data: []*pb.User{
					&pb.User{
						Name: "abc",
					},
				},
			},
			ResponseHeader: func() *http.Header {
				h := &http.Header{}
				h.Set("Content-Type", "application/json")
				return h
			},
			Status: http.StatusOK,
		},
		{
			Method: "GET",
			Path:   URL + "/stream?name=abc",
			Response: &pb.Users{
				Data: []*pb.User{
					&pb.User{
						Name: "abc",
					},
				},
			},
			ResponseHeader: func() *http.Header {
				h := &http.Header{}
				h.Set("Content-Type", "application/json")
				return h
			},
			Status: http.StatusOK,
		},
	}

	users := &pb.Users{}
	for _, test := range tests {
		req, err := http.NewRequest(test.Method, test.Path, strings.NewReader(`{}`))
		if err != nil {
			t.Fatal(err)
		}
		// call handler
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		// decode body into user struct
		users = &pb.Users{}
		err = json.NewDecoder(resp.Body).Decode(users)
		resp.Body.Close()
		// assertions
		assert.Equal(t, test.Status, resp.StatusCode)
		assert.Equal(t, test.ResponseHeader().Get("Content-Type"), resp.Header.Get("Content-Type"))
		assert.Equal(t, len(test.Response.Data), len(users.Data))
		// fmt.Println(users)
	}
	// teardown
	for _, u := range users.Data {
		err := helperDeleteUser(u.Id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
