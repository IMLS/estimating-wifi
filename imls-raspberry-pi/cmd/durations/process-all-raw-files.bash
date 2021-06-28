#!/bin/bash
OUT=./output
mkdir -p $OUT
rm -rf $OUT/*.sqlite
go build
EXE=./durations
SRC=../cache-to-sqlite/output/*.sqlite
for i in $SRC ; do
    echo "Processing $i"
    $EXE --config config.yaml --dest $OUT --sqlite $i --swap true
done

