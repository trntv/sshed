sshed - ssh connections editor and bookmarks manager
---
Visual cross-platform editor created to manage list of ssh hosts in ssh config file.
sshed uses native ``ssh_config`` format to store connections information and supports all available ssh options.

[![asciicast](https://asciinema.org/a/164783.png)](https://asciinema.org/a/164783)

# Installation
download binary [here](https://github.com/trntv/sshed/releases) 
or run command (make sure to change X.X.X to real version)
```
curl -L -s https://github.com/trntv/sshed/releases/download/X.X.X/sshed-X.X.X-linux-amd64
```
or install with ``go get``
```
go get -u github.com/trntv/sshed
```
or compile it from source
```
git clone https://github.com/trntv/sshed.git
cd sshed
make
```
install with brew
```
brew install https://raw.githubusercontent.com/trntv/sshed/master/sshed.rb
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

# Make local 
```
find . -name "*.go" | xargs sed -i 's|github.com/trntv/sshed|..|g'
make
```

# TODO
 - [x] ``sshed at`` - executes command on server
 - [x] batch commands
 - [x] ssh_config integration
 - [ ] ssh options (-c, -E, -f, -T, -t)
 - [ ] key, password generation
 - [x] bind address
 - [x] replace sshpass with native go implementation
 - [ ] scp
 - [x] ssh bin flag
 - [x] autocompletion
 - [ ] backup
 - [ ] restore
 - [x] jump hosts
