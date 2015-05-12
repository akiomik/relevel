package main

import (
	"os"

	"github.com/codegangsta/cli"
	log "github.com/golang/glog"
	"github.com/jmhodges/levigo"
)

var Commands = []cli.Command{
	commandConsole,
	commandExec,
	commandNew,
	commandVersion,
}

var commandConsole = cli.Command{
	Name:      "console",
	ShortName: "c",
	Usage:     "Starts leveldb console",
	Action:    doConsole,
}

var commandExec = cli.Command{
	Name:      "exec",
	ShortName: "e",
	Usage:     "Executes leveldb queries",
	Action:    doExec,
}

var commandNew = cli.Command{
	Name:      "new",
	ShortName: "n",
	Usage:     "Creates new db",
	Action:    doNew,
}

var commandVersion = cli.Command{
	Name:      "version",
	ShortName: "v",
	Usage:     "Shows version",
	Action:    cli.ShowVersion,
}

func doConsole(c *cli.Context) {
	dbName := c.Args().First()
	if dbName == "" {
		log.Error("DB name is not provided", dbName)
		os.Exit(1)
	}

	opts := levigo.NewOptions()
	err := operateDb(dbName, opts, func(db *levigo.DB) error {
		err := startConsole(func(line string) {
			query, args := queryParser(line)
			queryHandler(query, args, db)
		})
		if err != nil {
			log.Error("Error reading line: ", err)
			os.Exit(1)
		}

		return nil
	})
	if err != nil {
		log.Error("Error db operation: ", err)
		os.Exit(1)
	}
}

func doExec(c *cli.Context) {
	dbName := c.Args().Get(0)
	if dbName == "" {
		log.Error("DB is not provied", dbName)
		os.Exit(1)
	}

	line := c.Args().Get(1)
	if line == "" {
		log.Error("Query is not provied", line)
		os.Exit(1)
	}

	opts := levigo.NewOptions()
	err := operateDb(dbName, opts, func(db *levigo.DB) error {
		query, args := queryParser(line)
		queryHandler(query, args, db)
		return nil
	})
	if err != nil {
		log.Error("Error db operation: ", err)
		os.Exit(1)
	}
}

func doNew(c *cli.Context) {
	dbName := c.Args().First()
	if dbName == "" {
		log.Error("DB is not provied", dbName)
		os.Exit(1)
	}

	opts := levigo.NewOptions()
	opts.SetCreateIfMissing(true)

	err := operateDb(dbName, opts, func(db *levigo.DB) error {
		return nil
	})
	if err != nil {
		log.Error("Error db operation: ", err)
		os.Exit(1)
	}
}
