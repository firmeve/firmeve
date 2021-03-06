package http

import (
	"errors"
	"github.com/firmeve/firmeve/kernel/contract"
	render2 "github.com/firmeve/firmeve/render"
	testing2 "github.com/firmeve/firmeve/testing"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

var configPath = "../testdata/config/config.yaml"

//
type MockResponseWriter struct {
	mock.Mock
	Bytes      []byte
	StatusCode int
	Headers    http.Header
}

func (m *MockResponseWriter) DoSomething(number int) (bool, error) {

	args := m.Called(number)
	return args.Bool(0), args.Error(1)

}
func (m *MockResponseWriter) Header() http.Header {
	return m.Headers
}

func (m *MockResponseWriter) Write(p []byte) (int, error) {
	m.Bytes = p
	return len(p), nil
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.StatusCode = statusCode
}

//func TestSomethingElse(t *testing.T) {
//
//	// create an instance of our test object
//	testObj := new(MockResponseWriter)
//
//	// setup expectations with a placeholder in the argument list
//	testObj.On("DoSomething", mock.Anything).Return(true, nil)
//
//	// call the code we are testing
//	targetFuncThatDoesSomethingWithObj(testObj)
//
//	// assert that the expectations were met
//	testObj.AssertExpectations(t)
//
//}

func assertBaseRoute(t *testing.T, router *Router, method, path, name string, HandlerLen int) {
	key := router.routeKey(method, path)
	assert.NotNil(t, router.routes[key])
	assert.IsType(t, &Route{}, router.routes[key])
	assert.Equal(t, HandlerLen, len(router.routes[key].(*Route).handlers))
	assert.Equal(t, name, router.routes[key].(*Route).name)
}

func TestRouter_BaseRoute(t *testing.T) {
	router := New(testing2.ApplicationDefault())
	router.GET("/gets/1", func(ctx contract.Context) {
		ctx.RenderWith(200, render2.Plain, "Body")
		ctx.Next()
	}).Use(func(ctx contract.Context) {
		ctx.RenderWith(200, render2.Plain, "After 1")
		ctx.Next()
	}).Use(func(ctx contract.Context) {
		ctx.RenderWith(200, render2.Plain, "After 2")
		ctx.Next()
	}).Use(func(ctx contract.Context) {
		ctx.RenderWith(200, render2.Plain, "Before 1")
		ctx.Next()
	}).Name("gets.1")

	assertBaseRoute(t, router.(*Router), http.MethodGet, "/gets/1", "gets.1", 3)

	router.POST("/posts", func(ctx contract.Context) {
		ctx.RenderWith(200, render2.Plain, "Body")
		ctx.Next()
	}).Name("posts.1")
	assertBaseRoute(t, router.(*Router), http.MethodPost, "/posts", "posts.1", 0)

	router.PUT("/resources/1/put", func(ctx contract.Context) {
		ctx.RenderWith(200, render2.Plain, "Body")
		ctx.Next()
	})
	assertBaseRoute(t, router.(*Router), http.MethodPut, "/resources/1/put", "", 0)

	router.DELETE("/1/delete", func(ctx contract.Context) {
		ctx.RenderWith(200, render2.Plain, "Body")
		ctx.Next()
	})
	assertBaseRoute(t, router.(*Router), http.MethodDelete, "/1/delete", "", 0)

	router.PATCH("/patch", func(ctx contract.Context) {
		ctx.RenderWith(200, render2.Plain, "Body")
		ctx.Next()
	}).Name("patch")
	assertBaseRoute(t, router.(*Router), http.MethodPatch, "/patch", "patch", 0)

	router.OPTIONS("/options", func(ctx contract.Context) {
		ctx.RenderWith(200, render2.Plain, "Body")
		ctx.Next()
	})
	assertBaseRoute(t, router.(*Router), http.MethodOptions, "/options", "", 0)

	router.Handler("GET", "/original", func(writer http.ResponseWriter, request *http.Request) {

	})
	assertBaseRoute(t, router.(*Router), http.MethodGet, "/original", "", 0)

	//	http.ListenAndServe("127.0.0.1:28082",router)
}

func TestRouter_HttpRouter(t *testing.T) {
	router := New(testing2.ApplicationDefault())
	assert.IsType(t, &httprouter.Router{}, router.HttpRouter())
}

func TestRoute(t *testing.T) {
	route := newRoute("/route", func(c contract.Context) {

	}).Use(func(c contract.Context) {
		c.Next()
	})
	assert.Equal(t, 2, len(route.Handlers()))
}

func TestRouter_Group(t *testing.T) {
	router := New(testing2.ApplicationDefault())
	v1 := router.Group("/v1").Use(func(ctx contract.Context) {
		ctx.RenderWith(200, render2.Plain, "Group v1 After")
		ctx.Next()
	}).Use(Recovery, func(ctx contract.Context) {
		ctx.RenderWith(200, render2.Plain, "Group v1 Before")
		ctx.Next()
	})
	{
		v1.GET("/gets/1", func(ctx contract.Context) {
			ctx.RenderWith(200, render2.Plain, "bdc")
			ctx.Next()
		}).Name("gets.1")
		assertBaseRoute(t, router.(*Router), http.MethodGet, "/v1/gets/1", "gets.1", 3)

		v1.POST("/posts", func(ctx contract.Context) {
			ctx.Next()
		}).Name("v1.posts")
		assertBaseRoute(t, router.(*Router), http.MethodPost, "/v1/posts", "v1.posts", 3)

		//
		v1.DELETE("/delete", func(ctx contract.Context) {
		})
		assertBaseRoute(t, router.(*Router), http.MethodDelete, "/v1/delete", "", 3)

		v1.PUT("/put", func(ctx contract.Context) {
		})
		assertBaseRoute(t, router.(*Router), http.MethodPut, "/v1/put", "", 3)

		v1.PATCH("/patch", func(ctx contract.Context) {
		})
		assertBaseRoute(t, router.(*Router), http.MethodPatch, "/v1/patch", "", 3)

		v1.OPTIONS("/options", func(ctx contract.Context) {
		})
		assertBaseRoute(t, router.(*Router), http.MethodOptions, "/v1/options", "", 3)

		v1.Handler(http.MethodGet, "/handler", func(writer http.ResponseWriter, request *http.Request) {
		})
	}

	v1Dep := v1.Group("/dep").Use(func(ctx contract.Context) {
		ctx.RenderWith(200, render2.Plain, "Group v1--dep before")
		ctx.Next()
	})
	{
		v1Dep.GET("/gets/1", func(ctx contract.Context) {

		})
	}
	assertBaseRoute(t, router.(*Router), http.MethodGet, "/v1/dep/gets/1", "", 4)
}

func TestRouter_createRouter(t *testing.T) {
	r := New(testing2.ApplicationDefault())
	r.GET("/router-handler", func(c contract.Context) {
		c.RenderWith(200, render2.Plain, "hello")
		c.Next()
	})
	r.Handler(http.MethodGet, "/handler", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello"))
	})
	r.Static("/", "../testdata/")
	r.(*Router).NotFound(func(c contract.Context) {
		c.Error(404, errors.New("not found"))
	})
	s := httptest.NewServer(r)

	res, err := s.Client().Get(s.URL + "/router-handler")
	assert.Nil(t, err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(body))

	res1, err := s.Client().Get(s.URL + "/handler")
	assert.Nil(t, err)
	defer res1.Body.Close()
	body1, err := ioutil.ReadAll(res1.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(body1))

	res2, err := s.Client().Get(s.URL + "/404")
	assert.Nil(t, err)
	defer res2.Body.Close()

	assert.Equal(t, 404, res2.StatusCode)
}
