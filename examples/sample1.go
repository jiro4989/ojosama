package main

import (
	"fmt"

	"github.com/jiro4989/ojosama"
)

func main() {
	s := "ハーブがありました！"
	text, err := ojosama.Convert(s, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
