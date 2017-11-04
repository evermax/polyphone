package main

import (
	"log"
	"os"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/evermax/polyphone/protos"
)

func main() {
	host := os.Getenv("SERVER_HOST")
	connCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(connCtx, host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err, "unable to connect")
	}
	defer conn.Close()

	client := auth.NewAuthClient(conn)

	loginRequest := &auth.LoginRequest{Username: "maxime", Password: "password"}

	loginRes, err := client.Login(connCtx, loginRequest)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Received token after login: %s", loginRes.Token)

	<-time.After(5 * time.Second)
	refreshRes, err := client.RefreshToken(connCtx, &auth.RefreshRequest{Token: loginRes.Token})

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Received token after refresh: %s", refreshRes.Token)
}
