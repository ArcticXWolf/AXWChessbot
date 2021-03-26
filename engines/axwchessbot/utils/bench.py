import chess
from search import search
from tests import puzzles
from timeit import default_timer as timer
import yaml


class Benchmark:
    def __init__(
        self,
        abdepth: int = 5,
        qdepth: int = 6,
        num_positions: int = 20,
        output: bool = True,
    ):
        self.abdepth = abdepth
        self.qdepth = qdepth
        self.num_positions = num_positions
        self.output = output

    def run(self):
        bench_start = timer()
        single_tests = []
        for test in puzzles.wac.tests[: self.num_positions]:
            self.print(f"Running test: {test}")
            test_start = timer()
            self.run_test(test)
            test_end = timer()
            single_tests.append((test, test_end - test_start))
        bench_end = timer()
        self.print(f"Full runtime: {bench_end - bench_start:.3f} sec")

        return bench_end - bench_start, single_tests

    def run_test(self, test):
        board = chess.Board()
        board.set_epd(test)
        search_obj = search.Search(board, self.abdepth, self.qdepth, 0)
        search_obj.next_move_by_engine()
        self.print(yaml.dump(search_obj.get_measurements(show_cache=True), indent=4))

    def print(self, *args, **kwargs):
        if self.output:
            print(*args, **kwargs)
