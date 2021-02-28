from .base_event import BaseEvent
import chess


class GameEndEvent(BaseEvent):
    def __init__(self):
        pass

    def is_triggering(self, board: chess.Board, game):
        return board.is_game_over() or game.state["status"] != "started"

    def react_with(self, board, game):
        return f"This was a nice game. Feel free to challenge me anytime. Thanks for playing!"