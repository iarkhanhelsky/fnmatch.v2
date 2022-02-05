#!/bin/bash

GO_REPORT=$(go test ./... | grep FAIL: | sort | grep -o -e 'test_[^ ]*')
RB_REPORT=$(ruby fnmatch_test.rb | grep Failure | sort | grep -o -e 'test_[^(]*')

echo "diff RB GO"
diff -U 8 <(echo "$RB_REPORT") <(echo "$GO_REPORT") | sed 's|--- /dev/[^ ]*|--- ruby  |' | sed 's|+++ /dev/[^ ]*|+++ golang|'
