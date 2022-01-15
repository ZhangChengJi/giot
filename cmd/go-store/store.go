package main

import (
	"fmt"
	"giot/cmd/go-store/app"

	"os"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download
func main() {
	if err := app.NewStoreCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
