from typing import Tuple
import chess
import score_tables


class Evaluation:
    board = chess.Board()

    def __init__(self, board: chess.Board):
        self.board = board

    def evaluate(self) -> int:
        finished, score = self.evaluate_if_endstate_is_reached()

        if finished:
            return score

        score = self.evaluate_pst_and_material()

        if self.board.turn:
            return score
        else:
            return -score

    def evaluate_if_endstate_is_reached(self) -> Tuple[bool, int]:
        if self.board.is_checkmate():
            if self.board.turn:
                return (True, -9999)
            else:
                return (True, 9999)

        if self.board.is_stalemate():
            return (True, 0)

        if self.board.is_insufficient_material():
            return (True, 0)

        return (False, None)

    def evaluate_pst_and_material(self) -> int:
        m_score = self.evaluate_material_score()
        pst_score = {
            i: self.evaluate_piece_square_table(chess.WHITE, i)
            + self.evaluate_piece_square_table(chess.BLACK, i)
            for i in chess.PIECE_TYPES
        }

        return m_score + sum(pst_score.values())

    def evaluate_piece_square_table(
        self, color: chess.Color, piece_type: chess.PieceType
    ) -> int:
        return sum(
            [
                self.evaluate_piece_square_table_single_piece(color, piece_type, i)
                for i in self.board.pieces(piece_type, color)
            ]
        )

    def evaluate_piece_square_table_single_piece(
        self, color: chess.Color, piece_type: chess.PieceType, position: chess.Square
    ):
        if color == chess.BLACK:
            return -score_tables.piece_square_tables[piece_type][
                chess.square_mirror(position)
            ]
        return score_tables.piece_square_tables[piece_type][position]

    def evaluate_material_score(self) -> int:
        white_piece_counts = self.get_pieces_counts(chess.WHITE)
        black_piece_counts = self.get_pieces_counts(chess.BLACK)

        piece_scores = {
            i: score_tables.piece_values[i]
            * (white_piece_counts[i] - black_piece_counts[i])
            for i in chess.PIECE_TYPES
            if i != chess.KING
        }

        return sum(piece_scores.values())

    def get_pieces_counts(self, color: chess.Color) -> dict:
        return {
            i: len(self.board.pieces(i, color))
            for i in chess.PIECE_TYPES
            if i != chess.KING
        }

    def move_value(self, move: chess.Move):
        if move.promotion is not None:
            return -float("inf")

        from_piece = self.board.piece_at(move.from_square)
        to_piece = self.board.piece_at(move.to_square)
        from_value = self.evaluate_piece_square_table_single_piece(
            from_piece.color, from_piece.piece_type, move.from_square
        )
        to_value = self.evaluate_piece_square_table_single_piece(
            from_piece.color, from_piece.piece_type, move.to_square
        )
        position_value = to_value - from_value

        capture_value = 0
        if self.board.is_capture(move):
            if self.board.is_en_passant(move):
                capture_value = score_tables.piece_values[chess.PAWN]
            else:
                capture_value = (
                    score_tables.piece_values[to_piece.piece_type]
                    - score_tables.piece_values[from_piece.piece_type]
                )

        return -(capture_value + position_value)
