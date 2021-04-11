package uci

import (
	"errors"
	"fmt"
	"strings"

	"go.janniklasrichter.de/axwchessbot/game"
)

type UciOption struct {
	name  string
	value string
}

type UciProtocol struct {
	name        string
	author      string
	version     string
	options     []UciOption
	currentGame *game.Game
}

func New(name, author, version string, options []UciOption) *UciProtocol {
	return &UciProtocol{
		name:    name,
		author:  author,
		version: version,
		options: options,
	}
}

func (p *UciProtocol) HandleInput(message string) error {
	messageParts := strings.Fields(message)
	if len(messageParts) <= 0 {
		return nil
	}

	command := messageParts[0]
	messageParts = messageParts[1:]

	switch command {
	case "uci":
		return p.uciCmd(messageParts)
	case "isready":
		return p.isReadyCmd(messageParts)
	case "setoption":
		return p.setOptionCmd(messageParts)
	case "position":
		return p.positionCmd(messageParts)
	case "go":
		return p.goCmd(messageParts)
	case "ucinewgame":
		return p.uciNewGameCmd(messageParts)
	default:
		return errors.New("unknown command")
	}
}

func (p *UciProtocol) uciCmd(messageParts []string) error {
	fmt.Printf("id name %s %s\n", p.name, p.version)
	fmt.Printf("id author %s\n", p.author)
	for _, option := range p.options {
		fmt.Printf("option name %v type string default", option.name)
	}
	fmt.Println("uciok")
	return nil
}

func (p *UciProtocol) isReadyCmd(messageParts []string) error {
	fmt.Printf("readyok")
	return nil
}

func (p *UciProtocol) setOptionCmd(messageParts []string) error {
	if len(messageParts) < 4 {
		return errors.New("wrong arguments for setoption command")
	}

	option := UciOption{messageParts[1], messageParts[3]}
	p.options = append(p.options, option)

	return nil
}

func (p *UciProtocol) positionCmd(messageParts []string) error {
	command := messageParts[0]
	messageParts = messageParts[1:]

	p.currentGame = game.New()
	if command == "fen" {
		fen := ""
		extracted_keys := 0
		for _, value := range messageParts {
			if value == "moves" {
				break
			}
			fen = fen + value
			extracted_keys++
		}
		messageParts = messageParts[extracted_keys:]
		p.currentGame = game.NewFromFen(fen)
	} else if command != "startpos" {
		return errors.New("unknown arguments for position command")
	}

	if len(messageParts) > 0 {
		if messageParts[0] != "moves" {
			return fmt.Errorf("unknown argument %v for position command", messageParts[0])
		}
		messageParts = messageParts[1:]

		for _, moveStr := range messageParts {
			err := p.currentGame.PushMove(moveStr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *UciProtocol) goCmd(messageParts []string) error {
	return nil
}

func (p *UciProtocol) uciNewGameCmd(messageParts []string) error {
	return nil
}
