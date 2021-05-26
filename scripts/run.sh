#!/bin/bash

SCRIPT=`realpath $0`
SCRIPTPATH=`dirname $SCRIPT`
cd $SCRIPTPATH/..

OK=0
while [ $OK -eq 0 ]
do
    go run cmd/kickstart.go -watch "*.go *.html *.css *.js"
    OK=$? # 1 = a file changed, 2 = build failed

    # if build failure, do not loop
    if [ $OK -eq 2 ]
    then
        exit 2
    fi
done
