package commands

import (
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/trntv/sshed/host"
	"github.com/trntv/sshed/keychain"
	"github.com/trntv/sshed/ssh"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"sort"
	"strings"
)

type Commands struct {
	bin string
}

type options struct {
	verbose bool
}

func RegisterCommands(app *cli.App) {
	commands := &Commands{}

	beforeFunc := app.Before
	app.Before = func(context *cli.Context) error {

		err := beforeFunc(context)
		if err != nil {
			return err
		}

		commands.bin = context.String("bin")

		if keychain.Bootstrapped == false {
			fmt.Println("Creating keychain...")

			var encrypt bool
			err = survey.AskOne(&survey.Confirm{
				Message: "Protect keychain with password?",
				Default: false,
			}, &encrypt, nil)

			if encrypt == true {
				key := commands.askPassword()
				err = keychain.EncryptDatabase(key)
				if err != nil {
					return err
				}
			}

			return nil
		}

		if keychain.Encrypted == true {
			key := commands.askPassword()
			keychain.Password = key
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
		commands.newConfigCommand(),
	}
}

func (cmds *Commands) completeWithServers() {
	hosts := ssh.Config.GetAll()
	for key := range hosts {
		fmt.Println(key)
	}
}

func (cmds *Commands) askPassword() string {
	key := ""
	prompt := &survey.Password{
		Message: "Please type your password:",
	}
	survey.AskOne(prompt, &key, nil)

	return key
}

func (cmds *Commands) askServerKey() (string, error) {
	var key string
	options := make([]string, 0)
	srvs := ssh.Config.GetAll()
	for key := range srvs {
		options = append(options, key)
	}

	sort.Strings(options)
	prompt := &survey.Select{
		Message:  "Choose server:",
		Options:  options,
		PageSize: 16,
	}
	err := survey.AskOne(prompt, &key, survey.Required)

	return key, err
}

func (cmds *Commands) askServersKeys() ([]string, error) {
	var keys []string
	options := make([]string, 0)
	srvs := ssh.Config.GetAll()
	for _, h := range srvs {
		options = append(options, h.Key)
	}

	sort.Strings(options)
	prompt := &survey.MultiSelect{
		Message:  "Choose servers:",
		Options:  options,
		PageSize: 16,
	}
	err := survey.AskOne(prompt, &keys, survey.Required)

	return keys, err
}

func (cmds *Commands) createCommand(c *cli.Context, srv *host.Host, options *options, command string) (cmd *exec.Cmd, err error) {
	var username string
	if srv.User == "" {
		u, err := user.Current()
		if err != nil {
			return nil, err
		}
		username = u.Username
	} else {
		username = srv.User
	}

	var args = make([]string, 0)
	if srv.Password() != "" {
		args = []string{
			"sshpass",
			fmt.Sprintf("-p %s", srv.Password()),
		}
	}

	args = append(args, cmds.bin)
	args = append(args, fmt.Sprintf("-F %s", ssh.Config.Path))

	if pk := srv.PrivateKey(); pk != "" {
		tf, err := ioutil.TempFile("", "")
		defer os.Remove(tf.Name())
		defer tf.Close()

		if err != nil {
			return nil, err
		}

		_, err = tf.Write([]byte(pk))
		if err != nil {
			return nil, err
		}

		err = tf.Chmod(os.FileMode(0600))
		if err != nil {
			return nil, err
		}

		srv.IdentityFile = tf.Name()
	}

	if srv.User != "" {
		args = append(args, fmt.Sprintf("%s@%s", username, srv.Hostname))
	} else {
		args = append(args, fmt.Sprintf("%s", srv.Hostname))
	}

	if srv.Port != "" {
		args = append(args, fmt.Sprintf("-p %s", srv.Port))
	}

	if srv.IdentityFile != "" {
		args = append(args, fmt.Sprintf("-i %s", srv.IdentityFile))
	}

	if options.verbose == true {
		args = append(args, "-v")
	}

	if command != "" {
		args = append(args, command)
	}

	if options.verbose == true {
		fmt.Printf("%s: %s\r\n", ansi.Color("Executing", "green"), strings.Join(args, " "))
	}

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", strings.Join(args, " "))
	} else {
		cmd = exec.Command("sh", "-c", strings.Join(args, " "))
	}

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd, err
}
