package cmd

import (
	// System packages.
	"fmt"
	"net/http"
	"golang.org/x/crypto/ssh/terminal"
	"strings"

	// Jotter packages.
	"github.com/dmiprops/jotter/modules/setting"
	"github.com/dmiprops/jotter/modules/auth"
	"github.com/dmiprops/jotter/modules/daemon"
	"github.com/dmiprops/jotter/modules/log"
	"github.com/dmiprops/jotter/handlers"

	// Vendor packages.
	"github.com/urfave/cli"
	"github.com/gorilla/mux"
)

var (
	// CmdState shows current state jotter service.
	CmdState = cli.Command{
		Name:      "state",
		Usage:     "Show current state jotter service",
		Action:    state,
		ArgsUsage: " ",
	}

	// CmdSetPass sets the administrative password jotter service.
	CmdSetPass = cli.Command{
		Name:      "setpass",
		Usage:     "Set the administrative password jotter service",
		Action:    setPass,
		ArgsUsage: " ",
	}

	// CmdSetAddr sets address for listening jotter service.
	CmdSetAddr = cli.Command{
		Name:      "setaddr",
		Usage:     "Set address for listening jotter service",
		Action:    setAddr,
		ArgsUsage: "[host]:port (default: \":" + setting.DefaultPort + ")\"",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "r",
				Usage: "restart the jotter service if running",
			},
		},
	}

	// CmdGetAddr shows set address and current address, if they are different.
	CmdGetAddr = cli.Command{
		Name:      "getaddr",
		Usage:     "Show set address and current address, if they are different",
		Action:    getAddr,
		ArgsUsage: " ",
	}

	// CmdSetDb sets database settings for usage jotter service.
	CmdSetDb = cli.Command{
		Name:      "setdb",
		Usage:     "Set database settings for usage jotter service",
		Action:    setDb,
		ArgsUsage: "user:password@host[:port]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "r",
				Usage: "restart the jotter service if running",
			},
		},
	}

	// CmdGetDb shows database settings for usage jotter service.
	CmdGetDb = cli.Command{
		Name:      "getdb",
		Usage:     "Show database settings for usage jotter service",
		Action:    getDb,
		ArgsUsage: " ",
	}

	// CmdStart starts jotter service.
	CmdStart = cli.Command{
		Name:      "start",
		Usage:     "Start jotter service",
		Action:    start,
		ArgsUsage: " ",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "addr",
				Value: ":" + setting.DefaultPort,
				Usage: "jotter service address [host]:ip",
			},
			cli.StringFlag{
				Name:  "db",
				Value: "",
				Usage: "database settings user:password@host:[port]",
			},
			cli.BoolFlag{
				Name:  "s",
				Usage: "save jotter service address and database settings",
			},
		},
	}

	// CmdStop stops jotter service.
	CmdStop = cli.Command{
		Name:      "stop",
		Usage:     "Stop jotter service",
		Action:    stop,
		ArgsUsage: " ",
	}

	// CmdInnerStart executes start web (hidden command for inner usage).
	CmdInnerStart = cli.Command{
		Name:      "inner-start",
		Usage:     "",
		Action:    innerStart,
		ArgsUsage: " ",
		Hidden:    true,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "addr",
				Value: "",
				Usage: "jotter service address [host]:ip",
			},
			cli.StringFlag{
				Name:  "db",
				Value: "",
				Usage: "database settings user:password@host:[port]",
			},
		},
	}
)

// User command handlers.

func state(ctx *cli.Context) error {
	fmt.Printf("Version\t\t%s\n", setting.AppVer)

	isRunning := daemon.CheckRun()
	fmt.Print("State\t\t")
	if isRunning {
		fmt.Println("Run")
	} else {
		fmt.Println("Shutdown")
	}

	setting.ReadCurrentAdminSettings()
	showAddr(isRunning)
	showDb(isRunning)

	return nil
}

func setPass(ctx *cli.Context) error {
	fmt.Print("Enter new administrative password: ")
	password, err := terminal.ReadPassword(1)
	if err != nil {
		return err
	}
	fmt.Println("")

	hash, err := auth.HashPassword(string(password))
	if err != nil {
		return err
	}
	setting.StoredAdminSettings.Password = hash

	return setting.SaveStoredAdminSettings()
}

