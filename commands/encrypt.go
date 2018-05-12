package commands

import (
	"fmt"

	"github.com/mgutz/ansi"
	"github.com/trntv/sshed/keychain"
	"github.com/urfave/cli"
)

func (cmds *Commands) newEncryptCommand() cli.Command {
	return cli.Command{
		Name:   "encrypt",
		Usage:  "Encrypts keychain",
		Action: cmds.encryptAction(),
	}
}

func (cmds *Commands) encryptAction() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		if keychain.Encrypted {
			fmt.Println(ansi.Color("Keychain is already encrypted", "green"))
		}

		password := cmds.askPassword()
		return keychain.EncryptDatabase(password)
	}
}
