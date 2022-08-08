package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type CmdArgs struct {
	Text        string
	OutFile     string
	Version     bool
	CharCode    string
	Completions string
	Args        []string
}

const (
	helpMsgHelp        = "print help"
	helpMsgText        = "input text"
	helpMsgOutFile     = "output file"
	helpMsgCharCode    = "input text file encoding. default is utf8. (utf8, sjis)"
	helpMsgVersion     = "print version"
	helpMsgCompletions = "print completions file. (bash, zsh)"
)

func ParseArgs() (*CmdArgs, error) {
	opts := CmdArgs{}

	flag.Usage = flagHelpMessage
	flag.StringVar(&opts.Text, "t", "", helpMsgText)
	flag.StringVar(&opts.OutFile, "o", "", helpMsgOutFile)
	flag.StringVar(&opts.CharCode, "charcode", "utf8", helpMsgCharCode)
	flag.BoolVar(&opts.Version, "v", false, helpMsgVersion)
	flag.StringVar(&opts.Completions, "completions", "", helpMsgCompletions)
	flag.Parse()
	opts.Args = flag.Args()

	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &opts, nil
}

func flagHelpMessage() {
	cmd := os.Args[0]
	fmt.Fprintln(os.Stderr, fmt.Sprintf("%s convert text to '%s' style.", cmd, appName))
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
	switch c.CharCode {
	case "utf8", "sjis":
		// 何もしない
	default:
		err := errors.New("charcode must be 'utf8' or 'sjis'.")
		return err
	}

	if c.Completions != "" && !isSupportedCompletions(c.Completions) {
		return fmt.Errorf("illegal completions. completions = %s", c.Completions)
	}

	return nil
}
