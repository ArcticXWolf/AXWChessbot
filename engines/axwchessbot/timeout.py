import signal
import time


class TimeOut:
    """
    TimeOut for *nix systems
    """

    class TimeOutException(Exception):
        pass

    def __init__(self, sec):
        self.sec = sec  # time we want an exception to be raised after
        self.start_time = None  # to measure actual time

    def start(self):
        self.start_time = time.time()
        signal.signal(signal.SIGALRM, self.raise_timeout)
        signal.alarm(self.sec)

    def raise_timeout(self, *args):
        signal.alarm(0)  # disable
        message = "Alarm set for {}s and actually passed {}s".format(
            self.sec, round(time.time() - self.start_time, 5)
        )
        raise TimeOut.TimeOutException(message)