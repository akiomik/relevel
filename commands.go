package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
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
	db_name := c.Args().First()
	if db_name == "" {
		log.Error("DB name is not provided", db_name)
		os.Exit(1)
	}

	opts := levigo.NewOptions()
	err := operateDb(db_name, opts, func(db *levigo.DB) error {
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
	db_name := c.Args().Get(0)
	if db_name == "" {
		log.Error("DB is not provied", db_name)
		os.Exit(1)
	}

	line := c.Args().Get(1)
	if line == "" {
		log.Error("Query is not provied", line)
		os.Exit(1)
	}

	opts := levigo.NewOptions()
	err := operateDb(db_name, opts, func(db *levigo.DB) error {
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
	db_name := c.Args().First()
	if db_name == "" {
		log.Error("DB is not provied", db_name)
		os.Exit(1)
	}

	opts := levigo.NewOptions()
	opts.SetCreateIfMissing(true)

	err := operateDb(db_name, opts, func(db *levigo.DB) error {
		return nil
	})
	if err != nil {
		log.Error("Error db operation: ", err)
		os.Exit(1)
	}
}
