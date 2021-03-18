from typing import Tuple
import chess
from . import score_tables
import functools
import operator


class EvaluationResult:
    def __init__(self):
        self.gamephase = 0
        self.piece_score = {piece_type: 0 for piece_type in chess.PIECE_TYPES}
        self.piece_square_score_midgame = {
            piece_type: 0 for piece_type in chess.PIECE_TYPES
        }
        self.piece_square_score_endgame = {
            piece_type: 0 for piece_type in chess.PIECE_TYPES
        }
        self.material_score_midgame = 0
        self.material_score_endgame = 0
        self.pair_bonus = 0
        self.rook_bonus = 0
        self.tempo_bonus = 0
        self.blocked_pieces_score = 0
        self.king_shield_bonus = 0
        self.passed_pawn_bonus = 0


class Evaluation:
    def __init__(self, board: chess.Board):
        self.board = board
        self.eval_result = {
            chess.WHITE: EvaluationResult(),
            chess.BLACK: EvaluationResult(),
        }

    def evaluate(self):
        for color in chess.COLORS:
            self.evaluate_material_score(color)
            self.evaluate_gamephase(color)
            self.evaluate_tempo(color)
            self.evaluate_pair_bonus(color)
            self.evaluate_rook_bonus(color)
            self.evaluate_blocked_pieces(color)
            self.evaluate_king_shield(color)
            self.evaluate_passed_pawns(color)

    def evaluate_material_score(self, color: chess.Color) -> None:
        for piece_type in chess.PIECE_TYPES:
            pieces = list(self.board.pieces(piece_type, color))
            if color == chess.BLACK:
                pieces = [chess.square_mirror(i) for i in pieces]

            self.eval_result[color].piece_square_score_midgame[piece_type] = [
                score_tables.piece_square_tables[piece_type][i] for i in list(pieces)
            ]
            self.eval_result[color].piece_square_score_endgame[piece_type] = [
                score_tables.piece_square_tables_endgame[piece_type][i]
                for i in list(pieces)
            ]
            self.eval_result[color].piece_score[piece_type] = (
                len(pieces) * score_tables.piece_values[piece_type]
            )

            self.eval_result[color].material_score_midgame += (
                sum(self.eval_result[color].piece_square_score_midgame[piece_type])
                + self.eval_result[color].piece_score[piece_type]
            )
            self.eval_result[color].material_score_endgame += (
                sum(self.eval_result[color].piece_square_score_endgame[piece_type])
                + self.eval_result[color].piece_score[piece_type]
            )

    def evaluate_gamephase(self, color: chess.Color) -> None:
        self.eval_result[color].gamephase = (
            len(self.board.pieces(chess.KNIGHT, color))
            + len(self.board.pieces(chess.BISHOP, color))
            + 2 * len(self.board.pieces(chess.ROOK, color))
            + 4 * len(self.board.pieces(chess.QUEEN, color))
        )

    def evaluate_tempo(self, color: chess.Color) -> None:
        if self.board.turn == color:
            self.eval_result[color].tempo_bonus += score_tables.additional_modifiers[
                "tempo"
            ]

    def evaluate_pair_bonus(self, color: chess.Color) -> None:
        if len(self.board.pieces(chess.BISHOP, color)) > 1:
            self.eval_result[color].pair_bonus += score_tables.additional_modifiers[
                "bishop_pair"
            ]
        if len(self.board.pieces(chess.KNIGHT, color)) > 1:
            self.eval_result[color].pair_bonus -= score_tables.additional_modifiers[
                "knight_pair"
            ]
        if len(self.board.pieces(chess.ROOK, color)) > 1:
            self.eval_result[color].pair_bonus -= score_tables.additional_modifiers[
                "rook_pair"
            ]

    def evaluate_rook_bonus(self, color: chess.Color) -> None:
        for piece in self.board.pieces(chess.ROOK, color):
            own_pawns_in_same_file = (
                self.board.pawns
                & self.board.occupied_co[color]
                & chess.BB_FILES[chess.square_file(piece)]
            )

            if own_pawns_in_same_file > 0:
                continue

            enemy_pawns_in_same_file = (
                self.board.pawns
                & self.board.occupied_co[color]
                & chess.BB_FILES[chess.square_file(piece)]
            )
            if enemy_pawns_in_same_file > 0:  # semi open file
                self.eval_result[color].rook_bonus += score_tables.additional_modifiers[
                    "half_rook"
                ]
            else:  # open file
                self.eval_result[color].rook_bonus += score_tables.additional_modifiers[
                    "open_rook"
                ]

    def evaluate_blocked_pieces(self, color: chess.Color) -> None:
        # king blocks rook
        side_rank = 0 if color == chess.WHITE else 7
        if (
            self.board.piece_type_at(chess.square(5, side_rank)) == chess.KING
            or self.board.piece_type_at(chess.square(6, side_rank)) == chess.KING
        ) and (
            self.board.piece_type_at(chess.square(6, side_rank)) == chess.ROOK
            or self.board.piece_type_at(chess.square(7, side_rank)) == chess.ROOK
        ):
            self.eval_result[
                color
            ].blocked_pieces_score += score_tables.additional_modifiers[
                "king_blocks_rook_penalty"
            ]

        if (
            self.board.piece_type_at(chess.square(1, side_rank)) == chess.KING
            or self.board.piece_type_at(chess.square(2, side_rank)) == chess.KING
        ) and (
            self.board.piece_type_at(chess.square(0, side_rank)) == chess.ROOK
            or self.board.piece_type_at(chess.square(1, side_rank)) == chess.ROOK
        ):
            self.eval_result[
                color
            ].blocked_pieces_score += score_tables.additional_modifiers[
                "king_blocks_rook_penalty"
            ]

    def evaluate_king_shield(self, color: chess.Color) -> None:
        rank_2_to_check = 2
        rank_3_to_check = 3
        king_position = self.board.pieces(chess.KING, color).pop()

        if color == chess.BLACK:
            rank_2_to_check = 7
            rank_3_to_check = 6

        if chess.square_file(king_position) > 4:
            pawn_count_2 = len(
                [
                    self.board.piece_at(chess.square(i, rank_2_to_check))
                    == chess.Piece(chess.PAWN, color)
                    for i in range(5, 8)
                ]
            )
            pawn_count_3 = len(
                [
                    self.board.piece_at(chess.square(i, rank_3_to_check))
                    == chess.Piece(chess.PAWN, color)
                    for i in range(5, 8)
                ]
            )
            self.eval_result[color].king_shield_bonus += (
                pawn_count_2 * score_tables.additional_modifiers["king_shield_rank_2"]
            )
            self.eval_result[color].king_shield_bonus += (
                pawn_count_3 * score_tables.additional_modifiers["king_shield_rank_3"]
            )
        elif chess.square_file(king_position) < 3:
            pawn_count_2 = len(
                [
                    self.board.piece_at(chess.square(i, rank_2_to_check))
                    == chess.Piece(chess.PAWN, color)
                    for i in range(0, 3)
                ]
            )
            pawn_count_3 = len(
                [
                    self.board.piece_at(chess.square(i, rank_3_to_check))
                    == chess.Piece(chess.PAWN, color)
                    for i in range(0, 3)
                ]
            )
            self.eval_result[color].king_shield_bonus += (
                pawn_count_2 * score_tables.additional_modifiers["king_shield_rank_2"]
            )
            self.eval_result[color].king_shield_bonus += (
                pawn_count_3 * score_tables.additional_modifiers["king_shield_rank_3"]
            )

    def evaluate_passed_pawns(self, color: chess.Color) -> None:
        bb_enemy_spans = chess.BB_EMPTY
        for pawn in self.board.pieces(chess.PAWN, not color):
            ranks_to_go = chess.BB_RANKS[chess.square_rank(pawn) + 1 : 8]
            if color != chess.BLACK:
                ranks_to_go = chess.BB_RANKS[0 : chess.square_rank(pawn)]

            files_to_go = chess.BB_FILES[
                max(0, chess.square_file(pawn) - 1) : min(
                    8, chess.square_file(pawn) + 2
                )
            ]

            bb_ranks = functools.reduce(operator.or_, ranks_to_go)
            bb_files = functools.reduce(operator.or_, files_to_go)
            bb_front_of_pawn = bb_ranks & bb_files
            print("------------------------")
            print(chess.SquareSet(bb_front_of_pawn))
            bb_enemy_spans |= bb_front_of_pawn

        passed_pawns = chess.SquareSet(
            self.board.pawns & self.board.occupied_co[color] & ~bb_enemy_spans
        )
        print("------------------------")
        print(chess.SquareSet(bb_enemy_spans))
        print("------------------------")
        print(passed_pawns)
        print("------------------------")
        print("------------------------")
        if color == chess.BLACK:
            passed_pawns = [chess.square_mirror(pawn) for pawn in list(passed_pawns)]
        self.eval_result[color].passed_pawn_bonus += sum(
            [
                score_tables.additional_piece_square_tables["passed_pawn"][pawn]
                for pawn in list(passed_pawns)
            ]
        )