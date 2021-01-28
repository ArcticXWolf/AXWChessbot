import datetime
import chess
import chess.pgn
import evaluation
import search


def play_self():
    board = chess.Board()
    movehistory = []

    while not board.is_game_over(claim_draw=True):
        move = chess.Move.null()
        if board.turn:
            move = search.Search(board, 3).next_move()
        else:
            move = search.Search(board, 3).next_move()

        board.push(move)
        movehistory.append(move)

        print(f"===================================")
        print(f'-------------- {"WHITE" if board.turn else "BLACK" } --------------')
        print(f"The move was {move}")
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
    print(game)


if __name__ == "__main__":
    play_self()