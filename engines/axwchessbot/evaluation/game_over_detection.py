import chess
import collections


# Unfortunately the gameover detection of python-chess is very slow
# because it treats the game as drawn if one player can claim draw
# AFTER their next move. This means the detection loops over every
# possible move and that calls the move generator which is SLOW.
# We dont need this claiming mechanic here, so we implement our own
# detection. (Speedup of 2x)
class GameOverDetection:
    def is_game_over(board) -> bool:
        if board.is_seventyfive_moves():
            return True

        # Insufficient material.
        if board.is_insufficient_material():
            return True

        # Stalemate or checkmate.
        if not any(board.generate_legal_moves()):
            return True

        # Fivefold repetition.
        if board.is_fivefold_repetition():
            return True

        # Fifty move rule
        if board.halfmove_clock > 100:
            return True

        # Threefold repetition
        if GameOverDetection.is_threefold_repetition(board):
            return True

        return False

    def is_threefold_repetition(board) -> bool:
        transposition_key = board._transposition_key()
        transpositions = collections.Counter()
        transpositions.update((transposition_key,))

        # Count positions.
        switchyard = []
        while board.move_stack:
            move = board.pop()
            switchyard.append(move)

            if board.is_irreversible(move):
                break

            transpositions.update((board._transposition_key(),))

        while switchyard:
            board.push(switchyard.pop())

        # Threefold repetition occured.
        if transpositions[transposition_key] >= 3:
            return True

        return False