package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dgrijalva/jwt-go"
	"github.com/evermax/polyphone/protos"
	"github.com/evermax/polyphone/twil"
	"golang.org/x/net/context"
)

func main() {
	host := os.Getenv("SERVER_HOST")

	client := twil.NewTwilioClient("AC123", "123", "", "", "")
	s := newServer(host, client)

	s.Run(SignalContext(context.Background()))
}

type server struct {
	Host       string
	TwilClient *twil.TwilioClient
	mutex      *sync.RWMutex
	clients    map[string]bool
}

func newServer(host string, client *twil.TwilioClient) *server {
	return &server{
		Host:       host,
		TwilClient: client,
		mutex:      &sync.RWMutex{},
		clients:    make(map[string]bool),
	}
}

func (s *server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Printf("Starting server on host %s\n", s.Host)
	s.printBindings()

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
		log.Println("Password is incorrect")
		return nil, status.Error(codes.Unauthenticated, "password is incorrect")
	case req.Username == "":
		log.Println("Username is required")
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	// Identify user using username password

	tkn, err := s.genToken(req.Username)

	if err != nil {
		log.Printf("Token generation error %s\n", err)
		return nil, status.Error(codes.Internal, "token generation error")
	}
	s.loginClient(req.Username)

	log.Printf("New user login: %s, generated token %s\n", req.Username, tkn)

	return &auth.LoginResponse{Token: tkn}, nil
}

func (s *server) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	client, err := s.extractClient(req.Token)
	if err != nil {
		return nil, err
	}

	s.logoutClient(client)
	return &auth.LogoutResponse{}, nil
}

func (s *server) RefreshToken(ctx context.Context, req *auth.RefreshRequest) (*auth.RefreshResponse, error) {
	client, err := s.extractClient(req.Token)
	if err != nil {
		return nil, err
	}

	if !s.isClientLoggedIn(client) {
		return nil, fmt.Errorf("Not logged in")
	}

	token, err := s.genToken(client)

	if err != nil {
		return nil, err
	}

	log.Printf("New token: %s, for client %s\n", token, client)

	return &auth.RefreshResponse{Token: token}, nil
}

func (s *server) genToken(client string) (string, error) {
	return s.TwilClient.GenerateToken(client, twil.InOutScope, time.Now().Add(1*time.Hour))
}

func (s *server) verifyToken(token string) error {
	return s.TwilClient.VerifyTokenSignature(token)
}

func (s *server) extractClient(token string) (string, error) {
	tkn, err := s.TwilClient.Parse(token)
	if err != nil {
		return "", err
	}
	claims := tkn.Claims.(jwt.MapClaims)
	return claims["client"].(string), nil
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

func (s *server) printBindings() {
	host := strings.Split(s.Host, ":")

	if host[0] == "" || host[0] == "0.0.0.0" {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			log.Println(err)
		}
		for _, addr := range addrs {
			log.Printf("Listening on %v:%s\n", addr, host[1])
		}
	} else {
		log.Printf("Listening on %s", s.Host)
	}
}
