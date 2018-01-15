package commands

import (
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/trntv/sshdb/db"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
)

type Commands struct {
	database *db.DB
}

func RegisterCommands(app *cli.App) {
	commands := &Commands{}

	app.Before = func(context *cli.Context) error {

		if context.Args().First() == "help" {
			return nil
		}

		dbpath := context.String("db")

		database, err := db.NewDB(dbpath)
		if err != nil {
			return err
		}

		commands.database = database

		if database.Bootstrapped == false {
			fmt.Println("Creating database...")
			res := struct {
				Encrypt bool
			}{}

			err = survey.Ask([]*survey.Question{
				{
					Name: "Encrypt",
					Prompt: &survey.Confirm{
						Message: "Protect database with password?",
						Default: false,
					},
				},
			}, &res)

			if res.Encrypt == true {
				key := commands.askPassword()
				err = database.EncryptDatabase(key)
				if err != nil {
					return err
				}
			}

			return nil
		}

		isEncrypted, err := database.IsEncrypted()
		if err != nil {
			return err
		}

		if isEncrypted == true {
			key := commands.askPassword()
			database.Password = key
		}

		return nil
	}

	app.After = func(context *cli.Context) error {
		if commands.database != nil {
			return commands.database.Close()
		}

		return nil
	}

	app.Commands = []cli.Command{
		commands.newShowCommand(),
		commands.newListCommand(),
		commands.newAddCommand(),
		commands.newRemoveCommand(),
		commands.newToCommand(),
		commands.newEncryptCommand(),
	}
}

func (cmds *Commands) askPassword() string {
	key := ""
	prompt := &survey.Password{
		Message: "Please type your password",
	}
	survey.AskOne(prompt, &key, nil)

	return key
}

func (cmds *Commands) askServerKey() (string, error) {
	var key string
	options := make([]string, 0)
	srvs, err := cmds.database.GetAll()
	if err != nil {
		return key, err
	}
	for key := range srvs {
		options = append(options, key)
	}
	prompt := &survey.Select{
		Message: "Choose a server:",
		Options: options,
	}
	err = survey.AskOne(prompt, &key, nil)

	return key, err
}

func (cmds *Commands) printServer(key string, srv *db.Server) {
	fmt.Printf("  %s: %s\r\n", ansi.Color("Server", "cyan"), ansi.Color(key, "white"))

	f := "  %s: %s\r\n"

	fmt.Printf(f, ansi.Color("Host", "green"), ansi.Color(srv.Host, "white"))
	fmt.Printf(f, ansi.Color("Port", "green"), ansi.Color(srv.Port, "white"))
	fmt.Printf(f, ansi.Color("User", "green"), ansi.Color(srv.User, "white"))
	if srv.KeyFile != "" {
		fmt.Printf(f, ansi.Color("Password File", "green"), ansi.Color(srv.KeyFile, "white"))
	}
}
