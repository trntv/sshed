#!/usr/bin/env bash

: ${PROG:=sshed}

if [ -n "$ZSH_VERSION" ]; then
  autoload -U compinit && compinit
  autoload -U bashcompinit && bashcompinit
fi


_cli_bash_autocomplete() {
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
}

complete -F _cli_bash_autocomplete $PROG

unset PROG
