from . import events
import logging


class Eventmanager:
    def __init__(self, conversation):
        self.conversation = conversation

    def registered_events(self):
        return events.BaseEvent.__subclasses__()

    def handle_events(self, board, game):
        logging.info(f"Checking events: {self.registered_events()}")
        for event_class in self.registered_events():
            event = event_class()
            if event.is_triggering(board, game):
                logging.info(f"triggering event {event}")
                if event.is_sending_to_player():
                    self.conversation.send_to_player(event.react_with(board, game))
                else:
                    self.conversation.send_to_spectator(event.react_with(board, game))
