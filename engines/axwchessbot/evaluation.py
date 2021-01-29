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
        is_endgame = self.check_if_endgame()
        m_score = self.evaluate_material_score()
        pst_score = {
            i: self.evaluate_piece_square_table(chess.WHITE, i, is_endgame)
            + self.evaluate_piece_square_table(chess.BLACK, i, is_endgame)
            for i in chess.PIECE_TYPES
        }

        return m_score + sum(pst_score.values())

    def evaluate_piece_square_table(
        self, color: chess.Color, piece_type: chess.PieceType, is_endgame: bool
    ) -> int:
        return sum(
            [
                self.evaluate_piece_square_table_single_piece(
                    color, piece_type, i, is_endgame
                )
                for i in self.board.pieces(piece_type, color)
            ]
        )

    def evaluate_piece_square_table_single_piece(
        self,
        color: chess.Color,
        piece_type: chess.PieceType,
        position: chess.Square,
        is_endgame: bool,
    ):
        pst = score_tables.piece_square_tables
        if is_endgame:
            pst = score_tables.piece_square_tables_endgame

        if color == chess.BLACK:
            return -pst[piece_type][chess.square_mirror(position)]
        return pst[piece_type][position]

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

    def check_if_endgame(self) -> bool:
        queens_white = len(self.board.pieces(chess.QUEEN, chess.WHITE))
        queens_black = len(self.board.pieces(chess.QUEEN, chess.BLACK))
        minors_white = len(self.board.pieces(chess.BISHOP, chess.WHITE)) + len(
            self.board.pieces(chess.KNIGHT, chess.WHITE)
        )
        minors_black = len(self.board.pieces(chess.BISHOP, chess.BLACK)) + len(
            self.board.pieces(chess.KNIGHT, chess.BLACK)
        )

        return (queens_white + queens_black == 0) or (
            (queens_white == 0 or minors_white <= 1)
            and (queens_black == 0 or minors_black <= 1)
        )

    def move_value(self, move: chess.Move, is_endgame: bool):
        if move.promotion is not None:
            return -float("inf")

        from_piece = self.board.piece_at(move.from_square)
        to_piece = self.board.piece_at(move.to_square)
        from_value = self.evaluate_piece_square_table_single_piece(
            from_piece.color, from_piece.piece_type, move.from_square, is_endgame
        )
        to_value = self.evaluate_piece_square_table_single_piece(
            from_piece.color, from_piece.piece_type, move.to_square, is_endgame
        )
        position_value = to_value - from_value

        capture_value = self.capture_value(move)

        return -(capture_value + position_value)

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
