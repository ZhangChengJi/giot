package main

import (
	"fmt"
	"giot/cmd/scheduler/app"

	"os"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download
func main() {
	if err := app.NewSchedulerCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
