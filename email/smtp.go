package email

import (
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"golang.org/x/net/html/charset"
	"io"
	"mime"
	"net/mail"
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
	from                          *mail.Address
	to                            *mail.Address
	subject                       string
	body                          string
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
	if session.username == "" {
		return fmt.Errorf("not authenticated")
	}
	fromAddr, err := mail.ParseAddress(from)
	if err != nil {
		return err
	}
	if fromAddr.Address != session.username {
		return fmt.Errorf("authentication mismatch")
	}
	session.from = fromAddr
	return nil
}

func (session *SmtpSession) Rcpt(to string, opts *smtp.RcptOptions) error {
	if session.username == "" {
		return fmt.Errorf("not authenticated")
	}
	toAddr, err := mail.ParseAddress(to)
	if err != nil {
		return err
	}
	session.to = toAddr
	return nil
}

func (session *SmtpSession) Data(r io.Reader) error {
	if session.username == "" {
		return fmt.Errorf("not authenticated")
	}
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return err
	}
	mediaType, contentTypeParams, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if mediaType == "" {
		mediaType = "text/plain"
	}
	if mediaType != "text/plain" {
		return fmt.Errorf("unsupported media type: %s", mediaType)
	}
	charsetName := contentTypeParams["charset"]
	if charsetName == "" {
		charsetName = "US-ASCII"
	}
	cs, _ := charset.Lookup(charsetName)
	if cs == nil {
		return fmt.Errorf("unknown charset: %s", charsetName)
	}

	session.subject = msg.Header.Get("Subject")

	bodyBytes, err := io.ReadAll(cs.NewDecoder().Reader(msg.Body))
	if err != nil {
		return err
	}
	session.body = string(bodyBytes)

	fmt.Printf("Subject: %s\n", session.subject)
	fmt.Println()
	fmt.Println(session.body)

	return nil
}

func (session *SmtpSession) Reset() {
	session.to = nil
	session.from = nil
	session.subject = ""
	session.body = ""
}

func (session *SmtpSession) Logout() error {
	return nil
}
