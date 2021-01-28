import chess


class Chatter:
    conversation_obj = None

    def __init__(self, conversation_obj):
        self.conversation_obj = conversation_obj

    def handle_gamestart(self):
        pass

    def handle_move(self, board):
        pass