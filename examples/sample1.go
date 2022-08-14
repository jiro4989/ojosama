package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jiro4989/ojosama"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	s := "ハーブがありました！"
	text, err := ojosama.Convert(s, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
