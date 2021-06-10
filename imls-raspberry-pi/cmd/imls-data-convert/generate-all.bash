#!/bin/bash

go run main.go --fcfs_seq_id GA0058-003 --device_tag rpi01 --tzoffset -5 $1
go run main.go --fcfs_seq_id GA0058-003 --device_tag rpi02 --tzoffset -5 $1
go run main.go --fcfs_seq_id GA0058-003 --device_tag rpi03 --tzoffset -5 $1
go run main.go --fcfs_seq_id GA0014-002 --device_tag behind-desk --tzoffset -5 $1
go run main.go --fcfs_seq_id GA0027-004 --device_tag in-ops --tzoffset -5 $1

go run main.go --fcfs_seq_id MA0352-002 --device_tag securitydesk --tzoffset -5 $1

go run main.go --fcfs_seq_id KY0069-002 --device_tag richmond0 --tzoffset -5 $1
go run main.go --fcfs_seq_id KY0069-003 --device_tag berea1 --tzoffset -5 $1

go run main.go --fcfs_seq_id AR0012-004 --device_tag dardanelle-pac-desk --tzoffset -6 $1
