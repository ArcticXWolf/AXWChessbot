import chess
import evaluation
import search
from tests import tests as test_modules
from timeit import default_timer as timer


def print_result(
    test_module, id, time, move, result, operation, solution, info: dict, board
):
    info.pop("moves_analysis")
    print(
        f"[{test_module}#{id+1:03d}] {time :05.2f}s {result}: {str(board.san(move)):>6s} vs {operation} {solution} - {str(info)}"
    )


def test_engine():
    print("Starting tests.")
    start_whole = timer()
    for test_module, tests in test_modules.items():
        if test_module in []:
            continue
        print(f"Starting module {test_module}")
        solved = 0

        start_module = timer()
        for i in range(len(tests)):
            board = chess.Board()
            operations = board.set_epd(tests[i])
            if "am" in operations.keys():
                operations["am"] = [str(board.san(x)) for x in operations["am"]]
            if "bm" in operations.keys():
                operations["bm"] = [str(board.san(x)) for x in operations["bm"]]

            start_search = timer()
            move, info = search.Search(board, 10, 4, 15).next_move()
            end_search = timer()

            if "am" in operations.keys() and str(board.san(move)) in operations["am"]:
                print_result(
                    test_module,
                    i,
                    end_search - start_search,
                    move,
                    "ERRO",
                    "am",
                    operations["am"],
                    info,
                    board,
                )
            elif (
                "bm" in operations.keys()
                and str(board.san(move)) not in operations["bm"]
            ):
                print_result(
                    test_module,
                    i,
                    end_search - start_search,
                    move,
                    "ERRO",
                    "bm",
                    operations["bm"],
                    info,
                    board,
                )
            elif "am" not in operations.keys() and "bm" not in operations.keys():
                print_result(
                    test_module,
                    i,
                    end_search - start_search,
                    move,
                    "ERRO",
                    "",
                    "N/A",
                    info,
                    board,
                )
            else:
                solved += 1
                print_result(
                    test_module,
                    i,
                    end_search - start_search,
                    move,
                    "OKAY",
                    "am" if "am" in operations.keys() else "bm",
                    operations["am"] if "am" in operations.keys() else operations["bm"],
                    info,
                    board,
                )

        end_module = timer()
        print(
            f"Ending module {test_module}, solved {solved}/{len(tests)} in {end_module - start_module :05.2f}s"
        )
    end_whole = timer()
    print(f"Ending tests in {end_whole - start_whole :05.2f}s")
