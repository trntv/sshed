package commands

import (
	"fmt"

	"github.com/trntv/sshed/ssh"
	"github.com/mgutz/ansi"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
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

	fmt.Printf(f, ansi.Color("Hostname", "green"), ansi.Color(srv.Hostname, "white"))
	fmt.Printf(f, ansi.Color("Port", "green"), ansi.Color(srv.Port, "white"))
	fmt.Printf(f, ansi.Color("User", "green"), ansi.Color(srv.User, "white"))
	if srv.IdentityFile != "" {
		fmt.Printf(f, ansi.Color("IdentityFile", "green"), ansi.Color(srv.IdentityFile, "white"))
	}
	if srv.Options["ProxyJump"] != "" {
		fmt.Printf(f, ansi.Color("ProxyJump", "green"), ansi.Color(srv.Options["ProxyJump"], "white"))
	}

	return nil
}
