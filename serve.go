package main

import (
	"context"

	"github.com/emersion/go-smtp"
	"github.com/ngyewch/notification-server/email"
	"github.com/urfave/cli/v3"
)

func doServe(ctx context.Context, cmd *cli.Command) error {
	usernamePasswordAuthenticator := email.NewUsernamePasswordAuthenticator([]email.UsernamePassword{
		{
			Username: "bob@test.com",
			Password: "bobspassword",
		},
	})
	smtpBackend := email.NewSmtpBackend(usernamePasswordAuthenticator)

	smtpServer := smtp.NewServer(smtpBackend)
	smtpServer.Addr = ":5555"
	smtpServer.AllowInsecureAuth = true
	return smtpServer.ListenAndServe()
}
