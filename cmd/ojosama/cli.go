package main

import (
	"flag"
	"fmt"
	"os"
)

type CmdArgs struct {
	Text    string
	OutFile string
	Args    []string
}

func ParseArgs() (*CmdArgs, error) {
	opts := CmdArgs{}

	flag.Usage = flagHelpMessage
	flag.StringVar(&opts.Text, "t", "", "input text")
	flag.StringVar(&opts.OutFile, "o", "", "output file")
	flag.Parse()
	opts.Args = flag.Args()

	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &opts, nil
}

func flagHelpMessage() {
	cmd := os.Args[0]
	fmt.Fprintln(os.Stderr, fmt.Sprintf("%s changes file permissions with a slot", cmd))
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, fmt.Sprintf("  %s [OPTIONS] [files...]", cmd))
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, fmt.Sprintf("  %s sample.txt", cmd))
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Options:")

	flag.PrintDefaults()
}

func (c *CmdArgs) Validate() error {
	return nil
}
