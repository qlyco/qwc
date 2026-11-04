package main

import (
	"github.com/qlyco/qwc/cli"
	"github.com/qlyco/qwc/internal/res"
)

var Version = res.VersionDefault

func main() {
	cli.Run(cli.App{Name: res.AppName, Version: Version})
}
