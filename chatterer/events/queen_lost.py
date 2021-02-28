from .base_event import BaseEvent
import chess
import logging


class QueenLostEvent(BaseEvent):
    def __init__(self):
        pass

    def is_triggering(self, board: chess.Board, game):
        old_board = board.copy()
        move = old_board.pop()
        if len(old_board.pieces(chess.QUEEN, game.is_white)) > len(
            board.pieces(chess.QUEEN, game.is_white)
        ):
            return True
        return False

    def react_with(self, board: chess.Board, game):
        return f"(ノಠ益ಠ)ノ彡┻━┻ NO, MY QUEEN!"