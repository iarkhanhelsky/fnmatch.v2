#!/bin/bash

echo ":: GO"
go test ./... | grep FAIL: | sort

echo ":: Ruby"
ruby fnmatch_test.rb | grep Failure | sort
