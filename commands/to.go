package commands

import (
	"github.com/trntv/sshed/host"
	"github.com/trntv/sshed/sshconn"
	"github.com/trntv/sshed/sshf"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func (cmds *Commands) newToCommand() cli.Command {
	return cli.Command{
		Name:      "to",
		Usage:     "Connects to host",
		ArgsUsage: "<key>",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "verbose ssh output",
			},
		},
		BashComplete: func(c *cli.Context) {
			// This will complete if no args are passed
			if c.NArg() > 0 {
				return
			}
			cmds.completeWithServers()
		},
		Action: cmds.toAction,
	}
}

func (cmds *Commands) toAction(c *cli.Context) (err error) {
	var key string
	var srv *host.Host

	if c.NArg() == 0 {
		key, err = cmds.askServerKey()
		if err != nil {
			return err
		}
	} else {
		key = c.Args().First()
	}

	srv = sshf.Config.Get(key)
	if srv == nil {
		return errors.New("host not found")
	}

	conn, ses := sshconn.Conn(srv)
	defer conn.Close()
	defer ses.Close()
	sshconn.Shell(ses)

	return nil
}
