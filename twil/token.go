package twil

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	// scope formats
	incomingOutgoingScopeFormat = "scope:client:incoming?clientName=%s scope:client:outgoing?appSid=%s&clientName=%s"

	// InOutScope is the scope for a client accepting both incoming and outgoing calls
	InOutScope = "inoutscope"
)

var scopeFormatMap = map[string]string{
	InOutScope: incomingOutgoingScopeFormat,
}

// GenerateToken generate a JWT token based off of the parameters
func (c *TwilioClient) GenerateToken(client string, scopeType string, expirationDate time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	scope, err := c.buildScope(scopeType, client)
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["client"] = client
	claims["scope"] = scope
	claims["iss"] = c.accountSid
	claims["exp"] = expirationDate.Unix()

	return token.SignedString([]byte(c.accountSecret))
}

func (c *TwilioClient) buildScope(scope, client string) (string, error) {
	switch scope {
	case InOutScope:
		return fmt.Sprintf(incomingOutgoingScopeFormat, client, c.appSid, client), nil
	default:
		return "", fmt.Errorf("Unknown format %s", scope)
	}
}

func (c *TwilioClient) VerifyTokenSignature(token string) error {
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) { return []byte(c.accountSecret), nil })
	if err != nil {
		return err
	}
	sstr, err := tkn.SigningString()
	if err != nil {
		return err
	}
	sig, err := tkn.Method.Sign(sstr, []byte(c.accountSecret))
	if err != nil {
		return err
	}
	pieces := strings.Split(token, ".")
	if len(pieces) != 3 {
		return errors.New("Malformed in token")
	}
	if sig != pieces[2] {
		return errors.New("Not signed properly")
	}
	return nil
}

func (c *TwilioClient) Parse(token string) (*jwt.Token, error) {
	err := c.VerifyTokenSignature(token)
	if err != nil {
		return nil, err
	}
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) { return []byte(c.accountSecret), nil })
}
