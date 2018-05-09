package commands

import (
	"fmt"
	"sync"

	"github.com/trntv/sshed/sshconn"
	"github.com/trntv/sshed/sshf"
	"github.com/mgutz/ansi"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
)

func (cmds *Commands) newAtCommand() cli.Command {
	return cli.Command{
		Name:      "at",
		Usage:     "Executes commands",
		ArgsUsage: "[key] [command]",
		Action:    cmds.atAction,
		BashComplete: func(c *cli.Context) {
			if c.NArg() > 0 {
				return
			}
			cmds.completeWithServers()
		},
	}
}
func (cmds *Commands) atAction(c *cli.Context) (err error) {
	keys := []string{c.Args().First()}
	if keys[0] == "" {
		keys, err = cmds.askServersKeys()
		if err != nil {
			return err
		}
	}

	command := c.Args().Get(1)
	if command == "" {

		err = survey.AskOne(&survey.Input{Message: "Command:"}, &command, nil)
		if err != nil {
			return err
		}

		fmt.Println("")
	}

	var wg sync.WaitGroup
	for _, key := range keys {
		var srv = sshf.Config.Get(key)
		if srv == nil {
			return errors.New("host not found")
		}

		if err != nil {
			return err
		}

		wg.Add(1)
		go (func() {
			defer wg.Done()

			conn, ses := sshconn.Conn(srv)
			defer conn.Close()
			defer ses.Close()
			out, _ := sshconn.RunCmd(ses, command)

			fmt.Printf("%s:\r\n", ansi.Color(srv.Key, "yellow"))
			fmt.Println(string(out))
		})()
	}

	wg.Wait()

	return err
}
