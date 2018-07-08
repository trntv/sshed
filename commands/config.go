package commands

import (
	"fmt"

	"github.com/trntv/sshed/ssh"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
)

func (cmds *Commands) newConfigCommand() cli.Command {
	return cli.Command{
		Name:   "config",
		Usage:  "Shows SSH config",
		Action: cmds.configAction,
		Subcommands: []cli.Command{
			{
				Name:   "edit",
				Usage:  "edit SSH config",
				Action: cmds.configEditAction,
			},
		},
	}
}

func (cmds *Commands) configAction(ctx *cli.Context) error {
	fmt.Print(string(ssh.Config.Content))

	return nil
}

func (cmds *Commands) configEditAction(ctx *cli.Context) (err error) {
	var content string
	err = survey.AskOne(&survey.Editor{
		Message:       "",
		Default:       ssh.Config.String(),
		HideDefault:   true,
		AppendDefault: true,
	}, &content, nil)
	if err != nil {
		return err
	}

	err = ssh.Config.SaveContent([]byte(content))

	return err
}
