package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

var (
	version string

	smtpListenAddrFlag = &cli.StringFlag{
		Name:    "smtp-listen-addr",
		Sources: cli.EnvVars("SMTP_LISTEN_ADDR"),
	}

	app = &cli.Command{
		Name:    "notification server",
		Usage:   "Notification server",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:   "serve",
				Usage:  "serve",
				Action: doServe,
				Flags: []cli.Flag{
					smtpListenAddrFlag,
				},
			},
		},
	}
)

func main() {
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
