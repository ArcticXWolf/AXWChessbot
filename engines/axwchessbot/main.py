#!/usr/bin/env python3
import argparse
from selfplay import play_self
from bratko_kopec import test_engine
import uci

parser = argparse.ArgumentParser()
parser.add_argument("--depth", default=3, help="provide an integer (default: 3)")
parser.add_argument("--selfplay", default=False, type=bool)
parser.add_argument("--test", default=False, type=bool)
args = parser.parse_args()

if args.selfplay:
    play_self()
elif args.test:
    test_engine()
else:
    uci.Uci(args.depth).communicate()
