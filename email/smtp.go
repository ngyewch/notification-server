package email

import (
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"io"
)

type SmtpAuth interface {
	AuthMechanisms() []string
	Auth(mech string) (sasl.Server, error)
}

// ----

type SmtpBackend struct {
	usernamePasswordAuthenticator *UsernamePasswordAuthenticator
}

func NewSmtpBackend(usernamePasswordAuthenticator *UsernamePasswordAuthenticator) *SmtpBackend {
	return &SmtpBackend{
		usernamePasswordAuthenticator: usernamePasswordAuthenticator,
	}
}

func (backend *SmtpBackend) NewSession(conn *smtp.Conn) (smtp.Session, error) {
	return &SmtpSession{
		usernamePasswordAuthenticator: backend.usernamePasswordAuthenticator,
		conn:                          conn,
	}, nil
}

type SmtpSession struct {
	usernamePasswordAuthenticator *UsernamePasswordAuthenticator
	conn                          *smtp.Conn
	username                      string
}

func (session *SmtpSession) AuthMechanisms() []string {
	return []string{sasl.Plain, sasl.Login}
}

func (session *SmtpSession) Auth(mech string) (sasl.Server, error) {
	switch mech {
	case sasl.Plain:
		return sasl.NewPlainServer(func(identity string, username string, password string) error {
			err := session.usernamePasswordAuthenticator.Authenticate(username, password)
			if err != nil {
				return err
			}
			session.username = username
			return nil
		}), nil
	case sasl.Login:
		return NewLoginServer(func(username string, password string) error {
			err := session.usernamePasswordAuthenticator.Authenticate(username, password)
			if err != nil {
				return err
			}
			session.username = username
			return nil
		}), nil
	default:
		return nil, fmt.Errorf("mechanism '%s' is not supported", mech)
	}
}

func (session *SmtpSession) Mail(from string, opts *smtp.MailOptions) error {
	fmt.Println("**** Username:", session.username)
	fmt.Printf("MAIL From: %s\n", from)
	return nil
}

func (session *SmtpSession) Rcpt(to string, opts *smtp.RcptOptions) error {
	fmt.Printf("RCPT To: %s\n", to)
	return nil
}

func (session *SmtpSession) Data(r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	fmt.Printf("DATA %s\n", string(b))
	return nil
}

func (session *SmtpSession) Reset() {
	fmt.Println("**** Reset")
}

func (session *SmtpSession) Logout() error {
	fmt.Println("**** Logout")
	return nil
}
