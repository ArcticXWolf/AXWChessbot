import chess
import yaml
from search import search
from evaluation.game_over_detection import GameOverDetection


def test_if_search_runs_without_error():
    board = chess.Board()
    search_obj = search.Search(board, 1, 0, 0)
    search_obj.next_move_by_engine()


def test_if_engine_can_play_itself():
    board = chess.Board()
    while not GameOverDetection.is_game_over(board):
        search_obj = search.Search(board, 1, 0, 0)
        result = search_obj.next_move_by_engine()
        board.push(result)