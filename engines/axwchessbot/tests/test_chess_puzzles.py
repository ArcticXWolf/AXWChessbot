import pytest
import chess
from search import search
from tests import puzzles


@pytest.mark.slow
def test_chess_puzzle_wac():
    result = puzzle_execution(puzzles.wac.tests, True)
    assert result == len(puzzles.wac.tests)


def puzzle_execution(puzzle_list: list, output: bool):
    count_right = 0

    for test in puzzle_list:
        result_str = ""
        board = chess.Board()
        operations = board.set_epd(test)
        search_obj = search.Search(board, 40, 10, 30)
        result = search_obj.next_move_by_engine()
        if "bm" in operations:
            if board.san(result[0]) in [board.san(move) for move in operations["bm"]]:
                count_right += 1
                result_str += "OKAY | "
            else:
                result_str += "ERRO | "
        elif "am" in operations:
            if board.san(result[0]) not in [
                board.san(move) for move in operations["am"]
            ]:
                count_right += 1
                result_str += "OKAY | "
            else:
                result_str += "ERRO | "
        result_str += f"{str(result[1]['depth_reached']):>4s} | {board.san(result[0]):>6s} | {test}"
        if output:
            print(result_str)

    return count_right