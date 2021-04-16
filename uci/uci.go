package uci

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"go.janniklasrichter.de/axwchessbot/evaluation"
	"go.janniklasrichter.de/axwchessbot/game"
	"go.janniklasrichter.de/axwchessbot/search"
)

type UciOption struct {
	name  string
	value string
}

type UciProtocol struct {
	name        string
	author      string
	version     string
	logger      *log.Logger
	options     []UciOption
	currentGame *game.Game
}

func New(name, author, version string, options []UciOption, logger *log.Logger) *UciProtocol {
	return &UciProtocol{
		name:    name,
		author:  author,
		version: version,
		options: options,
		logger:  logger,
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
	fmt.Printf("readyok\n")
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
	p.logger.Printf("Position: %v", messageParts)
	messageParts = messageParts[1:]

	p.currentGame = game.New()
	if command == "fen" {
		fen := ""
		extracted_keys := 0
		for _, value := range messageParts {
			if value == "moves" {
				break
			}
			fen = fmt.Sprintf("%s %s", fen, value)
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
			err := p.currentGame.PushMoveStr(moveStr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *UciProtocol) goCmd(messageParts []string) error {
	timingInfo := NewTimingInfo(messageParts)
	context, cancel := timingInfo.calculateTimeoutContext(context.Background(), p.currentGame)
	defer cancel()

	tTable := search.NewTranspositionTable(1000000)
	evaluator := evaluation.Evaluation{}
	searchObj := search.New(p.currentGame, p.logger, tTable, &evaluator, 40, 10)
	bestMove, score := searchObj.SearchBestMove(context)

	fmt.Printf("bestmove %v\n", bestMove.String())
	infoStr := fmt.Sprintf("info depth %d", searchObj.SearchInfo.MaxDepthCompleted)
	infoStr += fmt.Sprintf(" score cp %d", int(score))
	infoStr += fmt.Sprintf(" nodes %d", searchObj.SearchInfo.NodesTraversed)
	infoStr += fmt.Sprintf(" nps %d", int(float64(searchObj.SearchInfo.NodesTraversed)/searchObj.SearchInfo.TotalSearchTime.Seconds()))
	infoStr += fmt.Sprintf(" time %d\n", searchObj.SearchInfo.TotalSearchTime.Milliseconds())
	fmt.Print(infoStr)
	p.logger.Printf("Move: %v, Info: %v", bestMove.String(), infoStr)

	return nil
}

func (p *UciProtocol) uciNewGameCmd(messageParts []string) error {
	return nil
}
