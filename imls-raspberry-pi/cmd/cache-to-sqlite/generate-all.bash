#!/bin/bash
DEST=./output
mkdir -p $DEST
# Clear old stuff
rm -rf $DEST/*.sqlite
go build
EXE=./cache-to-sqlite

$EXE --fcfs_seq_id GA0058-003 --device_tag rpi01 --tzoffset -5 --dest $DEST
$EXE --fcfs_seq_id GA0058-003 --device_tag rpi02 --tzoffset -5 --dest $DEST  
$EXE --fcfs_seq_id GA0014-002 --device_tag behind-desk --tzoffset -5  --dest $DEST 
$EXE --fcfs_seq_id GA0027-004 --device_tag in-ops --tzoffset -5  --dest $DEST 

$EXE --fcfs_seq_id MA0352-002 --device_tag securitydesk --tzoffset -5  --dest $DEST
$EXE --fcfs_seq_id MA0269-002 --device_tag mnspear-1 --tzoffset -5  --dest $DEST 

$EXE --fcfs_seq_id KY0069-002 --device_tag richmond0 --tzoffset -5  --dest $DEST 
$EXE --fcfs_seq_id KY0069-003 --device_tag berea1 --tzoffset -5  --dest $DEST 

$EXE --fcfs_seq_id AR0012-004 --device_tag dardanelle-pac-desk --tzoffset -6  --dest $DEST

