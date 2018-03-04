package commands

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/trntv/sshed/ssh"
	"github.com/urfave/cli"
	"github.com/fatih/color"
)

func (cmds *Commands) newShowCommand() cli.Command {
	return cli.Command{
		Name:      "show",
		Usage:     "Shows host",
		ArgsUsage: "<key>",
		Action:    cmds.showAction,
		BashComplete: func(c *cli.Context) {
			// This will complete if no args are passed
			if c.NArg() > 0 {
				return
			}
			cmds.completeWithServers()
		},
	}
}

func (cmds *Commands) showAction(c *cli.Context) (err error) {
	var key string

	if c.NArg() == 0 {
		key, err = cmds.askServerKey()
		if err != nil {
			return err
		}
	} else {
		key = c.Args().First()
	}

	srv := ssh.Config.Get(key)
	if srv == nil {
		return errors.New("host not found")
	}

	f := "%s: %s\r\n"

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf(f, green("Hostname"), srv.Hostname)
	fmt.Printf(f, green("Port"), srv.Port)
	fmt.Printf(f, green("User"), srv.User)
	if srv.IdentityFile != "" {
		fmt.Printf(f, green("IdentityFile"), srv.IdentityFile)
	}

	return nil
}
