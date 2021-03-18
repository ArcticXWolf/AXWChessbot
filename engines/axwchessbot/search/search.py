from typing import Tuple
import chess
import chess.polyglot
import chess.syzygy
from evaluation import evaluation
import os
from .cache import TranspositionTable
from .timeout import TimeOut

LOWER = -1
EXACT = 0
UPPER = 1


class Search:
    board = chess.Board()
    alpha_beta_depth = 2
    quiesce_depth = 10
    opening_db_path = os.path.join(
        os.path.dirname(os.path.abspath(__file__)), "../opening_db/jnr-combine.bin"
    )
    ending_db_path = os.path.join(
        os.path.dirname(os.path.abspath(__file__)), "../ending_db"
    )
    ending_piece_count = 5
    cache = TranspositionTable(1e7)

    def __init__(
        self,
        board: chess.Board,
        alpha_beta_depth: int = 10,
        quiesce_depth: int = 10,
        timeout: int = 180,
        cache: TranspositionTable = None,
    ):
        self.board = board.copy()
        self.alpha_beta_depth = alpha_beta_depth
        self.quiesce_depth = quiesce_depth
        self.cache = TranspositionTable(1e7)
        self.timeout = timeout

        if cache:
            self.cache = cache

    def next_move(self):
        debug_info = {"moves_analysis": [], "positions_analyzed": 0, "cache_hits": 0}
        move = self.next_move_by_opening_db()
        if move is not None:
            debug_info["moves"] = [move.uci()]
            debug_info["eval"] = "opening_db"
            debug_info["depth_reached"] = 0
            return (move, debug_info)

        move = self.next_move_by_ending_db()
        if move is not None:
            debug_info["moves"] = [move.uci()]
            debug_info["eval"] = "ending_db"
            debug_info["depth_reached"] = 0
            return (move, debug_info)

        return self.next_move_by_engine()

    def next_move_by_engine(self):
        moves, score, debug_info = self.iterative_deepening()
        debug_info["moves"] = [move.uci() for move in moves]
        eval = evaluation.Evaluation(self.board).evaluate()
        debug_info["current_eval"] = (
            eval.eval_result[chess.WHITE].__dict__,
            eval.eval_result[chess.BLACK].__dict__,
        )
        return moves[-1], debug_info

    def iterative_deepening(self):
        depth_reached = 1
        board_copy = self.board.copy()
        moves, score, debug_info = self.alpha_beta_search(1)
        timeout = TimeOut(self.timeout)
        try:
            timeout.start()
            for i in range(2, self.alpha_beta_depth + 1):
                moves, score, debug_info = self.alpha_beta_search(
                    i, previous_moves=moves
                )
                depth_reached += 1
        except TimeOut.TimeOutException as e:
            self.board = board_copy
        finally:
            timeout.disable_timeout()

        debug_info["depth_reached"] = depth_reached

        return moves, score, debug_info

    def alpha_beta_search(
        self,
        depth_left: int = 0,
        alpha: float = float("-inf"),
        beta: float = float("inf"),
        move=None,
        previous_moves=None,
    ) -> Tuple:
        best_score = float("-inf")
        best_move = None
        alpha_orig = alpha
        moves = []
        debug_info = {"moves_analysis": [], "positions_analyzed": 0, "cache_hits": 0}

        cached = self.cache[self.board]
        if cached and cached.entry_depth >= depth_left:
            if cached.flag == EXACT:
                move = cached.move if not move else move
                debug_info["positions_analyzed"] += 1
                debug_info["cache_hits"] += 1
                moves.append(move)
                return moves, cached.val, debug_info
            elif cached.flag == LOWER:
                alpha = max(alpha, cached.val)
            elif cached.flag == UPPER:
                beta = min(beta, cached.val)
            if alpha >= beta:
                move = cached.move if not move else move
                debug_info["positions_analyzed"] += 1
                debug_info["cache_hits"] += 1
                moves.append(move)
                return moves, cached.val, debug_info

        if depth_left <= 0 or self.board.is_game_over():
            moves.append(move)
            debug_info["positions_analyzed"] += 1
            return (
                moves,
                self.quiesce_search(alpha, beta, self.quiesce_depth - 1),
                debug_info,
            )

        move_list_to_choose_from = evaluation.Evaluation(self.board).move_order()

        if (
            previous_moves
            and len(previous_moves) > depth_left
            and previous_moves[depth_left - 1] in move_list_to_choose_from
        ):
            move_list_to_choose_from.insert(0, previous_moves[depth_left - 1])

        for m in move_list_to_choose_from:
            san = self.board.san(m)
            self.board.push(m)

            new_moves, score, new_debug_info = self.alpha_beta_search(
                depth_left - 1, -beta, -alpha, m, previous_moves
            )
            score = -score
            debug_info["moves_analysis"].append(
                (
                    str(san),
                    str(evaluation.Evaluation(self.board).evaluate()),
                    score,
                )
            )
            debug_info["positions_analyzed"] += new_debug_info["positions_analyzed"]
            debug_info["cache_hits"] += new_debug_info["cache_hits"]

            self.board.pop()

            if score > best_score:
                moves = new_moves
                best_score = score
                best_move = m
            if score > alpha:
                alpha = score
            if alpha >= beta:
                break

        if best_score <= alpha_orig:
            flag = UPPER
        elif best_score >= beta:
            flag = LOWER
        else:
            flag = EXACT

        if not best_move:
            best_move = m
        self.cache.store(
            self.board, best_score, flag, depth_left, best_move, debug_info
        )
        moves.append(best_move)
        return (moves, best_score, debug_info)

    def quiesce_search(self, alpha: float, beta: float, depth_left: int = 0):

        stand_pat = evaluation.Evaluation(self.board).evaluate().total_score_perspective
        if stand_pat >= beta:
            return beta
        if alpha < stand_pat:
            alpha = stand_pat

        if depth_left > 0 or not self.board.is_game_over():
            for move in self.get_captures_by_value():
                if self.board.is_capture(move):
                    self.board.push(move)
                    score = -self.quiesce_search(-beta, -alpha, depth_left - 1)
                    self.board.pop()

                    if score >= beta:
                        return beta
                    if score > alpha:
                        alpha = score
        return alpha

    def next_move_by_opening_db(self):
        try:
            move = (
                chess.polyglot.MemoryMappedReader(self.opening_db_path)
                .weighted_choice(self.board)
                .move
            )
            # fix error in which move is a function.
            if callable(move):
                move = move()
            return move
        except IndexError:
            return None

    def next_move_by_ending_db(self):
        if chess.popcount(self.board.occupied) > self.ending_piece_count:
            return None

        tablebase = chess.syzygy.open_tablebase(self.ending_db_path)
        chosen_move = None

        try:
            current_wdl = tablebase.probe_wdl(self.board)
            dtz_moves = {}

            for move in self.board.legal_moves:
                self.board.push(move)
                new_dtz = -tablebase.probe_dtz(self.board)
                if new_dtz not in dtz_moves:
                    dtz_moves[new_dtz] = []

                dtz_moves[new_dtz].append(move)

                self.board.pop()

            if current_wdl >= 0:
                best_index = min([i for i in dtz_moves.keys() if i > 0], default=0)
                chosen_move = dtz_moves[best_index].pop()
            else:
                if 0 in dtz_moves.keys():
                    chosen_move = dtz_moves[0].pop()
                else:
                    chosen_move = dtz_moves[min(dtz_moves.keys())].pop()
        except:
            tablebase.close()
            return None

        tablebase.close()
        return chosen_move

    def get_captures_by_value(self):
        def sort_function(move):
            return evaluation.Evaluation(self.board).capture_value(move)

        captures = [
            move for move in self.board.legal_moves if self.board.is_capture(move)
        ]
        captures_ordered = sorted(captures, key=sort_function, reverse=True)
        return list(captures_ordered)