package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalln("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = "/etc/golang-example/.tls/server.crt"
		}
		if *keyFile == "" {
			*keyFile = "/etc/golang-example/.tls/server.key"
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)
	//pb.RegisterRouteGuideServer(grpcServer, newServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
