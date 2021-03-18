import chess
from ..evaluation import evaluation


def test_if_evaluation_runs_without_error():
    board = chess.Board()
    eval = evaluation.Evaluation(board)
    eval.evaluate()


def test_if_evaluation_calculates_piece_score_correctly():
    board = chess.Board()
    eval = evaluation.Evaluation(board)
    eval.evaluate_material_score(chess.WHITE)
    eval.evaluate_material_score(chess.BLACK)

    assert sum(eval.eval_result[chess.WHITE].piece_score.values()) == 4000
    assert sum(eval.eval_result[chess.BLACK].piece_score.values()) == 4000


def test_if_evaluation_calculates_piece_score_correctly1():
    board = chess.Board("r2qkb1r/p6p/8/8/8/8/PP5P/R3K1NR w KQkq - 0 1")
    eval = evaluation.Evaluation(board)
    eval.evaluate_material_score(chess.WHITE)
    eval.evaluate_material_score(chess.BLACK)

    assert sum(eval.eval_result[chess.WHITE].piece_score.values()) == 1620
    assert sum(eval.eval_result[chess.BLACK].piece_score.values()) == 2430