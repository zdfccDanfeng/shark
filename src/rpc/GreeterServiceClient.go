package rpc

import (
	"context"
	"fmt"
	"github.com/shark/src/rpc/proto/greeter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

const (
	defaultName = "rokety"
)

func GreeterSerClient() {

	resolver.SetDefaultScheme("dns")
	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s","MethodConfig": [{"Name": [{"Service": "helloworld.Greeter"}], "RetryPolicy": {"MaxAttempts":2, "InitialBackoff": "0.1s", "MaxBackoff": "1s", "BackoffMultiplier": 2.0, "RetryableStatusCodes": ["UNAVAILABLE", "CANCELLED"]}}], "HealthCheckConfig": {"ServiceName": "helloworld.Greeter"}}`, roundrobin.Name)),
		grpc.WithBlock(), grpc.WithBackoffMaxDelay(time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := helloworld.NewGreeterClient(conn)

	// Contact the server and print out its response.
	for range time.Tick(time.Second) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)
		r, err := c.Say2Hello(ctx, &helloworld.HelloRequest{Name: defaultName})
		if err != nil {
			log.Printf("could not greet: %v\n", err)
		} else {
			log.Printf("Greeting: %s", r.Message)
		}
		cancel()
	}

}
