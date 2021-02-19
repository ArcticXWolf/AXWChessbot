import sys
import chess
import search
from timeit import default_timer as timer


class GoCommandArgs:
    has_args = False
    has_timing_info = False
    time = {chess.WHITE: None, chess.BLACK: None}
    inc = {chess.WHITE: None, chess.BLACK: None}

    def __init__(self, line: str):
        if len(line) <= 2:
            return

        self.has_args = True
        args = iter(line[3:].split(" "))
        try:
            for arg in args:
                if arg == "wtime":
                    self.time[chess.WHITE] = int(next(args, None))
                if arg == "btime":
                    self.time[chess.BLACK] = int(next(args, None))
                if arg == "winc":
                    self.inc[chess.WHITE] = int(next(args, None))
                if arg == "binc":
                    self.inc[chess.BLACK] = int(next(args, None))
        except:
            pass

        self.has_timing_info = (
            self.time[chess.WHITE] is not None
            and self.time[chess.BLACK] is not None
            and self.inc[chess.WHITE] is not None
            and self.inc[chess.BLACK] is not None
        )


class Uci:
    def __init__(self, abdepth, qdepth=10):
        self.abdepth = abdepth
        self.qdepth = qdepth
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
            go_args = GoCommandArgs(msg)
            if go_args.has_args and go_args.has_timing_info:
                self.set_depth_by_timing(go_args)

            start_search = timer()
            move, num_positions = search.Search(
                self.board, self.abdepth, self.qdepth
            ).next_move()
            end_search = timer()
            self.debug(
                f"Analyzed {num_positions} positions in {end_search - start_search :.2f} sec at depths ({self.abdepth}, {self.qdepth})"
            )
            self.output(f"bestmove {move.uci()}")
            return

    def output(self, msg):
        print(msg)
        print(f"OUTPUT: {msg}", file=sys.stderr)

    def debug(self, msg):
        print(f"DEBUG: {msg}", file=sys.stderr)

    def set_depth_by_timing(self, go_args: GoCommandArgs):
        # if go_args.time[self.board.turn] > 300000:
        #    self.abdepth = 4
        #    self.qdepth = 10
        # if go_args.time[self.board.turn] > 120000:
        #    self.abdepth = 3
        #    self.qdepth = 10
        if go_args.time[self.board.turn] > 10000:
            self.abdepth = 2
            self.qdepth = 5
        else:
            self.abdepth = 1
            self.qdepth = 2
