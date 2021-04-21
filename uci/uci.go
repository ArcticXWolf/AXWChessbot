package uci

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot/evaluation"
	"go.janniklasrichter.de/axwchessbot/game"
	"go.janniklasrichter.de/axwchessbot/search"
)

const (
	TranspositionTableSize = 268435456
)

type UciOption struct {
	name  string
	value string
}

type UciProtocol struct {
	name               string
	author             string
	version            string
	logger             *log.Logger
	options            []UciOption
	transpositionTable *search.TranspositionTable
	currentGame        *game.Game
}

func New(name, author, version string, options []UciOption, logger *log.Logger) *UciProtocol {
	return &UciProtocol{
		name:               name,
		author:             author,
		version:            version,
		options:            options,
		transpositionTable: search.NewTranspositionTable(TranspositionTableSize),
		logger:             logger,
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
	fmt.Printf("option name Hash type spin default 256 min 1 max 2048\n")
	fmt.Printf("option name Move Overhead type spin default 200 min 1 max 1000\n")
	fmt.Printf("option name Max Time type spin default 30 min 2 max 300\n")
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

	option := UciOption{}
	nameTokenFound := false
	valueTokenFound := false
	for _, content := range messageParts {
		if !nameTokenFound {
			if content == "name" {
				nameTokenFound = true
			}
			continue
		}

		if !valueTokenFound {
			if content == "value" {
				valueTokenFound = true
			} else {
				option.name = fmt.Sprintf("%v %v", option.name, strings.TrimSpace(content))
			}
			continue
		}

		option.value = fmt.Sprintf("%v %v", option.value, strings.TrimSpace(content))
	}

	option.name = strings.TrimSpace(option.name)
	option.value = strings.TrimSpace(option.value)
	p.options = append(p.options, option)

	p.recreateTranspositionTable()
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
	context, cancel := timingInfo.calculateTimeoutContext(context.Background(), p.currentGame, p.options)
	defer cancel()

	evaluator := evaluation.Evaluation{}
	searchObj := search.New(p.currentGame, p, p.logger, p.transpositionTable, &evaluator, 40, 10)
	bestMove, _ := searchObj.SearchBestMove(context)

	fmt.Printf("bestmove %v\n", bestMove.String())
	p.logger.Printf("TranspositionTable: %v / %v", len(p.transpositionTable.Entries), p.transpositionTable.MaxSizeInEntries)

	return nil
}

func (p *UciProtocol) uciNewGameCmd(messageParts []string) error {
	p.recreateTranspositionTable()
	return nil
}

func (p *UciProtocol) recreateTranspositionTable() {
	transpositionTableSize := TranspositionTableSize
	for _, option := range p.options {
		if option.name == "Hash" {
			optionsTTSize, err := strconv.Atoi(option.value)
			if err == nil {
				transpositionTableSize = optionsTTSize * 1048576
			}
		}
	}

	p.transpositionTable = search.NewTranspositionTable(transpositionTableSize)
}

func (p *UciProtocol) SendInfo(depth, score, nodes, nps int, time time.Duration, pv []dragontoothmg.Move) {
	infoStr := fmt.Sprintf("info depth %d", depth)
	infoStr += fmt.Sprintf(" score cp %d", score)
	infoStr += fmt.Sprintf(" nodes %d", nodes)
	infoStr += fmt.Sprintf(" nps %d", nps)
	infoStr += fmt.Sprintf(" time %d", time)
	if len(pv) > 0 {
		infoStr += " pv"
		for i := len(pv) - 1; i >= 0; i-- {
			infoStr += fmt.Sprintf(" %s", &(pv[i]))
		}
	}
	infoStr += "\n"
	fmt.Print(infoStr)
}
