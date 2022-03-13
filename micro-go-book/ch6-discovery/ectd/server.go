package main

import (
	"context"
	service "example.com/micro-go-book/ch7-rpc/basic/string-service"
	"fmt"
	"github.com/go-kit/kit/sd/etcdv3"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	kitlog "github.com/go-kit/log"
	"google.golang.org/grpc"
)

func main() {
	// Let's say this is a service that means to register itself.
	// First, we will set up some context.
	var (
		etcdServer = "192.168.137.3:2379" // in the change from v2 to v3, the schema is no longer necessary if connecting directly to an etcd v3 instance
		prefix     = "/services/foosvc/"  // known at compile time
		instance   = "127.0.0.1:8080"     // taken from runtime or platform, somehow
		key        = prefix + instance    // should be globally unique
		value      = instance             // based on our transport
		ctx        = context.Background()
	)

	options := etcdv3.ClientOptions{
		// Path to trusted ca file
		CACert: "",

		// Path to certificate
		Cert: "",

		// Path to private key
		Key: "",

		// Username if required
		Username: "",

		// Password if required
		Password: "",

		// If DialTimeout is 0, it defaults to 3s
		DialTimeout: time.Second * 3,

		// If DialKeepAlive is 0, it defaults to 3s
		DialKeepAlive: time.Second * 3,

		// If passing `grpc.WithBlock`, dial connection will block until success.
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}

	// Build the client.
	client, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
	if err != nil {
		panic(err)
	}

	// Build the registrar.
	registrar := etcdv3.NewRegistrar(client, etcdv3.Service{
		Key:   key,
		Value: value,
	}, kitlog.NewNopLogger())

	// Register our instance.
	registrar.Register()

	// At the end of our service lifecycle, for example at the end of func main,
	// we should make sure to deregister ourselves. This is important! Don't
	// accidentally skip this step by invoking a log.Fatal or os.Exit in the
	// interim, which bypasses the defer stack.
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
