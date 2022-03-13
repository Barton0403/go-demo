package main

import (
	service "example.com/micro-go-book/ch7-rpc/basic/string-service"
	"fmt"
	"github.com/go-kit/kit/sd/zk"
	kitlog "github.com/go-kit/log"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

func main() {
	var (
		server   = "192.168.137.3:2181" // in the change from v2 to v3, the schema is no longer necessary if connecting directly to an etcd v3 instance
		instance = "127.0.0.1:8080"     // taken from runtime or platform, somehow
		path     = "/services"          // known at compile time
		name     = "foosvc"             // taken from runtime or platform, somehow
		value    = instance             // based on our transport
	)

	zk.ConnectTimeout(time.Second * 3)
	logger := kitlog.NewNopLogger()
	client, err := zk.NewClient([]string{server}, logger)
	if err != nil {
		panic(err)
	}

	registrar := zk.NewRegistrar(client, zk.Service{
		Path: path,
		Name: name,
		Data: []byte(value),
	}, logger)
	registrar.Register()
	defer registrar.Deregister()

	stringService := new(service.StringService)
	registerError := rpc.Register(stringService)
	if registerError != nil {
		log.Fatal("Register error: ", registerError)
	}

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", instance)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	fmt.Println("listen " + instance)
	http.Serve(l, nil)
}
