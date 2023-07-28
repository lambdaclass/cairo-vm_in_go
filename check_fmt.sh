#!/usr/bin/env bash

files=$(gofmt -l pkg)
if [[ $files ]]; then
    echo -e "Some files are not correctly formatted:\n${files}"
    exit 1
else 
    exit 0
fi
