from .base_event import BaseEvent
import chess


class GameStartEvent(BaseEvent):
    def __init__(self):
        pass

    def is_triggering(self, board: chess.Board, game):
        return len(board.move_stack) <= 1

    def react_with(self, board, game):
        return "Hello, nice to meet you! Have fun and good luck! :)"