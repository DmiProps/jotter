package cmd

import (
	// System packages.
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"golang.org/x/crypto/ssh/terminal"

	// Jotter packages.
	"github.com/dmiprops/jotter/modules/setting"
	"github.com/dmiprops/jotter/modules/auth"

	// Vendor packages.
	"github.com/urfave/cli"
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
				Name:  "r",
				Usage: "restart the jotter service if running",
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
	}
)

// Command hundlers.

func state(ctx *cli.Context) error {

	fmt.Printf("Jotter version %s\n", setting.AppVer)

	err := checkRun()
	if err != nil {
		fmt.Printf("State %s\n", err.Error())
	} else {
		fmt.Println("State is done")
	}

	fmt.Printf("Listening address %s\n", setting.StoredAdminSettings.Address)
	fmt.Printf("Using databases %s\n", setting.StoredAdminSettings.Database)

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

	fmt.Println("Run command setaddr...")
	fmt.Println("New address: " + ctx.Args().First())

	return nil
}

func getAddr(ctx *cli.Context) error {

	fmt.Println("Run command getaddr...")

	return nil
}

func setDb(ctx *cli.Context) error {

	fmt.Println("Run command setdb...")
	fmt.Printf("New database: %s\n", ctx.Args().First())
	fmt.Printf("Need restart: %t\n", ctx.Bool("r"))

	return nil
}

func getDb(ctx *cli.Context) error {

	fmt.Println("Run command getdb...")

	return nil
}

func start(ctx *cli.Context) error {

	fmt.Println("Run command start...")
	if ctx.IsSet("addr") {
		fmt.Println("By address: %s\n" + ctx.String("addr"))
	}
	if ctx.IsSet("db") {
		fmt.Println("Using database: %s\n" + ctx.String("db"))
	}
	fmt.Printf("Need restart: %t\n", ctx.Bool("r"))
	fmt.Printf("Need save settings: %t\n", ctx.Bool("s"))

	cmd := exec.Command("./jotter", "inner-start")
	cmd.Start()

	return nil
}

func stop(ctx *cli.Context) error {

	fmt.Println("Run command stop...")

	data := []byte(`{"method":"stop"}`)

	rdata := bytes.NewReader(data)
	w, err := http.Post("http://localhost:"+setting.DefaultPort, "application/json", rdata)
	if err != nil {
		fmt.Println("Service stoped")
	} else {
		wdata := responseData{}
		json.NewDecoder(w.Body).Decode(&wdata)
		w.Body.Close()

		fmt.Println("Response: " + wdata.Response)
	}

	if err != nil {
		return err
	}

	return nil
}

// ********* INNER COMMAND ********* //

type requestData struct {
	Method string `json:"method"`
}

type responseData struct {
	Response string `json:"response"`
}

func checkRun() error {
	w, err := http.Get(setting.Protocol + "://" + setting.StoredAdminSettings.Address + "/checkrun")
	if err == nil {
		return json.NewDecoder(w.Body).Decode(&responseData{})
	}
	return err
}

func innerStart(ctx *cli.Context) error {
	handler := func(w http.ResponseWriter, r *http.Request) {
		flusher, _ := w.(http.Flusher)

		rdata := requestData{}
		wdata := responseData{}
		json.NewDecoder(r.Body).Decode(&rdata)
		r.Body.Close()

		if rdata.Method == "" {
			// Request from browser.
			io.WriteString(w, "We working!")
		} else if rdata.Method == "stop" {
			// Request method "stop".
			wdata.Response = "OK, we closed."
			json.NewEncoder(w).Encode(wdata)
			flusher.Flush()
			os.Exit(0)
		} else {
			// Unknown method.
			wdata.Response = "Unknown method: " + rdata.Method + ". We running."
			json.NewEncoder(w).Encode(wdata)
		}
	}

	http.ListenAndServe(setting.StoredAdminSettings.Address, http.HandlerFunc(handler))

	return nil
}
