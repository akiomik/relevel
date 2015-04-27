package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/jmhodges/levigo"
	"github.com/peterh/liner"
)

var (
	historyFile = "/tmp/." + appName + "_history"
	queries     = []string{"GET", "PUT", "DELETE", "KEYS", "HELP"}
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
	app.Commands = []cli.Command{
		{
			Name:      "console",
			ShortName: "c",
			Usage:     "Starts leveldb console",
			Action:    consoleCommand,
		},
		{
			Name:      "exec",
			ShortName: "e",
			Usage:     "Executes leveldb queries",
			Action:    execCommand,
		},
		{
			Name:      "new",
			ShortName: "n",
			Usage:     "Creates new db",
			Action:    newCommand,
		},
		{
			Name:      "version",
			ShortName: "v",
			Usage:     "Shows version",
			Action:    cli.ShowVersion,
		},
	}
	app.Run(os.Args)
}

func consoleCommand(c *cli.Context) {
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

func queryParser(line string) (string, []string) {
	input := strings.Split(line, " ")
	query := strings.ToLower(input[0])
	args := input[1:]
	return query, args
}

func queryHandler(query string, args []string, db *levigo.DB) {
	switch {
	case query == "help":
		showQueryHelp()
	case query == "keys":
		executeKeys(db)
	case query == "get":
		if len(args) < 1 {
			println("Key must be provided.")
			return
		}

		executeGet(args[0], db)
	case query == "put":
		if len(args) < 2 {
			println("Key and value must be provided.")
			return
		}

		executePut(args[0], args[1], db)
	case query == "delete":
		if len(args) < 1 {
			println("Key must be provided.")
			return
		}

		executeDelete(args[0], db)
	case query == "":
		// nothing to do
	default:
		println("Unknown query: ", query)
	}
}

func showQueryHelp() {
	println("keys")
	println("get <key>")
	println("put <key> <value>")
	println("delete <key>")
	println("help")
}

func execCommand(c *cli.Context) {
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

func newCommand(c *cli.Context) {
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

func startConsole(exec func(string)) error {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetCompleter(consoleCompleter)

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

func consoleCompleter(line string) (c []string) {
	for _, q := range queries {
		if strings.HasPrefix(strings.ToLower(q), line) {
			c = append(c, strings.ToLower(q)+" ")
		} else if strings.HasPrefix(q, strings.ToUpper(line)) {
			c = append(c, q+" ")
		}
	}
	return
}

func executeGet(key string, db *levigo.DB) {
	ro := levigo.NewReadOptions()
	defer ro.Close()

	value, err := db.Get(ro, []byte(key))
	if err != nil {
		log.Print("Error db operation:", err)
		return
	}

	if value == nil {
		println("Not found:", key)
		return
	}

	println(string(value[:]))
	println("")
	println("Got:", key)
}

func executePut(key string, value string, db *levigo.DB) {
	wo := levigo.NewWriteOptions()
	defer wo.Close()

	err := db.Put(wo, []byte(key), []byte(value))
	if err != nil {
		log.Print("Error db operation:", err)
		return
	}

	println("Created:", key)
}

func executeKeys(db *levigo.DB) {
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)

	it := db.NewIterator(ro)
	defer it.Close()

	count := 0
	for it.SeekToFirst(); it.Valid(); it.Next() {
		key := it.Key()
		println(string(key[:]))
		count++
	}

	println("")
	fmt.Printf("Found: %d keys\n", count)
}

func executeDelete(key string, db *levigo.DB) {
	wo := levigo.NewWriteOptions()
	defer wo.Close()

	err := db.Delete(wo, []byte(key))
	if err != nil {
		log.Print("Error db operation:", err)
		return
	}

	println("Deleted:", key)
}
