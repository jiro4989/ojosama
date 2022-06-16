package main

import (
	"fmt"
	"os"
)

func Err(err error) {
	fmt.Fprintln(os.Stderr, "[ERR] "+err.Error())
}
