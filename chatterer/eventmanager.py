from . import events


class Eventmanager:
    def __init__(self, conversation):
        self.conversation = conversation

        self.handle_event_setups()

    def handle_event_setups(self):
        for event_class in self.registered_events():
            event_class.setup_class()

    def registered_events(self):
        return events.BaseEvent.__subclasses__()

    def handle_events(self, board, game):
        for event_class in self.registered_events():
            event = event_class()
            if event.is_triggering(board, game):
                if event.is_sending_to_player():
                    self.conversation.send_to_player(event.react_with(board, game))
                else:
                    self.conversation.send_to_spectator(event.react_with(board, game))
