#!/usr/bin/env python3
import argparse
from utils import uci

parser = argparse.ArgumentParser()
parser.add_argument("--depth", default=4, help="provide an integer (default: 3)")
args = parser.parse_args()

uci.Uci(args.depth).communicate()
