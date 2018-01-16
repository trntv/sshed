sshdb - SSH connections manager
---
Simple program created to manage list of ssh connections.

![Interface](gui.gif)

# Installation
download binary
```
curl -L -s https://github.com/trntv/sshdb/releases/download/0.3.0/sshdb -o sshdb
```
or install with ``go get``
```
go get -u github.com/trntv/sshdb
```
or compile it from source
```
git clone https://github.com/trntv/sshdb.git
cd sshdb
make
```
install with brew
```
brew install https://raw.githubusercontent.com/trntv/sshdb/master/sshdb.rb
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
   sshdb - SSH connections manager

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
   --database value, --db value  Path to database file (default: "$HOME/.sshdbb") [$SSHED_DB_PATH]
   --help, -h                    show help
   --version, -v                 print the version

```

# Tips
to store passwords you need to install sshpass that allows to 
offer a password via SSH

to install it with brew use
```
brew install http://git.io/sshpass.rb
```
for other package managers see: [https://github.com/kevinburke/sshpass](https://github.com/kevinburke/sshpass)

# Similar projects
Searching for such tool i've found some similar projects but ended up writing my own solution:
 - [https://github.com/mmeyer724/sshmenu](https://github.com/mmeyer724/sshmenu)    
 - [https://github.com/vaniacer/sshto](https://github.com/vaniacer/sshto)
 - [https://github.com/xiongharry/sshtoy](https://github.com/xiongharry/sshtoy)
 - [https://github.com/sciancio/connectionmanager2](https://github.com/sciancio/connectionmanager2)
 - etc.
 
# TODO
 - [ ] ``sshed at`` - executes command on server
 - [ ] backup
 - [ ] restore
 - [ ] manage ssh config (view, edit)
 - [ ] additional arguments to ssh command
 - [ ] key, password generation
 - [ ] bind address
 - [ ] replace sshpass with native go implementation
 - [ ] scp
 - [x] executes command on given server
