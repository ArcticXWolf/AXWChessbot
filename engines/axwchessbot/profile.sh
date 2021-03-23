#!/bin/bash

python -m cProfile -o tests/results/profiling main.py --bench
gprof2dot -f pstats tests/results/profiling | dot -Tpng -o tests/results/profile.png