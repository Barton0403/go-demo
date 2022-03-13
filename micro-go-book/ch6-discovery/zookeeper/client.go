package main

import (
	"context"
	service "example.com/micro-go-book/ch7-rpc/basic/string-service"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/sd/zk"
	kitlog "github.com/go-kit/log"
	"io"
	"log"
	"net/rpc"
	"time"
)

func main() {
	var (
		server = "192.168.137.3:2181" // in the change from v2 to v3, the schema is no longer necessary if connecting directly to an etcd v3 instance
		path   = "/services"          // known at compile time
		name   = "foosvc"             // taken from runtime or platform, somehow
		ctx    = context.Background()
	)

	logger := kitlog.NewNopLogger()
	client, err := zk.NewClient([]string{server}, logger)
	if err != nil {
		panic(err)
	}

	instancer, err := zk.NewInstancer(client, path+"/"+name, logger)
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
