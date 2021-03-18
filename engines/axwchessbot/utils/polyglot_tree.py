#!/usr/bin/env python3

"""Print a Polyglot opening book in tree form."""

import argparse

from typing import Set

import chess
import chess.polyglot

DEFAULT_RETURN = {"max_move_depth": 0}


def print_tree(args: argparse.Namespace, tree: list, level: int = 0) -> None:
    if level >= args.depth:
        return

    for node in tree:
        entry = node["entry"]

        print(
            "{}├─ \033[1m{}\033[0m (weight: {}, learn: {}, depth: {})".format(
                "|  " * level,
                node["move_san"],
                entry.weight,
                entry.learn,
                node["max_depth"],
            )
        )

        print_tree(args, node["children"], level + 1)


def gen_tree(args: argparse.Namespace, visited: Set[int], level: int = 0) -> list:
    result = []

    zobrist_hash = chess.polyglot.zobrist_hash(args.board)
    if zobrist_hash in visited:
        return []

    visited.add(zobrist_hash)
    for entry in args.book.find_all(zobrist_hash):
        current_result = {
            "move": entry.move(),
            "entry": entry,
            "move_san": args.board.san(entry.move()),
            "children": [],
            "max_depth": level,
        }

        args.board.push(entry.move())
        current_result["children"] = gen_tree(args, visited, level + 1)
        args.board.pop()

        for child in current_result["children"]:
            if child["max_depth"] > current_result["max_depth"]:
                current_result["max_depth"] = child["max_depth"]

        result.append(current_result)

    return result


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("book", type=chess.polyglot.open_reader)
    parser.add_argument("--depth", type=int, default=5)
    parser.add_argument("--fen", type=chess.Board, default=chess.Board(), dest="board")
    args = parser.parse_args()
    tree = gen_tree(args, visited=set())
    print_tree(args, tree)