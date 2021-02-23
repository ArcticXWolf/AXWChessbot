import chess
import chess.polyglot


class Entry:
    def __init__(self, val, flag, entry_depth, move):
        self.val = val
        self.flag = flag
        self.entry_depth = entry_depth
        self.move = move


class TranspositionTable:
    def __init__(self, size):
        self.size = size
        self.basic_cache = {}

    def __getitem__(self, position):
        return self.basic_cache.get(chess.polyglot.zobrist_hash(position), None)

    def store(self, position, value, flag, entry_depth, move):
        if len(self.basic_cache) > self.size:
            print("Trans Table is full!")
            return False
        self.basic_cache[chess.polyglot.zobrist_hash(position)] = Entry(
            value, flag, entry_depth, move
        )
        if len(self.basic_cache) % 5000 == 0:
            print("cache size:", len(self.basic_cache))

        return True

    def empty_cache(self):
        self.basic_cache = {}