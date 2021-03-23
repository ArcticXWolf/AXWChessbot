import chess
from evaluation import evaluation


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


def test_if_evaluation_calculates_piece_square_score_correctly():
    board = chess.Board()
    eval = evaluation.Evaluation(board)
    eval.evaluate_material_score(chess.WHITE)
    eval.evaluate_material_score(chess.BLACK)

    assert eval.eval_result[chess.WHITE].piece_square_score_midgame[chess.PAWN] == -66
    assert eval.eval_result[chess.BLACK].piece_square_score_midgame[chess.PAWN] == -66


def test_if_evaluation_calculates_king_shield_bonus_correctly():
    board = chess.Board("6k1/5ppp/8/8/8/8/PPP5/1K6 w - - 0 1")
    eval = evaluation.Evaluation(board)
    eval.evaluate_king_shield(chess.WHITE)
    eval.evaluate_king_shield(chess.BLACK)

    assert eval.eval_result[chess.WHITE].king_shield_bonus == 30
    assert eval.eval_result[chess.BLACK].king_shield_bonus == 30


def test_if_evaluation_calculates_king_shield_bonus_correctly1():
    board = chess.Board("6k1/5p1p/6p1/8/8/1P6/P1P5/1K6 w - - 0 1")
    eval = evaluation.Evaluation(board)
    eval.evaluate_king_shield(chess.WHITE)
    eval.evaluate_king_shield(chess.BLACK)

    assert eval.eval_result[chess.WHITE].king_shield_bonus == 25
    assert eval.eval_result[chess.BLACK].king_shield_bonus == 25


def test_if_evaluation_calculates_passed_pawn_bonus_correctly():
    board = chess.Board("rnbqkbnr/ppppp3/8/8/8/8/PPPPPPP1/RNBQKBNR w KQkq - 0 1")
    eval = evaluation.Evaluation(board)
    eval.evaluate()

    assert eval.eval_result[chess.WHITE].passed_pawn_bonus == 20
    assert eval.eval_result[chess.BLACK].passed_pawn_bonus == 0


def test_if_evaluation_calculates_passed_pawn_bonus_correctly1():
    board = chess.Board("rnbqkbnr/ppppp2p/7p/8/8/8/PPPPPP2/RNBQKBNR w KQkq - 0 1")
    eval = evaluation.Evaluation(board)
    eval.evaluate()

    assert eval.eval_result[chess.WHITE].passed_pawn_bonus == 0
    assert eval.eval_result[chess.BLACK].passed_pawn_bonus == 20


def test_if_evaluation_calculates_passed_pawn_bonus_correctly2():
    board = chess.Board("rnbqkbnr/p4p1p/7p/8/8/8/PPPPPP2/RNBQKBNR w KQkq - 0 1")
    eval = evaluation.Evaluation(board)
    eval.evaluate()

    assert eval.eval_result[chess.WHITE].passed_pawn_bonus == 40
    assert eval.eval_result[chess.BLACK].passed_pawn_bonus == 20