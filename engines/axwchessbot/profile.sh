#!/bin/bash
dt=$(date '+%Y-%m-%d-%H-%M-%S')
python -m cProfile -o "tests/results/profiling-$dt" main.py --bench
gprof2dot -f pstats "tests/results/profiling-$dt" | dot -Tpng -o "tests/results/profile-$dt.png"