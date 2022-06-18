package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jiro4989/ojosama"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// CIでビルド時に値を埋め込む。
// 埋め込む値の設定は .goreleaser.yaml を参照。
var (
	version  = "dev"
	revision = "dev"
)

const (
	exitStatusOK = iota
	exitStatusCLIError
	exitStatusConvertError
	exitStatusInputFileError
	exitStatusOutputError
)

func main() {
	args, err := ParseArgs()
	if err != nil {
		Err(err)
		os.Exit(exitStatusCLIError)
	}

	if args.Version {
		msg := fmt.Sprintf("ojosama %s (%s)", version, revision)
		fmt.Println(msg)
		fmt.Println("")
		fmt.Println("author:     jiro")
		fmt.Println("repository: https://github.com/jiro4989/ojosama")
		os.Exit(exitStatusOK)
	}

	if args.Text != "" {
		exitStatus, err := run(args.Text, args)
		if err != nil {
			Err(err)
			os.Exit(exitStatus)
		}
		os.Exit(exitStatus)
	}

	if len(args.Args) < 1 {
		// SJIS指定の時だけSJISとして読み込む
		var r io.Reader = os.Stdin
		if args.CharCode == "sjis" {
			r = transform.NewReader(os.Stdin, japanese.ShiftJIS.NewDecoder())
		}

		b, err := io.ReadAll(r)
		if err != nil {
			Err(err)
			os.Exit(exitStatusInputFileError)
		}

		s := string(b)
		exitStatus, err := run(s, args)
		if err != nil {
			Err(err)
			os.Exit(exitStatus)
		}
		os.Exit(exitStatus)
	}

	for _, f := range args.Args {
		f, err := os.Open(f)
		if err != nil {
			Err(err)
			os.Exit(exitStatusOutputError)
		}
		defer f.Close()

		// SJIS指定の時だけSJISとして読み込む
		var r io.Reader = f
		if args.CharCode == "sjis" {
			r = transform.NewReader(f, japanese.ShiftJIS.NewDecoder())
		}

		b, err := io.ReadAll(r)
		if err != nil {
			Err(err)
			os.Exit(exitStatusInputFileError)
		}

		s := string(b)
		exitStatus, err := run(s, args)
		if err != nil {
			Err(err)
			os.Exit(exitStatus)
		}
	}

	os.Exit(exitStatusOK)
}

func run(s string, args *CmdArgs) (int, error) {
	text, err := ojosama.Convert(s, nil)
	if err != nil {
		return exitStatusConvertError, err
	}

	out := os.Stdout
	if args.OutFile != "" {
		out, err = os.Create(args.OutFile)
		if err != nil {
			return exitStatusOutputError, err
		}
	}
	out.WriteString(text)
	return exitStatusOK, nil
}
