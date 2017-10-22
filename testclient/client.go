package main

import (
	"log"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/evermax/polyphone/protos"
)

func main() {
	connCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(connCtx, "10.0.1.50:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println(err, "unable to connect")
	}
	defer conn.Close()

	client := auth.NewAuthClient(conn)

	loginRequest := &auth.LoginRequest{Username: "maxime", Password: "password"}

	client.Login(connCtx, loginRequest)
}
