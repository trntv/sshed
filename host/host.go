package host

import (
	"strings"

	"github.com/kevinburke/ssh_config"
	"github.com/trntv/sshed/keychain"
)

type Host struct {
	Key          string
	Hostname     string
	Port         string
	User         string
	IdentityFile string
	ProxyJump    string

	Options map[string]string

	creds *keychain.Record
}

func CreateFromConfig(h *ssh_config.Host) *Host {
	hh := &Host{
		Options: make(map[string]string),
	}
	for _, node := range h.Nodes {
		switch node.(type) {
		case *ssh_config.KV:
			c := node.(*ssh_config.KV)
			switch strings.ToLower(c.Key) {
			case "hostname":
				hh.Hostname = c.Value
			case "port":
				hh.Port = c.Value
			case "user":
				hh.User = c.Value
			case "identityfile":
				hh.IdentityFile = c.Value
			case "ProxyJump":
				hh.ProxyJump = c.Value
			default:
				hh.Options[c.Key] = c.Value
			}
		}
	}

	return hh
}

func (h *Host) Password() (password string) {
	h.readKeychain()
	if h.creds != nil {
		password = h.creds.Password
	}
	return password
}

func (h *Host) PrivateKey() (pk string) {
	h.readKeychain()
	if h.creds != nil {
		pk = h.creds.PrivateKey
	}
	return pk
}

func (h *Host) readKeychain() error {
	var err error
	if h.creds != nil {
		return nil
	}
	h.creds, err = keychain.Get(h.Key)
	return err
}
