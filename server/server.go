package main

import (
	"fmt"
	"html"
	"net"
	"net/http"
	"os"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	rpc "github.com/mikerjacobi/grpc/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	if err := runServer(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run cache server: %s\n", err)
		os.Exit(1)
	}
}
func runServer() error {
	tlsCreds, err := credentials.NewServerTLSFromFile("tls.crt", "tls.key")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(grpc.Creds(tlsCreds))
	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	rpc.RegisterCacheServer(grpcServer, newCacheService())
	rpc.RegisterAccountServiceServer(grpcServer, newAccountService())
	l, err := net.Listen("tcp", "localhost:5051")
	if err != nil {
		return err
	}

	logrus.Infof("starting server")
	go func() {
		grpcServer.Serve(l)
	}()

	http.HandleFunc("/grpc", grpcHandler(wrappedGrpc))
	http.HandleFunc("/hello", helloHandler)
	fs := http.FileServer(http.Dir("../web-client"))
	http.Handle("/", fs)

	return http.ListenAndServeTLS(":443", "tls.crt", "tls.key", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func grpcHandler(wrappedGrpc *grpcweb.WrappedGrpcServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(r) {
			wrappedGrpc.ServeHTTP(w, r)
		}
		http.DefaultServeMux.ServeHTTP(w, r)
	}
}
