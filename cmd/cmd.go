package main

import (
	"os"

	"fmt"
	"github.com/mgutz/ansi"
	"github.com/trntv/sshdb/commands"
	"github.com/urfave/cli"
)

var version, build string

func main() {
	app := cli.NewApp()

	app.Name = "sshdb"
	app.Usage = "SSH connections manager"
	app.Author = "Eugene Terentev"
	app.Email = "eugene@terentev.net"

	if version != "" && build != "" {
		app.Version = fmt.Sprintf("%s (build %s)", version, build)
	}

	app.HelpName = "help"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "database, db",
			EnvVar: "SSHED_DB_PATH",
			Value:  fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".sshdb"),
			Usage:  "Path to database file",
		},
	}

	commands.RegisterCommands(app)

	err := app.Run(os.Args)

	if err != nil {
		fmt.Println(ansi.Red, fmt.Sprintf("Error: %s", err))
		os.Exit(1)
	}
}
