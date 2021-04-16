package main

import (
	"log"
	"os"
	"runtime/pprof"

	"go.janniklasrichter.de/axwchessbot/uci"
)

const (
	engineName   = "AXWChessBot"
	engineAuthor = "Jan Niklas Richter"
)

var (
	version = "undefined"
	date    = "undefined"
	commit  = "undefined"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	f, perr := os.Create("cpu.pprof")
	if perr != nil {
		logger.Fatal(perr)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	logger.Println(engineName, "Version", version, "BuildDate", date, "GitCommitHash", commit)

	uci.StartProtocol(logger, uci.New(engineName, engineAuthor, version, []uci.UciOption{}, logger))
}
