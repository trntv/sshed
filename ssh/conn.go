package ssh

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"time"

	"github.com/trntv/sshed/host"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
)

type jsonconf struct {
	Cmd   string
	Bargs string
	Aargs string
	Hosts []string
}

func getKeyFile(keypath string) (gossh.Signer, error) {
	buf, err := ioutil.ReadFile(keypath)
	if err != nil {
		return nil, err
	}

	pubkey, err := gossh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

func getSSHConfig(config *host.Host) *gossh.ClientConfig {
	auths := []gossh.AuthMethod{}

	if config.Password() != "" {
		auths = append(auths, gossh.Password(config.Password()))
	}

	if config.IdentityFile != "" {
		if pubkey, err := getKeyFile(config.IdentityFile); err != nil {
			panic(err)
		} else {
			auths = append(auths, gossh.PublicKeys(pubkey))
		}
	}

	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, gossh.PublicKeysCallback(agent.NewClient(sshAgent).Signers))
		defer sshAgent.Close()
	}

	return &gossh.ClientConfig{
		Timeout:         10 * time.Second,
		User:            config.User,
		Auth:            auths,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
	}
}

// Conn Connection Method
func Conn(srv *host.Host) ([]*gossh.Client, *gossh.Session) {
	var client *gossh.Client
	var err error

	config := getSSHConfig(srv)

	var Hosts []*host.Host
	var Clients []*gossh.Client

	key := srv.Key
	Hosts = append(Hosts, srv)

	for {
		c := Config.Get(key).Options
		if _, ok := c["ProxyJump"]; ok {
			key = c["ProxyJump"]
		} else {
			break
		}
		Hosts = append(Hosts, Config.Get(key))
	}

	if len(Hosts) > 2 {
		dialFunc := net.Dial

		for i := len(Hosts) - 1; i >= 0; i-- {
			proxyConfig := getSSHConfig(Hosts[i])

			tcpConn, err := dialFunc("tcp", net.JoinHostPort(Hosts[i].Hostname, Hosts[i].Port))
			if err != nil {
				panic(err)
			}

			sshConn, chans, reqs, err := gossh.NewClientConn(tcpConn, net.JoinHostPort(Hosts[i].Hostname, Hosts[i].Port), proxyConfig)
			if err != nil {
				panic(err)
			}
			client = gossh.NewClient(sshConn, chans, reqs)
			Clients = append(Clients, client)
			dialFunc = client.Dial
		}
	} else if len(Hosts) == 2 {
		proxyConfig := getSSHConfig(Hosts[1])

		proxyClient, err := gossh.Dial("tcp", net.JoinHostPort(Hosts[1].Hostname, Hosts[1].Port), proxyConfig)
		if err != nil {
			panic(err)
		}

		conn, err := proxyClient.Dial("tcp", net.JoinHostPort(srv.Hostname, srv.Port))
		if err != nil {
			panic(err)
		}

		ncc, chans, reqs, err := gossh.NewClientConn(conn, net.JoinHostPort(srv.Hostname, srv.Port), config)
		if err != nil {
			panic(err)
		}

		client = gossh.NewClient(ncc, chans, reqs)
		Clients = append(Clients, client)
	} else {
		client, err = gossh.Dial("tcp", net.JoinHostPort(srv.Hostname, srv.Port), config)
		if err != nil {
			panic(err)
		}
		Clients = append(Clients, client)
	}

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}

	return Clients, session
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// Shell Remote Terminal Attachment
func Shell(session *gossh.Session, host string) {
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	user, _ := user.Current()
	var jsonConf jsonconf
	file, err := ioutil.ReadFile(user.HomeDir + "/.sshed.config")
	if err == nil {
		json.Unmarshal(file, &jsonConf)
		if contains(jsonConf.Hosts, host) {
			exec.Command(jsonConf.Cmd, jsonConf.Bargs).Run()
		}
	}
	modes := gossh.TerminalModes{
		gossh.ECHO:          1,
		gossh.ECHOCTL:       0,
		gossh.TTY_OP_ISPEED: 115200,
		gossh.TTY_OP_OSPEED: 115200,
	}

	termFD := int(os.Stdin.Fd())

	w, h, _ := terminal.GetSize(termFD)

	termState, _ := terminal.MakeRaw(termFD)
	defer terminal.Restore(termFD, termState)

	session.RequestPty(os.Getenv("TERM"), h, w, modes)
	session.Shell()
	session.Wait()
	if err == nil {
		if contains(jsonConf.Hosts, host) {
			exec.Command(jsonConf.Cmd, jsonConf.Aargs).Run()
		}
	}
}

// RunCmd Run Command in Remote Host
func RunCmd(session *gossh.Session, cmd string) (string, error) {
	var stdout, stderr bytes.Buffer
	var buf bytes.Buffer

	session.Stdout = &stdout
	session.Stderr = &stderr

	err := session.Run(cmd)
	if err != nil {
		buf.Write(stderr.Bytes())
	}

	buf.Write(stdout.Bytes())

	return buf.String(), err
}
