from .base_event import BaseEvent
import chess
import random

STRINGS_MSG = [
    "This was a nice game. Feel free to challenge me anytime. Thanks for playing!",
    "GG! Do you have time for another game?",
    "Well played! I hope you want to play again in the future.",
    "That was fun! If you know of any bugs, improvements or are just curious, you can find the link to the sourcecode in my profile.",
]


class GameEndEvent(BaseEvent):
    def __init__(self):
        pass

    def is_triggering(self, board: chess.Board, game):
        return board.is_game_over() or game.state["status"] != "started"

    def react_with(self, board, game):
        return random.choice(STRINGS_MSG)