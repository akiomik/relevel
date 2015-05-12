package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/jmhodges/levigo"
)

var (
	appName    = "relevel"
	appVersion = "0.0.1"
	appAuthor  = "Akiomi Kamakura"
	appEmail   = "akiomik@gmail.com"
	appUsage   = "REPL for leveldb"
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

func operateDb(dbName string, opts *levigo.Options, f func(db *levigo.DB) error) error {
	db, err := levigo.Open(dbName, opts)
	if err != nil {
		return err
	}

	defer db.Close()
	return f(db)
}
