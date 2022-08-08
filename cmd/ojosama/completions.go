package main

import (
	"fmt"
	"strings"
)

var (
	completionsBash = strings.ReplaceAll(`# {{APPNAME}}(1) completion                                       -*- shell-script -*-

_{{APPNAME}}_module() {
  local cur prev cword
  _get_comp_words_by_ref -n : cur prev cword

  case "${cword}" in
    1)
      local opts="-h -help -t -o -charcode -v -completions"
      COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
      ;;
    2)
      case "${prev}" in
        -o)
          COMPREPLY=($(compgen -f -- "${cur}"))
          ;;
        -charcode)
          local opts="utf8 sjis"
          COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
          ;;
        -completions)
          local opts="bash zsh"
          COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
          ;;
      esac
      ;;
  esac
}

complete -F _{{APPNAME}}_module {{APPNAME}}`, "{{APPNAME}}", appName)

	completionsZsh = strings.ReplaceAll(`#compdef {{APPNAME}}

_{{APPNAME}}() {
  _arguments \
    {-h,-help}'[print help]: :->etc' \
    {-t}'[input text]: :->etc' \
    {-o}'[output file]:file:_files' \
    {-charcode}'[input text file encoding. (utf8, sjis)]: :->charcode' \
    {-v}'[print version]: :->etc' \
    {-completions}'[print completions file. (bash, zsh)]: :->completions'

  case "$state" in
    charcode)
      _values 'charcode' utf8 sjis
      ;;
    completions)
      _values 'completions' bash zsh
      ;;
    etc)
      # nothing to do
      ;;
  esac
}

compdef _{{APPNAME}} {{APPNAME}}

# vim: ft=zsh`, "{{APPNAME}}", appName)

	completionsMap = map[string]string{
		"bash": completionsBash,
		"zsh": completionsZsh,
	}
)

func isSupportedCompletions(sh string) bool {
	_, ok := completionsMap[sh]
	return ok
}

func printCompletions(sh string) {
	sh = strings.ToLower(sh)
	fmt.Println(completionsMap[sh])
}
