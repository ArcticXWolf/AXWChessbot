from .base_event import BaseEvent
import chess
import random

STRINGS_MSG = [
    "Hello, nice to meet you! Have fun and good luck! :)",
    "A new challenger! May the best player (or bot) win :)",
    "A new face! GL and HF!",
    "Heyho, I'm a chessbot developed by ArcticXWolf, who instructed me to send you this message: 'Crush the bot!'",
]


class GameStartEvent(BaseEvent):
    def __init__(self):
        pass

    def is_triggering(self, board: chess.Board, game):
        return len(board.move_stack) <= 1

    def react_with(self, board, game):
        return random.choice(STRINGS_MSG)