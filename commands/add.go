package commands

import (
	"os/user"

	"github.com/trntv/sshed/host"
	"github.com/trntv/sshed/keychain"
	"github.com/trntv/sshed/sshf"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
)

type answers struct {
	Key            string
	Host           string
	Port           string
	User           string
	Password       string
	KeyFile        string
	KeyFileContent string
	ProxyJump      string
}

func (cmds *Commands) newAddCommand() cli.Command {
	return cli.Command{
		Name:      "add",
		Usage:     "Add or edit host",
		ArgsUsage: "[key]",
		Action:    cmds.addAction,
		BashComplete: func(c *cli.Context) {
			// This will complete if no args are passed
			if c.NArg() > 0 {
				return
			}
			cmds.completeWithServers()
		},
	}
}

func (cmds *Commands) addAction(c *cli.Context) error {
	var h *host.Host
	var err error
	var usr, _ = user.Current()
	var key = c.Args().First()

	if key != "" {
		h = sshf.Config.Get(key)
	}

	if h == nil {
		h = &host.Host{
			Key:  key,
			Port: "22",
			User: usr.Username,
		}
	}

	var qs = []*survey.Question{
		{
			Name: "key",
			Prompt: &survey.Input{
				Message: "Alias:",
				Default: h.Key,
			},
			Validate: survey.Required,
		},
		{
			Name: "host",
			Prompt: &survey.Input{
				Message: "Hostname:",
				Default: h.Hostname,
			},
			Validate:  survey.Required,
			Transform: survey.ToLower,
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "Port:",
				Default: h.Port,
			},
			Transform: survey.ToLower,
		},
		{
			Name: "user",
			Prompt: &survey.Input{
				Message: "User:",
				Help:    "paste single space to leave this field empty (active user will be used when connecting)",
				Default: h.User,
			},
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Password (optional):",
			},
		},
	}

	answers := &answers{}

	// perform the questions
	err = survey.Ask(qs, answers)
	if err != nil {
		return err
	}

	askForIdentityFile(answers, h)

	askForJumpHost(answers, h)

	h = &host.Host{
		Key:          answers.Key,
		Hostname:     answers.Host,
		Port:         answers.Port,
		User:         answers.User,
		IdentityFile: answers.KeyFile,
		ProxyJump:    answers.ProxyJump,
		Options:      make(map[string]string),
	}

	isOptions := false
	for {
		err = survey.AskOne(&survey.Confirm{
			Message: "Add additional SSH options?",
		}, &isOptions, nil)
		if err != nil {
			return err
		}
		if isOptions == true {
			option := struct {
				Key   string
				Value string
			}{}
			optionQuestions := []*survey.Question{
				{
					Name: "key",
					Prompt: &survey.Input{
						Message: "Option:",
					},
					Validate: survey.Required,
				},
				{
					Name: "value",
					Prompt: &survey.Input{
						Message: "Value:",
					},
					Validate: survey.Required,
				},
			}
			err = survey.Ask(optionQuestions, &option)
			if err != nil {
				return err
			}
			h.Options[option.Key] = option.Value
		} else {
			break
		}
	}

	err = keychain.Put(h.Key, &keychain.Record{
		Password:   answers.Password,
		PrivateKey: answers.KeyFileContent,
	})
	if err != nil {
		return err
	}

	sshf.Config.Add(h)

	return sshf.Config.Save()
}

func askForIdentityFile(answers *answers, srv *host.Host) (err error) {
	const OPTION_SKIP = "Leave empty"
	const OPTION_SELECT = "Select known key"
	const OPTION_INPUT = "Input custom path"
	const OPTION_EDITOR = "Paste file contents"

	if srv.IdentityFile != "" {
		var change bool
		err = survey.AskOne(&survey.Confirm{
			Message: "Do you want to change key information?",
		}, &change, nil)
		if err != nil {
			return err
		}

		if change == false {
			answers.KeyFile = srv.IdentityFile
			return nil
		}
	}

	var options = []string{
		OPTION_INPUT,
		OPTION_EDITOR,
		OPTION_SKIP,
	}

	if len(sshf.Config.Keys) > 0 {
		options = append(options, OPTION_SELECT)
	}

	var choice string

	err = survey.AskOne(&survey.Select{
		Options: options,
		Message: "How do you want to provide key file?",
	}, &choice, survey.Required)
	if err != nil {
		return err
	}

	switch choice {
	case OPTION_SKIP:
		answers.KeyFile = srv.IdentityFile
		return
	case OPTION_SELECT:
		err = survey.AskOne(&survey.Select{
			Options: sshf.Config.Keys,
			Message: "Choose private key:",
			Default: srv.IdentityFile,
		}, &answers.KeyFile, nil)
	case OPTION_INPUT:
		err = survey.AskOne(&survey.Input{
			Message: "Private key path:",
			Default: srv.IdentityFile,
		}, &answers.KeyFile, nil)
	case OPTION_EDITOR:
		err = survey.AskOne(&survey.Editor{
			Message: "Private key content:",
		}, &answers.KeyFileContent, nil)
	}

	return err
}

func askForJumpHost(answers *answers, srv *host.Host) (err error) {
	options := make([]string, 0)

	options = append(options, "Without ProxyJump")

	srvs := sshf.Config.GetAll()
	for key := range srvs {
		options = append(options, key)
	}

	var choice string

	err = survey.AskOne(&survey.Select{
		Options: options,
		Message: "Choose ProxyJump",
	}, &choice, survey.Required)

	answers.ProxyJump = choice

	if err != nil {
		return err
	}

	return err
}
