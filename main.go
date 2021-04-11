package main

import (
	"log"
	"os"

	"go.janniklasrichter.de/axwchessbot/uci"
)

const (
	engineName   = "AXWChessBot"
	engineAuthor = "Jan Niklas Richter"
)

var (
	engineVersion = "undefined"
	buildDate     = "undefined"
	gitCommit     = "undefined"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	logger.Println(engineName, "Version", engineVersion, "BuildDate", buildDate, "GitCommitHash", gitCommit)

	uci.StartProtocol(logger, uci.New(engineName, engineAuthor, engineVersion, []uci.UciOption{}))
}
