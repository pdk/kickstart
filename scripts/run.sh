#!/bin/bash

SCRIPT=`realpath $0`
SCRIPTPATH=`dirname $SCRIPT`
cd $SCRIPTPATH/..

OK=0
while [ $OK -eq 0 ]
do
    go run cmd/kickstart.go -watch "*.go *.html *.css *.js"
    OK=$?

    # if build failure, exit
    if [ $OK -eq 2 ]
    then
        exit 2
    fi
done
