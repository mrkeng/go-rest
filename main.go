package main

import (
	"encoding/json"
	"fmt"
	"go-rest/rest"
	"go-rest/server"
	"go-rest/server/context"
	"net/http"
	"net/url"
	"os"
)

type Consumer struct {
	Key    string
	Secret string
}

func (c Consumer) Authorize(urlStr string, requestType string, form url.Values) url.Values {
	return form
}

type Foo struct {
	Foo string
	Bar float64
}

func (f Foo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"foo": f.Foo,
		"bar": f.Bar,
	})
}

type FooHandler struct {
	server.BaseResourceHandler
}

func (f FooHandler) ResourceName() string {
	return "foos"
}

func (f FooHandler) ReadResource(ctx context.RequestContext, id string, version string) (server.Resource, error) {
	if id == "42" {
		return &Foo{"hello", 42}, nil
	}

	return nil, server.ResourceNotFound()
}

func (f FooHandler) ReadResourceList(ctx context.RequestContext, limit int, version string) ([]server.Resource, string, error) {
	resources := make([]server.Resource, 0)
	resources = append(resources, &Foo{Foo: "hello", Bar: 42})
	resources = append(resources, &Foo{Foo: "world", Bar: 100})
	return resources, "cursor123", nil
}

func (f FooHandler) CreateResource(ctx context.RequestContext, data server.Payload, version string) (server.Resource, error) {
	foo := &Foo{Foo: data["foo"].(string), Bar: data["bar"].(float64)}
	return foo, nil
}

func (f FooHandler) UpdateResource(ctx context.RequestContext, id string, data server.Payload, version string) (server.Resource, error) {
	foo := &Foo{Foo: data["foo"].(string), Bar: data["bar"].(float64)}
	return foo, nil
}

func (f FooHandler) DeleteResource(ctx context.RequestContext, id string, version string) (server.Resource, error) {
	foo := &Foo{}
	return foo, nil
}

func MyMiddleware(wrapped http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var scheme string
		scheme = r.URL.Scheme
		if scheme == "" {
			scheme = "http"
		}

		fmt.Println(scheme + "://" + r.Host + r.RequestURI)
		wrapped(w, r)
	}
}

func main() {
	if os.Args[1] == "1" {
		api := server.NewRestApi()
		api.RegisterResourceHandler(FooHandler{}, MyMiddleware)
		api.Start(":8080")
	}

	rc := rest.Client{Consumer{"key", "secret"}}
	params := map[string]string{
		"foo": "bar",
	}
	resp, err := rc.Get("http://localhost:8080/api/v0.1/foos/1", params)
	fmt.Println(resp)
	fmt.Println(err)
}
