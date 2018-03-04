package commands

import (
	"github.com/trntv/sshed/keychain"
	"github.com/urfave/cli"
	"github.com/fatih/color"
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
			color.New(color.FgGreen).Println("Keychain is already encrypted", "green")
		}

		password := cmds.askPassword()
		return keychain.EncryptDatabase(password)
	}
}
