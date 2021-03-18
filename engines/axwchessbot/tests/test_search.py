import chess
import yaml
from search import search


def test_if_search_runs_without_error():
    board = chess.Board()
    search_obj = search.Search(board, 1, 0, 0)
    search_obj.next_move_by_engine()


# def test_if_engine_can_play_itself():
#    board = chess.Board()
#    while not board.is_game_over():
#        search_obj = search.Search(board, 2, 4, 0)
#        result = search_obj.next_move_by_engine()
#        board.push(result[0])
#        print(board.unicode())
#        print("----------------")
#
#    raise Exception()