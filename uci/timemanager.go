package uci

import (
	"context"
	"strconv"
	"time"

	"go.janniklasrichter.de/axwchessbot/game"
)

const (
	DefaultMovesToGo = 30
	MoveOverhead     = 200 * time.Millisecond
	MaxTime          = 30000 * time.Millisecond
	MinTimeLeft      = 3000 * time.Millisecond
)

type UciTimingInfo struct {
	StartTimestamp time.Time
	TimeWhite      int
	TimeBlack      int
	IncrementWhite int
	IncrementBlack int
	MovesToGo      int
	MoveTime       int
}

func NewTimingInfo(messageParts []string) (timingInfo *UciTimingInfo) {
	timingInfo = &UciTimingInfo{
		StartTimestamp: time.Now(),
	}
	for i, token := range messageParts {
		switch token {
		case "wtime":
			timingInfo.TimeWhite, _ = strconv.Atoi(messageParts[i+1])
		case "btime":
			timingInfo.TimeBlack, _ = strconv.Atoi(messageParts[i+1])
		case "winc":
			timingInfo.IncrementWhite, _ = strconv.Atoi(messageParts[i+1])
		case "binc":
			timingInfo.IncrementBlack, _ = strconv.Atoi(messageParts[i+1])
		case "movestogo":
			timingInfo.MovesToGo, _ = strconv.Atoi(messageParts[i+1])
		case "movetime":
			timingInfo.MoveTime, _ = strconv.Atoi(messageParts[i+1])
		}
	}
	return
}

func (timingInfo *UciTimingInfo) calculateTimeoutContext(ctx context.Context, g *game.Game, options []UciOption) (context.Context, func()) {
	if timingInfo.MovesToGo <= 0 && timingInfo.TimeWhite <= 0 && timingInfo.TimeBlack <= 0 {
		return context.WithCancel(ctx)
	}

	if timingInfo.MovesToGo <= 0 {
		timingInfo.MovesToGo = DefaultMovesToGo
	}

	timeLeft, increment := time.Duration(timingInfo.TimeWhite)*time.Millisecond, time.Duration(timingInfo.IncrementWhite)*time.Millisecond
	if !g.Position.Wtomove {
		timeLeft, increment = time.Duration(timingInfo.TimeBlack)*time.Millisecond, time.Duration(timingInfo.IncrementBlack)*time.Millisecond
	}

	moveOverhead := MoveOverhead
	maxTime := MaxTime
	for _, option := range options {
		if option.name == "Move Overhead" {
			optionsMO, err := strconv.Atoi(option.value)
			if err == nil {
				moveOverhead = time.Duration(optionsMO) * time.Millisecond
			}
		}
		if option.name == "Max Time" {
			optionsMT, err := strconv.Atoi(option.value)
			if err == nil {
				maxTime = time.Duration(optionsMT) * time.Second
			}
		}
	}

	timeLeft -= moveOverhead
	if timeLeft <= 0 {
		timeLeft = 0
	}

	total := float64(timeLeft) + float64(timingInfo.MovesToGo-1)*float64(increment)
	limit := time.Duration(total / float64(timingInfo.MovesToGo-1))

	if limit > timeLeft-MinTimeLeft {
		limit = timeLeft - MinTimeLeft
	}
	if limit > maxTime {
		limit = maxTime
	}

	return context.WithDeadline(ctx, timingInfo.StartTimestamp.Add(limit))
}
