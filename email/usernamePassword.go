package email

import "fmt"

type UsernamePassword struct {
	Username string
	Password string
}

type UsernamePasswordAuthenticator struct {
	entries []UsernamePassword
}

func NewUsernamePasswordAuthenticator(entries []UsernamePassword) *UsernamePasswordAuthenticator {
	return &UsernamePasswordAuthenticator{
		entries: entries,
	}
}

func (authenticator *UsernamePasswordAuthenticator) Authenticate(username string, password string) error {
	for _, entry := range authenticator.entries {
		if (entry.Username == username) && (entry.Password == password) {
			return nil
		}
	}
	return fmt.Errorf("authentication failure")
}
