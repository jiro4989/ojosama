package main

import (
	"fmt"
	"os"

	"github.com/jiro4989/ojosama"
)

const (
	exitStatusOK = iota
	exitStatusCLIError
	exitStatusConvertError
)

func main() {
	args, err := ParseArgs()
	if err != nil {
		Err(err)
		os.Exit(exitStatusCLIError)
	}

	if args.Text != "" {
		text, err := ojosama.Convert(args.Text, nil)
		if err != nil {
			Err(err)
			os.Exit(exitStatusConvertError)
		}
		fmt.Println(text)
		os.Exit(exitStatusOK)
	}
}
