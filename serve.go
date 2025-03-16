package main

import (
	"github.com/emersion/go-smtp"
	"github.com/ngyewch/notification-server/email"
	"github.com/urfave/cli/v2"
)

func doServe(cCtx *cli.Context) error {
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
