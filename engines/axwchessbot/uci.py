import sys
import chess
import search


class Uci:
    def __init__(self, depth):
        self.depth = depth
        self.board = chess.Board()

    def communicate(self):
        while True:
            msg = input()
            print(f"INPUT: {msg}", file=sys.stderr)
            self.command(msg)

    def command(self, msg):
        if msg == "quit":
            quit()

        if msg == "uci":
            self.output("id name AXWChess")
            self.output("id author Jan Niklas Richter")
            self.output("uciok")
            return

        if msg == "isready":
            self.output("readyok")
            return

        if msg == "ucinewgame":
            return

        if "position startpos moves" in msg:
            moves = msg.split(" ")[3:]
            self.board.clear()
            self.board.set_fen(chess.STARTING_FEN)
            for move in moves:
                self.board.push(chess.Move.from_uci(move))
            return

        if "position fen" in msg:
            fen = " ".join(msg.split(" ")[2:])
            self.board.set_fen(fen)
            return

        if msg[0:2] == "go":
            move = search.Search(self.board, self.depth).next_move()
            self.output(f"bestmove {move.uci()}")
            return

    def output(self, msg):
        print(msg)
        print(f"OUTPUT: {msg}", file=sys.stderr)