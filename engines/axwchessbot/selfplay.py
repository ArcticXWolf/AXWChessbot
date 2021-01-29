import datetime
import chess
import chess.pgn
import evaluation
import search
from timeit import default_timer as timer


def play_self():
    board = chess.Board("3r2k1/ppp3pp/8/3bn1K1/2q5/8/4r3/8 b - - 9 44")
    movehistory = []

    while not board.is_game_over(claim_draw=True):
        move = None
        start_search = end_search = None
        if board.turn:
            start_search = timer()
            move, _ = search.Search(board, 3).next_move()
            end_search = timer()
        else:
            start_search = timer()
            move, _ = search.Search(board, 3).next_move()
            end_search = timer()

        board.push(move)
        movehistory.append(move)

        print(f"===================================")
        print(f'-------------- {"WHITE" if board.turn else "BLACK" } --------------')
        print(f"The move was {move} (took {end_search - start_search :.2f} sec)")
        print(board)
        print(
            f"Evaluation result in this position: {evaluation.Evaluation(board).evaluate()}"
        )
        print(f"===================================")
        print("")

    game = chess.pgn.Game()
    game.headers["Event"] = "Example"
    game.headers["Site"] = "None"
    game.headers["Date"] = str(datetime.datetime.now().date())
    game.headers["Round"] = 1
    game.headers["White"] = "JNR Chess Engine"
    game.headers["Black"] = "JNR Chess Engine"
    game.add_line(movehistory)
    game.headers["Result"] = str(board.result(claim_draw=True))
    # print(game)
    print(f"===================================")
    print(f"RESULT (Turn {board.turn})")
    print(board)
    print(
        f"Evaluation result in this position: {evaluation.Evaluation(board).evaluate()}"
    )
    print(f"===================================")
    print("")


if __name__ == "__main__":
    play_self()