#!/usr/bin/env python3
import argparse
import utils.uci

parser = argparse.ArgumentParser()
parser.add_argument("--depth", default=4, help="provide an integer (default: 3)")
args = parser.parse_args()

utils.uci.Uci(args.depth).communicate()
