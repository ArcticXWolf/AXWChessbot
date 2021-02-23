from typing import Tuple
import chess
import score_tables
import itertools

NORMALIZATION_VALUE = 3900


class Evaluation:
    board = chess.Board()

    def __init__(self, board: chess.Board):
        self.board = board

    def evaluate(self) -> int:
        finished, score = self.evaluate_if_endstate_is_reached()

        if finished:
            # evals are from whites perspective, convert to current perspective
            if self.board.turn:
                return score
            else:
                return -score

        score = self.evaluate_pst_and_material()

        if self.board.turn:
            # evals are from whites perspective, convert to current perspective
            return score
        else:
            return -score

    def evaluate_if_endstate_is_reached(self) -> Tuple[bool, int]:
        if not self.board.is_game_over():
            return (False, None)

        result = self.board.result()
        if result == "0-1":
            return (True, -1.0)
        elif result == "1-0":
            return (True, 1.0)
        return (True, 0)

    def evaluate_pst_and_material(self) -> int:
        pst = score_tables.piece_square_tables
        if self.check_if_endgame():
            pst = score_tables.piece_square_tables_endgame

        score = 0
        for color in chess.COLORS:
            color_score = 0
            for piece in chess.PIECE_TYPES:
                pieces = self.board.pieces(piece, color)
                pieces_pst_scores = (
                    itertools.compress(pst[piece], list(pieces))
                    if len(pieces) > 0
                    else [0]
                )
                color_score += (
                    sum(pieces_pst_scores)
                    + len(pieces) * score_tables.piece_values[piece]
                )
            if color == chess.WHITE:
                score += color_score
            else:
                score -= color_score
        return max(min(score / NORMALIZATION_VALUE, 1.0), -1.0)

    def check_if_endgame(self) -> bool:
        queens_white = len(self.board.pieces(chess.QUEEN, chess.WHITE))
        queens_black = len(self.board.pieces(chess.QUEEN, chess.BLACK))
        rooks_white = len(self.board.pieces(chess.ROOK, chess.WHITE))
        rooks_black = len(self.board.pieces(chess.ROOK, chess.BLACK))
        minors_white = len(self.board.pieces(chess.BISHOP, chess.WHITE)) + len(
            self.board.pieces(chess.KNIGHT, chess.WHITE)
        )
        minors_black = len(self.board.pieces(chess.BISHOP, chess.BLACK)) + len(
            self.board.pieces(chess.KNIGHT, chess.BLACK)
        )

        return (
            queens_white
            + queens_black
            + rooks_white
            + rooks_black
            + minors_white
            + minors_black
            <= 4
        )

    def attacked_by_inferior_piece(
        self, move: chess.Move, evaluate_to_square: bool
    ) -> bool:
        checked_square = move.to_square if evaluate_to_square else move.from_square
        for square in self.board.attackers(not self.board.turn, checked_square):
            our_piece_value = int(
                score_tables.piece_values[self.board.piece_type_at(move.from_square)]
                / 100
            )
            attacker_value = int(
                score_tables.piece_values[self.board.piece_type_at(square)] / 100
            )
            if our_piece_value > attacker_value:
                return True
        return False

    def defenders_of_square(self, square: chess.Square):
        return self.board.attackers(self.board.turn, square)

    def attackers_of_square(self, square: chess.Square):
        return self.board.attackers(not self.board.turn, square)

    def move_order(self):
        good_moves = []
        medium_moves = []
        bad_moves = []

        move: chess.Move
        for move in self.board.legal_moves:
            move_text = self.board.san(move)

            if "#" in move_text:
                return [move]

            if self.board.is_capture(move):
                if self.board.piece_at(move.from_square) == chess.PAWN:
                    good_moves.insert(0, move)
                    continue
                elif not self.board.is_attacked_by(not self.board.turn, move.to_square):
                    good_moves.insert(0, move)
                    continue
                else:
                    medium_moves.insert(0, move)
                    continue

            if self.board.piece_at(move.from_square) == chess.QUEEN:
                if self.board.is_attacked_by(not self.board.turn, move.to_square):
                    bad_moves.insert(0, move)
                    continue

            if self.attacked_by_inferior_piece(move, False):
                if self.attacked_by_inferior_piece(move, True):
                    bad_moves.insert(0, move)
                    continue
                else:
                    if len(self.defenders_of_square(move.to_square)) >= len(
                        self.attackers_of_square(move.to_square)
                    ):
                        good_moves.insert(0, move)
                        continue
                    else:
                        bad_moves.insert(0, move)
                        continue
            elif self.attacked_by_inferior_piece(move, True):
                bad_moves.insert(0, move)
                continue

            medium_moves.insert(0, move)
        return good_moves + medium_moves + bad_moves

    def capture_value(self, move: chess.Move):
        if not self.board.is_capture(move):
            return 0

        from_piece = self.board.piece_at(move.from_square)
        to_piece = self.board.piece_at(move.to_square)

        if self.board.is_en_passant(move):
            return score_tables.piece_values[chess.PAWN]
        else:
            return (
                score_tables.piece_values[to_piece.piece_type]
                - score_tables.piece_values[from_piece.piece_type]
            )
