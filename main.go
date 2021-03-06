package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/sergeyfast/stubsrv/app"
	"github.com/sergeyfast/stubsrv/db"
	"gopkg.in/pg.v5"
	"io/ioutil"
	"log"
	"os"
)

const appName = "stubsrv"

var (
	flVerbose    = flag.Bool("verbose", false, "enable debug output")
	flConfigPath = flag.String("config", "config.cfg", "Path to config file")
	cfg          config
	version      string
)

// Application config
type config struct {
	Server   app.HttpConfig
	Database *pg.Options
}

func main() {
	flag.Parse()
	fixStdLog(*flVerbose)

	log.Printf("starting %v version=%v", appName, version)
	if _, err := toml.DecodeFile(*flConfigPath, &cfg); err != nil {
		die(err)
	}

	dbc := pg.Connect(cfg.Database)
	v, err := db.New(dbc).Version()
	die(err)
	log.Println(v)

	a := app.New(appName, *flVerbose, cfg.Server, dbc)
	die(a.Run())
}

// fixStdLog sets additional params to std logger (prefix D, filename & line).
func fixStdLog(verbose bool) {
	log.SetPrefix("D")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if verbose {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

// die calls log.Fatal if err wasn't nil.
func die(err error) {
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
}
