#!/bin/bash

OUT=./output
mkdir -p $OUT
rm -rf $OUT/*
go build
EXE=./waterfall
for i in ../durations/output/*.sqlite
do
    echo "Processing $i"
    $EXE --config config.sqlite --type sqlite --data $OUT --src $i

done
