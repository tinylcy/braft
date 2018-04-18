#!/bin/bash
echo -n "" > data

for k in $( seq 1 1)
do
    go test
done

# python average.py
python committed_elapsed_time.py

rm *_data
