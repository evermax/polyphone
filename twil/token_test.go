package twil

import (
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestGenerateTokenTestSuccess(t *testing.T) {
	client := NewTwilioClient(
		"Max",
		"twilioAccountSid",
		"twilioAccountSecret",
		"appSid",
		"8008008888",
	)
	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnQiOiJtYXhpbWUiLCJleHAiOjE0ODkzNjQ3MDQsImlzcyI6Ik1heCIsInNjb3BlIjoic2NvcGU6Y2xpZW50OmluY29taW5nP2NsaWVudE5hbWU9bWF4aW1lIHNjb3BlOmNsaWVudDpvdXRnb2luZz9hcHBTaWQ9dHdpbGlvQWNjb3VudFNlY3JldFx1MDAyNmNsaWVudE5hbWU9bWF4aW1lIn0.g_YoUOv8hG_JtFVYv7Czp3gr0kWo2nADwAi1JL_94QI"
	clientName := "maxime"
	token, err := client.GenerateToken(clientName, InOutScope, time.Unix(1489364704, 0))
	if err != nil {
		t.Fatalf("An error occured on valid token generation: %v", err)
	}

	if expectedToken != token {
		t.Fatalf("Expected value: %s, got value: %s", expectedToken, token)
	}

	components := strings.Split(token, ".")
	sig, err := jwt.SigningMethodHS256.Sign(strings.Join(components[0:2], "."), []byte(client.accountSecret))
	if err != nil {
		t.Fatalf("An error occured on signature of token: %v", err)
	}

	if sig != components[2] {
		t.Fatalf("Expected value: %s, got value: %s", sig, components[2])
	}
}

func TestVerifyFunc(t *testing.T) {
	client := NewTwilioClient(
		"Max",
		"twilioAccountSid",
		"twilioAccountSecret",
		"appSid",
		"8008008888",
	)

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnQiOiJtYXhpbWUiLCJleHAiOjE0ODkzNjQ3MDQsImlzcyI6Ik1heCIsInNjb3BlIjoic2NvcGU6Y2xpZW50OmluY29taW5nP2NsaWVudE5hbWU9bWF4aW1lIHNjb3BlOmNsaWVudDpvdXRnb2luZz9hcHBTaWQ9dHdpbGlvQWNjb3VudFNlY3JldFx1MDAyNmNsaWVudE5hbWU9bWF4aW1lIn0.g_YoUOv8hG_JtFVYv7Czp3gr0kWo2nADwAi1JL_94QI"
	err := client.VerifyTokenSignature(token)
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("An unexpected error: %v, expected %v", err, "Token is expired")
	}

	token = "a.b.c"
	err = client.VerifyTokenSignature(token)
	if err == nil || !strings.Contains(err.Error(), "illegal") {
		t.Fatalf("An unexpected error: %v, expected %v", err, "illegal base64 data at input byte 1")
	}

	token = "a.b"
	err = client.VerifyTokenSignature(token)
	if err == nil || !strings.Contains(err.Error(), "invalid number of segments") {
		t.Fatalf("An unexpected error: %v, expected %v", err, "token contains an invalid number of segments")
	}

	clientName := "maxime"
	token, err = client.GenerateToken(clientName, InOutScope, time.Now())
	if err != nil {
		t.Fatalf("An unexpected error happened: %v", err)
	}

	err = client.VerifyTokenSignature(token)
	if err != nil {
		t.Fatalf("An unexpected error happened: %v", err)
	}
}
