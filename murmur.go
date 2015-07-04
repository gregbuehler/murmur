package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"

	"github.com/gregbuehler/murmur/cmd"
	"github.com/gregbuehler/murmur/modules/setting"
)

const (
	APP_VER = "0.0.0"
)

func init() {
	// TODO: This might be excisive and should probably be configurable
	runtime.GOMAXPROCS(runtime.NumCPU())

	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "Murmur"
	app.Usage = "A vector driven time series database"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.CmdServer,
		// TODO: cmd.CmdDump,
		// TODO: cmd.CmdRestore,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
