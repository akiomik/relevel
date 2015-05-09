package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jmhodges/levigo"
)

type Query struct {
	Name   string
	Help   string
	Action func([]string, *levigo.DB)
}

var queries = []Query{
	getQuery,
	putQuery,
	deleteQuery,
	keysQuery,
	// avoid initialization loop
	// helpQuery,
}

var allQuery = append(queries, helpQuery)

var getQuery = Query{
	Name:   "get",
	Help:   "get <key>",
	Action: executeGet,
}

var putQuery = Query{
	Name:   "put",
	Help:   "put <key> <value>",
	Action: executePut,
}

var deleteQuery = Query{
	Name:   "delete",
	Help:   "delete <key>",
	Action: executeDelete,
}

var keysQuery = Query{
	Name:   "keys",
	Help:   "keys",
	Action: executeKeys,
}

var helpQuery = Query{
	Name:   "help",
	Help:   "help",
	Action: showQueryHelp,
}

func queryParser(line string) (string, []string) {
	input := strings.Split(line, " ")
	query := strings.ToLower(input[0])
	args := input[1:]
	return query, args
}

func queryHandler(query string, args []string, db *levigo.DB) {
	if query == "" {
		return
	}

	for _, q := range allQuery {
		if query != q.Name {
			continue
		}

		q.Action(args, db)
		return
	}

	println("Unknown query: ", query)
}

func showQueryHelp(_ []string, _ *levigo.DB) {
	println("help")

	for _, q := range queries {
		println(q.Help)
	}
}

func queryCompleter(line string) (c []string) {
	for _, q := range allQuery {
		if strings.HasPrefix(strings.ToUpper(q.Name), line) {
			c = append(c, strings.ToLower(q.Name)+" ")
		} else if strings.HasPrefix(q.Name, strings.ToLower(line)) {
			c = append(c, q.Name+" ")
		}
	}
	return
}

func executeGet(args []string, db *levigo.DB) {
	if len(args) < 1 {
		println("Key must be provided.")
		return
	}
	key := args[0]

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

func executePut(args []string, db *levigo.DB) {
	if len(args) < 2 {
		println("Key and value must be provided.")
		return
	}
	key := args[0]
	value := args[1]

	wo := levigo.NewWriteOptions()
	defer wo.Close()

	err := db.Put(wo, []byte(key), []byte(value))
	if err != nil {
		log.Print("Error db operation:", err)
		return
	}

	println("Created:", key)
}

func executeKeys(_ []string, db *levigo.DB) {
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

func executeDelete(args []string, db *levigo.DB) {
	if len(args) < 1 {
		println("Key must be provided.")
		return
	}
	key := args[0]

	wo := levigo.NewWriteOptions()
	defer wo.Close()

	err := db.Delete(wo, []byte(key))
	if err != nil {
		log.Print("Error db operation:", err)
		return
	}

	println("Deleted:", key)
}
