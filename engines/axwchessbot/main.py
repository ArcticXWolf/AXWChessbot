#!/usr/bin/env python3
import argparse
from utils import uci
from utils import bench

parser = argparse.ArgumentParser(description="Start the AXWChessBot chess engine")
parser.add_argument(
    "--uci",
    dest="mode",
    action="store_const",
    const="uci",
    help="Start in uci mode (default)",
)
parser.add_argument(
    "--bench",
    dest="mode",
    action="store_const",
    const="bench",
    help="Start in bench mode",
)
parser.set_defaults(mode="uci")
args = parser.parse_args()

if args.mode == "uci":
    uci.Uci().communicate()
elif args.mode == "bench":
    bench.Benchmark(4, 6, 10, True).run()
