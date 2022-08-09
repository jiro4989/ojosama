package main

import (
	"fmt"
	"strings"
)

var (
	paramCharCodes = "utf8 sjis"
	paramCompletions = "bash zsh fish"

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
          local opts="`+paramCharCodes+`"
          COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
          ;;
        -completions)
          local opts="`+paramCompletions+`"
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
    {-h,-help}'[`+helpMsgHelp+`]: :->etc' \
    -t'[`+helpMsgText+`]: :->etc' \
    -o'[`+helpMsgOutFile+`]:file:_files' \
    -charcode'[`+helpMsgCharCode+`]: :->charcode' \
    -v'[`+helpMsgVersion+`]: :->etc' \
    -completions'[`+helpMsgCompletions+`]: :->completions'

  case "$state" in
    charcode)
      _values 'charcode' `+paramCharCodes+`
      ;;
    completions)
      _values 'completions' `+paramCompletions+`
      ;;
    etc)
      # nothing to do
      ;;
  esac
}

compdef _{{APPNAME}} {{APPNAME}}

# vim: ft=zsh`, "{{APPNAME}}", appName)

	// -x 引数を受け取るけれどファイルを指定できない
	// -r 引数としてファイルを指定する
	// -a 入力可能な文字列を指定する
	// -o 古いロングオプション(-helpとか)を指定
	completionsFish = strings.ReplaceAll(`complete -c {{APPNAME}} -r

complete -c {{APPNAME}} -o h -d '`+helpMsgHelp+`'
complete -c {{APPNAME}} -o help -d '`+helpMsgHelp+`'
complete -c {{APPNAME}} -o t -x -d '`+helpMsgText+`'
complete -c {{APPNAME}} -o o -r -d '`+helpMsgOutFile+`'
complete -c {{APPNAME}} -o charcode -x -a '`+paramCharCodes+`' -d '`+helpMsgCharCode+`'
complete -c {{APPNAME}} -o v -d '`+helpMsgVersion+`'
complete -c {{APPNAME}} -o completions -x -a '`+paramCompletions+`' -d '`+helpMsgCompletions+`'`,
		"{{APPNAME}}", appName)

	completionsMap = map[string]string{
		"bash": completionsBash,
		"zsh":  completionsZsh,
		"fish": completionsFish,
	}
)

func isSupportedCompletions(sh string) bool {
	sh = strings.ToLower(sh)
	_, ok := completionsMap[sh]
	return ok
}

func printCompletions(sh string) {
	sh = strings.ToLower(sh)
	fmt.Println(completionsMap[sh])
}
