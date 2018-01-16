package commands

import (
	"github.com/trntv/sshdb/db"
	"github.com/urfave/cli"
)

func (cmds *Commands) newToCommand() cli.Command {
	return cli.Command{
		Name:      "to",
		Usage:     "connects to server",
		ArgsUsage: "<key>",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "verbose ssh output",
			},
		},
		Action: cmds.toAction,
	}
}

func (cmds *Commands) toAction(c *cli.Context) (err error) {
	var key string
	var srv *db.Server

	if c.NArg() == 0 {
		key, err = cmds.askServerKey()
		if err != nil {
			return err
		}
	} else {
		key = c.Args().First()
	}

	srv, err = cmds.database.Get(key)
	if err != nil {
		return err
	}

	err = cmds.database.Close()
	if err != nil {
		return err
	}

	cmds.exec(srv, &options{verbose: c.Bool("verbose")}, "")

	return err
}
