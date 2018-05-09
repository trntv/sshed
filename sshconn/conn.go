package sshconn

import (
	"bytes"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/trntv/sshed/host"
	"github.com/trntv/sshed/sshf"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
)

func getKeyFile(keypath string) (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(keypath)
	if err != nil {
		return nil, err
	}

	pubkey, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

func getSSHConfig(config *host.Host) *ssh.ClientConfig {
	auths := []ssh.AuthMethod{}

	if config.Password() != "" {
		auths = append(auths, ssh.Password(config.Password()))
	}

	if config.IdentityFile != "" {
		if pubkey, err := getKeyFile(config.IdentityFile); err != nil {
			panic(err)
		} else {
			auths = append(auths, ssh.PublicKeys(pubkey))
		}
	}

	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers))
		defer sshAgent.Close()
	}

	return &ssh.ClientConfig{
		Timeout:         10 * time.Second,
		User:            config.User,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

// Conn Connection Method
func Conn(srv *host.Host) (*ssh.Client, *ssh.Session) {
	var client *ssh.Client
	var err error

	config := getSSHConfig(srv)

	if srv.Options["ProxyJump"] != "" {
		proxyConfig := getSSHConfig(sshf.Config.Get(srv.Options["ProxyJump"]))

		proxyClient, err := ssh.Dial("tcp", net.JoinHostPort(sshf.Config.Get(srv.Options["ProxyJump"]).Hostname, sshf.Config.Get(srv.Options["ProxyJump"]).Port), proxyConfig)
		if err != nil {
			panic(err)
		}

		conn, err := proxyClient.Dial("tcp", net.JoinHostPort(srv.Hostname, srv.Port))
		if err != nil {
			panic(err)
		}

		ncc, chans, reqs, err := ssh.NewClientConn(conn, net.JoinHostPort(srv.Hostname, srv.Port), config)
		if err != nil {
			panic(err)
		}

		client = ssh.NewClient(ncc, chans, reqs)
	} else {
		client, err = ssh.Dial("tcp", net.JoinHostPort(srv.Hostname, srv.Port), config)
		if err != nil {
			panic(err)
		}
	}

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}

	return client, session
}

// Shell Remote Terminal Attachment
func Shell(session *ssh.Session) {
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.ECHOCTL:       0,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}

	termFD := int(os.Stdin.Fd())

	w, h, _ := terminal.GetSize(termFD)

	termState, _ := terminal.MakeRaw(termFD)
	defer terminal.Restore(termFD, termState)

	session.RequestPty(os.Getenv("TERM"), h, w, modes)
	session.Shell()
	session.Wait()
}

// RunCmd Run Command in Remote Host
func RunCmd(session *ssh.Session, cmd string) (string, error) {
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
