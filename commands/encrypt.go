package commands

import (
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/urfave/cli"
)

func (cmds *Commands) newEncryptCommand() cli.Command {
	return cli.Command{
		Name:   "encrypt",
		Usage:  "encrypt database",
		Action: cmds.encryptAction(),
	}
}

func (cmds *Commands) encryptAction() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		if isEncrypted, _ := cmds.database.IsEncrypted(); isEncrypted {
			fmt.Println(ansi.Color("Database already encrypted", "red"))
		}

		key := cmds.askPassword()
		return cmds.database.EncryptDatabase(key)
	}
}
