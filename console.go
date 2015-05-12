package main

import (
	"io"
	"os"

	log "github.com/golang/glog"
	"github.com/peterh/liner"
)

var (
	historyFile = "/tmp/." + appName + "_history"
	prompt      = appName + "> "
)

func showConsoleHeader() {
	println(appName + " - " + appVersion)
	println("")
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
