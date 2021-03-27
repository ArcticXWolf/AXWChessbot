from typing import Tuple
import chess
import chess.polyglot
import chess.syzygy
from evaluation import evaluation
import os

from evaluation.game_over_detection import GameOverDetection
from .cache import TranspositionTable
from .timeout import TimeOut
from timeit import default_timer as timer

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
        self.killer_moves = {}

        self.nodes_traversed = 0
        self.q_nodes_traversed = 0
        self.cache_hits = 0
        self.cache_cutoffs = 0
        self.max_depth_used = 0
        self.search_finished = False
        self.total_search_time = 0.0
        self.time_spent_per_depth = {}

        if cache:
            self.cache = cache

    def next_move(self):
        if self.search_finished:
            raise Exception("reused old search object!")

        move = self.next_move_by_opening_db()
        if move is not None:
            return move

        move = self.next_move_by_ending_db()
        if move is not None:
            return move

        return self.next_move_by_engine()

    def next_move_by_engine(self):
        start = timer()
        moves, _ = self.iterative_deepening()
        self.total_search_time = timer() - start
        self.search_finished = True
        return moves[-1]

    def iterative_deepening(self):
        board_copy = self.board.copy()
        start_depth_time = timer()
        moves, score = self.alpha_beta_search(1)
        self.time_spent_per_depth[1] = timer() - start_depth_time
        timeout = TimeOut(self.timeout)
        try:
            timeout.start()
            for i in range(2, self.alpha_beta_depth + 1):
                start_depth_time = timer()
                moves, score = self.alpha_beta_search(i, previous_moves=moves)
                self.time_spent_per_depth[i] = timer() - start_depth_time
                self.max_depth_used = i
        except TimeOut.TimeOutException as e:
            self.board = board_copy
        finally:
            timeout.disable_timeout()

        return moves, score

    def alpha_beta_search(
        self,
        depth_left: int = 0,
        alpha: float = float("-inf"),
        beta: float = float("inf"),
        move=None,
        previous_moves=None,
        ply=0,
    ) -> Tuple:
        best_score = float("-inf")
        best_move = None
        alpha_orig = alpha
        moves = []

        if depth_left <= 0 or GameOverDetection.is_game_over(self.board):
            moves.append(move)
            return (
                moves,
                self.quiesce_search(alpha, beta, self.quiesce_depth - 1),
            )

        self.nodes_traversed += 1

        cached = self.cache[self.board]
        if cached:
            self.cache_hits += 1
            if cached.entry_depth >= depth_left:
                if cached.flag == EXACT:
                    move = cached.move if not move else move
                    moves.append(move)
                    self.cache_cutoffs += 1
                    return moves, cached.val
                elif cached.flag == LOWER:
                    alpha = max(alpha, cached.val)
                elif cached.flag == UPPER:
                    beta = min(beta, cached.val)
                if alpha >= beta:
                    move = cached.move if not move else move
                    moves.append(move)
                    self.cache_cutoffs += 1
                    return moves, cached.val

        move_list_to_choose_from = evaluation.Evaluation(self.board).move_order()

        if ply in self.killer_moves:
            if self.killer_moves[ply][1] is not None:
                try:
                    move_list_to_choose_from.remove(self.killer_moves[ply][1])
                    move_list_to_choose_from.insert(0, self.killer_moves[ply][1])
                except ValueError:
                    pass
            if self.killer_moves[ply][0] is not None:
                try:
                    move_list_to_choose_from.remove(self.killer_moves[ply][0])
                    move_list_to_choose_from.insert(0, self.killer_moves[ply][0])
                except ValueError:
                    pass

        if (
            previous_moves
            and len(previous_moves) > depth_left
            and previous_moves[depth_left - 1] in move_list_to_choose_from
        ):
            try:
                move_list_to_choose_from.remove(previous_moves[depth_left - 1])
                move_list_to_choose_from.insert(0, previous_moves[depth_left - 1])
            except ValueError:
                pass

        if cached:
            try:
                move_list_to_choose_from.remove(cached.move)
                move_list_to_choose_from.insert(0, cached.move)
            except ValueError:
                pass

        for m in move_list_to_choose_from:
            self.board.push(m)

            new_moves, score = self.alpha_beta_search(
                depth_left - 1, -beta, -alpha, m, previous_moves, ply + 1
            )
            score = -score

            self.board.pop()

            if score > best_score:
                moves = new_moves
                best_score = score
                best_move = m
            if score > alpha:
                alpha = score
            if alpha >= beta:
                self.update_killer_moves(ply, m)
                break

        if best_score <= alpha_orig:
            flag = UPPER
        elif best_score >= beta:
            flag = LOWER
        else:
            flag = EXACT

        if not best_move:
            best_move = m
        self.cache.store(self.board, best_score, flag, depth_left, best_move)
        moves.append(best_move)
        return (moves, best_score)

    def quiesce_search(self, alpha: float, beta: float, depth_left: int = 0):
        self.nodes_traversed += 1
        self.q_nodes_traversed += 1

        stand_pat = evaluation.Evaluation(self.board).evaluate().total_score_perspective
        if stand_pat >= beta:
            return beta
        if alpha < stand_pat:
            alpha = stand_pat

        if depth_left > 0 or not GameOverDetection.is_game_over(self.board):
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

    def update_killer_moves(self, ply: int, new_move: chess.Move) -> None:
        if not self.board.is_capture(new_move):
            if ply not in self.killer_moves:
                self.killer_moves[ply] = [new_move, None]
                return

            if new_move != self.killer_moves[ply][0]:
                self.killer_moves[ply][1] = self.killer_moves[ply][0]
            self.killer_moves[ply][0] = new_move

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
        return evaluation.Evaluation(self.board).capture_order()

    def get_measurements(
        self, show_exact_timings: bool = False, show_cache: bool = False
    ):
        result = {
            "finished": self.search_finished,
            "max_depth_used": self.max_depth_used,
            "nodes_traversed": self.nodes_traversed,
            "q_nodes_traversed": self.q_nodes_traversed,
            "total_search_time": self.total_search_time,
        }
        if self.nodes_traversed >= 0 and self.total_search_time >= 1.0:
            result["nps"] = float(self.nodes_traversed) / self.total_search_time
        if show_exact_timings:
            result["time_spent_per_depth"] = self.time_spent_per_depth
        if show_cache:
            result["cache_hits"] = self.cache_hits
            result["cache_cutoffs"] = self.cache_cutoffs
            result["cache_length"] = self.cache.get_length()
        return result

    def __str__(self):
        measurements = [f"{k}={str(v)}" for k, v in self.get_measurements().items()]
        return f"<Search {' '.join(measurements)}>"