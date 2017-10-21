package twil

import (
	"fmt"
	"time"
)

const (
	jwtMarshalFormat = `{"token":"%s"}`
)

// JWT
type JWT struct {
	tokenID        int
	Token          string
	userID         int
	creationDate   time.Time
	expirationDate time.Time
	lastUseDate    time.Time
}

// MarshalJSON
func (t JWT) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(jwtMarshalFormat, t.Token)), nil
}
