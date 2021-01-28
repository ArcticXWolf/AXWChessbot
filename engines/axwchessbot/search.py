import chess
import chess.polyglot
import chess.syzygy
import evaluation
import os


class Search:
    board = chess.Board()
    max_depth = 2
    opening_db_path = os.path.join(
        os.path.dirname(os.path.abspath(__file__)), "opening_db/Perfect2017-SF12.bin"
    )
    ending_db_path = os.path.join(
        os.path.dirname(os.path.abspath(__file__)), "ending_db"
    )
    ending_piece_count = 5

    def __init__(self, board: chess.Board, max_depth: int = 2):
        self.board = board
        self.max_depth = max_depth

    def next_move(self):
        move = self.next_move_by_opening_db()
        if move is not None:
            return move

        move = self.next_move_by_ending_db()
        if move is not None:
            return move

        move = self.next_move_by_engine()
        return move

    def next_move_by_engine(self):
        best_score = -99999
        best_move = None
        alpha = -100000
        beta = 100000

        for move in self.get_moves_by_value():
            self.board.push(move)
            score = -self.alpha_beta_search(-beta, -alpha, self.max_depth - 1)
            if score > best_score:
                best_score = score
                best_move = move
            if score > alpha:
                alpha = score
            self.board.pop()

        return best_move

    def alpha_beta_search(self, alpha: int, beta: int, depth_left: int = 0):
        best_score = -99999

        if depth_left <= 0 or self.board.is_game_over():
            return evaluation.Evaluation(self.board).evaluate()

        for move in self.get_moves_by_value():
            self.board.push(move)

            score = -self.alpha_beta_search(-beta, -alpha, depth_left - 1)

            self.board.pop()

            if score >= beta:
                return score
            if score > best_score:
                best_score = score
            if score > alpha:
                alpha = score

        return best_score

    def quiesce_search(self, alpha: int, beta: int):
        stand_pat = evaluation.Evaluation(self.board).evaluate()
        if stand_pat >= beta:
            return beta
        if alpha < stand_pat:
            alpha = stand_pat

        for move in self.get_moves_by_value():
            if self.board.is_capture(move):
                self.board.push(move)
                score = -self.quiesce_search(-beta, -alpha)
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
            current_dtz = tablebase.probe_dtz(self.board)
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
        def sort_function(move):
            return evaluation.Evaluation(self.board).move_value(move)

        moves_ordered = sorted(self.board.legal_moves, key=sort_function, reverse=True)
        return list(moves_ordered)