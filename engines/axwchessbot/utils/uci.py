import sys
import chess
from search import search
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
    def __init__(self, abdepth=40, qdepth=6, timeout=180):
        self.abdepth = abdepth
        self.qdepth = qdepth
        self.timeout = timeout
        self.board = chess.Board()
        self.cache = None

    def communicate(self):
        while True:
            msg = input()
            print(f"> {msg}", file=sys.stderr)
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
            fen = " ".join(msg.split(" ")[2:8])
            self.board.set_fen(fen)
            if len(msg.split(" ")) > 8:
                moves = msg.split(" ")[9:]
                for move in moves:
                    self.board.push(chess.Move.from_uci(move))
            return

        if msg[0:2] == "go":
            go_args = GoCommandArgs(msg)
            if go_args.has_args and go_args.has_timing_info:
                self.set_depth_by_timing(go_args)
            start_search = timer()
            search_obj = search.Search(
                self.board, self.abdepth, self.qdepth, self.timeout, self.cache
            )
            move, info = search_obj.next_move()
            self.cache = search_obj.cache
            end_search = timer()
            info.pop("moves_analysis")
            self.debug(f"[{end_search - start_search :.2f}] {str(info)}")
            self.output(f"bestmove {move.uci()}")
            return

    def output(self, msg):
        print(msg)
        print(f"< {msg}", file=sys.stderr)

    def debug(self, msg):
        print(f"# {msg}", file=sys.stderr)

    def set_depth_by_timing(self, go_args: GoCommandArgs):
        # if go_args.time[self.board.turn] > 300000:
        #    self.abdepth = 4
        #    self.qdepth = 10
        # if go_args.time[self.board.turn] > 120000:
        #    self.abdepth = 3
        #    self.qdepth = 10
        suggested_time = (go_args.time[self.board.turn] / 1000) / 30
        suggested_time = max(go_args.inc[self.board.turn] / 1000, suggested_time)
        suggested_time = min(go_args.time[self.board.turn] / 1000 / 2, suggested_time)
        suggested_time = min(suggested_time, 30)
        self.abdepth = 40
        self.qdepth = 6
        self.timeout = int(suggested_time)
