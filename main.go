//
// Jotter is a system for logging changes to settings, data, and user requests made to clients.
//
package main

import (
	// System packages.
	"os"

	// Jotter packages.
	"github.com/dmiprops/jotter/cmd"
	"github.com/dmiprops/jotter/modules/log"
	"github.com/dmiprops/jotter/modules/setting"

	// Vendor packages.
	"github.com/urfave/cli"
)

// Version holds the current Jotter version.
const Version = "0.1.0-dev"

func init() {
	setting.AppVer = Version
	setting.Protocol = "http"
}

func main() {
	app := cli.NewApp()
	app.Name = "Jotter"
	app.Author = "DmiProps"
	app.Email = "dmi.develop@gmail.com"
	app.Usage = "A system for logging changes to settings, data, and user requests."
	app.Description = `By default, Jotter will start serving using the webserver with no
arguments - which can alternatively be run by running the subcommand web.`
	app.Version = Version
	app.Commands = []cli.Command{
		cmd.CmdState,
		cmd.CmdSetPass,
		cmd.CmdSetAddr,
		cmd.CmdGetAddr,
		cmd.CmdSetDb,
		cmd.CmdGetDb,
		cmd.CmdStart,
		cmd.CmdStop,
		cmd.CmdInnerStart,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Failed to run app with %s: %v", os.Args, err)
	}
}
