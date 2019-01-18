#!/bin/sh

if [ $# = 1 ] ; then
    echo "shutting down previous workers..."
    killall worker
    echo "launching $1 workers..."
    [ -e address.config ] && rm address.config
    [ ! -e worker ] && go build worker.go
    left=$1
    while [ $left -gt 0 ] ; do 
        ./worker &
        left=`expr $left - 1`
    done
    echo "done"
    ps | grep worker
else
    echo "Usage: ./launch_workers <number of workers>"
fi
