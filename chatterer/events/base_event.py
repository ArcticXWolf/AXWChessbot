class BaseEvent:
    def __init__(self):
        pass

    def is_triggering(self, board, game):
        return False

    def is_sending_to_player(self):
        return True

    def react_with(self, board, game):
        return ""