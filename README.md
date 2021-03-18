# AXWChessbot

This is a simple chess engine written in python. Its main code can be found
under [engines/axwchessbot/](engines/axwchessbot), while the surrounding code
is a fork of [lichess-bot](https://github.com/ShailChoksi/lichess-bot) for
easy integration with [lichess](https://lichess.org).

You can play against the engine [here](https://lichess.org/@/AXWChessBot).

## Strength

I have not been able to conduct a good comparison against other engines yet
and furthermore I change the engine quite a bit, so the strength may vary
from time to time. However the engine is available on lichess and there it
beats stockfish on difficulty levels 1-4 (lichess rating levels of 800-1700)
confidently on rapid 10+10 rules. Level 5 (rating 2000) is dominated by
stockfish, but AXWChessBot manages to secure some wins as well.

So its strength should be between 1700 and 2000. Remember that this is a
lichess rating which does not correspond to the FIDE elo ratings and that its
strength also depends on the machine it is running on. More computing
resources mean more depth in the search tree. The server running the bot on
lichess provides enough resources for an average search depth of 3-4 plys.

## Features

* Written in python
* Basic UCI interface
* Simple timemanagement
* Simple evaluation function using
    * Gamephase detection of mid- and endgame
    * Piece value
    * Piece square tables
    * Tempo evaluation
    * Pair bonus
    * Rook on (half-)open files bonus
    * Pawnshield bonus
    * Blocked piece penalty
    * Passed pawn bonus
* Simple search algorithm using
    * Negamax with alpha-beta-pruning
    * Quiescence search
    * Iterative deepening
    * Move ordering
    * Transposition tables as cache
    * Opening books
    * Endgame tables
* Integrated into lichess-bot
    * Bot reacts to game events via chat

## Features ideas

* Killer move heuristics
* Multiprocessing (Lazy SMP?)

## Used materials / Research

* [Chessprogramming-Wiki](https://www.chessprogramming.org)
* Stockfish source
* [python-chess](https://github.com/niklasf/python-chess)
* [cpw-engine](https://github.com/nescitus/cpw-engine)
* [chess-ai](https://github.com/xtreemtg/Chess_AI)
* [sunfish](https://github.com/thomasahle/sunfish)
* [Adam Berent Blog](https://adamberent.com/2019/03/02/chess-board-evaluation/)
* Several blogposts/stackoverflow questions found via Google
