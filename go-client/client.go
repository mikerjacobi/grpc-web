package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	rpc "github.com/mikerjacobi/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	if err := runClient(); err != nil {
		fmt.Fprintf(os.Stderr, "failed: %v\n", err)
		os.Exit(1)
	}
}
func runClient() error {
	// connnect
	tlsCreds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial("localhost:5051", grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		return fmt.Errorf("failed to dial server: %v", err)
	}

	cache := rpc.NewCacheClient(conn)
	asCli := rpc.NewAccountServiceClient(conn)
	// store
	_, err = cache.Store(context.Background(), &rpc.StoreReq{Key: "gopher", Val: []byte("con")})
	if err != nil {
		return fmt.Errorf("failed to store: %v", err)
	}
	// get
	resp, err := cache.Get(context.Background(), &rpc.GetReq{Key: "gopher"})
	if err != nil {
		return fmt.Errorf("failed to get: %v", err)
	}
	fmt.Printf("Got cached value %s\n", resp.Val)

	newAcct := &rpc.Account{Username: "abcd", Password: "def"}
	a, err := asCli.Create(context.Background(), newAcct)
	if err != nil {
		return fmt.Errorf("fail to ascli create: %+v", err)
	}
	a2, err := asCli.Get(context.Background(), &rpc.Account{AccountID: a.AccountID})
	if err != nil {
		return fmt.Errorf("fail to ascli get1: %+v", err)
	}
	a3, err := asCli.Get(context.Background(), &rpc.Account{Username: a.Username})
	if err != nil {
		return fmt.Errorf("fail to ascli get2: %+v", err)
	}
	a4, err := asCli.Authenticate(context.Background(), newAcct)
	if err != nil {
		return fmt.Errorf("fail to ascli get2: %+v", err)
	}
	fmt.Println(a)
	fmt.Println(a2)
	fmt.Println(a3)
	fmt.Println(a4)

	return nil
}
