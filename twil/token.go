package twil

import (
	"errors"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GenerateToken generate a JWT token based off of the parameters
func (c *TwilioClient) GenerateToken(expirationDate time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["scope"] = c.buildScope()
	claims["iss"] = c.accountSid
	claims["exp"] = expirationDate.Unix()

	return token.SignedString([]byte(c.accountSecret))
}

func (c *TwilioClient) buildScope() string {
	return "scope:client:incoming?clientName=" + c.name + " scope:client:outgoing?appSid=" + c.appSid + "&clientName=" + c.name
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
