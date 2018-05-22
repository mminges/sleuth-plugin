#!/bin/bash

GOOS=darwin go build -o sleuth-plugin-osx
#GOOS=linux go build -o sleuth-plugin-linux
#GOOS=windows GOARCH=amd64 go build -o sleuth-plugin.exe

if [ $? != 0 ]; then
   printf "Error when executing compile\n"
   exit 1
fi
cf uninstall-plugin sleuth
cf install-plugin -f ./sleuth-plugin-osx
