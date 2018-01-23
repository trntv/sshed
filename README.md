sshme - ssh bookmarks manager
---
Simple program created to manage bookmarks for ssh connections. Alternative to native ``~/.ssh/config`` approach

![Interface](gui.gif)

# Installation
download binary [here](https://github.com/trntv/sshme/releases) 
or run command (make sure to change X.X.X to real version)
```
curl -L -s https://github.com/trntv/sshme/releases/download/X.X.X/sshme-X.X.X-linux-amd64
```
or install with ``go get``
```
go get -u github.com/trntv/sshme
```
or compile it from source
```
git clone https://github.com/trntv/sshme.git
cd sshme
make
```
install with brew
```
brew install https://raw.githubusercontent.com/trntv/sshme/master/sshme.rb
```

# Features
- add, show, list, remove ssh connections
- supported fields
    - Host
    - Port
    - User
    - Password
    - Key File
- connect to server by key
- database encryption

# Usage
```
NAME:
   sshme - SSH connections manager

USAGE:
   help [global options] command [command options] [arguments...]

VERSION:
   X.X.X (build xxxxxx)

AUTHOR:
   Eugene Terentev <eugene@terentev.net>

COMMANDS:
     show     show server information
     list     list all servers from database
     add      adds server to database
     remove   removes server from database
     to       connects to server
     at       executes command on given server
     encrypt  encrypt database
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --database value, --db value  Path to database file (default: "$HOME/.sshdb") [$SSHME_DB_PATH]
   --help, -h                    show help
   --version, -v                 print the version

```

# Bash (ZSH) autocomplete
to enable autocomplete run
```
PROG=sshme source completions/sshme.bash
```

# Tips
to store passwords you need to install sshpass that allows to 
offer a password via SSH

to install it with brew use
```
brew install http://git.io/sshpass.rb
```
for other package managers see: [https://github.com/kevinburke/sshpass](https://github.com/kevinburke/sshpass)

# Native ssh built-in bookmarks
[https://blog.viktorpetersson.com/2010/12/05/ssh-tips-how-to-create-ssh-bookmarks.html](https://blog.viktorpetersson.com/2010/12/05/ssh-tips-how-to-create-ssh-bookmarks.html)

# Similar projects
Searching for such tool i've found some similar projects but ended up writing my own solution:
 - [https://github.com/mmeyer724/sshmenu](https://github.com/mmeyer724/sshmenu)    
 - [https://github.com/vaniacer/sshto](https://github.com/vaniacer/sshto)
 - [https://github.com/xiongharry/sshtoy](https://github.com/xiongharry/sshtoy)
 - [https://github.com/sciancio/connectionmanager2](https://github.com/sciancio/connectionmanager2)
 - [https://github.com/andreyantipov/ssh-cli-bookmarks](https://github.com/andreyantipov/ssh-cli-bookmarks)
 - etc.
 
# TODO
 - [x] ``sshme at`` - executes command on server
 - [ ] backup
 - [ ] restore
 - [ ] manage ssh config (view, edit)
 - [ ] additional arguments to ssh command
 - [ ] key, password generation
 - [ ] bind address
 - [ ] replace sshpass with native go implementation
 - [ ] scp
 - [ ] batch commands
 - [x] autocompletion
