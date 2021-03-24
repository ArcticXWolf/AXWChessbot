from .base_event import BaseEvent
from .opening_db import opening_db
import os
import chess
import random

STRINGS_UNKNOWN_OPENING = [
    "You hit me with an opening, that I do not know! Congratz!",
    "Oh, I do not know this opening. Guess I need to learn more theory then.",
    "Wow, that opening is not in my book. Did you invent it?",
    "Care to teach me this opening line? I have never seen it before.",
]

STRINGS_KNOWN_OPENING = [
    "Oh, the {}! This will be an exciting game!",
    "Hmm, my theory on the {} is a bit hazy...",
    "I studied the {} a lot, prepare for my attack!",
    "Wonderful, the {} is my favorite line!",
]


class OpeningDetectionEvent(BaseEvent):
    @staticmethod
    def setup_class():
        folder_path = os.path.join(
            os.path.dirname(os.path.abspath(__file__)), "opening_db"
        )
        OpeningDetectionEvent.opening_db_obj = opening_db.OpeningDB(folder_path)

    def __init__(self):
        pass

    def is_triggering(self, board: chess.Board, game):
        return len(board.move_stack) == 6

    def react_with(self, board, game):
        opening = OpeningDetectionEvent.opening_db_obj.find_opening_by_board(board)
        if opening:
            return random.choice(STRINGS_KNOWN_OPENING).format(opening.name)
        return random.choice(STRINGS_UNKNOWN_OPENING)