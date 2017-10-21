package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/evermax/polyphone/protos"
	"github.com/evermax/polyphone/twil"
	"golang.org/x/net/context"
)

func main() {
	fmt.Println("Hello, world")

	host := os.Getenv("SERVER_HOST")

	client := twil.NewTwilioClient("AC123", "123", "", "", "")
	s := Server(host, client)

	s.Run(SignalContext(context.Background()))
}

type server struct {
	Host       string
	TwilClient *twil.TwilioClient
}

func Server(host string, client *twil.TwilioClient) *server {
	return &server{
		Host:       host,
		TwilClient: client,
	}
}

func (s *server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Printf("Starting server on host %s\n", s.Host)

	srv := grpc.NewServer()
	auth.RegisterAuthServer(srv, s)

	l, err := net.Listen("tcp", s.Host)
	if err != nil {
		return err
	}

	go func() {
		srv.Serve(l)
		cancel()
	}()

	<-ctx.Done()

	log.Println("shutting down")
	srv.GracefulStop()

	return nil
}

func (s *server) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	switch {
	case req.Password != "password":
		return nil, status.Error(codes.Unauthenticated, "password is incorrect")
	case req.Username == "":
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	// Identify user using username password

	tkn, err := s.genToken()
	if err != nil {
		return nil, status.Error(codes.Internal, "token generation error")
	}
	s.linkUser(req.Username, tkn)

	log.Printf("New user login: %s, generated token %s\n", req.Username, tkn)

	return &auth.LoginResponse{Token: tkn}, nil
}

func (s *server) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	//
	return &auth.LogoutResponse{}, nil
}

func (s *server) RefreshToken(ctx context.Context, req *auth.RefreshRequest) (*auth.RefreshResponse, error) {
	//
	return &auth.RefreshResponse{}, nil
}

func (s *server) genToken() (string, error) {
	return s.TwilClient.GenerateToken(time.Now().Add(1 * time.Hour))
}

func (s *server) linkUser(username, token string) {

}

func (s *server) verifyToken(token string) error {
	s.TwilClient.VerifyTokenSignature(token)
	return nil
}

func SignalContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Print("Listening for shutdown signal")
		<-sigs
		log.Print("shutdown signal received")
		signal.Stop(sigs)
		close(sigs)
		cancel()
	}()

	return ctx
}
