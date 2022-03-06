# Stand with Ukraine!

---
<p align="center">
   <img align="center" src="https://github.blog/wp-content/uploads/2022/03/1200x630-GitHub-1.png?resize=320%2C240">
</p>
While Russia is destroying my home and my country, killing my friends and neighbors - any russian company, organization, or citizen, who do nothing about it,
is prohibited from using this package.
For others - please, pray for us, share information about war crimes Russia is conducting in Ukraine, do everything you can
to urge your governments to be on the right side of history.
Ukraine will prevail! Good triumph over evil! Русский военный корабль, иди нах#й!

---

SSHed - ssh connections manager and config editor
---
Visual cross-platform editor created to manage list of SSH hosts in ssh config file and execute commands on those hosts.
sshed uses native ``ssh_config`` format to store connections information and supports all available ssh options.

[![asciicast](https://asciinema.org/a/164783.svg)](https://asciinema.org/a/164783)

# Installation
install with brew
```
brew install trntv/sshed/sshed
```
install with scoop
```
scoop bucket add taskctl https://github.com/trntv/scoop-sshed.git
scoop install sshed
```
download binary [here](https://github.com/trntv/sshed/releases)
or run in console
```
curl -sf https://gobinaries.com/trntv/sshed | sh
```
or install with ``go get``
```
go get -u github.com/trntv/sshed
```

# Features
- add, show, list, remove ssh hosts in ssh_config file
- show, edit ssh config via preferred text editor
- connect to host by key
- execute commands via ssh (on single or multiple hosts)
- encrypted keychain to store ssh passwords and private keys

# Usage
```
NAME:
   sshed - SSH config editor and hosts manager

USAGE:
   help [global options] command [command options] [arguments...]

VERSION:
   X.X.X

AUTHOR:
   Eugene Terentev <eugene@terentev.net>

COMMANDS:
     show     Shows host
     list     Lists all hosts
     add      Add or edit host
     remove   Removes host
     to       Connects to host
     at       Executes commands
     encrypt  Encrypts keychain
     config   Shows SSH config
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --keychain value  path to keychain database (default: "/Users/e.terentev/.sshed") [$SSHED_KEYCHAIN]
   --config value    path to SSH config file (default: "/Users/e.terentev/.ssh/config") [$SSHED_CONFIG_FILE]
   --bin value       path to SSH binary (default: "ssh") [$SSHED_BIN]
   --help, -h        show help
   --version, -v     print the version
```

# Bash (ZSH) autocomplete
to enable autocomplete run
```
PROG=sshed source completions/autocomplete.sh
```
if installed with brew, just add those lines to ``.bash_profile`` (``.zshrc``) file
```
PROG=sshed source $(brew --prefix sshed)/autocomplete.sh
```

# Tips
1. to store passwords you need to install sshpass that allows to offer a password via SSH

    to install it with brew use
    ```
    brew install http://git.io/sshpass.rb
    ```
    for other options see: [https://github.com/kevinburke/sshpass](https://github.com/kevinburke/sshpass)

2. To see all available ssh options run ``man ssh_config``

# TODO
 - [x] ``sshed at`` - executes command on server
 - [x] batch commands
 - [x] ssh_config integration
 - [ ] ssh options (-c, -E, -f, -T, -t)
 - [ ] key, password generation
 - [x] bind address
 - [ ] replace sshpass with native go implementation
 - [ ] scp
 - [x] ssh bin flag
 - [x] autocompletion
 - [ ] backup
 - [ ] restore
 - [ ] jump hosts
