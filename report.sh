#!/bin/bash

echo ":: GO"
go test ./... | grep FAIL: | grep -v 30-range | sort

echo ":: Ruby"
ruby fnmatch_test.rb | grep Failure | sort
