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

func ParseArgs() (*CmdArgs, error) {
	opts := CmdArgs{}

	flag.Usage = flagHelpMessage
	flag.StringVar(&opts.Text, "t", "", "input text")
	flag.StringVar(&opts.OutFile, "o", "", "output file")
	flag.StringVar(&opts.CharCode, "charcode", "utf8", "input text file encoding. (utf8 or sjis)")
	flag.BoolVar(&opts.Version, "v", false, "print version")
	flag.StringVar(&opts.Completions, "completions", "", "print completions file. (bash)")
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
