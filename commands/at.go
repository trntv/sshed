package commands

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/trntv/sshed/ssh"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"io/ioutil"
	"log"
	"sync"
	"github.com/fatih/color"
)

func (cmds *Commands) newAtCommand() cli.Command {
	return cli.Command{
		Name:      "at",
		Usage:     "Executes commands",
		ArgsUsage: "[key] [command]",
		Action:    cmds.atAction,
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
		var srv = ssh.Config.Get(key)
		if srv == nil {
			return errors.New("host not found")
		}

		if err != nil {
			return err
		}

		wg.Add(1)
		go (func() {
			defer wg.Done()

			cmd, err := cmds.createCommand(c, srv, &options{verbose: true}, command)
			if err != nil {
				log.Panicln(err)
			}

			var buf []byte
			w := bytes.NewBuffer(buf)
			cmd.Stdout = w

			err = cmd.Run()
			if err != nil {
				log.Panicln(err)
			}

			sr, err := ioutil.ReadAll(w)
			if err != nil {
				log.Panicln(err)
			}

			color.New(color.FgYellow).Printf("%s:\r\n", srv.Key, "yellow")
			fmt.Println(string(sr))
		})()
	}

	wg.Wait()

	return err
}
