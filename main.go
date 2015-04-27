package main

import (
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/jmhodges/levigo"
	"github.com/peterh/liner"
)

var (
	historyFile = "/tmp/." + appName + "_history"
	appName     = "relevel"
	appVersion  = "0.0.1"
	appAuthor   = "Akiomi Kamakura"
	appEmail    = "akiomik@gmail.com"
	appUsage    = "REPL for leveldb"
	prompt      = appName + "> "
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = appVersion
	app.Author = appAuthor
	app.Email = appEmail
	app.Usage = appUsage
	app.Commands = Commands
	app.Run(os.Args)
}

func startConsole(exec func(string)) error {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetCompleter(queryCompleter)

	f, err := os.Open(historyFile)
	if err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	showConsoleHeader()

	// main loop
	for {
		name, err := line.Prompt(prompt)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		} else {
			line.AppendHistory(name)
			exec(name)
			println("")
		}

		f, err := os.Create(historyFile)
		if err != nil {
			log.Error("Error writing history file: ", err)
		}

		line.WriteHistory(f)
		f.Close()
	}

	return nil
}

func operateDb(db_name string, opts *levigo.Options, f func(db *levigo.DB) error) error {
	db, err := levigo.Open(db_name, opts)
	if err != nil {
		return err
	}

	defer db.Close()
	return f(db)
}

func showConsoleHeader() {
	println(appName + " - " + appVersion)
	println("")
}
