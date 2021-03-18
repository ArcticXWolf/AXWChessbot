from __future__ import annotations
from typing import Tuple
import chess
from . import score_tables
import functools
import itertools
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
        self.total_score = 0
        self.total_score_perspective = 0
        self.evaluated = False

    def evaluate(self) -> Evaluation:
        if self.evaluated:
            return self

        self.evaluated = True
        if self.evaluate_gameover() is not None:
            self.total_score = self.evaluate_gameover()
            self.total_score_perspective = self.total_score
            if self.board.turn == chess.BLACK:
                self.total_score_perspective = -self.total_score
            return self

        for color in chess.COLORS:
            self.evaluate_material_score(color)
            self.evaluate_gamephase(color)
            self.evaluate_pair_bonus(color)
            self.evaluate_rook_bonus(color)
            self.evaluate_tempo(color)
            self.evaluate_blocked_pieces(color)
            self.evaluate_king_shield(color)
            self.evaluate_passed_pawns(color)

        self.combine_results()

        return self

    def combine_results(self) -> None:
        gamephase = (
            self.eval_result[chess.WHITE].gamephase
            + self.eval_result[chess.BLACK].gamephase
        )
        self.total_score += (
            float(
                (gamephase * self.eval_result[chess.WHITE].material_score_midgame)
                + (
                    (24 - gamephase)
                    * self.eval_result[chess.WHITE].material_score_endgame
                )
            )
            / 24.0
        )

        self.total_score -= (
            float(
                (gamephase * self.eval_result[chess.BLACK].material_score_midgame)
                + (
                    (24 - gamephase)
                    * self.eval_result[chess.BLACK].material_score_endgame
                )
            )
            / 24.0
        )

        self.total_score += self.eval_result[chess.WHITE].pair_bonus
        self.total_score -= self.eval_result[chess.BLACK].pair_bonus
        self.total_score += self.eval_result[chess.WHITE].rook_bonus
        self.total_score -= self.eval_result[chess.BLACK].rook_bonus
        self.total_score += self.eval_result[chess.WHITE].tempo_bonus
        self.total_score -= self.eval_result[chess.BLACK].tempo_bonus
        self.total_score += self.eval_result[chess.WHITE].blocked_pieces_score
        self.total_score -= self.eval_result[chess.BLACK].blocked_pieces_score
        # kingshield bonus is included in material_score_midgame
        self.total_score += self.eval_result[chess.WHITE].passed_pawn_bonus
        self.total_score -= self.eval_result[chess.BLACK].passed_pawn_bonus

        self.total_score_perspective = self.total_score
        if self.board.turn == chess.BLACK:
            self.total_score_perspective = -self.total_score

    def evaluate_gameover(self) -> float:
        if not self.board.is_game_over():
            return None

        result = self.board.result()
        if result == "0-1":
            return float("-inf")
        elif result == "1-0":
            return float("inf")
        return 0.0

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
            ].blocked_pieces_score -= score_tables.additional_modifiers[
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
            ].blocked_pieces_score -= score_tables.additional_modifiers[
                "king_blocks_rook_penalty"
            ]

    def evaluate_king_shield(self, color: chess.Color) -> None:
        rank_2_to_check = chess.BB_RANK_2
        rank_3_to_check = chess.BB_RANK_3
        king_position = self.board.king(color)

        if color == chess.BLACK:
            if chess.square_rank(king_position) != 7:
                return
            rank_2_to_check = chess.BB_RANK_7
            rank_3_to_check = chess.BB_RANK_6
        else:
            if chess.square_rank(king_position) != 0:
                return

        if chess.square_file(king_position) > 4:
            bb_side_of_board = chess.BB_FILE_F | chess.BB_FILE_G | chess.BB_FILE_H
            pawn_count_2 = len(
                chess.SquareSet(
                    self.board.pawns
                    & self.board.occupied_co[color]
                    & bb_side_of_board
                    & rank_2_to_check
                )
            )
            pawn_count_3 = len(
                chess.SquareSet(
                    self.board.pawns
                    & self.board.occupied_co[color]
                    & bb_side_of_board
                    & rank_3_to_check
                )
            )
            self.eval_result[color].king_shield_bonus += (
                pawn_count_2 * score_tables.additional_modifiers["king_shield_rank_2"]
            )
            self.eval_result[color].king_shield_bonus += (
                pawn_count_3 * score_tables.additional_modifiers["king_shield_rank_3"]
            )
        elif chess.square_file(king_position) < 3:
            bb_side_of_board = chess.BB_FILE_A | chess.BB_FILE_B | chess.BB_FILE_C
            pawn_count_2 = len(
                chess.SquareSet(
                    self.board.pawns
                    & self.board.occupied_co[color]
                    & bb_side_of_board
                    & rank_2_to_check
                )
            )
            pawn_count_3 = len(
                chess.SquareSet(
                    self.board.pawns
                    & self.board.occupied_co[color]
                    & bb_side_of_board
                    & rank_3_to_check
                )
            )
            self.eval_result[color].king_shield_bonus += (
                pawn_count_2 * score_tables.additional_modifiers["king_shield_rank_2"]
            )
            self.eval_result[color].king_shield_bonus += (
                pawn_count_3 * score_tables.additional_modifiers["king_shield_rank_3"]
            )

        self.eval_result[color].material_score_midgame += self.eval_result[
            color
        ].king_shield_bonus

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
            bb_enemy_spans |= bb_front_of_pawn

        passed_pawns = chess.SquareSet(
            self.board.pawns & self.board.occupied_co[color] & ~bb_enemy_spans
        )
        if color == chess.BLACK:
            passed_pawns = [chess.square_mirror(pawn) for pawn in list(passed_pawns)]
        self.eval_result[color].passed_pawn_bonus += sum(
            [
                score_tables.additional_piece_square_tables["passed_pawn"][pawn]
                for pawn in list(passed_pawns)
            ]
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

    def __str__(self) -> str:
        result = "---------- EVAL ---------- \n"
        result += str(self.board) + "\n"
        result += "--------- WHITE ---------- \n"
        for k, v in self.eval_result[chess.WHITE].__dict__.items():
            if isinstance(v, dict):
                result += f"{k}: \n"
                for piece, value in v.items():
                    result += f"    {chess.PIECE_NAMES[piece]}: {value} \n"
                if isinstance(value, int):
                    result += f"    #: {sum(v.values())} \n"
                elif isinstance(value, list):
                    result += f"    #: {sum([x_value for x_list in v.values() for x_value in x_list ])} \n"
            else:
                result += f"{k}: {v} \n"
        result += "--------- BLACK ---------- \n"
        for k, v in self.eval_result[chess.BLACK].__dict__.items():
            if isinstance(v, dict):
                result += f"{k}: \n"
                for piece, value in v.items():
                    result += f"    {chess.PIECE_NAMES[piece]}: {value} \n"
                if isinstance(value, int):
                    result += f"    #: {sum(v.values())} \n"
                elif isinstance(value, list):
                    result += f"    #: {sum([x_value for x_list in v.values() for x_value in x_list ])} \n"
            else:
                result += f"{k}: {v} \n"
        result += "--------- RESULT --------- \n"
        result += f"FEN: {self.board.fen()}\n"
        result += f"Total: {self.total_score} \n"
        result += "-------------------------- \n"
        return result