package main

import (
	"fmt"
	"giot/cmd/giot/app"
	"os"
)

func main() {
	if err := app.NewGiotCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
