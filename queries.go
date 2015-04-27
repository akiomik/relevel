package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jmhodges/levigo"
)

var (
	queries = []string{"GET", "PUT", "DELETE", "KEYS", "HELP"}
)

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

func queryCompleter(line string) (c []string) {
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
