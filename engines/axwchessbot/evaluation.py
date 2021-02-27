from typing import Tuple
import chess
import score_tables
import itertools


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

        score = int(self.evaluate_material_score())
        score += self.evaluate_pair_bonus()
        score += self.evaluate_tempo()
        score += self.evaluate_blocked_pieces(chess.WHITE)
        score += self.evaluate_blocked_pieces(chess.BLACK)

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
            return (True, float("-inf"))
        elif result == "1-0":
            return (True, float("inf"))
        return (True, 0)

    def evaluate_material_score(self) -> float:
        midgame_score = float(self.evaluate_material_score_phase(False))
        endgame_score = float(self.evaluate_material_score_phase(True))
        gamestate = float(max(self.evaluate_gamephase(), 24))

        return (
            (gamestate * midgame_score) + ((24.0 - gamestate) * endgame_score)
        ) / 24.0

    def evaluate_material_score_phase(self, is_endgame: bool) -> int:
        pst = score_tables.piece_square_tables
        if is_endgame:
            pst = score_tables.piece_square_tables_endgame

        score = 0
        for color in chess.COLORS:
            color_score = 0
            for piece in chess.PIECE_TYPES:
                pieces = self.board.pieces(piece, color)
                if color == chess.WHITE:
                    pieces_pst_scores = [pst[piece][i] for i in list(pieces)]
                else:
                    pieces_pst_scores = [
                        pst[piece][chess.square_mirror(i)] for i in list(pieces)
                    ]

                color_score += (
                    sum(pieces_pst_scores)
                    + len(pieces) * score_tables.piece_values[piece]
                )
            if color == chess.WHITE:
                score += color_score
            else:
                score -= color_score

        if not is_endgame:
            score += self.evaluate_king_shield(chess.WHITE)
            score += self.evaluate_king_shield(chess.BLACK)

        return score

    def evaluate_gamephase(self):
        return (
            len(self.board.pieces(chess.KNIGHT, chess.WHITE))
            + len(self.board.pieces(chess.BISHOP, chess.WHITE))
            + 2 * len(self.board.pieces(chess.ROOK, chess.WHITE))
            + 4 * len(self.board.pieces(chess.QUEEN, chess.WHITE))
            + len(self.board.pieces(chess.KNIGHT, chess.BLACK))
            + len(self.board.pieces(chess.BISHOP, chess.BLACK))
            + 2 * len(self.board.pieces(chess.ROOK, chess.BLACK))
            + 4 * len(self.board.pieces(chess.QUEEN, chess.BLACK))
        )

    def evaluate_pair_bonus(self):
        score = 0

        if len(self.board.pieces(chess.BISHOP, chess.WHITE)) > 1:
            score += score_tables.additional_modifiers["bishop_pair"]
        if len(self.board.pieces(chess.KNIGHT, chess.WHITE)) > 1:
            score -= score_tables.additional_modifiers["knight_pair"]
        if len(self.board.pieces(chess.ROOK, chess.WHITE)) > 1:
            score -= score_tables.additional_modifiers["rook_pair"]

        if len(self.board.pieces(chess.BISHOP, chess.BLACK)) > 1:
            score -= score_tables.additional_modifiers["bishop_pair"]
        if len(self.board.pieces(chess.KNIGHT, chess.BLACK)) > 1:
            score += score_tables.additional_modifiers["knight_pair"]
        if len(self.board.pieces(chess.ROOK, chess.BLACK)) > 1:
            score += score_tables.additional_modifiers["rook_pair"]

        return score

    def evaluate_tempo(self):
        return (
            score_tables.additional_modifiers["tempo"]
            if self.board.turn == chess.WHITE
            else -score_tables.additional_modifiers["tempo"]
        )

    def evaluate_passed_pawns(self):
        score = 0
        for color in chess.COLORS:
            for pawn in self.board.pieces(chess.PAWN, color):
                if self.is_passed_pawn(pawn, color):
                    if color == chess.WHITE:
                        score += score_tables.additional_piece_square_tables[
                            "passed_pawn"
                        ][pawn]
                    else:
                        score -= score_tables.additional_piece_square_tables[
                            "passed_pawn"
                        ][pawn]
        return score

    def is_passed_pawn(self, square: chess.Square, color: chess.Color) -> bool:
        piece = self.board.piece_type_at(square)
        if self.board.piece_type_at(square) != chess.PAWN:
            return False

        ranks_to_go = range(chess.square_rank(square) + 1, 8)
        if color == chess.BLACK:
            ranks_to_go = range(0, chess.square_rank(square))

        amount_pawns_on_same_file = sum(
            [
                self.board.piece_type_at(chess.square(chess.square_file(square), i))
                == chess.PAWN
                for i in ranks_to_go
            ]
        )
        if amount_pawns_on_same_file > 0:
            return False

        amount_enemy_pawns_on_adjacent_files = 0
        if chess.square_file(square) - 1 >= 0:
            amount_enemy_pawns_on_adjacent_files += sum(
                [
                    self.board.piece_at(chess.square(chess.square_file(square) - 1, i))
                    == chess.Piece(chess.PAWN, not color)
                    for i in ranks_to_go
                ]
            )
        if chess.square_file(square) + 1 <= 8:
            amount_enemy_pawns_on_adjacent_files += sum(
                [
                    self.board.piece_at(chess.square(chess.square_file(square) + 1, i))
                    == chess.Piece(chess.PAWN, not color)
                    for i in ranks_to_go
                ]
            )

        if amount_enemy_pawns_on_adjacent_files > 0:
            return False

        return True

    def evaluate_king_shield(self, color: chess.Color):
        score = 0
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
            score += (
                pawn_count_2 * score_tables.additional_modifiers["king_shield_rank_2"]
            )
            score += (
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
            score += (
                pawn_count_2 * score_tables.additional_modifiers["king_shield_rank_2"]
            )
            score += (
                pawn_count_3 * score_tables.additional_modifiers["king_shield_rank_3"]
            )

        if color == chess.BLACK:
            return -score
        return score

    def evaluate_blocked_pieces(self, color: chess.Color):
        score = 0

        # king blocks rook
        side_rank = 0 if color == chess.WHITE else 7
        if (
            self.board.piece_type_at(chess.square(5, side_rank)) == chess.KING
            or self.board.piece_type_at(chess.square(6, side_rank)) == chess.KING
        ) and (
            self.board.piece_type_at(chess.square(6, side_rank)) == chess.ROOK
            or self.board.piece_type_at(chess.square(7, side_rank)) == chess.ROOK
        ):
            score -= score_tables.additional_modifiers["king_blocks_rook_penalty"]

        if (
            self.board.piece_type_at(chess.square(1, side_rank)) == chess.KING
            or self.board.piece_type_at(chess.square(2, side_rank)) == chess.KING
        ) and (
            self.board.piece_type_at(chess.square(0, side_rank)) == chess.ROOK
            or self.board.piece_type_at(chess.square(1, side_rank)) == chess.ROOK
        ):
            score -= score_tables.additional_modifiers["king_blocks_rook_penalty"]

        if color == chess.BLACK:
            return -score
        return score

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