func setAddr(ctx *cli.Context) error {
	setting.StoredAdminSettings.Address = ctx.Args().First()

	err := setting.SaveStoredAdminSettings()
	if err != nil {
		return err
	}

	if ctx.Bool("r") {
		return daemon.RestartDaemon("inner-start")
	}
	return nil
}

func getAddr(ctx *cli.Context) error {
	setting.ReadCurrentAdminSettings()
	showAddr(daemon.CheckRun())
	return nil
}

func setDb(ctx *cli.Context) error {
	setting.StoredAdminSettings.Database = ctx.Args().First()

	err := setting.SaveStoredAdminSettings()
	if err != nil {
		return err
	}

	if ctx.Bool("r") {
		return daemon.RestartDaemon("inner-start")
	}
	return nil
}

func getDb(ctx *cli.Context) error {
	setting.ReadCurrentAdminSettings()
	showDb(daemon.CheckRun())
	return nil
}

func start(ctx *cli.Context) error {
	if daemon.CheckRun() {
		err := daemon.StopDaemon()
		if err != nil {
			return err
		}
	}

	args := []string{"inner-start"}
	save := ctx.Bool("s")

	if ctx.IsSet("addr") {
		args = append(args, []string{"--addr", ctx.String("addr")}...)
		if save {
			setting.StoredAdminSettings.Address = ctx.String("addr")
		}
	}
	if ctx.IsSet("db") {
		args = append(args, []string{"--db", ctx.String("db")}...)
		if save {
			setting.StoredAdminSettings.Database = ctx.String("db")
		}
	}

	if save {
		err := setting.SaveStoredAdminSettings()
		if err != nil {
			return err
		}
	}

	return daemon.StartDaemon(args ...)
}

func stop(ctx *cli.Context) error {
	return daemon.StopDaemon()
}

// Inner command handlers.

func innerStart(ctx *cli.Context) error {
	// Prepare arguments.
	log.Info("inner-start: ", ctx.Args())
	if ctx.IsSet("addr") {
		setting.CurrentAdminSettings.Address = ctx.String("addr")
	} else {
		setting.CurrentAdminSettings.Address = setting.StoredAdminSettings.Address
	}
	if ctx.IsSet("db") {
		setting.CurrentAdminSettings.Database = ctx.String("db")
	} else {
		setting.CurrentAdminSettings.Database = setting.StoredAdminSettings.Database
	}
	err := setting.SaveCurrentAdminSettings()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Output log information.
	log.Info(`Start jotter service:
- listening address %s
- using database %s`,
		setting.CurrentAdminSettings.Address,
		setting.CurrentAdminSettings.Database,
	)

	// Start listener.
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.MainPage)

	http.Handle("/", r)
	http.ListenAndServe(setting.CurrentAdminSettings.Address, nil)

	return nil
}

// Service functions.

func showAddr(isRunning bool) {
	fmt.Printf("Address\t\t%s\n", setting.StoredAdminSettings.Address)
	if (isRunning && strings.ToLower(setting.StoredAdminSettings.Address) != strings.ToLower(setting.CurrentAdminSettings.Address)) {
		fmt.Printf("- used\t\t%s\n", setting.CurrentAdminSettings.Address)
	}
}

func showDb(isRunning bool) {
	storedDatabase := setting.ConnectionStringWithoutPassword(setting.StoredAdminSettings.Database)
	if storedDatabase == "" {
		storedDatabase = "<undefined>"
	}
	currentDatabase := setting.ConnectionStringWithoutPassword(setting.CurrentAdminSettings.Database)
	if currentDatabase == "" {
		currentDatabase = "<undefined>"
	}

	fmt.Printf("Database\t%s\n", storedDatabase)

	if (isRunning && strings.ToLower(storedDatabase) != strings.ToLower(currentDatabase)) {
		fmt.Printf("- used\t\t%s\n", currentDatabase)
	}
}