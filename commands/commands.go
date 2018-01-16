package commands

import (
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/trntv/sshdb/db"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

type Commands struct {
	database *db.DB
}

type options struct {
	verbose bool
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
		commands.newAtCommand(),
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

func (cmds *Commands) exec(srv *db.Server, options *options, command string) error {
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

	if options.verbose == true {
		args = append(args, "-v")
	}

	if command != "" {
		args = append(args, command)
	}

	cmd := exec.Command("sh", "-c", strings.Join(args, " "))

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	return err
}
