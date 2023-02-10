package main

import (
	"context"
	"flag"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"log"
)

func main() {
	var endpoint = flag.String("endpoint", "opc.tcp://localhost:55000", "OPC UA Endpoint URL")
	flag.BoolVar(&debug.Enable, "debug", false, "enable debug logging")
	flag.Parse()
	log.SetFlags(0)

	eps, err := opcua.GetEndpoints(context.Background(), *endpoint)
	if err != nil {
		log.Fatal(err)
	}

	for _, ep := range eps {
		log.Println(ep.EndpointURL, ep.SecurityPolicyURI, ep.SecurityMode)
	}
}
