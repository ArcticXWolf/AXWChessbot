import chess
import chess.polyglot
import chess.syzygy
import evaluation
import os


class Search:
    board = chess.Board()
    alpha_beta_depth = 2
    quiesce_depth = 10
    opening_db_path = os.path.join(
        os.path.dirname(os.path.abspath(__file__)), "opening_db/Perfect2017-SF12.bin"
    )
    ending_db_path = os.path.join(
        os.path.dirname(os.path.abspath(__file__)), "ending_db"
    )
    ending_piece_count = 5

    def __init__(
        self, board: chess.Board, alpha_beta_depth: int = 2, quiesce_depth: int = 10
    ):
        self.board = board
        self.alpha_beta_depth = alpha_beta_depth
        self.quiesce_depth = quiesce_depth

    def next_move(self):
        move = self.next_move_by_opening_db()
        if move is not None:
            return (move, -1)

        move = self.next_move_by_ending_db()
        if move is not None:
            return (move, -2)

        move, num_positions = self.next_move_by_engine()
        return (move, num_positions)

    def next_move_by_engine(self):
        alpha = -100000
        beta = 100000
        analyzed_positions = 0
        moves = []

        for move in self.get_moves_by_value():
            self.board.push(move)
            score, num_positions, depth = self.alpha_beta_search(
                -beta, -alpha, self.alpha_beta_depth - 1
            )
            score = -score
            analyzed_positions += num_positions
            moves.append((move, score, depth))
            if score > alpha:
                alpha = score
            self.board.pop()

        sort_depth = self.alpha_beta_depth

        def sort_function(move_eval):
            if move_eval[1] == 9999:
                return (move_eval[1], sort_depth - move_eval[2])
            return (move_eval[1], 1)

        moves_ordered = sorted(moves, reverse=True, key=sort_function)
        print(moves_ordered)
        return (moves_ordered.pop(0)[0], analyzed_positions)

    def alpha_beta_search(self, alpha: int, beta: int, depth_left: int = 0):
        best_score = -99999
        analyzed_positions = 0
        best_move_depth = 0

        if depth_left <= 0 or self.board.is_game_over():
            return (self.quiesce_search(alpha, beta, self.quiesce_depth - 1), 1, 0)
            # return (evaluation.Evaluation(self.board).evaluate(), 1)

        for move in self.get_moves_by_value():
            self.board.push(move)

            score, num_positions, depth = self.alpha_beta_search(
                -beta, -alpha, depth_left - 1
            )
            score = -score
            analyzed_positions += num_positions

            self.board.pop()

            if score >= beta:
                return (score, analyzed_positions, depth)
            if score > best_score:
                best_score = score
                best_move_depth = depth
            if score > alpha:
                alpha = score

        return (best_score, analyzed_positions, best_move_depth + 1)

    def quiesce_search(self, alpha: int, beta: int, depth_left: int = 0):

        stand_pat = evaluation.Evaluation(self.board).evaluate()
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

    def get_moves_by_value(self):
        is_endgame = evaluation.Evaluation(self.board).check_if_endgame()

        def sort_function(move):
            return evaluation.Evaluation(self.board).move_value(move, is_endgame)

        moves_ordered = sorted(self.board.legal_moves, key=sort_function, reverse=True)
        return list(moves_ordered)

    def get_captures_by_value(self):
        def sort_function(move):
            return evaluation.Evaluation(self.board).capture_value(move)

        captures = [
            move for move in self.board.legal_moves if self.board.is_capture(move)
        ]
        captures_ordered = sorted(captures, key=sort_function, reverse=True)
        return list(captures_ordered)