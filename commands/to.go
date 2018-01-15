package commands

import (
	"fmt"
	"github.com/trntv/sshdb/db"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"os/user"
	"strings"
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

	var username string
	if srv.User == "" {
		u, err := user.Current()
		if err != nil {
			return err
		}
		username = u.Username
	} else {
		username = srv.User
	}

	var args = make([]string, 0)
	if srv.Password != "" {
		args = []string{
			"sshpass",
			fmt.Sprintf("-p %s", srv.Password),
		}
	}

	args = append(args, []string{
		"ssh",
		fmt.Sprintf("%s@%s", username, srv.Host),
		fmt.Sprintf("-p %s", srv.Port),
	}...)

	if srv.KeyFile != "" {
		args = append(args, fmt.Sprintf("-i %s", srv.KeyFile))
	}

	if c.Bool("verbose") == true {
		args = append(args, "-v")
	}

	cmd := exec.Command("sh", "-c", strings.Join(args, " "))

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err != nil {
		return err
	}

	err = cmd.Run()

	return err
}
