package ssh

import (
	"bytes"
	"fmt"
	"github.com/kevinburke/ssh_config"
	"github.com/trntv/sshed/host"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
)

var Config *sshConfig
var maskPatternRegexp = regexp.MustCompile(`[\*!\?]`)

type ErrNotFound struct {
	Key string
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("Host %s not found", err.Key)
}

type sshConfig struct {
	Hosts   []string
	Keys    []string
	Content []byte

	Path string
	cfg  *ssh_config.Config
}

func Parse(path string, nonStrict bool) (conf *sshConfig, err error) {
	conf = &sshConfig{Path: path}

	if _, err := os.Stat(conf.Path); os.IsNotExist(err) == false {
		conf.Content, err = ioutil.ReadFile(conf.Path)
		if err != nil {
			return nil, err
		}
	}

	conf.cfg, err = ssh_config.Decode(bytes.NewReader(conf.Content))
	if err != nil {
		return nil, err
	}

	for _, h := range conf.cfg.Hosts {
		if h.EOLComment == " -- added by sshed" || nonStrict == true {
			for _, pattern := range h.Patterns {
				if maskPatternRegexp.MatchString(pattern.String()) == false {
					conf.Hosts = append(conf.Hosts, pattern.String())
				}
			}

			for _, node := range h.Nodes {
				switch node.(type) {
				case *ssh_config.KV:
					key := node.(*ssh_config.KV).Key
					if key != "IdentityFile" {
						continue
					}

					path = node.(*ssh_config.KV).Value
					exists := false
					for _, v := range conf.Keys {
						v = convertTilde(v)
						path = convertTilde(path)
						if v == path {
							exists = true
							break
						}
					}
					if exists == false {
						conf.Keys = append(conf.Keys, path)
					}
				}
			}
		}
	}

	return conf, nil
}

func (s *sshConfig) Get(k string) (h *host.Host) {
	for _, v := range s.cfg.Hosts {
		for _, pattern := range v.Patterns {
			if pattern.String() == k {
				h := host.CreateFromConfig(v)
				h.Key = k
				return h
			}
		}
	}

	return
}

func (s *sshConfig) GetAll() map[string]*host.Host {
	hs := make(map[string]*host.Host)
	for _, v := range s.cfg.Hosts {
		for _, pattern := range v.Patterns {
			if maskPatternRegexp.MatchString(pattern.String()) == false {
				hh := host.CreateFromConfig(v)
				hh.Key = pattern.String()
				hs[hh.Key] = hh
			}
		}
	}

	return hs
}

func (s *sshConfig) Add(h *host.Host) error {
	s.Remove(h.Key)

	pattern, err := ssh_config.NewPattern(h.Key)
	if err != nil {
		return err
	}

	nodes := make([]ssh_config.Node, 0)

	if h.Hostname != "" {
		nodes = append(nodes, &ssh_config.KV{Key: "  HostName", Value: h.Hostname})
	}

	if h.Port != "" {
		nodes = append(nodes, &ssh_config.KV{Key: "  Port", Value: h.Port})
	}

	if h.User != "" {
		nodes = append(nodes, &ssh_config.KV{Key: "  User", Value: h.User})
	}

	if h.IdentityFile != "" {
		nodes = append(nodes, &ssh_config.KV{Key: "  IdentityFile", Value: h.IdentityFile})
	}

	for key, option := range h.Options {
		if option == "" {
			continue
		}
		nodes = append(nodes, &ssh_config.KV{Key: fmt.Sprintf("  %s", key), Value: option})
	}

	nodes = append(nodes, &ssh_config.Empty{})

	s.cfg.Hosts = append(s.cfg.Hosts, &ssh_config.Host{
		Patterns:   []*ssh_config.Pattern{pattern},
		Nodes:      nodes,
		EOLComment: " -- added by sshed",
	})

	return nil
}

func (s *sshConfig) Remove(p string) {
	var hosts = make([]*ssh_config.Host, 0)

	for _, h := range s.cfg.Hosts {
		skip := false
		for _, pattern := range h.Patterns {
			if maskPatternRegexp.MatchString(pattern.String()) != false {
				continue
			}

			if pattern.String() == p {
				skip = true
				break
			}
		}
		if skip != true {
			hosts = append(hosts, h)
		}
	}

	s.cfg.Hosts = hosts
}

func (s *sshConfig) Save() (err error) {
	return s.SaveContent([]byte(s.cfg.String()))
}

func (s *sshConfig) SaveContent(data []byte) (err error) {
	err = ioutil.WriteFile(s.Path, data, os.FileMode(0644))
	if err != nil {
		return err
	}

	return err
}

func (s *sshConfig) String() string {
	return s.cfg.String()
}

func convertTilde(path string) string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	if len(path) > 2 && path[:2] == "~/" {
		path = filepath.Join(homeDir, path[2:])
	}

	return path
}
