#!/bin/bash

while read -u 10 line; do
  if [ -n "$line" ]; then
    echo "go get -u $line"
    go get -u $line
  fi
done 10<Goopfile
