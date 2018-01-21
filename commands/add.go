package commands

import (
	"github.com/trntv/sshme/db"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"os/user"
)

func (cmds *Commands) newAddCommand() cli.Command {
	return cli.Command{
		Name:      "add",
		Usage:     "adds server to database",
		ArgsUsage: "[key]",
		Action:    cmds.addAction,
	}
}

func (cmds *Commands) addAction(c *cli.Context) error {
	var srv *db.Server
	var err error
	var usr, _ = user.Current()

	if c.NArg() == 1 {
		srv, err = cmds.database.Get(c.Args().First())
		if err != nil {
			if _, ok := err.(db.ErrNotFound); ok == false {
				return err
			}
		}
	}

	if srv == nil {
		srv = &db.Server{
			Port: "22",
			User: usr.Username,
		}
	}

	var qs = []*survey.Question{
		{
			Name: "key",
			Prompt: &survey.Input{
				Message: "Server alias:",
				Default: c.Args().First(),
			},
			Validate: survey.Required,
		},
		{
			Name: "host",
			Prompt: &survey.Input{
				Message: "Host:",
				Default: srv.Host,
			},
			Validate:  survey.Required,
			Transform: survey.ToLower,
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "Port:",
				Default: srv.Port,
			},
			Validate:  survey.Required,
			Transform: survey.ToLower,
		},
		{
			Name: "user",
			Prompt: &survey.Input{
				Message: "User:",
				Help:    "paste single space to leave this field empty (active user will be used when connecting)",
				Default: srv.User,
			},
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Password (optional):",
			},
		},
		{
			Name: "keyfile",
			Prompt: &survey.Input{
				Message: "Key file (optional):",
				Default: srv.KeyFile,
			},
		},
	}

	answers := struct {
		Key      string
		Host     string
		Port     string
		User     string
		Password string
		KeyFile  string
	}{}

	// perform the questions
	err = survey.Ask(qs, &answers)
	if err != nil {
		return err
	}

	err = cmds.database.Put(answers.Key, &db.Server{
		Host:     answers.Host,
		Port:     answers.Port,
		User:     answers.User,
		Password: answers.Password,
		KeyFile:  answers.KeyFile,
	})

	return err
}
