package main

import (
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/trntv/sshed/commands"
	"github.com/trntv/sshed/keychain"
	"github.com/trntv/sshed/ssh"
	"github.com/urfave/cli"
	"os"
	"os/user"
	"path/filepath"
)

var version, build string

func main() {

	app := cli.NewApp()

	app.Name = "sshed"
	app.Usage = "SSH config editor and hosts manager"
	app.Author = "Eugene Terentev"
	app.Email = "eugene@terentev.net"

	if version != "" && build != "" {
		app.Version = fmt.Sprintf("%s (build %s)", version, build)
	}

	usr, _ := user.Current()
	homeDir := usr.HomeDir

	app.HelpName = "help"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "keychain",
			EnvVar: "SSHED_KEYCHAIN",
			Value:  filepath.Join(homeDir, ".sshed"),
			Usage:  "path to keychain database",
		},
		cli.StringFlag{
			Name:   "config",
			EnvVar: "SSHED_CONFIG_FILE",
			Value:  filepath.Join(homeDir, ".ssh", "config"),
			Usage:  "path to SSH config file",
		},
		cli.StringFlag{
			Name:   "bin",
			EnvVar: "SSHED_BIN",
			Value:  "ssh",
			Usage:  "path to SSH binary",
		},
	}

	app.EnableBashCompletion = true

	app.Before = func(context *cli.Context) error {
		if context.Command.Name == "help" {
			return nil
		}

		var err error
		ssh.Config, err = ssh.Parse(context.String("config"))
		if err != nil {
			return err
		}

		dbpath := context.String("keychain")

		err = keychain.Open(dbpath)
		return err
	}

	commands.RegisterCommands(app)

	err := app.Run(os.Args)

	if err != nil {
		fmt.Println(ansi.Red, fmt.Sprintf("Error: %s", err))
		os.Exit(1)
	}
}
