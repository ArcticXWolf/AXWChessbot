import chess
import os
import glob
import csv


class Opening:
    def __init__(self, code: str, name: str, moves: str):
        self.code = code
        self.name = name
        board = chess.Board()
        for move in moves.split(" "):
            board.push(chess.Move.from_uci(move))
        self.fen = board.board_fen()
        self.moves = moves

    def __str__(self) -> str:
        return f"{self.code}: {self.name}"


class OpeningDB:
    def __init__(self, folderpath):
        self.db = {}

        self.load_db_from_folder(folderpath)

    def load_db_from_folder(self, path: str) -> None:
        for filename in glob.glob(os.path.join(path, "*.tsv")):
            with open(filename) as fd:
                rd = csv.DictReader(fd, delimiter="\t", quotechar='"')
                for row in rd:
                    opening = Opening(row["eco"], row["name"], row["moves"])
                    self.db[opening.fen] = opening

    def find_opening_by_board(self, board: chess.Board) -> Opening:
        move_string = " ".join(move.uci() for move in board.move_stack)

        if board.board_fen() in self.db:
            return self.db[board.board_fen()]

        for fen, opening in self.db.items():
            if move_string.startswith(opening.moves):
                return opening

        if len(board.move_stack) > 0:
            board_tmp = chess.Board()
            detected_opening = None

            for move in board.move_stack:
                board_tmp.push(move)
                if board_tmp.board_fen() in self.db:
                    detected_opening = self.db.get(board.board_fen(), detected_opening)

            return detected_opening

        return None