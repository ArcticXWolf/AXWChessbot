package uci

import (
	"bufio"
	"log"
	"os"
)

func StartProtocol(logger *log.Logger, protocol *UciProtocol) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "quit" {
			logger.Println("exiting")
			break
		}
		err := protocol.HandleInput(message)

		if err != nil {
			logger.Println(err)
		}
	}
}
