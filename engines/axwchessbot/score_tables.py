import chess

piece_values = {
    chess.PAWN: 100,
    chess.KNIGHT: 320,
    chess.BISHOP: 330,
    chess.ROOK: 500,
    chess.QUEEN: 900,
    chess.KING: 2000 # has no value, will not be used
}

additional_modifiers = {
    "bishop_pair": 30,
    "knight_pair": 8,
    "rook_pair": 16,
    "open_rook": 10,
    "half_rook": 5,
    "tempo": 10,
    "king_shield_rank_2": 10,
    "king_shield_rank_3": 5,
    "king_blocks_rook_penalty": 24
}

additional_settings = {
    "endgame_threshold": 1300
}

piece_square_tables = {
    chess.PAWN: [ 0,   0,   0,   0,    0,  0,   0,   0,
                 -6,  -4,   1, -24,  -24,  1,  -4,  -6,
                 -4,  -4,   1,   5,   5,   1,  -4,  -4,
                 -6,  -4,   5,  10,  10,   5,  -4,  -6,
                 -6,  -4,   2,   8,   8,   2,  -4,  -6,
                 -6,  -4,   1,   2,   2,   1,  -4,  -6,
                 -6,  -4,   1,   1,   1,   1,  -4,  -6,
                  0,   0,   0,   0,   0,   0,   0,   0
                 ],
    chess.KNIGHT: [
        -8, -12,  -8,  -8,  -8,  -8, -12,  -8,
        -8,   0,   1,   2,   2,   1,   0,  -8,
        -8,   0,   4,   4,   4,   4,   0,  -8,
        -8,   0,   4,   8,   8,   4,   0,  -8,
        -8,   0,   4,   8,   8,   4,   0,  -8,
        -8,   0,   4,   4,   4,   4,   0,  -8,
        -8,   0,   0,   0,   0,   0,   0,  -8,
        -8,  -8,  -8,  -8,  -8,  -8,  -8,  -8,
    ],
    chess.BISHOP: [
        -4,  -4, -12,  -4,  -4, -12,  -4,  -4,
        -4,   2,   1,   1,   1,   1,   2,  -4,
        -4,   1,   2,   4,   4,   2,   1,  -4,
        -4,   0,   4,   6,   6,   4,   0,  -4,
        -4,   0,   4,   6,   6,   4,   0,  -4,
        -4,   0,   2,   4,   4,   2,   0,  -4,
        -4,   0,   0,   0,   0,   0,   0,  -4,
        -4,  -4,  -4,  -4,  -4,  -4,  -4,  -4,
    ],
    chess.ROOK: [
         0,   0,   0,   2,   2,   0,   0,   0,
        -5,   0,   0,   0,   0,   0,   0,  -5,
        -5,   0,   0,   0,   0,   0,   0,  -5,
        -5,   0,   0,   0,   0,   0,   0,  -5,
        -5,   0,   0,   0,   0,   0,   0,  -5,
        -5,   0,   0,   0,   0,   0,   0,  -5,
        20,  20,  20,  20,  20,  20,  20,  20,
        5,   5,   5,   5,   5,   5,   5,   5,
    ],
    chess.QUEEN: [
        -5, -5, -5, -5, -5, -5, -5, -5,
        0, 0, 1, 1, 1, 1, 0, 0,
        0, 0, 1, 2, 2, 1, 0, 0,
        0, 0, 2, 3, 3, 2, 0, 0,
        0, 0, 2, 3, 3, 2, 0, 0,
        0, 0, 1, 2, 2, 1, 0, 0,
        0, 0, 1, 1, 1, 1, 0, 0,
        0, 0, 0, 0, 0, 0, 0, 0,
    ],
    chess.KING: [
        40,  50,  30,  10,  10,  30,  50,  40,
        30,  40,  20,   0,   0,  20,  40,  30,
        10,  20,   0, -20, -20,   0,  20,  10,
        0,  10, -10, -30, -30, -10,  10,   0,
        -10,   0, -20, -40, -40, -20,   0, -10,
        -20, -10, -30, -50, -50, -30, -10, -20,
        -30, -20, -40, -60, -60, -40, -20, -30,
        -40, -30, -50, -70, -70, -50, -30, -40,
    ]
}

piece_square_tables_endgame = piece_square_tables.copy()
piece_square_tables_endgame[chess.KING] = [
  -72, -48, -36, -24, -24, -36, -48, -72,
  -48, -24, -12,   0,   0, -12, -24, -48,
  -36, -12,   0,  12,  12,   0, -12, -36,
  -24,   0,  12,  24,  24,  12,   0, -24,
  -24,   0,  12,  24,  24,  12,   0, -24,
  -36, -12,   0,  12,  12,   0, -12, -36,
  -48, -24, -12,   0,   0, -12, -24, -48,
  -72, -48, -36, -24, -24, -36, -48, -72,
]

additional_piece_square_tables = {
    "weak_pawn": [
        0,   0,   0,   0,   0,   0,   0,   0,
        -10, -12, -14, -16, -16, -14, -12, -8,
        -10, -12, -14, -16, -16, -14, -12, -8,
        -10, -12, -14, -16, -16, -14, -12, -10,
        -10, -12, -14, -16, -16, -14, -12, -10,
        -10, -12, -14, -16, -16, -14, -12, -10,
        -10, -12, -14, -16, -16, -14, -12, -10,
        0,   0,   0,   0,   0,   0,   0,   0
    ],
    "passed_pawn": [
        0,   0,   0,   0,   0,   0,   0,   0,
        20,  20,  20,  20,  20,  20,  20,  20,
        20,  20,  20,  20,  20,  20,  20,  20,
        40,  40,  40,  40,  40,  40,  40,  40,
        60,  60,  60,  60,  60,  60,  60,  60,
        80,  80,  80,  80,  80,  80,  80,  80,
        100, 100, 100, 100, 100, 100, 100, 100,
        0,   0,   0,   0,   0,   0,   0,   0
    ]
}