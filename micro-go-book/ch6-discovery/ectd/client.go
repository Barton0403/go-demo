package main

import (
	"context"
	service "example.com/micro-go-book/ch7-rpc/basic/string-service"
	"fmt"
	"github.com/go-kit/kit/sd/etcdv3"
	"io"
	"log"
	"net/rpc"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	kitlog "github.com/go-kit/log"
	"google.golang.org/grpc"
)

func main() {
	// Let's say this is a service that means to register itself.
	// First, we will set up some context.
	var (
		etcdServer = "192.168.137.3:2379" // in the change from v2 to v3, the schema is no longer necessary if connecting directly to an etcd v3 instance
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

	// It's likely that we'll also want to connect to other services and call
	// their methods. We can build an Instancer to listen for changes from etcd,
	// create Endpointer, wrap it with a load-balancer to pick a single
	// endpoint, and finally wrap it with a retry strategy to get something that
	// can be used as an endpoint directly.
	barPrefix := "/services/foosvc"
	logger := kitlog.NewNopLogger()
	instancer, err := etcdv3.NewInstancer(client, barPrefix, logger)
	if err != nil {
		panic(err)
	}
	endpointer := sd.NewEndpointer(instancer, barFactory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, 3*time.Second, balancer)

	// And now retry can be used like any other endpoint.
	req := struct{}{}
	if _, err = retry(ctx, req); err != nil {
		panic(err)
	}
}

func barFactory(address string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		client, err := rpc.DialHTTP("tcp", address)
		if err != nil {
			log.Fatal("dialing:", err)
		}

		stringReq := &service.StringRequest{"A", "B"}
		var reply string
		err = client.Call("StringService.Concat", stringReq, &reply)
		if err != nil {
			log.Fatal("StringService error:", err)
		}
		fmt.Printf("StringService Concat: %s concat %s = %s\n", stringReq.A, stringReq.B, reply)

		stringReq = &service.StringRequest{"ACD", "BDF"}
		// 异步调用
		call := client.Go("StringService.Diff", stringReq, &reply, nil)
		_ = <-call.Done
		fmt.Printf("StringService Diff: %s diff %s = %s\n", stringReq.A, stringReq.B, reply)

		return nil, nil
	}, nil, nil
}
